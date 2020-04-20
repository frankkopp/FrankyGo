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

package search

import (
	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

func see(p *position.Position, move Move) Value {

	// enpassant moves are ignored in a sense that it will be winning
	// capture and therefore should lead to no cut-offs when using see()
	if move.MoveType() == EnPassant {
		return 100
	}

	// prepare short array to store the captures - max 32 pieces
	// TODO: avoid allocation
	gain := make([]Value, 32, 32)

	ply := 0
	toSquare := move.To()
	fromSquare := move.From()
	movedPiece := p.GetPiece(fromSquare)
	nextPlayer := p.NextPlayer()

	// get a bitboard of all occupied squares to remove single pieces later
	// to reveal hidden attacks (x-ray)
	occupiedBitboard := p.OccupiedAll()

	// get all attacks to the square as a bitboard
	remainingAttacks := AttacksTo(p, toSquare, White) | AttacksTo(p, toSquare, Black)

	// log := myLogging.GetLog()
	// log.Debugf("Determine gain for %s %s", p.StringFen(), move.StringUci())

	// initial value of the first capture
	capturedValue := p.GetPiece(toSquare).ValueOf()
	gain[ply] = capturedValue

	// log.Debugf("gain[%d] = %s | %s", ply, gain[ply].String(), move.StringUci());

	// loop through all remaining attacks/captures
	for {
		ply++
		nextPlayer = nextPlayer.Flip()

		// speculative store, if defended
		if move.MoveType() == Promotion {
			gain[ply] = move.PromotionType().ValueOf() - Pawn.ValueOf() - gain[ply-1]
		} else {
			gain[ply] = movedPiece.ValueOf() - gain[ply-1]
		}

		// pruning if defended - will not change final see score
		if max(-gain[ply-1], gain[ply]) < 0 {
			break
		}

		remainingAttacks.PopSquare(fromSquare) // reset bit in set to traverse
		occupiedBitboard.PopSquare(fromSquare) // reset bit in temporary occupancy (for x-Rays)

		// reevaluate attacks to reveal attacks after removing the moving piece
		remainingAttacks |= revealedAttacks(p, toSquare, occupiedBitboard, White) |
			revealedAttacks(p, toSquare, occupiedBitboard, Black)

		// determine next capture
		fromSquare = getLeastValuablePiece(p, remainingAttacks, nextPlayer)

		// break if no more attackers
		if fromSquare == SqNone {
			break
		}

		// log.Debugf("gain[%d] = %s | %s%s", ply, gain[ply].String(), fromSquare.String(), toSquare.String())
		movedPiece = p.GetPiece(fromSquare)
	}

	ply--
	for ply > 0 {
		gain[ply-1] = -max(-gain[ply-1], gain[ply])
		ply--
	}

	return gain[0]
}

// AttacksTo determine all attacks for SEE. EnPassant is not included as this is not
// relevant for SEE as the move preceding enpassant is always non capturing.
func AttacksTo(p *position.Position, square Square, color Color) Bitboard {
	occupiedAll := p.OccupiedAll()
	return (GetPawnAttacks(color.Flip(), square) & p.PiecesBb(color, Pawn)) |
		// Knight
		(GetAttacksBb(Knight, square, occupiedAll) & p.PiecesBb(color, Knight)) |
		// King
		(GetAttacksBb(King, square, occupiedAll) & p.PiecesBb(color, King)) |
		// Sliding rooks and queens
		(GetAttacksBb(Rook, square, occupiedAll) & (p.PiecesBb(color, Rook) | p.PiecesBb(color, Queen))) |
		// Sliding bishops and queens
		(GetAttacksBb(Bishop, square, occupiedAll) & (p.PiecesBb(color, Bishop) | p.PiecesBb(color, Queen)))
}

// Returns sliding attacks after a piece has been removed to reveal new attacks.
// It is only necessary to look at slider pieces as only their attacks can be revealed
func revealedAttacks(p *position.Position, square Square, occupied Bitboard, color Color) Bitboard {
	// Sliding rooks and queens
	return (GetAttacksBb(Rook, square, occupied) & (p.PiecesBb(color, Rook) | p.PiecesBb(color, Queen)) & occupied) |
		// Sliding bishops and queens
		(GetAttacksBb(Bishop, square, occupied) & (p.PiecesBb(color, Bishop) | p.PiecesBb(color, Queen)) & occupied)
}

// Returns a square with the least valuable attacker. When several of same
// type are available it uses the least significant bit of the bitboard.
func getLeastValuablePiece(position *position.Position, bitboard Bitboard, color Color) Square {
	// check all piece types with increasing value
	switch {
	case (bitboard & position.PiecesBb(color, Pawn)) != 0:
		return (bitboard & position.PiecesBb(color, Pawn)).Lsb()
	case (bitboard & position.PiecesBb(color, Knight)) != 0:
		return (bitboard & position.PiecesBb(color, Knight)).Lsb()
	case (bitboard & position.PiecesBb(color, Bishop)) != 0:
		return (bitboard & position.PiecesBb(color, Bishop)).Lsb()
	case (bitboard & position.PiecesBb(color, Rook)) != 0:
		return (bitboard & position.PiecesBb(color, Rook)).Lsb()
	case (bitboard & position.PiecesBb(color, Queen)) != 0:
		return (bitboard & position.PiecesBb(color, Queen)).Lsb()
	case (bitboard & position.PiecesBb(color, King)) != 0:
		return (bitboard & position.PiecesBb(color, King)).Lsb()
	default:
		return SqNone
	}
}

func max(x, y Value) Value {
	if x > y {
		return x
	}
	return y
}
