package field

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"
)

const FieldWidth = 12
const FieldHeight = 21
const clearBoardANSII = "\033[22A"

type Field struct {
	Val          *big.Int
	CurrentPiece *Piece
	Score        *int
	CleanCount   *int
}

func MakeDefaultField() Field {
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
	return MakeField(fieldVal)
}

func MakeField(fieldVal *big.Int) Field {
	score := 0
	cleanCount := 0
	return Field{Val: fieldVal, Score: &score, CleanCount: &cleanCount}
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
	currentCleanCount := 0
	for i := 0; i < FieldHeight-1; i++ {
		curRange := uint(i * FieldWidth)
		lineMask := big.NewInt(0).Lsh(fullLine, curRange)
		lineIsFilled := big.NewInt(0).And(lineMask, gameField.Val).Cmp(lineMask) == 0

		if lineIsFilled {
			// add empy line to end of field
			restField.Lsh(restField, FieldWidth)
			restField.Or(restField, emptyLine)
			currentCleanCount += 1
		} else {
			// add current line to start of field
			lineMask.And(lineMask, gameField.Val)
			restField.Or(lineMask, restField)
		}
	}
	*gameField.CleanCount += currentCleanCount
	*gameField.Score += currentCleanCount * gameField.GetSpeed() * 10 / (5 - currentCleanCount)
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
			"000000000000", 2)
	gameField.Val.Or(gameField.Val, restField)
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
	result := ""
	fieldStr := field.String()
	for i := 20; i >= 0; i-- {
		line := fieldStr[i*12 : i*12+12]
		line = strings.ReplaceAll(line, "1", " Ж ")
		line = strings.ReplaceAll(line, "0", "   ")
		result += line + "\n"
	}
	result += "Score: " + strconv.Itoa(*field.Score) + " | Speed: " + strconv.Itoa(field.GetSpeed()) + " | Cleaned: " + strconv.Itoa(*field.CleanCount)
	fmt.Println(clearBoardANSII + result)
}

func (gameField *Field) GetSpeed() int {
	return *gameField.CleanCount/4 + 1
}
