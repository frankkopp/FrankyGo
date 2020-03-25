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

package uci

import (
	"strconv"
	"strings"

	"github.com/frankkopp/FrankyGo/config"
)

// init will define all available uci options and store them into the uciOption map
func init() {
	uciOptions = map[string]*uciOption{
		"Clear Hash": {NameID: "Clear Hash", HandlerFunc: clearCache, OptionType: Button},
		"Use_Hash":   {NameID: "Use_Hash", HandlerFunc: useCache, OptionType: Check, DefaultValue: "true", CurrentValue: strconv.FormatBool(config.Settings.Search.UseTT)},
		"Hash":       {NameID: "Hash", HandlerFunc: cacheSize, OptionType: Spin, DefaultValue: "64", CurrentValue: string(config.Settings.Search.TTSize), MinValue: "0", MaxValue: "65000"},
		"Use_Book":   {NameID: "Use_Book", HandlerFunc: useBook, OptionType: Check, DefaultValue: "true", CurrentValue: strconv.FormatBool(config.Settings.Search.UseBook)},
	}
}

// GetOptions returns all available uci options a slice of strings
// to be send to the UCI user interface during the initialization
// phase of the UCI protocol
func (o *optionMap) GetOptions() *[]string {
	var options []string
	for _, opt := range uciOptions {
		options = append(options, opt.String())
	}
	return &options
}

// String for uciOption will return a representation of the uci option as required by
// the UCI protocol during the initialization phase of the UCI protocol
func (o *uciOption) String() string {
	var os strings.Builder
	os.WriteString("option name ")
	os.WriteString(o.NameID)
	os.WriteString(" type ")
	switch o.OptionType {
	case Check:
		os.WriteString("check ")
		os.WriteString("default ")
		os.WriteString(o.DefaultValue)
	case Spin:
		os.WriteString("spin ")
		os.WriteString("default ")
		os.WriteString(o.DefaultValue)
		os.WriteString(" min ")
		os.WriteString(o.MinValue)
		os.WriteString(" max ")
		os.WriteString(o.MaxValue)
	case Combo:
		os.WriteString("combo ")
		os.WriteString("default ")
		os.WriteString(o.DefaultValue)
		os.WriteString(" var ")
		os.WriteString(o.VarValue)
	case Button:
		os.WriteString("button")
	case String:
		os.WriteString("string ")
		os.WriteString("default ")
		os.WriteString(o.DefaultValue)
	}

	return os.String()
}

// uciOptionType is a enum representing the different UCI Option types
type uciOptionType int

// uci option types constants
const (
	Check  uciOptionType = 0
	Spin   uciOptionType = 1
	Combo  uciOptionType = 2
	Button uciOptionType = 3
	String uciOptionType = 4
)

// optionHandler is a function type to by used as function pointer
// in each uci option defined. This is called when the uci option
// is changed by the "setoption" command
type optionHandler func(*UciHandler, *uciOption)

// uciOption defines UCI Options as described in the UCI protocol.
// Each options has a function pointer to a handler which will be
// called when the "setoption" command changes the option.
type uciOption struct {
	NameID       string
	HandlerFunc  optionHandler
	OptionType   uciOptionType
	DefaultValue string
	MinValue     string
	MaxValue     string
	VarValue     string
	CurrentValue string
}

// optionMap convenience type for a map of pointers to uci options
type optionMap map[string]*uciOption

// uciOptions stores all available uci options
var uciOptions optionMap

// ////////////////////////////////////////////////////////////////
// HandlerFunc for uci options changes
// ////////////////////////////////////////////////////////////////

func clearCache(u *UciHandler, o *uciOption) {
	u.mySearch.ClearHash()
	log.Debug("Cleared Cache")
}

func useCache(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseTT = v
	log.Debugf("Set Use Hash to %v", config.Settings.Search.UseTT)
}

func cacheSize(u *UciHandler, o *uciOption) {
	v, _ := strconv.Atoi(o.CurrentValue)
	config.Settings.Search.TTSize = v
	u.mySearch.ResizeCache()
}

func useBook(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseBook = v
	log.Debugf("Set Use Book to %v", config.Settings.Search.UseBook)
}
