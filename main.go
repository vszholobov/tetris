package main

import (
	"math/rand"
	"sync"
	"tetris/field"
	"time"
)

func main() {
	extField := field.MakeField()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(gameField *field.Field) {
		piece := SelectNextPiece(gameField)
		for {
			field.PrintField(gameField)
			moveType := rand.Intn(4)
			switch moveType {
			case 0:
				piece.MoveLeft()
			case 1:
				piece.MoveRight()
			case 2:
				piece.MoveDown()
			case 3:
				piece.Rotate(field.Left)
			}
			if !piece.MoveDown() {
				gameField.Val.Or(gameField.Val, piece.GetVal())
				piece = SelectNextPiece(gameField)
			}
			gameField.Clean()

			time.Sleep(time.Second / 10)
		}
	}(&extField)
	wg.Wait()
}

func SelectNextPiece(gameField *field.Field) field.Piece {
	pieceTypeRnd := rand.Intn(5)
	var pieceType field.PieceType
	if pieceTypeRnd == 0 {
		pieceType = field.TShape
	} else if pieceTypeRnd == 1 {
		pieceType = field.ZigZagLeft
	} else if pieceTypeRnd == 2 {
		pieceType = field.ZigZagRight
	} else if pieceTypeRnd == 3 {
		pieceType = field.IShape
	} else if pieceTypeRnd == 4 {
		pieceType = field.LShape
	}
	piece := field.MakePiece(gameField, pieceType)
	gameField.CurrentPiece = &piece
	return piece
}
