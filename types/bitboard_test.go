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

package types

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"math/bits"
	"testing"
)

func TestBitboardType(t *testing.T) {
	log.Printf("Testing Bitboards")
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
			t.Logf("Bit count %d of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardStr(t *testing.T) {
	log.Printf("Testing Bitboards String conversion")
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
		got := test.value.str()
		if got != test.expected {
			t.Errorf("Bit str of %d should be %s. Got %s", test.value, test.expected, got)
		} else {
			t.Logf("Bit str %s of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardOps(t *testing.T) {
	log.Printf("Testing Bitboards Operations")
	tests := []struct {
		value    Bitboard
		expected string
	}{
		{SqA1.bitboard_(), "0000000000000000000000000000000000000000000000000000000000000001"},
		{SqH8.bitboard_(), "1000000000000000000000000000000000000000000000000000000000000000"},
		{BbZero.put(SqA1), "0000000000000000000000000000000000000000000000000000000000000001"},
		{BbZero.put(SqH8), "1000000000000000000000000000000000000000000000000000000000000000"},
		{BbZero.put(SqE5), "0000000000000000000000000001000000000000000000000000000000000000"},
		{BbZero.put(SqE4), "0000000000000000000000000000000000010000000000000000000000000000"},
		{BbZero.put(SqE4).remove(SqE4), "0000000000000000000000000000000000000000000000000000000000000000"},
		{BbZero.put(SqA1).remove(SqA1), "0000000000000000000000000000000000000000000000000000000000000000"},
		{BbZero.remove(SqA1), "0000000000000000000000000000000000000000000000000000000000000000"},
	}
	for _, test := range tests {
		got := test.value.str()
		if got != test.expected {
			t.Errorf("Bit str of %d should be %s. Got %s", test.value, test.expected, got)
		} else {
			t.Logf("Bit str %s of %d is correct.", got, test.value)
		}
	}
}

func TestBitboardStrBoard(t *testing.T) {
	Init()
	fmt.Println(BbZero.strBoard())
	fmt.Println(BbOne.strBoard())
	fmt.Println(BbAll.strBoard())
}

func TestBitboardStrGrp(t *testing.T) {
	Init()
	fmt.Println(BbZero.strGrp())
	fmt.Println(BbOne.strGrp())
	fmt.Println(BbAll.strGrp())

	assert.Equal(t, "10000000.00000000.00000000.00000000.00000000.00000000.00000000.00000000 (1)", BbOne.strGrp())
	assert.Equal(t, "0000000000000000000000000000000000000000000000000000000000000001", BbOne.str())
}

func TestBitboardDiagUp(t *testing.T) {
	Init()
	fmt.Println(DiagUpA1.strBoard())
	fmt.Println(DiagUpB1.strBoard())
	fmt.Println(DiagUpC1.strBoard())
	fmt.Println(DiagUpD1.strBoard())
	fmt.Println(DiagUpE1.strBoard())
	fmt.Println(DiagUpF1.strBoard())
	fmt.Println(DiagUpG1.strBoard())
	fmt.Println(DiagUpH1.strBoard())
	assert.Equal(t, "10000000.01000000.00100000.00010000."+
		"00001000.00000100.00000010.00000001 (9241421688590303745)", DiagUpA1.strGrp())
	assert.Equal(t, "00000010.00000001.00000000.00000000."+
		"00000000.00000000.00000000.00000000 (32832)", DiagUpG1.strGrp())

	fmt.Println(DiagUpA2.strBoard())
	fmt.Println(DiagUpA3.strBoard())
	fmt.Println(DiagUpA4.strBoard())
	fmt.Println(DiagUpA5.strBoard())
	fmt.Println(DiagUpA6.strBoard())
	fmt.Println(DiagUpA7.strBoard())
	fmt.Println(DiagUpA8.strBoard())
	assert.Equal(t, "00000000.10000000.01000000.00100000."+
		"00010000.00001000.00000100.00000010 (4620710844295151872)", DiagUpA2.strGrp())
	assert.Equal(t, "00000000.00000000.00000000.00000000."+
		"00000000.00000000.10000000.01000000 (144396663052566528)", DiagUpA7.strGrp())
}

func TestBitboardDiagDown(t *testing.T) {
	Init()
	fmt.Println(DiagDownH1.strBoard())
	fmt.Println(DiagDownH2.strBoard())
	fmt.Println(DiagDownH3.strBoard())
	fmt.Println(DiagDownH4.strBoard())
	fmt.Println(DiagDownH5.strBoard())
	fmt.Println(DiagDownH6.strBoard())
	fmt.Println(DiagDownH7.strBoard())
	fmt.Println(DiagDownH8.strBoard())
	assert.Equal(t, "00000001.00000010.00000100.00001000."+
		"00010000.00100000.01000000.10000000 (72624976668147840)", DiagDownH1.strGrp())
	assert.Equal(t, "00000000.00000000.00000000.00000000."+
		"00000000.00000001.00000010.00000100 (2323998145211531264)", DiagDownH6.strGrp())

	fmt.Println(DiagDownG1.strBoard())
	fmt.Println(DiagDownF1.strBoard())
	fmt.Println(DiagDownE1.strBoard())
	fmt.Println(DiagDownD1.strBoard())
	fmt.Println(DiagDownC1.strBoard())
	fmt.Println(DiagDownB1.strBoard())
	fmt.Println(DiagDownA1.strBoard())
	assert.Equal(t, "00000100.00001000.00010000.00100000."+
		"01000000.10000000.00000000.00000000 (1108169199648)", DiagDownF1.strGrp())
	assert.Equal(t, "01000000.10000000.00000000.00000000."+
		"00000000.00000000.00000000.00000000 (258)", DiagDownB1.strGrp())

}

func TestBitboardLsbMsb(t *testing.T) {
	Init()

	tests := []struct {
		bitboard Bitboard
		lsb      Square
		msb      Square
	}{
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
		fmt.Printf("Lsb of %s == %s (%d)\n", test.bitboard.str(), test.bitboard.Lsb().str(), test.bitboard.Lsb())
		// Msb
		assert.Equal(t, test.msb, test.bitboard.Msb())
		fmt.Printf("Msb of %s == %s (%d)\n", test.bitboard.str(), test.bitboard.Msb().str(), test.bitboard.Msb())
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
		{DiagUpA2, DiagUpA2.remove(SqA2), SqA2},
	}

	for _, test := range tests {
		// PopLsb
		fmt.Printf("Bitboard in %s \n", test.bbIn.str())
		got := test.bbIn.PopLsb()
		fmt.Printf("Square is %s \nBitboard out %s \n", got.str(), test.bbIn.str())
		assert.Equal(t, test.square, got)
		assert.Equal(t, test.bbOut, test.bbIn)
	}

	i := 0
	b := DiagDownH3
	var sq Square
	fmt.Printf("Bitboard %d = %s \n", i, b.str())
	for sq = b.PopLsb(); sq != SqNone; sq = b.PopLsb() {
		i++
		fmt.Printf("Bitboard %d = %s \n", i, b.str())
	}
	assert.Equal(t, 6, i)

}

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
