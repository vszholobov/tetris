package main

import (
	"math"
	"math/big"
)

type PieceType int

const (
	TShape = 0
	ZigZag = 1
	LShape = 2
)

var rotationsCntByType = map[PieceType]int{
	TShape: 4,
	ZigZag: 2,
	LShape: 2,
}

//000001110000
//000000100000 = 458784

//000000100000
//000001100000
//000000100000 = 537264160

//000000100000
//000001110000
//000000000000 = 537329664

//000000100000
//000000110000
//000000100000 = 537067552

var rotationsByType = map[PieceType][]*big.Int{
	TShape: {big.NewInt(458784), big.NewInt(537264160), big.NewInt(537329664), big.NewInt(537067552)},
}

type RotationType int

const (
	Left  RotationType = -1
	Right RotationType = 1
)

type Piece struct {
	rotationCount int
	pieceType     PieceType
	rotations     []*big.Int
	field         *Field
}

func MakePiece(field *Field, pieceType PieceType) Piece {
	rotations := rotationsByType[pieceType]
	rotationsCopy := copyRotations(rotations)
	return Piece{
		rotationCount: 0,
		pieceType:     pieceType,
		rotations:     rotationsCopy,
		field:         field,
	}
}

func copyRotations(rotations []*big.Int) []*big.Int {
	rotationsCopy := make([]*big.Int, len(rotations))
	copy(rotationsCopy, rotations)
	for i, rotation := range rotationsCopy {
		rotationsCopy[i] = big.NewInt(0).Set(rotation)
	}
	return rotationsCopy
}

func (piece *Piece) Rotate(rotationType RotationType) bool {
	var diff int
	if rotationType == Left {
		diff = -1
	} else {
		diff = 1
	}

	piece.changeRotationCount(diff)
	if !piece.field.intersects(piece.GetVal()) {
		return true
	}
	piece.changeRotationCount(-diff)
	return false
}

func (piece *Piece) changeRotationCount(diff int) {
	maxRotations := len(rotationsByType[piece.pieceType])
	piece.rotationCount += diff
	if piece.rotationCount < 0 {
		piece.rotationCount = maxRotations - 1
	} else if piece.rotationCount == maxRotations {
		piece.rotationCount = 0
	}
}

func (piece *Piece) GetVal() *big.Int {
	abs := int64(math.Abs(float64(piece.rotationCount % rotationsCntByType[piece.pieceType])))
	return piece.rotations[abs]
}

func (piece *Piece) moveLeft() bool {
	newPieceVal := big.NewInt(0).Set(piece.GetVal())
	newPieceVal.Rsh(newPieceVal, 1)
	if piece.field.intersects(newPieceVal) {
		return false
	}
	for i, _ := range piece.rotations {
		newRotation := big.NewInt(0).Set(piece.rotations[i])
		piece.rotations[i] = newRotation.Rsh(newRotation, 1)
	}
	return true
}

func (piece *Piece) moveRight() bool {
	newPieceVal := big.NewInt(0).Set(piece.GetVal())
	newPieceVal.Lsh(newPieceVal, 1)
	if piece.field.intersects(newPieceVal) {
		return false
	}
	for i, _ := range piece.rotations {
		newRotation := big.NewInt(0).Set(piece.rotations[i])
		piece.rotations[i] = newRotation.Lsh(newRotation, 1)
	}
	return true
}

func (piece *Piece) moveDown() bool {
	newPieceVal := big.NewInt(0).Set(piece.GetVal())
	newPieceVal.Lsh(newPieceVal, fieldWidth)
	if piece.field.intersects(newPieceVal) {
		return false
	}
	for i, _ := range piece.rotations {
		newRotation := big.NewInt(0).Set(piece.rotations[i])
		piece.rotations[i] = newRotation.Lsh(newRotation, fieldWidth)
	}
	return true
}

//func (piece Piece) dropPiece(field Field) Piece {
//	for !piece.moveDown().intersects() {
//		piece = piece.moveDown()
//	}
//	return piece
//}
