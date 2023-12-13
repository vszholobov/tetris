package field

import (
	"fmt"
	"math/big"
)

const FieldWidth = 12

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

// ! Нельзя доходить до последней линии. Она всегда заполнена

// fullLine  := 111...111
// emptyLine := 100...001
// range := fieldWidth * i
// line := ((fullLine << range) & field) >> range
// line is filled -> (fullLine << range) == (fullLine << range) & field

//restField := 0
// if line is filled then
//    restField := (restField << fieldWidth) | emptyLine
//    score += ...
// else
//    restField := (line << range) | restField

// field = (field >> range << range) | restField

// example
// field = 11111
// field >> 3 = 11
// field << 3 = 11000

// var fullLine = big.NewInt(111111111111)
func (gameField *Field) Clean() {
	restField := big.NewInt(0)

	fullLine, _ := big.NewInt(0).SetString("111111111111", 2)
	emptyLine, _ := big.NewInt(0).SetString("100000000001", 2)
	for i := 0; i < 21; i++ {
		curRange := uint(i * FieldWidth)
		lineMask := big.NewInt(0).Lsh(fullLine, curRange)
		//lineMask.And(lineMask, gameField.Val)
		lineIsFilled := big.NewInt(0).And(lineMask, gameField.Val).Cmp(lineMask) == 0

		if lineIsFilled {
			// добавляем пустую линию в конец поля
			restField.Lsh(restField, FieldWidth)
			restField.Or(restField, emptyLine)
			// TODO: score += награда за соженную линию
		} else {
			lineMask.And(lineMask, gameField.Val)
			restField.Or(lineMask, restField)

			// Expr!
			//restField.Rsh(restField, FieldWidth)
		}

		// field = (field >> range << range) | restField
		// Обнуляем остаток поля и проставляем туда restField
		gameField.Val.Rsh(gameField.Val, curRange)
		gameField.Val.Lsh(gameField.Val, curRange)
		gameField.Val.Or(gameField.Val, restField)

		//gameField.Val.Rsh(gameField.Val, curRange+FieldWidth)
		//gameField.Val.Lsh(gameField.Val, curRange+FieldWidth)
		//gameField.Val.Or(gameField.Val, restField)
		//fmt.Println(gameField.Val)
	}

	//gameField.Val.Rsh(gameField.Val, FieldWidth)

	// 22 строки. Больше на одну строку, чем надо, поэтому сдвигаем вправо на длину поля в конце
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
