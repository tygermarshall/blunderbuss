package main

import (
	"github.com/gin-gonic/gin"
	"github.com/tygermarshall/blunderbuss/shared/board"
	"log"
)

func startNewGame(c *gin.Context) {

}

func main() {
	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong\n",
		})
	})
	router.Run() // listens on 0.0.0.0:8080 by default
}

func hold() {
	log.SetPrefix("Server: ")
	log.SetFlags(0)
	var gameBoard = board.CreateDefaultBoard()

	PrettyPrint(gameBoard)
	gameBoard, err := gameBoard.MovePiece(board.Coordinate{X: 6, Y: 0}, board.Coordinate{X: 4, Y: 0})
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(gameBoard)

	gameBoard, err = gameBoard.MovePiece(board.Coordinate{X: 1, Y: 4}, board.Coordinate{X: 0, Y: 4})
	if err != nil {
		log.Fatal(err)
	}
	PrettyPrint(gameBoard)

}
