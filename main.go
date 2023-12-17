package main

import (
	"fmt"
	"github.com/mattn/go-tty"
	"log"
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
		//gameField.CurrentPiece = SelectRandomPiece(gameField)
		//nextPiece := SelectRandomPiece(gameField)
		for {
			inputControl(keyboardInputChannel, gameField)

			if !gameField.CurrentPiece.MoveDown() {
				gameField.Val.Or(gameField.Val, gameField.CurrentPiece.GetVal())
				gameField.SelectNextPiece()
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
