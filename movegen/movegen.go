/*
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

// Package movegen contains functionality to create moves on a
// chess position. It implements several variants like
// pseudo legal move list, legal move list or on demand move
// generation of pseudo legal moves.
package movegen

import (
	"github.com/frankkopp/FrankyGo/assert"
	"github.com/frankkopp/FrankyGo/movearray"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

type movegen struct {
	pseudoLegalMoves   movearray.MoveArray
	legalMoves         movearray.MoveArray
	onDemandMoves      movearray.MoveArray
	killerMoves        movearray.MoveArray
	pvMove             Move
	currentODStage     int
	currentIteratorKey position.Key
	maxNumberOfKiller  int
}

// States for the on demand move generator
const (
	odNew = iota
	odPv  = iota
	od1   = iota
	od2   = iota
	od3   = iota
	od4   = iota
	od5   = iota
	od6   = iota
	od7   = iota
	od8   = iota
	odEnd = iota
)

// //////////////////////////////////////////////////////
// // Public functions
// //////////////////////////////////////////////////////

// GenMode generation modes for on demand move generation
type GenMode int

// GenMode generation modes for on demand move generation
const (
	GenZero   GenMode = 0b00
	GenCap    GenMode = 0b01
	GenNonCap GenMode = 0b10
	GenAll    GenMode = 0b11
)

// New creates a new instance of a move generator
func New() movegen {
	tmpMg := movegen{
		pseudoLegalMoves:   movearray.New(MaxMoves),
		legalMoves:         movearray.New(MaxMoves),
		onDemandMoves:      movearray.New(MaxMoves),
		killerMoves:        movearray.New(4),
		pvMove:             MoveNone,
		currentODStage:     odNew,
		currentIteratorKey: 0,
		maxNumberOfKiller:  2, // default
	}
	// tmpMg.pseudoLegalMoves.SetMinCapacity(6)
	// tmpMg.legalMoves.SetMinCapacity(6)
	// tmpMg.onDemandMoves.SetMinCapacity(6)
	return tmpMg
}

// GeneratePseudoLegalMoves generates pseudo moves for the next player. Does not check if
// king is left in check or passes an attacked square when castling or has been in check
// before castling. Disregards PV moves and Killer moves. They need to be handled after
// the returned MoveList. Or just use the OnDemand Generator.
func (mg *movegen) GeneratePseudoLegalMoves(position *position.Position, mode GenMode) *movearray.MoveArray {
	mg.pseudoLegalMoves.Clear()
	if mode&GenCap != 0 {
		mg.generatePawnMoves(position, GenCap, &mg.pseudoLegalMoves)
		mg.generateCastling(position, GenCap, &mg.pseudoLegalMoves)
		mg.generateKingMoves(position, GenCap, &mg.pseudoLegalMoves)
		mg.generateMoves(position, GenCap, &mg.pseudoLegalMoves)
	}
	if mode&GenNonCap != 0 {
		mg.generatePawnMoves(position, GenNonCap, &mg.pseudoLegalMoves)
		mg.generateCastling(position, GenNonCap, &mg.pseudoLegalMoves)
		mg.generateKingMoves(position, GenNonCap, &mg.pseudoLegalMoves)
		mg.generateMoves(position, GenNonCap, &mg.pseudoLegalMoves)
	}
	// sort.Stable(&mg.pseudoLegalMoves)
	mg.pseudoLegalMoves.Sort()
	// remove internal sort value
	mg.pseudoLegalMoves.ForEach(func(i int) {
		mg.pseudoLegalMoves.Set(i, mg.pseudoLegalMoves.At(i).MoveOf())
	})
	return &mg.pseudoLegalMoves
}

func (mg *movegen) String() string {
	return "movegen instance"
}

// //////////////////////////////////////////////////////
// // Private functions
// //////////////////////////////////////////////////////

func (mg *movegen) generatePawnMoves(position *position.Position, mode GenMode, ml *movearray.MoveArray) {

	nextPlayer := position.NextPlayer()
	myPawns := position.PiecesBb(nextPlayer, Pawn)
	oppPieces := position.OccupiedBb(nextPlayer.Flip())
	gamePhase := position.GamePhase()
	piece := MakePiece(nextPlayer, Pawn)

	// captures
	if mode&GenCap != 0 {

		// This algorithm shifts the own pawn bitboard in the direction of pawn captures
		// and ANDs it with the opponents pieces. With this we get all possible captures
		// and can easily create the moves by using a loop over all captures and using
		// the backward shift for the from-Square.
		// All moves get stable_sort values so that stable_sort order should be:
		//   captures: most value victim least value attacker - promotion piece value
		//   non captures: killer (TBD), promotions, castling, normal moves (position value)
		// Values for sorting are descending - the most valuable move has the highest value
		// values are not compatible to position evaluation values.

		var tmpCaptures, promCaptures Bitboard

		for _, dir := range []Direction{West, East} {
			// normal pawn captures - promotions first
			tmpCaptures = ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North+dir) & oppPieces
			promCaptures = tmpCaptures & nextPlayer.PromotionRankBb()
			// promotion captures
			for promCaptures != 0 {
				toSquare := promCaptures.PopLsb()
				fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North - dir)
				// value is the delta of values from the two pieces involved plus the positional value
				value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
					PosValue(piece, toSquare, gamePhase)
				// add the possible promotion moves to the move list and also add value of the promoted piece type
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, value+Queen.ValueOf()))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, value+Knight.ValueOf()))
				// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
				// therefore we give them lower sort order
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, value+Rook.ValueOf()-Value(2000)))
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, value+Bishop.ValueOf()-Value(2000)))
			}
			tmpCaptures &= ^nextPlayer.PromotionRankBb()
			for tmpCaptures != 0 {
				toSquare := tmpCaptures.PopLsb()
				fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection())*North - dir)
				// value is the delta of values from the two pieces involved plus the positional value
				value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
					PosValue(piece, toSquare, gamePhase)
				ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
			}
		}

		// en passant captures
		enPassantSquare := position.GetEnPassantSquare()
		if enPassantSquare != SqNone {
			for _, dir := range []Direction{West, East} {
				tmpCaptures = ShiftBitboard(enPassantSquare.Bb(),
					Direction(nextPlayer.Flip().MoveDirection())*North+dir) & myPawns
				if tmpCaptures != 0 {
					fromSquare := tmpCaptures.PopLsb()
					toSquare := fromSquare.To(Direction(nextPlayer.MoveDirection())*North - dir)
					// value is the positional value of the piece at this game phase
					value := PosValue(piece, toSquare, gamePhase)
					ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
				}
			}
		}
	}

	// non captures
	if mode&GenNonCap != 0 {

		//  Move my pawns forward one step and keep all on not occupied squares
		//  Move pawns now on rank 3 (rank 6) another square forward to check for pawn doubles.
		//  Loop over pawns remaining on unoccupied squares and add moves.

		// pawns - check step one to unoccupied squares
		tmpMoves := ShiftBitboard(myPawns, Direction(nextPlayer.MoveDirection())*North) & ^position.OccupiedAll()
		// pawns double - check step two to unoccupied squares
		tmpMovesDouble := ShiftBitboard(tmpMoves&nextPlayer.PawnDoubleRank(), Direction(nextPlayer.MoveDirection())*North) & ^position.OccupiedAll()

		// single pawn steps - promotions first
		promMoves := tmpMoves & nextPlayer.PromotionRankBb()
		for promMoves != 0 {
			toSquare := promMoves.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			// value for non captures is lowered by 10k
			value := Value(-10_000)
			// add the possible promotion moves to the move list and also add value of the promoted piece type
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Queen, value+Queen.ValueOf()))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Knight, value+Knight.ValueOf()))
			// rook and bishops are usually redundant to queen promotion (except in stale mate situations)
			// therefore we give them lower sort order
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Rook, value+Rook.ValueOf()-Value(2000)))
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Promotion, Bishop, value+Bishop.ValueOf()-Value(2000)))
		}
		// double pawn steps
		for tmpMovesDouble != 0 {
			toSquare := tmpMovesDouble.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North).
				To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
		// normal single pawn steps
		tmpMoves &= ^nextPlayer.PromotionRankBb()
		for tmpMoves != 0 {
			toSquare := tmpMoves.PopLsb()
			fromSquare := toSquare.To(Direction(nextPlayer.Flip().MoveDirection()) * North)
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}
}

func (mg *movegen) generateCastling(position *position.Position, mode GenMode, ml *movearray.MoveArray) {
	nextPlayer := position.NextPlayer()
	occupiedBB := position.OccupiedAll()

	// castling - pseudo castling - we will not check if we are in check after the move
	// or if we have passed an attacked square with the king or if the king has been in check

	if mode&GenNonCap != 0 && position.CastlingRights() != CastlingNone {
		cr := position.CastlingRights()
		if nextPlayer == White { // white
			if cr.Has(CastlingWhiteOO) && Intermediate(SqE1, SqH1)&occupiedBB == 0 {
				if assert.DEBUG {
					assert.Assert(position.KingSquare(White) == SqE1, "MoveGen Castling: White King not on e1")
					assert.Assert(position.GetPiece(SqH1) == WhiteRook, "MoveGen Castling: White Rook not on h1")
				}
				ml.PushBack(CreateMoveValue(SqE1, SqG1, Castling, PtNone, Value(-5000)))
			}
			if cr.Has(CastlingWhiteOOO) && Intermediate(SqE1, SqA1)&occupiedBB == 0 {
				if assert.DEBUG {
					assert.Assert(position.KingSquare(White) == SqE1, "MoveGen Castling: White King not on e1")
					assert.Assert(position.GetPiece(SqA1) == WhiteRook, "MoveGen Castling: White Rook not on a1")
				}
				ml.PushBack(CreateMoveValue(SqE1, SqC1, Castling, PtNone, Value(-5000)))
			}
		} else { // black
			if cr.Has(CastlingBlackOO) && Intermediate(SqE8, SqH8)&occupiedBB == 0 {
				if assert.DEBUG {
					assert.Assert(position.KingSquare(Black) == SqE8, "MoveGen Castling: Black King not on e8")
					assert.Assert(position.GetPiece(SqH8) == BlackRook, "MoveGen Castling: Black Rook not on h8")
				}
				ml.PushBack(CreateMoveValue(SqE8, SqG8, Castling, PtNone, Value(-5000)))
			}
			if cr.Has(CastlingBlackOOO) && Intermediate(SqE8, SqA8)&occupiedBB == 0 {
				if assert.DEBUG {
					assert.Assert(position.KingSquare(Black) == SqE8, "MoveGen Castling: Black King not on e8")
					assert.Assert(position.GetPiece(SqA8) == BlackRook, "MoveGen Castling: Black Rook not on a8")
				}
				ml.PushBack(CreateMoveValue(SqE8, SqC8, Castling, PtNone, Value(-5000)))
			}
		}
	}
}

func (mg *movegen) generateKingMoves(position *position.Position, mode GenMode, ml *movearray.MoveArray) {
	nextPlayer := position.NextPlayer()
	piece := MakePiece(nextPlayer, King)
	gamePhase := position.GamePhase()
	kingSquareBb := position.PiecesBb(nextPlayer, King)
	if assert.DEBUG {
		assert.Assert(kingSquareBb.PopCount() == 1,
			"Chess always needs exactly one king. Found=%d ", kingSquareBb.PopCount())
	}
	fromSquare := kingSquareBb.PopLsb()

	// pseudo attacks include all moves no matter if the king would be in check
	pseudoMoves := GetPseudoAttacks(King, fromSquare)

	// captures
	if mode&GenCap != 0 {
		captures := pseudoMoves & position.OccupiedBb(nextPlayer.Flip())
		for captures != 0 {
			toSquare := captures.PopLsb()
			value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
				PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}

	// non captures
	if mode&GenNonCap != 0 {
		nonCaptures := pseudoMoves &^ position.OccupiedAll()
		for nonCaptures != 0 {
			toSquare := nonCaptures.PopLsb()
			value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
			ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
		}
	}
}

func (mg *movegen) generateMoves(position *position.Position, mode GenMode, ml *movearray.MoveArray) {
	nextPlayer := position.NextPlayer()
	gamePhase := position.GamePhase()
	occupiedBb := position.OccupiedAll()

	for pt := Knight; pt <= Queen; pt++ {
		pieces := position.PiecesBb(nextPlayer, pt)
		piece := MakePiece(nextPlayer, pt)

		for pieces != 0 {
			fromSquare := pieces.PopLsb()
			pseudoMoves := GetPseudoAttacks(pt, fromSquare)

			// captures
			if mode&GenCap != 0 {
				captures := pseudoMoves & position.OccupiedBb(nextPlayer.Flip())
				for captures != 0 {
					toSquare := captures.PopLsb()
					if pt > Knight { // sliding pieces
						if Intermediate(fromSquare, toSquare)&occupiedBb == 0 {
							value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
								PosValue(piece, toSquare, gamePhase)
							ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
						}
					} else { // king and knight cannot be blocked
						value := position.GetPiece(toSquare).ValueOf() - position.GetPiece(fromSquare).ValueOf() +
							PosValue(piece, toSquare, gamePhase)
						ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
					}
				}
			}

			// non captures
			if mode&GenNonCap != 0 {
				nonCaptures := pseudoMoves &^ occupiedBb
				for nonCaptures != 0 {
					toSquare := nonCaptures.PopLsb()
					if pt > Knight { // sliding pieces
						if Intermediate(fromSquare, toSquare)&occupiedBb == 0 {
							value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
							ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
						}
					} else { // king and knight cannot be blocked
						value := Value(-10_000) + PosValue(piece, toSquare, gamePhase)
						ml.PushBack(CreateMoveValue(fromSquare, toSquare, Normal, PtNone, value))
					}
				}
			}
		}
	}
}
