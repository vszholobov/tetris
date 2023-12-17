package main

import (
	"github.com/mattn/go-tty"
	"log"
	"math/rand"
	"sync"
	"tetris/field"
	"time"
)

func main() {
	field.InitClear()
	field.CallClear()
	extField := field.MakeDefaultField()

	keyboardInputChannel := make(chan rune)
	// input
	go func(keyboardInputChannel chan<- rune) {
		keyPressedChannel, err := tty.Open()
		if err != nil {
			log.Fatal(err)
		}
		defer keyPressedChannel.Close()

		for {
			r, err := keyPressedChannel.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("Key press => " + string(r))
			keyboardInputChannel <- r
		}
	}(keyboardInputChannel)
	// game
	go func(gameField *field.Field, keyboardInputChannel chan rune) {
		piece := SelectNextPiece(gameField)
		for {
			inputControl(keyboardInputChannel, gameField)

			if !piece.MoveDown() {
				gameField.Val.Or(gameField.Val, piece.GetVal())
				piece = SelectNextPiece(gameField)
				gameField.Clean()
			}
		}
	}(&extField, keyboardInputChannel)
	var wg sync.WaitGroup
	wg.Add(1)
	wg.Wait()
}

func inputControl(
	keyboardInputChannel chan rune,
	gameField *field.Field,
) {
	timeout := time.After(time.Second / time.Duration(*gameField.CleanCount/2+2))
	for {
		field.PrintField(gameField)
		select {
		case moveType := <-keyboardInputChannel:
			switch moveType {
			case 100:
				// d
				gameField.CurrentPiece.MoveLeft()
			case 97:
				// a
				gameField.CurrentPiece.MoveRight()
			case 115:
				// s
				gameField.CurrentPiece.MoveDown()
			case 113:
				// q
				gameField.CurrentPiece.Rotate(field.Left)
			case 101:
				// e
				gameField.CurrentPiece.Rotate(field.Right)
			}
		case <-timeout:
			return
		}
	}
}

func SelectNextPiece(gameField *field.Field) *field.Piece {
	pieceTypeRnd := rand.Intn(6)
	var pieceType field.PieceType
	if pieceTypeRnd == 0 {
		pieceType = field.IShape
	} else if pieceTypeRnd == 1 {
		pieceType = field.LShape
	} else if pieceTypeRnd == 2 {
		pieceType = field.TShape
	} else if pieceTypeRnd == 3 {
		pieceType = field.ZigZagRight
	} else if pieceTypeRnd == 4 {
		pieceType = field.ZigZagLeft
	} else if pieceTypeRnd == 5 {
		pieceType = field.SquareShape
	}
	piece := field.MakePiece(gameField, pieceType)
	gameField.CurrentPiece = &piece
	return &piece
}
