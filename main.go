package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

func main() {
	extField := MakeField()
	var wg sync.WaitGroup
	wg.Add(1)
	go func(field *Field) {
		piece := SelectNextPiece(field)
		for {
			printField(field)
			moveType := rand.Intn(4)
			switch moveType {
			case 0:
				piece.moveLeft()
			case 1:
				piece.moveRight()
			case 2:
				piece.moveDown()
			case 3:
				piece.Rotate(Left)
			}
			if !piece.moveDown() {
				field.val.Or(field.val, piece.GetVal())
				piece = SelectNextPiece(field)
			}
			time.Sleep(time.Second / 4)
		}
	}(&extField)
	wg.Wait()
}

func printField(field *Field) {
	s := field.String()
	for i := 20; i >= 0; i-- {
		fmt.Println(s[i*12 : i*12+12])
	}
	fmt.Println()
}

func (field Field) intersects(pieceVal *big.Int) bool {
	newField := big.NewInt(0).Set(field.val)
	newShape := big.NewInt(0).Set(pieceVal)
	return newField.And(newField, newShape).Cmp(big.NewInt(0)) != 0
}

func SelectNextPiece(field *Field) Piece {
	pieceTypeRnd := rand.Intn(5)
	var pieceType PieceType
	if pieceTypeRnd == 0 {
		pieceType = TShape
	} else if pieceTypeRnd == 1 {
		pieceType = ZigZagLeft
	} else if pieceTypeRnd == 2 {
		pieceType = ZigZagRight
	} else if pieceTypeRnd == 3 {
		pieceType = IShape
	} else if pieceTypeRnd == 4 {
		pieceType = LShape
	}
	piece := MakePiece(field, pieceType)
	field.currentPiece = &piece
	return piece
}
