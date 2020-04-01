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
		"Use_Hash":   {NameID: "Use_Hash", HandlerFunc: useCache, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseTT), CurrentValue: strconv.FormatBool(config.Settings.Search.UseTT)},
		"Hash":       {NameID: "Hash", HandlerFunc: cacheSize, OptionType: Spin, DefaultValue: strconv.Itoa(config.Settings.Search.TTSize), CurrentValue: strconv.Itoa(config.Settings.Search.TTSize), MinValue: "0", MaxValue: "65000"},
		"Use_Book":   {NameID: "Use_Book", HandlerFunc: useBook, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseBook), CurrentValue: strconv.FormatBool(config.Settings.Search.UseBook)},
		"Ponder":     {NameID: "Ponder", HandlerFunc: usePonder, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UsePonder), CurrentValue: strconv.FormatBool(config.Settings.Search.UsePonder)},
		"Quiescence": {NameID: "Quiescence", HandlerFunc: useQuiescence, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseQuiescence), CurrentValue: strconv.FormatBool(config.Settings.Search.UseQuiescence)},
		"Use_QHash":  {NameID: "Use_QHash", HandlerFunc: useQSHash, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseQSTT), CurrentValue: strconv.FormatBool(config.Settings.Search.UseQSTT)},
		"Use_PVS":    {NameID: "Use_PVS", HandlerFunc: usePvs, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UsePVS), CurrentValue: strconv.FormatBool(config.Settings.Search.UsePVS)},
		"Use_Mdp":    {NameID: "Use_Mdp", HandlerFunc: useMdp, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseMDP), CurrentValue: strconv.FormatBool(config.Settings.Search.UseMDP)},
		"Use_Mpp":    {NameID: "Use_Mpp", HandlerFunc: useMpp, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseMPP), CurrentValue: strconv.FormatBool(config.Settings.Search.UseMPP)},
		"Use_Killer": {NameID: "Use_Killer", HandlerFunc: useKiller, OptionType: Check, DefaultValue: strconv.FormatBool(config.Settings.Search.UseKiller), CurrentValue: strconv.FormatBool(config.Settings.Search.UseKiller)},
	}
	sortOrderUciOptions = []string{
		"Clear Hash",
		"Use_Hash",
		"Hash",
		"Use_Book",
		"Ponder",
		"Quiescence",
		"Use_QHash",
		"Use_PVS",
		"Use_Mdp",
		"Use_Mpp",
	}

}

// GetOptions returns all available uci options a slice of strings
// to be send to the UCI user interface during the initialization
// phase of the UCI protocol
func (o *optionMap) GetOptions() *[]string {
	var options []string
	for _, opt := range sortOrderUciOptions {
		options = append(options, uciOptions[opt].String())
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

// to control the sort order of all options
var sortOrderUciOptions []string

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

func usePonder(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UsePonder = v
	log.Debugf("Set Use Ponder to %v", config.Settings.Search.UsePonder)
}

func useQuiescence(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseQuiescence = v
	log.Debugf("Set Use Quiescence to %v", config.Settings.Search.UseQuiescence)
}

func useQSHash(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseQSTT = v
	log.Debugf("Set Use Hash in Quiescence to %v", config.Settings.Search.UseQSTT)
}

func usePvs(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UsePVS = v
	log.Debugf("Set Use PVS to %v", config.Settings.Search.UsePVS)
}

func useMdp(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseMDP = v
	log.Debugf("Set Use MDP to %v", config.Settings.Search.UseMDP)
}

func useMpp(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseMPP = v
	log.Debugf("Set Use MPP to %v", config.Settings.Search.UseMPP)
}

func useKiller(u *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	config.Settings.Search.UseKiller = v
	log.Debugf("Set Use Killer Moves to %v", config.Settings.Search.UseKiller)
}
