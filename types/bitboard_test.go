/*
 * MIT License
 *
 * Copyright (c) 2018-2020 Frank Kopp
 *
 * Permission is hereby granted, free of charge, To any person obtaining a copy
 * of this software and associated documentation files (the "Software"), To deal
 * in the Software without restriction, including without limitation the rights
 * To use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and To permit persons To whom the Software is
 * furnished To do so, subject To the following conditions:
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

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"math/bits"
	"testing"
)

// set To true for printing output during tests
const verbose bool = true

func TestBitboardType(t *testing.T) {
	Init()
	tests := []struct {
		value    Bitboard
		expected int
	}{
		{BbZero, 0},
		{BbAll, 64},
		{BbOne, 1},
		{Bitboard(128), 1},
		{Bitboard(7), 3},
	}
	for _, test := range tests {
		got := bits.OnesCount64(uint64(test.value))
		if got != test.expected {
			t.Errorf("Bit count of %d should be %d. Got %d", test.value, test.expected, got)
		} else {
			//t.Logf("Bit count %d of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardStr(t *testing.T) {
	Init()
	tests := []struct {
		value    Bitboard
		expected string
	}{
		{BbZero, "0000000000000000000000000000000000000000000000000000000000000000"},
		{BbAll, "1111111111111111111111111111111111111111111111111111111111111111"},
		{BbOne, "0000000000000000000000000000000000000000000000000000000000000001"},
		{FileA_Bb, "0000000100000001000000010000000100000001000000010000000100000001"},
		{Rank1_Bb, "0000000000000000000000000000000000000000000000000000000011111111"},
		{FileH_Bb, "1000000010000000100000001000000010000000100000001000000010000000"},
		{Rank8_Bb, "1111111100000000000000000000000000000000000000000000000000000000"},
	}
	for _, test := range tests {
		got := test.value.Str()
		if got != test.expected {
			t.Errorf("Bit Str of %d should be %s. Got %s", test.value, test.expected, got)
		} else {
			//t.Logf("Bit Str %s of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardPutRemove(t *testing.T) {
	Init()
	tests := []struct {
		value    Bitboard
		expected string
	}{
		{SqA1.bitboard_(), "0000000000000000000000000000000000000000000000000000000000000001"},
		{SqH8.bitboard_(), "1000000000000000000000000000000000000000000000000000000000000000"},
		{PushSquare(BbZero, SqA1), "0000000000000000000000000000000000000000000000000000000000000001"},
		{PushSquare(BbZero, SqH8), "1000000000000000000000000000000000000000000000000000000000000000"},
		{PushSquare(BbZero, SqE5), "0000000000000000000000000001000000000000000000000000000000000000"},
		{PushSquare(BbZero, SqE4), "0000000000000000000000000000000000010000000000000000000000000000"},
		{PopSquare(PushSquare(BbZero, SqE4), SqE4), "0000000000000000000000000000000000000000000000000000000000000000"},
		{PopSquare(PushSquare(BbZero, SqA1), SqA1), "0000000000000000000000000000000000000000000000000000000000000000"},
		{PopSquare(BbZero, SqA1), "0000000000000000000000000000000000000000000000000000000000000000"},
	}
	for _, test := range tests {
		got := test.value.Str()
		if got != test.expected {
			t.Errorf("Bit Str of %d should be %s. Got %s", test.value, test.expected, got)
		} else {
			//t.Logf("Bit Str %s of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardStrBoard(t *testing.T) {
	Init()
	if verbose {
		fmt.Println(BbZero.StrBoard())
	}
	if verbose {
		fmt.Println(BbOne.StrBoard())
	}
	if verbose {
		fmt.Println(BbAll.StrBoard())
	}
}

func TestBitboardStrGrp(t *testing.T) {
	Init()
	if verbose {
		fmt.Println(BbZero.StrGrp())
	}
	if verbose {
		fmt.Println(BbOne.StrGrp())
	}
	if verbose {
		fmt.Println(BbAll.StrGrp())
	}

	assert.Equal(t, "10000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000 (1)", BbOne.StrGrp())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", BbOne.Str())
}

func TestBitboardDiagUp(t *testing.T) {
	Init()
	if verbose {
		fmt.Println(DiagUpA1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpB1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpC1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpD1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpE1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpF1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpG1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpH1.StrBoard())
	}
	assert.Equal(t, "10000000.01000000.00100000.00010000."+
		"00001000.00000100.00000010.00000001 (9241421688590303745)", DiagUpA1.StrGrp())
	assert.Equal(t, "00000010.00000001.00000000.00000000."+
		"00000000.00000000.00000000.00000000 (32832)", DiagUpG1.StrGrp())

	if verbose {
		fmt.Println(DiagUpA2.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA3.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA4.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA5.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA6.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA7.StrBoard())
	}
	if verbose {
		fmt.Println(DiagUpA8.StrBoard())
	}
	assert.Equal(t, "00000000.10000000.01000000.00100000."+
		"00010000.00001000.00000100.00000010 (4620710844295151872)", DiagUpA2.StrGrp())
	assert.Equal(t, "00000000.00000000.00000000.00000000."+
		"00000000.00000000.10000000.01000000 (144396663052566528)", DiagUpA7.StrGrp())
}

func TestBitboardDiagDown(t *testing.T) {
	Init()
	if verbose {
		fmt.Println(DiagDownH1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH2.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH3.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH4.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH5.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH6.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH7.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownH8.StrBoard())
	}
	assert.Equal(t, "00000001.00000010.00000100.00001000."+
		"00010000.00100000.01000000.10000000 (72624976668147840)", DiagDownH1.StrGrp())
	assert.Equal(t, "00000000.00000000.00000000.00000000."+
		"00000000.00000001.00000010.00000100 (2323998145211531264)", DiagDownH6.StrGrp())

	if verbose {
		fmt.Println(DiagDownG1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownF1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownE1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownD1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownC1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownB1.StrBoard())
	}
	if verbose {
		fmt.Println(DiagDownA1.StrBoard())
	}
	assert.Equal(t, "00000100.00001000.00010000.00100000."+
		"01000000.10000000.00000000.00000000 (1108169199648)", DiagDownF1.StrGrp())
	assert.Equal(t, "01000000.10000000.00000000.00000000."+
		"00000000.00000000.00000000.00000000 (258)", DiagDownB1.StrGrp())

}

func TestBitboardLsbMsb(t *testing.T) {
	Init()

	tests := []struct {
		bitboard Bitboard
		lsb      Square
		msb      Square
	}{
		{BbZero, SqNone, SqNone},
		{SqA1.Bitboard(), SqA1, SqA1},
		{SqH8.Bitboard(), SqH8, SqH8},
		{SqE5.Bitboard(), SqE5, SqE5},
		{DiagUpA2, SqA2, SqG8},
		{DiagDownH3, SqH3, SqC8},
		{FileB_Bb, SqB1, SqB8},
		{Rank3_Bb, SqA3, SqH3},
	}

	for _, test := range tests {
		// Lsb
		assert.Equal(t, test.lsb, test.bitboard.Lsb())
		if verbose {
			fmt.Printf("Lsb of %s == %s (%d)\n", test.bitboard.Str(), test.bitboard.Lsb().Str(), test.bitboard.Lsb())
		}
		// Msb
		assert.Equal(t, test.msb, test.bitboard.Msb())
		if verbose {
			fmt.Printf("Msb of %s == %s (%d)\n", test.bitboard.Str(), test.bitboard.Msb().Str(), test.bitboard.Msb())
		}
	}
}

func TestBitboardPopLsb(t *testing.T) {
	Init()

	tests := []struct {
		bbIn   Bitboard
		bbOut  Bitboard
		square Square
	}{
		{SqA1.Bitboard(), BbZero, SqA1},
		{SqH8.Bitboard(), BbZero, SqH8},
		{DiagUpA2, PopSquare(DiagUpA2, SqA2), SqA2},
	}

	for _, test := range tests {
		// PopLsb
		if verbose {
			fmt.Printf("Bitboard in %s \n", test.bbIn.Str())
		}
		got := test.bbIn.PopLsb()
		if verbose {
			fmt.Printf("Square is %s \nBitboard out %s \n", got.Str(), test.bbIn.Str())
		}
		assert.Equal(t, test.square, got)
		assert.Equal(t, test.bbOut, test.bbIn)
	}

	i := 0
	b := DiagDownH3
	var sq Square
	if verbose {
		fmt.Printf("Bitboard %d = %s \n", i, b.Str())
	}
	for sq = b.PopLsb(); sq != SqNone; sq = b.PopLsb() {
		i++
		if verbose {
			fmt.Printf("Bitboard %d = %s \n", i, b.Str())
		}
	}
	assert.Equal(t, 6, i)

}

func TestBitboardShift(t *testing.T) {
	Init()

	tests := []struct {
		preShift  Bitboard
		shift     Direction
		postShift Bitboard
	}{
		//Vertical and horizontal shifts
		{DiagUpA2, North, DiagUpA3},
		{DiagUpA3, North, DiagUpA4},
		{DiagUpB1, South, DiagUpC1},
		{DiagUpC1, South, DiagUpD1},
		{DiagUpD1, South, DiagUpE1},
		{DiagDownH1, North, DiagDownH2},
		{DiagDownH2, North, DiagDownH3},
		{DiagDownH3, North, DiagDownH4},
		{DiagDownH4, North, DiagDownH5},
		{DiagDownH1, East, DiagDownH2},
		{DiagDownH2, East, DiagDownH3},
		{DiagDownH3, East, DiagDownH4},
		{DiagDownH4, East, DiagDownH5},
		{DiagDownH1, South, DiagDownG1},
		{DiagDownG1, South, DiagDownF1},
		{DiagDownF1, South, DiagDownE1},
		{DiagDownE1, South, DiagDownD1},
		{DiagDownH1, West, DiagDownG1},
		{DiagDownG1, West, DiagDownF1},
		{DiagDownF1, West, DiagDownE1},
		{DiagDownE1, West, DiagDownD1},
		{Rank8_Bb | FileH_Bb, East, PopSquare(Rank8_Bb, SqA8)},

		// diagonal shifts
		{Rank8_Bb | FileH_Bb, Northeast, BbZero},
		{Rank1_Bb | FileA_Bb, Northeast, Bitboard(0x20202020202fe00)},
		{Rank1_Bb | FileA_Bb, Southwest, BbZero},
		{Rank8_Bb | FileH_Bb, Southwest, Bitboard(0x7f404040404040)},
		{Rank8_Bb | FileA_Bb, Northwest, BbZero},
		{Rank1_Bb | FileH_Bb, Northwest, Bitboard(0x4040404040407f00)},
		{Rank1_Bb | FileH_Bb, Southeast, BbZero},
		{Rank8_Bb | FileA_Bb, Southeast, Bitboard(0xfe020202020202)},

		// single square all directions
		{SqE4.Bitboard(), North, SqE5.Bitboard()},
		{SqE4.Bitboard(), Northeast, SqF5.Bitboard()},
		{SqE4.Bitboard(), East, SqF4.Bitboard()},
		{SqE4.Bitboard(), Southeast, SqF3.Bitboard()},
		{SqE4.Bitboard(), South, SqE3.Bitboard()},
		{SqE4.Bitboard(), Southwest, SqD3.Bitboard()},
		{SqE4.Bitboard(), West, SqD4.Bitboard()},
		{SqE4.Bitboard(), Northwest, SqD5.Bitboard()},

		// single square at edge all directions
		{SqA4.Bitboard(), North, SqA5.Bitboard()},
		{SqA4.Bitboard(), Northeast, SqB5.Bitboard()},
		{SqA4.Bitboard(), East, SqB4.Bitboard()},
		{SqA4.Bitboard(), Southeast, SqB3.Bitboard()},
		{SqA4.Bitboard(), South, SqA3.Bitboard()},
		{SqA4.Bitboard(), Southwest, BbZero},
		{SqA4.Bitboard(), West, BbZero},
		{SqA4.Bitboard(), Northwest, BbZero},

		// single square at corner all directions
		{SqA1.Bitboard(), North, SqA2.Bitboard()},
		{SqA1.Bitboard(), Northeast, SqB2.Bitboard()},
		{SqA1.Bitboard(), East, SqB1.Bitboard()},
		{SqA1.Bitboard(), Southeast, BbZero},
		{SqA1.Bitboard(), South, BbZero},
		{SqA1.Bitboard(), Southwest, BbZero},
		{SqA1.Bitboard(), West, BbZero},
		{SqA1.Bitboard(), Northwest, BbZero},

		// single square at corner all directions
		{SqH8.Bitboard(), North, BbZero},
		{SqH8.Bitboard(), Northeast, BbZero},
		{SqH8.Bitboard(), East, BbZero},
		{SqH8.Bitboard(), Southeast, BbZero},
		{SqH8.Bitboard(), South, SqH7.Bitboard()},
		{SqH8.Bitboard(), Southwest, SqG7.Bitboard()},
		{SqH8.Bitboard(), West, SqG8.Bitboard()},
		{SqH8.Bitboard(), Northwest, BbZero},
	}

	for _, test := range tests {
		got := ShiftBitboard(test.preShift, test.shift)
		if verbose {
			fmt.Printf("Bitboard in  \n%s \n", test.preShift.StrBoard())
		}
		if verbose {
			fmt.Printf("Bitboard out \n%s \n", got.StrBoard())
		}
		assert.Equal(t, test.postShift, got)
	}
}

func TestBitboardInit(t *testing.T) {
	Init()

	// Square bitboards
	assert.Equal(t, SqA1.bitboard_().Str(), "0000000000000000000000000000000000000000000000000000000000000001")
	assert.Equal(t, SqH8.bitboard_().Str(), "1000000000000000000000000000000000000000000000000000000000000000")

	// square To file index
	assert.Equal(t, sqToFileBb[SqA2], FileA_Bb)
	assert.Equal(t, sqToFileBb[SqC5], FileC_Bb)
	assert.Equal(t, sqToFileBb[SqF6], FileF_Bb)
	assert.Equal(t, sqToFileBb[SqH8], FileH_Bb)

	// square To rank index
	assert.Equal(t, sqToRankBb[SqA2], Rank2_Bb)
	assert.Equal(t, sqToRankBb[SqC5], Rank5_Bb)
	assert.Equal(t, sqToRankBb[SqF6], Rank6_Bb)
	assert.Equal(t, sqToRankBb[SqH8], Rank8_Bb)

	// square To diag up index
	assert.Equal(t, sqDiagUpBb[SqA2], DiagUpA2)
	assert.Equal(t, sqDiagUpBb[SqC5], DiagUpA3)
	assert.Equal(t, sqDiagUpBb[SqF6], DiagUpA1)
	assert.Equal(t, sqDiagUpBb[SqH8], DiagUpA1)

	// square To diag down index
	assert.Equal(t, sqDiagDownBb[SqA2], DiagDownB1)
	assert.Equal(t, sqDiagDownBb[SqC5], DiagDownG1)
	assert.Equal(t, sqDiagDownBb[SqF6], DiagDownH4)
	assert.Equal(t, sqDiagDownBb[SqH8], DiagDownH8)
}

func TestBitboardFileDistance(t *testing.T) {
	Init()

	tests := []struct {
		f1   File
		f2   File
		dist int
	}{
		{FileA, FileA, 0},
		{FileA, FileB, 1},
		{FileB, FileA, 1},
		{FileA, FileH, 7},
		{FileH, FileA, 7},
		{FileC, FileF, 3},
		{FileF, FileC, 3},
	}

	for _, test := range tests {
		// PopLsb
		got := FileDistance(test.f1, test.f2)
		if verbose {
			fmt.Printf("File distance between %s and %s is %d \n", test.f1.Str(), test.f2.Str(), got)
		}
		assert.Equal(t, test.dist, got)
	}
}

func TestBitboardSquareDistance(t *testing.T) {
	Init()

	tests := []struct {
		s1   Square
		s2   Square
		dist int
	}{
		{SqA1, SqA1, 0},
		{SqA1, SqA2, 1},
		{SqA1, SqB1, 1},
		{SqA1, SqB2, 1},
		{SqA1, SqH8, 7},
		{SqA8, SqH1, 7},
		{SqD4, SqA1, 3},
		{SqE5, SqD4, 1},
	}

	for _, test := range tests {
		// PopLsb
		got := squareDistance[test.s1][test.s2]
		if verbose {
			fmt.Printf("Square distance between %s and %s is %d \n", test.s1.Str(), test.s2.Str(), got)
		}
		assert.Equal(t, test.dist, got)
	}
}

func TestBitboardRotateBb(t *testing.T) {
	Init()

	bitboard := FileA_Bb | Rank8_Bb | DiagDownH1

	rotatedBb := RotateR90(bitboard)
	if verbose {
		fmt.Printf("%s\n%s\n", bitboard.StrBoard(), bitboard.StrGrp())
	}
	if verbose {
		fmt.Printf("%s\n%s\n", rotatedBb.StrBoard(), rotatedBb.StrGrp())
	}
	assert.Equal(t, Bitboard(18428906217826189953), rotatedBb)

	rotatedBb = rotateL90(bitboard)
	if verbose {
		fmt.Printf("%s\n%s\n", bitboard.StrBoard(), bitboard.StrGrp())
	}
	if verbose {
		fmt.Printf("%s\n%s\n", rotatedBb.StrBoard(), rotatedBb.StrGrp())
	}
	assert.Equal(t, Bitboard(9313761861428380671), rotatedBb)

	bitboard = DiagUpA1
	rotatedBb = rotateR45(bitboard)
	if verbose {
		fmt.Printf("%s\n%s\n", bitboard.StrBoard(), bitboard.StrGrp())
	}
	if verbose {
		fmt.Printf("%s\n%s\n", rotatedBb.StrBoard(), rotatedBb.StrGrp())
	}
	assert.Equal(t, Bitboard(68451041280), rotatedBb)

	bitboard = DiagDownH1
	rotatedBb = rotateL45(bitboard)
	if verbose {
		fmt.Printf("%s\n%s\n", bitboard.StrBoard(), bitboard.StrGrp())
	}
	if verbose {
		fmt.Printf("%s\n%s\n", rotatedBb.StrBoard(), rotatedBb.StrGrp())
	}
	assert.Equal(t, Bitboard(68451041280), rotatedBb)
}

func TestBitboardRotateSq(t *testing.T) {
	Init()

	tests := []struct {
		rotation string
		square   Square
		expected Square
	}{
		{"R90", SqA1, SqA8},
		{"R90", SqD8, SqH5},
		{"L90", SqH8, SqA8},
		{"L90", SqH2, SqG8},
		{"R45", SqH8, SqD5},
		{"L45", SqH1, SqD5},
		{"R45", SqC7, SqA8},
		{"L45", SqB3, SqH1},
	}

	for _, test := range tests {
		bitboard := test.square.Bitboard()
		rotated := SqNone
		switch test.rotation {
		case "R90":
			rotated = RotateSquareR90(test.square)
		case "L90":
			rotated = RotateSquareL90(test.square)
		case "R45":
			rotated = RotateSquareR45(test.square)
		case "L45":
			rotated = RotateSquareL45(test.square)
		}
		rotatedBb := rotated.Bitboard()
		if verbose {
			fmt.Printf("Input   : %s\n%s\n%s\n", test.rotation, bitboard.StrBoard(), bitboard.StrGrp())
		}
		if verbose {
			fmt.Printf("Rotation: %s\n%s\n%s\n", test.rotation, rotatedBb.StrBoard(), rotatedBb.StrGrp())
		}
		assert.Equal(t, test.expected, rotated)
	}
}

// TODO implement getMoves functions and test

// //////////////////////////////////////////////////////////////////////////
// benchmarks

//noinspection GoUnusedGlobalVariable
var result Bitboard

func BenchmarkSqBb1(b *testing.B) {
	Init()
	var bb Bitboard
	for i := 0; i < b.N; i++ {
		for square := SqA1; square < SqNone; square++ {
			bb = square.bitboard_()
		}
	}
	result = bb
}

func BenchmarkSqBb2(b *testing.B) {
	Init()
	var bb Bitboard
	for i := 0; i < b.N; i++ {
		for square := SqA1; square < SqNone; square++ {
			bb = square.Bitboard()
		}
	}
	result = bb
}
