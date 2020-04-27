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

package types

// Magic holds all magic bitboards relevant for a single square
// Taken from Stockfish
// License see https://stockfishchess.org/about/
type Magic struct {
	Mask    Bitboard
	Magic   Bitboard
	Attacks []Bitboard
	Shift   uint
}

// init_magics() computes all rook and bishop attacks at startup. Magic
// bitboards are used to look up attacks of sliding pieces. As a reference see
// https://www.chessprogramming.org/Magic_Bitboards. In particular, here we use the so
// called "fancy" approach.
// Taken from Stockfish
func initMagics(table *[]Bitboard, magics *[64]Magic, directions *[4]Direction) {

	// Optimal PrnG seeds to pick the correct magics in the shortest time
	seeds := [RankLength]uint64{728, 10316, 55013, 32803, 12281, 15100, 16645, 255}

	occupancy := [4096]Bitboard{}
	reference := [4096]Bitboard{}
	var edges, b Bitboard
	cnt := 0
	size := 0
	epoch := [4096]int{}

	for sq := SqA1; sq <= SqH8; sq++ {
		// log.Debugf("Init magic bitboards: Phase 1 for Square: %s", sq)

		// Board edges are not considered in the relevant occupancies
		edges = ((Rank1_Bb | Rank8_Bb) &^ sq.RankOf().Bb()) | ((FileA_Bb | FileH_Bb) &^ sq.FileOf().Bb())

		// Given a square 'sq', the mask is the bitboard of sliding attacks from
		// 'sq' computed on an empty board. The index must be big enough to contain
		// all the attacks for each possible subset of the mask and so is 2 power
		// the number of 1s of the mask. Hence we deduce the size of the shift to
		// apply to the 64 bits word to get the index.
		m := &(*magics)[sq]
		m.Mask = slidingAttack(directions, sq, BbZero) &^ edges
		m.Shift = uint(64 - m.Mask.PopCount())

		// Set the offset for the attacks table of the square. We have individual
		// table sizes for each square with "Fancy Magic Bitboards".
		if sq == SqA1 {
			m.Attacks = *table
		} else {
			m.Attacks = magics[sq-1].Attacks[size:] // instead of pointer offset use slice offset
		}

		// Use Carry-Rippler trick to enumerate all subsets of masks[s] and
		// store the corresponding sliding attack bitboard in reference[].
		// https://www.chessprogramming.org/Traversing_Subsets_of_a_Set
		b = 0
		size = 0
		for {
			occupancy[size] = b
			reference[size] = slidingAttack(directions, sq, b)
			size++
			b = (b - m.Mask) & m.Mask
			if b == 0 { // do - while(b)
				break
			}
		}

		// special random number generator
		rng := newPrnG(seeds[sq.RankOf()])

		// Find a magic for square 's' picking up an (almost) random number
		// until we find the one that passes the verification test.
		// TODO: understand this better
		for i := 0; i < size; {
			for m.Magic = 0;; {
				m.Magic = Bitboard(rng.sparseRand())
				if ((m.Magic * m.Mask) >> 56).PopCount() < 6 {
					break
				}
			}

			// A good magic must map every possible occupancy to an index that
			// looks up the correct sliding attack in the attacks[s] database.
			// Note that we build up the database for square 's' as a side
			// effect of verifying the magic. Keep track of the attempt count
			// and save it in epoch[], little speed-up trick to avoid resetting
			// m.attacks[] after every failed attempt.
			cnt++
			for i = 0; i < size; i++ {
				idx := m.index(occupancy[i])
				if epoch[idx] < cnt {
					epoch[idx] = cnt
					m.Attacks[idx] = reference[i]
				} else if m.Attacks[idx] != reference[i] {
					break
				}
			}
		}
	}
}

// slidingAttack calculate sliding attacks along the given directions for the given square
// and the given board occupation. Uses loop in loop and is not very efficient.
// Doesn't matter for pre-computing but should not be used during move gen or search
func slidingAttack(directions *[4]Direction, sq Square, occupied Bitboard) Bitboard {
	attack := BbZero
	for i := 0; i < 4; i++ {
		s := sq
		for {
			s = s.To(directions[i])
			if !s.IsValid() {
				break
			}
			attack.PushSquare(s)
			if occupied.Has(s) {
				break
			}
			if !s.To(directions[i]).IsValid() || SquareDistance(s, s.To(directions[i])) != 1 {
				break
			}
		}
	}
	return attack
}

// Index calculates the index in the table for the attacks
// https://www.chessprogramming.org/Magic_Bitboards
//  occ      &= mBishopTbl[sq].mask;
//  occ      *= mBishopTbl[sq].magic;
//  occ     >>= mBishopTbl[sq].shift;
func (m *Magic) index(occupied Bitboard) uint {
	occ := occupied & m.Mask
	occ *= m.Magic
	occ >>= m.Shift
	return uint(occ)
}

// PrnG random generator for magic bitboards
// from Stockfish
// xorshift64star Pseudo-Random Number Generator
// This class is based on original code written and dedicated
// to the public domain by Sebastiano Vigna (2014).
// It has the following characteristics:
//  -  Outputs 64-bit numbers
//  -  Passes Dieharder and SmallCrush test batteries
//  -  Does not require warm-up, no zeroland to escape
//  -  Internal state is a single 64-bit integer
//  -  Period is 2^64 - 1
//  -  Speed: 1.60 ns/call (Core i7 @3.40GHz)
// For further analysis see
//   <http://vigna.di.unimi.it/ftp/papers/xorshift.pdf>
type PrnG struct {
	s uint64
}

// newPrnG creates a new instance of the pseudo random generator
func newPrnG(seed uint64) *PrnG {
	return &PrnG{s: seed}
}

func (r *PrnG) rand64() uint64 {
	r.s ^= r.s >> 12
	r.s ^= r.s << 25
	r.s ^= r.s >> 27
	return r.s * 2685821657736338717
}

// Special generator used to fast init magic numbers.
// Output values only have 1/8th of their bits set on average.
func (r *PrnG) sparseRand() uint64 {
	return r.rand64() & r.rand64() & r.rand64()
}
