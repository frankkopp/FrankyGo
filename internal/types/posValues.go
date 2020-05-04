//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package types

import (
	"github.com/frankkopp/FrankyGo/internal/assert"
)

// PosMidValue returns the pre computed positional value
// for the piece on the given square in mid game
func PosMidValue(p Piece, sq Square) Value {
	if assert.DEBUG {
		assert.Assert(initialized, "Pos values have not been initialized. Please call types.Init() first.")
	}
	return posMidValue[p][sq]
}

// PosEndValue returns the pre computed positional value
// for the piece on the given square in end game
func PosEndValue(p Piece, sq Square) Value {
	if assert.DEBUG {
		assert.Assert(initialized, "Pos values have not been initialized. Please call types.Init() first.")
	}
	return posEndValue[p][sq]
}

// PosValue returns the pre computed positional value
// for the piece on the given square and given game phase
func PosValue(p Piece, sq Square, gp int) Value {
	if assert.DEBUG {
		assert.Assert(initialized, "Pos values have not been initialized. Please call types.Init() first.")
	}
	return posValue[p][sq][gp]
}

// initPosValues pre computes arrays containing the positional values of each piece
// for each square and game phase
func initPosValues() {
	// all game phases
	for pc := WhiteKing; pc <= BlackQueen; pc++ {
		for sq := SqA1; sq <= SqH8; sq++ {
			for gp := GamePhaseMax; gp >= 0; gp-- {
				switch pc {
				case WhiteKing:
					posMidValue[pc][sq] = kingMidGame[63-sq]
					posEndValue[pc][sq] = kingEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &kingMidGame, &kingEndGame)
				case WhitePawn:
					posMidValue[pc][sq] = pawnsMidGame[63-sq]
					posEndValue[pc][sq] = pawnsEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &pawnsMidGame, &pawnsEndGame)
				case WhiteKnight:
					posMidValue[pc][sq] = knightMidGame[63-sq]
					posEndValue[pc][sq] = knightEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &knightMidGame, &knightEndGame)
				case WhiteBishop:
					posMidValue[pc][sq] = bishopMidGame[63-sq]
					posEndValue[pc][sq] = bishopEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &bishopMidGame, &bishopEndGame)
				case WhiteRook:
					posMidValue[pc][sq] = rookMidGame[63-sq]
					posEndValue[pc][sq] = rookEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &rookMidGame, &rookEndGame)
				case WhiteQueen:
					posMidValue[pc][sq] = queenMidGame[63-sq]
					posEndValue[pc][sq] = queenEndGame[63-sq]
					posValue[pc][sq][gp] = calcPosValueWhite(sq, gp, &queenMidGame, &queenEndGame)
				case BlackKing:
					posMidValue[pc][sq] = kingMidGame[sq]
					posEndValue[pc][sq] = kingEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &kingMidGame, &kingEndGame)
				case BlackPawn:
					posMidValue[pc][sq] = pawnsMidGame[sq]
					posEndValue[pc][sq] = pawnsEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &pawnsMidGame, &pawnsEndGame)
				case BlackKnight:
					posMidValue[pc][sq] = knightMidGame[sq]
					posEndValue[pc][sq] = knightEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &knightMidGame, &knightEndGame)
				case BlackBishop:
					posMidValue[pc][sq] = bishopMidGame[sq]
					posEndValue[pc][sq] = bishopEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &bishopMidGame, &bishopEndGame)
				case BlackRook:
					posMidValue[pc][sq] = rookMidGame[sq]
					posEndValue[pc][sq] = rookEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &rookMidGame, &rookEndGame)
				case BlackQueen:
					posMidValue[pc][sq] = queenMidGame[sq]
					posEndValue[pc][sq] = queenEndGame[sq]
					posValue[pc][sq][gp] = calcPosValueBlack(sq, gp, &queenMidGame, &queenEndGame)
				default:
				}
			}
		}
	}
}

func calcPosValueBlack(sq Square, gamePhase int, posMidTable *[SqLength]Value, posEndTable *[SqLength]Value) Value {
	return (Value(gamePhase)*posMidTable[sq] + (Value(1-gamePhase) * posEndTable[sq])) / GamePhaseMax
}

func calcPosValueWhite(sq Square, gamePhase int, posMidTable *[SqLength]Value, posEndTable *[SqLength]Value) Value {
	return (Value(gamePhase)*posMidTable[63-sq] + (Value(GamePhaseMax-gamePhase))*posEndTable[63-sq]) / GamePhaseMax
}

