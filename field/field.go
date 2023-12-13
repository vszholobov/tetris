package field

import (
	"fmt"
	"math/big"
)

const FieldWidth = 12
const FieldHeight = 21

type Field struct {
	Val          *big.Int
	CurrentPiece *Piece
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
	return Field{Val: fieldVal}
}

func (gameField *Field) String() string {
	newField := big.NewInt(0).Set(gameField.Val)
	newShape := big.NewInt(0).Set(gameField.CurrentPiece.GetVal())
	newField.Or(newField, newShape)
	return fmt.Sprintf("%b", newField)
}

func (gameField *Field) Clean() {
	restField := big.NewInt(0)

	fullLine, _ := big.NewInt(0).SetString("111111111111", 2)
	emptyLine, _ := big.NewInt(0).SetString("100000000001", 2)
	for i := 0; i < FieldHeight; i++ {
		curRange := uint(i * FieldWidth)
		lineMask := big.NewInt(0).Lsh(fullLine, curRange)
		lineIsFilled := big.NewInt(0).And(lineMask, gameField.Val).Cmp(lineMask) == 0

		if lineIsFilled {
			// add empy line to end of field
			restField.Lsh(restField, FieldWidth)
			restField.Or(restField, emptyLine)
			// TODO: score += награда за соженную линию
		} else {
			// add current line to start of field
			lineMask.And(lineMask, gameField.Val)
			restField.Or(lineMask, restField)
		}
	}
	// 22 lines. One redundant line for correct or concatenation.
	// So shift to the right by the length of the field after concatenation to remove redundant empty line
	gameField.Val.SetString(
		"111111111111"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000"+
			"000000000000", 2)
	gameField.Val.Or(gameField.Val, restField)
	gameField.Val.Rsh(gameField.Val, FieldWidth)
}

func (gameField *Field) Intersects(pieceVal *big.Int) bool {
	newField := big.NewInt(0).Set(gameField.Val)
	newShape := big.NewInt(0).Set(pieceVal)
	return newField.And(newField, newShape).Cmp(big.NewInt(0)) != 0
}

func CopyBigInt(val *big.Int) *big.Int {
	return big.NewInt(0).Set(val)
}

func PrintField(field *Field) {
	s := field.String()
	for i := 20; i >= 0; i-- {
		fmt.Println(s[i*12 : i*12+12])
	}
	fmt.Println()
}
