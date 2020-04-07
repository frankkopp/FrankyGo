/*
 * FrankyGo - UCI chess engine in GO for learning purposes
 *
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package evaluator

import (
	"github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

var log = logging.GetLog()

// Attacks is a data structure to store all attacks and defends of a position.
type Attacks struct {
	// the position key for which the attacks have been calculated
	Zobrist position.Key
	// bitboards of attacked/defended squares for each color and each from square
	// to get attackers us &^ ownPieces or & ownPieces for defenders
	From    [ColorLength][SqLength]Bitboard
	// bitboards of attackers/defenders for each color and to square
	// to get attackers us &^ ownPieces or & ownPieces for defenders
	To      [ColorLength][SqLength]Bitboard
	// bitboards for all attacked/defended squares of a color
	// to get attackers us &^ ownPieces or & ownPieces for defenders
	All     [ColorLength]Bitboard
	// bitboards of attacked/defended squares for each color and each piece type
	// to get attackers us &^ ownPieces or & ownPieces for defenders
	Piece   [ColorLength][PtLength]Bitboard
	// sum of possible moves for each color (moves to ownPieces already excluded)
	Mobility[ColorLength]int
	// pawn attacks - squares attacked by pawn of the given color
	Pawns [ColorLength]Bitboard
	// pawn double - squares which are attacked twice by pawns of the given color
	PawnsDouble [ColorLength]Bitboard
}

// NewAttacks creates a new instance of Attacks
func NewAttacks() *Attacks {
	return &Attacks{}
}

// Clear resets all fields of the Attacks instance without
// new allocation by looping through all fields
// This is considerably faster than creating a new instance
// Benchmark/New_Instance-8   1.904.764  691.0 ns/op
// Benchmark/Clear-8         13.043.875   91.7 ns/op
func (a *Attacks) Clear() {
	a.Zobrist = 0
	for sq := 0; sq < SqLength; sq++ {
		a.From[White][sq] = BbZero
		a.From[Black][sq] = BbZero
		a.To[White][sq] = BbZero
		a.To[Black][sq] = BbZero
	}
	for pt := PtNone; pt < PtLength; pt++ {
		a.Piece[White][pt] = BbZero
		a.Piece[Black][pt] = BbZero
	}
	a.All[White] = BbZero
	a.All[Black] = BbZero
	a.Mobility[White] = 0
	a.Mobility[Black] = 0
	a.Pawns[White] = 0
	a.Pawns[Black] = 0
	a.PawnsDouble[White] = 0
	a.PawnsDouble[Black] = 0
}

// Compute calculates all attacks on the position.
// Stores the positions zobrist key to be able to
// check if the position is already computed.
// if a position is called twice the already
// stored attacks are untouched.
func (a *Attacks) Compute(p *position.Position) {
	if p.ZobristKey() == a.Zobrist {
		log.Debugf("attacks compute: position was already computed")
		return
	}
	a.Zobrist = p.ZobristKey()
	a.nonPawnAttacks(p)
	a.pawnAttacks(p) // TODO safe time with pawn has table?
}

// nonPawnAttacks calculates all attacks of non pawn pieces including king
func (a *Attacks) nonPawnAttacks(p *position.Position) {

	ptList := [5]PieceType{King, Knight, Bishop, Rook, Queen}
	var attacks Bitboard
	allPieces := p.OccupiedAll()

	// iterate over colors
	for c := White; c <= Black; c++ {
		myPieces := p.OccupiedBb(c)
		// iterate over all piece types
		for _, pt := range ptList {
			// iterate over pieces of piece type
			for pieces := p.PiecesBb(c, pt); pieces != BbZero; {
				psq := pieces.PopLsb() // piece square
				// attacks will include attacks to opponents pieces
				// and defending own pieces
				attacks = GetAttacksBb(pt, psq, allPieces)
				// accumulate all attacks of this piece type for the color
				a.From[c][psq] = attacks
				a.Piece[c][pt] |= attacks
				a.All[c] |= attacks
				// store all attacks to the square
				tmp := attacks
				for tmp != BbZero {
					toSq := tmp.PopLsb() // attacked square
					a.To[c][toSq].PushSquare(psq)
				}
				a.Mobility[c] += (attacks &^ myPieces).PopCount()
			}
		}
	}
}

// pawnAttacks calculate all attacks for pawns.
func (a *Attacks) pawnAttacks(p *position.Position) {
	a.Pawns[White] = ShiftBitboard(p.PiecesBb(White, Pawn), Northwest) | ShiftBitboard(p.PiecesBb(White, Pawn), Northeast)
	a.Pawns[Black] = ShiftBitboard(p.PiecesBb(Black, Pawn), Northwest) | ShiftBitboard(p.PiecesBb(Black, Pawn), Northeast)
	a.PawnsDouble[White] = ShiftBitboard(p.PiecesBb(White, Pawn), Northwest) & ShiftBitboard(p.PiecesBb(White, Pawn), Northeast)
	a.PawnsDouble[Black] = ShiftBitboard(p.PiecesBb(Black, Pawn), Northwest) & ShiftBitboard(p.PiecesBb(Black, Pawn), Northeast)
}


