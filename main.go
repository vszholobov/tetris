package main

import (
	"fmt"
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

	var wg sync.WaitGroup
	wg.Add(1)
	keyboardSendChannel := make(chan rune)
	keyboardChannel := initInputChannel()
	defer keyboardChannel.Close()
	// input
	go func(keyboardChannel *tty.TTY, keyboardSendChannel chan<- rune) {
		for {
			r, err := keyboardChannel.ReadRune()
			if err != nil {
				log.Fatal(err)
			}
			//fmt.Println("Key press => " + string(r))
			keyboardSendChannel <- r
		}
	}(keyboardChannel, keyboardSendChannel)
	// game
	go func(gameField *field.Field, keyboardInputChannel chan rune, wg *sync.WaitGroup) {
		piece := SelectNextPiece(gameField)
		for {
			inputControl(keyboardInputChannel, gameField)

			if !piece.MoveDown() {
				gameField.Val.Or(gameField.Val, piece.GetVal())
				piece = SelectNextPiece(gameField)
				if !gameField.CurrentPiece.CanMoveDown() {
					field.CallClear()
					fmt.Println("Game over. Stats:")
					fmt.Printf("Score: %d | Speed: %d | Lines Cleand: %d\n", *gameField.Score, gameField.GetSpeed(), *gameField.CleanCount)
					wg.Done()
					break
				}
				gameField.Clean()
			}
		}
	}(&extField, keyboardSendChannel, &wg)
	wg.Wait()
}

func initInputChannel() *tty.TTY {
	keyPressedChannel, err := tty.Open()
	if err != nil {
		log.Fatal(err)
	}
	return keyPressedChannel
}

func inputControl(
	keyboardInputChannel chan rune,
	gameField *field.Field,
) {
	timeout := time.After(time.Second / 4 / time.Duration(gameField.GetSpeed()))
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
	pieceTypeRnd := rand.Intn(7)
	var pieceType field.PieceType
	if pieceTypeRnd == 0 {
		pieceType = field.IShape
	} else if pieceTypeRnd == 1 {
		pieceType = field.RightLShape
	} else if pieceTypeRnd == 2 {
		pieceType = field.TShape
	} else if pieceTypeRnd == 3 {
		pieceType = field.ZigZagRight
	} else if pieceTypeRnd == 4 {
		pieceType = field.ZigZagLeft
	} else if pieceTypeRnd == 5 {
		pieceType = field.SquareShape
	} else if pieceTypeRnd == 6 {
		pieceType = field.LeftLShape
	}
	piece := field.MakePiece(gameField, pieceType)
	gameField.CurrentPiece = &piece
	return &piece
}
