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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/frankkopp/FrankyGo/internal/position"
	. "github.com/frankkopp/FrankyGo/internal/types"
)

func TestAttacksTo(t *testing.T) {
	var p *position.Position
	var attacksTo Bitboard

	p = position.NewPosition("2brr1k1/1pq1b1p1/p1np1p1p/P1p1p2n/1PNPPP2/2P1BNP1/4Q1BP/R2R2K1 w - -")
	attacksTo = AttacksTo(p, SqE5, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 740294656, attacksTo)

	attacksTo = AttacksTo(p, SqF1, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 20552, attacksTo)

	attacksTo = AttacksTo(p, SqD4, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 3407880, attacksTo)

	attacksTo = AttacksTo(p, SqD4, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 4483945857024, attacksTo)

	attacksTo = AttacksTo(p, SqD6, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 582090251837636608, attacksTo)

	attacksTo = AttacksTo(p, SqF8, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 5769111122661605376, attacksTo)

	p = position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3")
	attacksTo = AttacksTo(p, SqE5, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2339760743907840, attacksTo)

	attacksTo = AttacksTo(p, SqB1, Black)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 1280, attacksTo)

	attacksTo = AttacksTo(p, SqG3, White)
	logTest.Debug("\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 40960, attacksTo)
}

func TestRevealedAttacks(t *testing.T) {
	p := position.NewPosition("1k1r3q/1ppn3p/p4b2/4p3/8/P2N2P1/1PP1R1BP/2K1Q3 w - -")
	occ := p.OccupiedAll()

	sq := SqE5

	attacksTo := AttacksTo(p, sq, White) | AttacksTo(p, sq, Black)
	logTest.Debug("Direct\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2286984186302464, attacksTo)

	// take away bishop on f6
	attacksTo.PopSquare(SqF6)
	occ.PopSquare(SqF6)

	attacksTo |= revealedAttacks(p, sq, occ, White) | revealedAttacks(p, sq, occ, Black)
	logTest.Debug("Revealed\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, Bitboard(9225623836668989440), attacksTo)

	// take away rook on e2
	attacksTo.PopSquare(SqE2)
	occ.PopSquare(SqE2)

	attacksTo |= revealedAttacks(p, sq, occ, White) | revealedAttacks(p, sq, occ, Black)
	logTest.Debug("Revealed\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, Bitboard(9225623836668985360), attacksTo)
}

func TestLeastValuablePiece(t *testing.T) {
	p := position.NewPosition("r3k2r/1ppn3p/2q1q1n1/4P3/2q1Pp2/6R1/pbp2PPP/1R4K1 b kq e3")
	attacksTo := AttacksTo(p, SqE5, Black)

	logTest.Debug("All attackers\n", attacksTo.StringBoard())
	logTest.Debug(attacksTo.StringGrouped())
	assert.EqualValues(t, 2339760743907840, attacksTo)

	lva := getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqG6, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqD7, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqB2, lva)

	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqE6, lva)
	// remove the attacker
	attacksTo.PopSquare(lva)

	lva = getLeastValuablePiece(p, attacksTo, Black)
	logTest.Debug("Least valuable piece:", lva.String())
	assert.EqualValues(t, SqNone, lva)
}
