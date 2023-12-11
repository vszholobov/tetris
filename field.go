package main

import (
	"fmt"
	"math/big"
)

const FieldWidth = 12

type Field struct {
	val          *big.Int
	currentPiece *Piece
}

func MakeField() Field {
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
	return Field{val: fieldVal}
}

func (field Field) String() string {
	newField := big.NewInt(0).Set(field.val)
	newShape := big.NewInt(0).Set(field.currentPiece.GetVal())
	newField.Or(newField, newShape)
	return fmt.Sprintf("%b", newField)
}
