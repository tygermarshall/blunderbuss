package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/gorilla/websocket"
	"github.com/tygermarshall/blunderbuss/shared"
	"github.com/tygermarshall/blunderbuss/shared/board"
)

// WebSocket server URL; ensure the backend is running on this address.
const wsServerURL = "ws://localhost:8080/ws"

// Server base URL for REST API (e.g. create game).
const serverBaseURL = "http://localhost:8080"

// Model represents the application state
type model struct {
	choices       []string
	cursor        int
	selected      map[int]struct{}
	width         int
	heigh         int
	conn          *websocket.Conn
	connErr       error
	lastMessage   string
	status        string // "connecting", "connected", "disconnected", "error"
	gameId        string // last created game ID
	createGameErr error
	Board         board.Board
}

func debugLog(resp *http.Response) {
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		fmt.Println("fatal:", err)
	}
	log.Printf("%T\n", dump)

}

// WebSocket async messages
type wsConnectedMsg struct{ conn *websocket.Conn }
type wsErrorMsg struct{ err error }
type wsMessageMsg struct{ data []byte }

type gameCreatedMsg struct {
	GameId string
	Board  board.Board
}
type gameCreateErrMsg struct{ Err error }

// connectCmd dials the backend and returns wsConnectedMsg or wsErrorMsg.
func connectCmd() tea.Msg {
	conn, _, err := websocket.DefaultDialer.Dial(wsServerURL, nil)
	if err != nil {
		return wsErrorMsg{err}
	}
	return wsConnectedMsg{conn}
}

// readNextCmd reads one message from conn and returns wsMessageMsg or wsErrorMsg.
func readNextCmd(conn *websocket.Conn) tea.Cmd {
	return func() tea.Msg {
		_, data, err := conn.ReadMessage()
		if err != nil {
			return wsErrorMsg{err}
		}
		return wsMessageMsg{data}
	}
}

//func movePieceCmd(m model) tea.Msg {
//	resp, err := http.Post(serverBaseURL+"/games/"+m.gameId+"/move", "application/json", nil)
//	if err != nil {
//		return gameCreateErrMsg{Err: err}
//	}
//	defer resp.Body.Close()
//	debugLog(resp)
//	if resp.StatusCode != http.StatusCreated {
//		return gameCreateErrMsg{Err: fmt.Errorf("create game: status %d", resp.StatusCode)}
//	}
//	var out shared.CreateGameReponse
//
//	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
//		return gameCreateErrMsg{Err: err}
//	}
//	return gameCreatedMsg{GameId: out.GameId, Board: out.Board}
//}

func movePieceCmd(m model) tea.Cmd {
	return func() tea.Msg {
		resp, err := http.Post(serverBaseURL+"/games/"+m.gameId+"/move", "application/json", nil)
		if err != nil {
			return gameCreateErrMsg{Err: err}
		}
		defer resp.Body.Close()
		debugLog(resp)
		if resp.StatusCode != http.StatusCreated {
			return gameCreateErrMsg{Err: fmt.Errorf("create game: status %d", resp.StatusCode)}
		}
		var out shared.CreateGameReponse

		if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
			return gameCreateErrMsg{Err: err}
		}
		return gameCreatedMsg{GameId: out.GameId, Board: out.Board}
	}
}

// createGameCmd POSTs to /games and returns gameCreatedMsg or gameCreateErrMsg.
func createGameCmd() tea.Msg {
	resp, err := http.Post(serverBaseURL+"/games", "application/json", nil)
	if err != nil {
		return gameCreateErrMsg{Err: err}
	}
	defer resp.Body.Close()
	debugLog(resp)
	if resp.StatusCode != http.StatusCreated {
		return gameCreateErrMsg{Err: fmt.Errorf("create game: status %d", resp.StatusCode)}
	}
	var out shared.CreateGameReponse

	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return gameCreateErrMsg{Err: err}
	}
	return gameCreatedMsg{GameId: out.GameId, Board: out.Board}
}

// Init returns the initial command for the application to run
func (m model) Init() tea.Cmd {
	return connectCmd
}

// Update handles events and updates the model
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			if m.conn != nil {
				_ = m.conn.Close()
			}
			return m, tea.Quit
		case "p":
			if m.conn != nil {
				_ = m.conn.WriteMessage(websocket.TextMessage, []byte("ping"))
			}
			return m, nil
		case "g":
			return m, createGameCmd
		case "m":
			return m, movePieceCmd(m)
		}
	case gameCreatedMsg:
		m.gameId = msg.GameId
		m.Board = msg.Board
		m.createGameErr = nil
		log.Printf("game created msg")
		return m, nil
	case gameCreateErrMsg:
		m.createGameErr = msg.Err
		return m, nil
	case wsConnectedMsg:
		m.conn = msg.conn
		m.status = "connected"
		m.connErr = nil
		return m, readNextCmd(m.conn)
	case wsErrorMsg:
		m.connErr = msg.err
		if m.conn != nil {
			_ = m.conn.Close()
			m.conn = nil
		}
		m.status = "error"
		return m, nil
	case wsMessageMsg:
		m.lastMessage = string(msg.data)
		if m.conn != nil {
			return m, readNextCmd(m.conn)
		}
		return m, nil
	}
	return m, nil
}

// View renders the UI based on the model's state
func (m model) View() string {
	var output strings.Builder
	if m.gameId != "" {
		b := m.Board
		shared.CreatePrettyPrint(b, &output)
	} else {
		output.WriteString("no game started yet\n")
	}

	moveCountString := fmt.Sprintf("Move count: %d", m.Board.MoveCount)
	output.WriteString(moveCountString)

	output.WriteString("\n")

	if m.lastMessage != "" {
		output.WriteString(" Last: " + m.lastMessage + "\n")
	}
	if m.connErr != nil {
		output.WriteString(" (" + m.connErr.Error() + ")")
	}
	if m.gameId != "" {
		output.WriteString(" Game: " + m.gameId + "\n")
	}
	if m.createGameErr != nil {
		output.WriteString(" Create game: " + m.createGameErr.Error() + "\n")
	}
	output.WriteString(" [g] create game  [p] send ping  [q] quit [m] move \n")
	return output.String()
}

func main() {
	f, err := tea.LogToFile("debuglog.log", "debug")
	if err != nil {
		fmt.Println("fatal:", err)
	}
	defer f.Close()
	p := tea.NewProgram(model{status: "connecting"})
	if _, err := p.Run(); err != nil { // Run the program
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