var (
	posMidValue = [PieceLength][SqLength]Value{}
	posEndValue = [PieceLength][SqLength]Value{}
	posValue    = [PieceLength][SqLength][GamePhaseMax + 1]Value{}

	// positional values for pieces
	// @formatter:off
	// PAWN Tables
	pawnsMidGame = [SqLength]Value {
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  5,  5,  5,  5,  5,  5,  0,
	5,  5, 10, 30, 30, 10,  5,  5,
	0,  0,  0, 30, 30,  0,  0,  0,
	5, -5,-10,  0,  0,-10, -5,  5,
	5, 10, 10,-30,-30, 10, 10,  5,
	0,  0,  0,  0,  0,  0,  0,  0}

	pawnsEndGame = [SqLength]Value {
	0,  0,  0,  0,  0,  0,  0,  0,
	90, 90, 90, 90, 90, 90, 90, 90,
	40, 50, 50, 60, 60, 50, 50, 40,
	20, 30, 30, 40, 40, 30, 30, 20,
	10, 10, 20, 20, 20, 10, 10, 10,
	5, 10, 10, 10, 10, 10, 10,  5,
	5, 10, 10, 10, 10, 10, 10,  5,
	0,  0,  0,  0,  0,  0,  0,  0}

	// KNIGHT Tables
	knightMidGame = [SqLength]Value {
	-50,-40,-30,-30,-30,-30,-40,-50,
	-40,-20,  0,  0,  0,  0,-20,-40,
	-30,  0, 10, 15, 15, 10,  0,-30,
	-30,  5, 15, 20, 20, 15,  5,-30,
	-30,  0, 15, 20, 20, 15,  0,-30,
	-30,  5, 10, 15, 15, 10,  5,-30,
	-40,-20,  0,  5,  5,  0,-20,-40,
	-50,-25,-20,-30,-30,-20,-25,-50}

	knightEndGame = [SqLength]Value {
	-50,-40,-30,-30,-30,-30,-40,-50,
	-40,-20,  0,  0,  0,  0,-20,-40,
	-30,  0, 10, 15, 15, 10,  0,-30,
	-30,  0, 15, 20, 20, 15,  0,-30,
	-30,  0, 15, 20, 20, 15,  0,-30,
	-30,  0, 10, 15, 15, 10,  0,-30,
	-40,-20,  0,  0,  0,  0,-20,-40,
	-50,-40,-20,-30,-30,-20,-40,-50}

	// BISHOP Tables
	bishopMidGame = [SqLength]Value {
	-20,-10,-10,-10,-10,-10,-10,-20,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-10,  0,  5, 10, 10,  5,  0,-10,
	-10,  5,  5, 10, 10,  5,  5,-10,
	-10,  0, 10, 10, 10, 10,  0,-10,
	-10, 10, 10, 10, 10, 10, 10,-10,
	-10,  5,  0,  0,  0,  0,  5,-10,
	-20,-10,-40,-10,-10,-40,-10,-20}

	bishopEndGame = [SqLength]Value {
	-20,-10,-10,-10,-10,-10,-10,-20,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-10,  0,  5, 10, 10,  5,  0,-10,
	-10,  0,  5, 10, 10,  5,  0,-10,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-20,-10,-10,-10,-10,-10,-10,-20}

	// ROOK Tables
	rookMidGame  = [SqLength]Value {
	5,  5,  5,  5,  5,  5,  5,  5,
	10, 10, 10, 10, 10, 10, 10, 10,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	-15,-10, 15, 15, 15, 15,-10,-15}

	rookEndGame  = [SqLength]Value {
	5,  5,  5,  5,  5,  5,  5,  5,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0,
	0,  0,  0,  0,  0,  0,  0,  0}

	// Queen Tables
	queenMidGame = [SqLength]Value {
	-20,-10,-10, -5, -5,-10,-10,-20,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-5,  0,  2,  2,  2,  2,  0, -5,
	-5,  0,  5,  5,  5,  5,  0, -5,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-20,-10,-10, -5, -5,-10,-10,-20}

	queenEndGame = [SqLength]Value {
	-20,-10,-10, -5, -5,-10,-10,-20,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-5,  0,  5,  5,  5,  5,  0, -5,
	-5,  0,  5,  5,  5,  5,  0, -5,
	-10,  0,  5,  5,  5,  5,  0,-10,
	-10,  0,  0,  0,  0,  0,  0,-10,
	-20,-10,-10, -5, -5,-10,-10,-20}

	// King Tables
	kingMidGame  = [SqLength]Value {
	-30,-40,-40,-50,-50,-40,-40,-30,
	-30,-40,-40,-50,-50,-40,-40,-30,
	-30,-40,-40,-50,-50,-40,-40,-30,
	-30,-40,-40,-50,-50,-40,-40,-30,
	-20,-30,-30,-40,-40,-30,-30,-20,
	-10,-20,-20,-30,-30,-30,-20,-10,
	0,  0,-20,-20,-20,-20,  0,  0,
	20, 50,  0,-20,-20,  0, 50, 20}

	kingEndGame  = [SqLength]Value {
	-50,-30,-30,-20,-20,-30,-30,-50,
	-30,-20,-10,  0,  0,-10,-20,-30,
	-30,-10, 20, 30, 30, 20,-10,-30,
	-30,-10, 30, 40, 40, 30,-10,-30,
	-30,-10, 30, 40, 40, 30,-10,-30,
	-30,-10, 20, 30, 30, 20,-10,-30,
	-30,-30,  0,  0,  0,  0,-30,-30,
	-50,-30,-30,-30,-30,-30,-30,-50}
	// @formatter:on
)
