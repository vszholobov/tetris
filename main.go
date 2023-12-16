package main

import (
	"github.com/mattn/go-tty"
	"log"
	"math/big"
	"math/rand"
	"sync"
	"tetris/field"
	"time"
)

func main() {
	//fmt.Println("012345678910")
	//fmt.Println("012345678910")
	//fmt.Println("012345678910")
	//fmt.Printf("\033[20D\r%s", "ABC")
	//fmt.Printf("\033[1A\r%s", "ABC")
	//fmt.Printf("\033[2A[3C%s", "ABC")
	//fmt.Println("DDD")
	//for i := 10; i >= 0; i-- {
	//	//fmt.Printf("\033[2K\r%d", i)
	//	//fmt.Printf("\033[2K\r%d\n", i)
	//	//fmt.Printf("\033[2K\r%d", 999)
	//	fmt.Printf("\\033[10D\r%d", 999)
	//	time.Sleep(1 * time.Second)
	//}
	//fmt.Println()
	field.InitClear()
	field.CallClear()
	fieldVal, _ := big.NewInt(0).SetString(
		"111111111111"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001"+
			"100000000001", 2)
	extField := field.MakeField(fieldVal)

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
	pieceTypeRnd := rand.Intn(1)
	// TODO: square type
	var pieceType field.PieceType
	if pieceTypeRnd == 0 {
		pieceType = field.IShape
	}
	//} else if pieceTypeRnd == 1 {
	//	pieceType = field.LShape
	//} else if pieceTypeRnd == 2 {
	//	pieceType = field.TShape
	//} else if pieceTypeRnd == 3 {
	//	pieceType = field.ZigZagRight
	//} else if pieceTypeRnd == 4 {
	//	pieceType = field.ZigZagLeft
	//}
	piece := field.MakePiece(gameField, pieceType)
	gameField.CurrentPiece = &piece
	return &piece
}
