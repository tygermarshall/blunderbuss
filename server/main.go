package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/tygermarshall/blunderbuss/shared"
	"github.com/tygermarshall/blunderbuss/shared/board"
)

var gameStore = NewGameStore()

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // allow all origins; restrict in production
	},
}

func handleWebSocket(c *gin.Context) {
	conn, err := wsUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("websocket upgrade: %v", err)
		return
	}
	defer conn.Close()
	for {
		mt, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("websocket read: %v", err)
			}
			break
		}
		if err := conn.WriteMessage(mt, msg); err != nil {
			log.Printf("websocket write: %v", err)
			break
		}
	}
}

func movePiece(c *gin.Context) {
	id := c.Param("id")
	//var req moveRequest
	//if err := c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: " + err.Error()})
	//	return
	//}
	//err := gameStore.Move(id, req.From, req.To)
	//todo move this seeing if works
	//end todo
	//if err == ErrGameNotFound {
	//	c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
	//	return
	//}
	//if err != nil {
	//	c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	//	return
	//	}

	gameStore.Move(id, board.Coordinate{X: 6, Y: 0}, board.Coordinate{X: 4, Y: 0})
	g, ok := gameStore.Get(id)

	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get game"})
	}
	g.Board.MoveCount += 1
	log.Printf("move count is: %d", g.Board.MoveCount)
	body := shared.CreateGameReponse{GameId: id, Board: g.Board}

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize game"})
		return
	}

	c.Data(http.StatusCreated, "application/json", buf.Bytes())
}

func startNewGame(c *gin.Context) {
	id, err := gameStore.Create()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}
	g, ok := gameStore.Get(id)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create game"})
		return
	}

	body := shared.CreateGameReponse{GameId: id, Board: g.Board}

	shared.PrintBoard(&body.Board)

	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(body); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to serialize game"})
		return
	}
	c.Data(http.StatusCreated, "application/json", buf.Bytes())
}

func getGame(c *gin.Context) {
	id := c.Param("id")
	g, ok := gameStore.Get(id)
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "game not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"gameId":     id,
		"turnNumber": g.TurnNumber,
		"board":      g.Board.Squares,
	})
}

type moveRequest struct {
	From board.Coordinate `json:"from"`
	To   board.Coordinate `json:"to"`
}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong\n",
		})
	})
	router.POST("/games", startNewGame)
	router.GET("/games/:id", getGame)
	router.POST("/games/:id/move", movePiece)
	router.GET("/ws", handleWebSocket)
	router.Run() // listens on 0.0.0.0:8080 by default
}

func hold() {
	log.SetPrefix("Server: ")
	log.SetFlags(0)
	var gameBoard = board.CreateDefaultBoard()
	shared.PrintBoard(&gameBoard)

	shared.PrettyPrint(gameBoard)
	gameBoard, err := gameBoard.MovePiece(board.Coordinate{X: 6, Y: 0}, board.Coordinate{X: 4, Y: 0})
	if err != nil {
		log.Fatal(err)
	}
	shared.PrettyPrint(gameBoard)

	gameBoard, err = gameBoard.MovePiece(board.Coordinate{X: 1, Y: 4}, board.Coordinate{X: 0, Y: 4})
	if err != nil {
		log.Fatal(err)
	}
	shared.PrettyPrint(gameBoard)

}
