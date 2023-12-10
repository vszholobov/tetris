package main

import (
	"fmt"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

const fieldWidth = 12

type Field struct {
	val          *big.Int
	currentPiece *Piece
}

func main() {
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
	fmt.Println(fieldVal)

	var wg sync.WaitGroup
	extField := Field{val: fieldVal}
	wg.Add(1)
	go func(field *Field) {
		piece := SelectNextPiece(field)
		for {
			field.printField(piece)
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

func (field Field) printField(piece Piece) {
	newField := big.NewInt(0).Set(field.val)
	newShape := big.NewInt(0).Set(piece.GetVal())
	newField.Or(newField, newShape)
	s := fmt.Sprintf("%b", newField)
	// TODO: one PrintLn to reduce output calls
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
	piece := MakePiece(field, TShape)
	field.currentPiece = &piece
	return piece
}
