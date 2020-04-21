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

package config

// searchConfiguration is a data structure to hold the configuration of an
// instance of a search.
type searchConfiguration struct {
	UseBook    bool
	BookPath   string
	BookFile   string
	BookFormat string

	UsePonder     bool
	UseQuiescence bool
	UseQSStandpat bool

	UsePVS       bool
	UseKiller    bool
	UseIID       bool
	IIDDepth     int
	IIDReduction int

	UseTT      bool
	TTSize     int
	UseTTMove  bool
	UseTTValue bool
	UseQSTT    bool
	UseEvalTT  bool

	UseMDP      bool
	UseNullMove bool

	NmpDepth     int
	NmpReduction int

	UseLmr           bool
	LmrDepth         int
	LmrMovesSearched int

	UseLmp bool
	UseSEE bool

	UseExt       bool
	UseCheckExt  bool
	UseThreatExt bool
}

// sets defaults which might be overwritten by config file
func init() {
	Settings.Search.UseBook = true
	Settings.Search.BookPath = "./assets/books"
	Settings.Search.BookPath = "book.txt"
	Settings.Search.BookFormat = "Simple"
	Settings.Search.UsePonder = true
	Settings.Search.UseQuiescence = true
	Settings.Search.UseQSStandpat = true
	Settings.Search.UsePVS = true
	Settings.Search.UseKiller = true
	Settings.Search.UseTT = true
	Settings.Search.TTSize = 128
	Settings.Search.UseTTMove = true
	Settings.Search.UseTTValue = true
	Settings.Search.UseQSTT = true
	Settings.Search.UseEvalTT = false
	Settings.Search.UseMDP = true
	Settings.Search.UseNullMove = true
	Settings.Search.NmpDepth = 3
	Settings.Search.NmpReduction = 2
	Settings.Search.UseIID = true
	Settings.Search.IIDDepth = 6
	Settings.Search.IIDReduction = 2
	Settings.Search.UseLmr = true
	Settings.Search.LmrDepth = 3
	Settings.Search.LmrMovesSearched = 3
	Settings.Search.UseLmp = false
	Settings.Search.UseSEE = true
	Settings.Search.UseExt = true
	Settings.Search.UseCheckExt = true
	Settings.Search.UseThreatExt = false

}

// set defaults for configurations here in case a configuration
// is not available from the config file
func setupSearch() {

}
