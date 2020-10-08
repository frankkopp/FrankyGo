//
// FrankyGo - UCI chess engine in GO for learning purposes
//
// MIT License
//
// Copyright (c) 2018-2020 Frank Kopp
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
//

package uci

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	. "github.com/frankkopp/FrankyGo/internal/config"
)

// init will define all available uci options and store them into the uciOption map.
func init() {
	uciOptions = map[string]*uciOption{
		"Print Config": {NameID: "Print Config", HandlerFunc: printConfig, OptionType: Button},
		"Clear Hash":   {NameID: "Clear Hash", HandlerFunc: clearCache, OptionType: Button},
		"Use_Hash":     {NameID: "Use_Hash", HandlerFunc: useCache, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseTT), CurrentValue: strconv.FormatBool(Settings.Search.UseTT)},
		"Hash":         {NameID: "Hash", HandlerFunc: cacheSize, OptionType: Spin, DefaultValue: strconv.Itoa(Settings.Search.TTSize), CurrentValue: strconv.Itoa(Settings.Search.TTSize), MinValue: "0", MaxValue: "65000"},

		"OwnBook": {NameID: "OwnBook", HandlerFunc: useBook, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseBook), CurrentValue: strconv.FormatBool(Settings.Search.UseBook)},

		"Ponder": {NameID: "Ponder", HandlerFunc: usePonder, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UsePonder), CurrentValue: strconv.FormatBool(Settings.Search.UsePonder)},

		"Quiescence":      {NameID: "Quiescence", HandlerFunc: useQuiescence, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseQuiescence), CurrentValue: strconv.FormatBool(Settings.Search.UseQuiescence)},
		"Use_QHash":       {NameID: "Use_QHash", HandlerFunc: useQSHash, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseQSTT), CurrentValue: strconv.FormatBool(Settings.Search.UseQSTT)},
		"Use QS Standpat": {NameID: "Use QS Standpat", HandlerFunc: useQSStandpat, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseQSStandpat), CurrentValue: strconv.FormatBool(Settings.Search.UseQSStandpat)},
		"Use_SEE":         {NameID: "Use_SEE", HandlerFunc: useSee, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseSEE), CurrentValue: strconv.FormatBool(Settings.Search.UseSEE)},

		"Use_AlphaBeta": {NameID: "Use_AlphaBeta", HandlerFunc: useAlphaBeta, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseAlphaBeta), CurrentValue: strconv.FormatBool(Settings.Search.UseAlphaBeta)},
		"Use_PVS":       {NameID: "Use_PVS", HandlerFunc: usePvs, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UsePVS), CurrentValue: strconv.FormatBool(Settings.Search.UsePVS)},
		"Use_ASP":       {NameID: "Use_ASP", HandlerFunc: useAsp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseAspiration), CurrentValue: strconv.FormatBool(Settings.Search.UseAspiration)},

		"Use_TTMove":      {NameID: "Use_TTMove", HandlerFunc: useTTMove, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseTTMove), CurrentValue: strconv.FormatBool(Settings.Search.UseTTMove)},
		"Use_IID":         {NameID: "Use_IID", HandlerFunc: useIID, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseIID), CurrentValue: strconv.FormatBool(Settings.Search.UseIID)},
		"Use_Killer":      {NameID: "Use_Killer", HandlerFunc: useKiller, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseKiller), CurrentValue: strconv.FormatBool(Settings.Search.UseKiller)},
		"Use_HistCount":   {NameID: "Use_HistCount", HandlerFunc: useHC, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseHistoryCounter), CurrentValue: strconv.FormatBool(Settings.Search.UseHistoryCounter)},
		"Use_CounterMove": {NameID: "Use_CounterMove", HandlerFunc: useCM, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseCounterMoves), CurrentValue: strconv.FormatBool(Settings.Search.UseCounterMoves)},

		"Use_TTValue":  {NameID: "Use_TTValue", HandlerFunc: useTTValue, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseTTValue), CurrentValue: strconv.FormatBool(Settings.Search.UseTTValue)},
		"Use_Mdp":      {NameID: "Use_Mdp", HandlerFunc: useMdp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseMDP), CurrentValue: strconv.FormatBool(Settings.Search.UseMDP)},
		"Use_Razoring": {NameID: "Use_Razoring", HandlerFunc: useRazor, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseRazoring), CurrentValue: strconv.FormatBool(Settings.Search.UseRazoring)},
		"Razor Margin": {NameID: "Razor Margin", HandlerFunc: razorMargin, OptionType: Spin, DefaultValue: strconv.Itoa(Settings.Search.RazorMargin), CurrentValue: strconv.Itoa(Settings.Search.RazorMargin), MinValue: "-10000", MaxValue: "10000"},
		"Use_Rfp":      {NameID: "Use_Rfp", HandlerFunc: useRfp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseRFP), CurrentValue: strconv.FormatBool(Settings.Search.UseRFP)},
		"Use_NullMove": {NameID: "Use_NullMove", HandlerFunc: useNullMove, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseNullMove), CurrentValue: strconv.FormatBool(Settings.Search.UseNullMove)},

		"Use_Fp":  {NameID: "Use_Fp", HandlerFunc: useFp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseFP), CurrentValue: strconv.FormatBool(Settings.Search.UseFP)},
		"Use_Qfp": {NameID: "Use_Qfp", HandlerFunc: useQfp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseQFP), CurrentValue: strconv.FormatBool(Settings.Search.UseQFP)},
		"Use_Lmr": {NameID: "Use_Lmr", HandlerFunc: useLmr, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseLmr), CurrentValue: strconv.FormatBool(Settings.Search.UseLmr)},
		"Use_Lmp": {NameID: "Use_Lmp", HandlerFunc: useLmp, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseLmp), CurrentValue: strconv.FormatBool(Settings.Search.UseLmp)},

		"Use_Ext":         {NameID: "Use_Ext", HandlerFunc: useExt, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseExt), CurrentValue: strconv.FormatBool(Settings.Search.UseExt)},
		"Use_ExtAddDepth": {NameID: "Use_ExtAddDepth", HandlerFunc: useExtAddDepth, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseExtAddDepth), CurrentValue: strconv.FormatBool(Settings.Search.UseExtAddDepth)},
		"Use_CheckExt":    {NameID: "Use_CheckExt", HandlerFunc: useCheckExt, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseCheckExt), CurrentValue: strconv.FormatBool(Settings.Search.UseCheckExt)},
		"Use_ThreatExt":   {NameID: "Use_ThreatExt", HandlerFunc: useThreatExt, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseThreatExt), CurrentValue: strconv.FormatBool(Settings.Search.UseThreatExt)},

		"Use_EvalHash":       {NameID: "Use_EvalHash", HandlerFunc: useEvalHash, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Search.UseEvalTT), CurrentValue: strconv.FormatBool(Settings.Search.UseEvalTT)},
		"Eval_Lazy":          {NameID: "Eval_Lazy", HandlerFunc: evalLazy, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Eval.UseLazyEval), CurrentValue: strconv.FormatBool(Settings.Eval.UseLazyEval)},
		"Eval_Pawns":         {NameID: "Eval_Pawns", HandlerFunc: evalPawns, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Eval.UsePawnEval), CurrentValue: strconv.FormatBool(Settings.Eval.UsePawnEval)},
		"Eval_PawnCache":     {NameID: "Eval_PawnCache", HandlerFunc: evalPawnCache, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Eval.UsePawnCache), CurrentValue: strconv.FormatBool(Settings.Eval.UsePawnCache)},
		"Eval_PawnCacheSize": {NameID: "Eval_PawnCacheSize", HandlerFunc: pawnCacheSize, OptionType: Spin, DefaultValue: strconv.Itoa(Settings.Eval.PawnCacheSize), CurrentValue: strconv.Itoa(Settings.Eval.PawnCacheSize), MinValue: "0", MaxValue: "1024"},

		// TODO vvv
		"Eval_Mobility": {NameID: "Eval_Mobility", HandlerFunc: evalMob, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Eval.UseMobility), CurrentValue: strconv.FormatBool(Settings.Eval.UseMobility)},
		"Eval_AdvPiece": {NameID: "Eval_AdvPiece", HandlerFunc: evalAdv, OptionType: Check, DefaultValue: strconv.FormatBool(Settings.Eval.UseAdvancedPieceEval), CurrentValue: strconv.FormatBool(Settings.Eval.UseAdvancedPieceEval)},
	}

	// as Go does not let us sort maps manually (e.g. insertion order) we need a
	// list to control order of the entries
	sortOrderUciOptions = []string{
		"Print Config",

		"OwnBook",

		"Clear Hash",
		"Use_Hash",
		"Hash",
		"Ponder",

		"Quiescence",
		"Use_QHash",
		"Use_SEE",
		"Use QS Standpat",
		"Use_EvalHash",

		"Use_AlphaBeta",
		"Use_PVS",
		"Use_ASP",

		"Use_TTMove",
		"Use_IID",
		"Use_Killer",
		"Use_HistCount",
		"Use_CounterMove",

		"Use_TTValue",
		"Use_Mdp",
		"Use_Razoring",
		"Razor Margin",
		"Use_Rfp",

		"Use_NullMove",
		"Use_Fp",
		"Use_Qfp",
		"Use_Lmr",
		"Use_Lmp",

		"Use_Ext",
		"Use_ExtAddDepth",
		"Use_CheckExt",
		"Use_ThreatExt",

		"Eval_Pawns",
		"Eval_PawnCache",
		"Eval_PawnCacheSize",
		"Eval_Mobility",
		"Eval_AdvPiece",
	}
}

// GetOptions returns all available uci options as a slice of strings
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

func printConfig(handler *UciHandler, _ *uciOption) {
	s := reflect.ValueOf(&Settings.Eval).Elem()
	typeOfT := s.Type()
	for i := s.NumField() - 1; i >= 0; i-- {
		f := s.Field(i)
		handler.SendInfoString(fmt.Sprintf("%-2d: %-22s %-6s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface()))
	}
	handler.SendInfoString("Evaluation Config:\n")
	s = reflect.ValueOf(&Settings.Search).Elem()
	typeOfT = s.Type()
	for i := s.NumField() - 1; i >= 0; i-- {
		f := s.Field(i)
		handler.SendInfoString(fmt.Sprintf("%-2d: %-22s %-6s = %v\n", i, typeOfT.Field(i).Name, f.Type(), f.Interface()))
	}
	handler.SendInfoString("Search Config:\n")
	log.Debug(Settings.String())

}

func clearCache(u *UciHandler, _ *uciOption) {
	u.mySearch.ClearHash()
	log.Debug("Cleared Cache")
}

func useCache(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseTT = v
	log.Debugf("Set Use Hash to %v", Settings.Search.UseTT)
}

func cacheSize(u *UciHandler, o *uciOption) {
	v, _ := strconv.Atoi(o.CurrentValue)
	Settings.Search.TTSize = v
	u.mySearch.ResizeCache()
}

func useBook(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseBook = v
	log.Debugf("Set Use Book to %v", Settings.Search.UseBook)
}

func usePonder(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UsePonder = v
	log.Debugf("Set Use Ponder to %v", Settings.Search.UsePonder)
}

func useQuiescence(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseQuiescence = v
	log.Debugf("Set Use Quiescence to %v", Settings.Search.UseQuiescence)
}

func useQSHash(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseQSTT = v
	log.Debugf("Set Use Hash in Quiescence to %v", Settings.Search.UseQSTT)
}

func useAlphaBeta(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseAlphaBeta = v
	log.Debugf("Set Use AlphaBeta Search to %v", Settings.Search.UseAlphaBeta)
}

func usePvs(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UsePVS = v
	log.Debugf("Set Use Principal Variation Search to %v", Settings.Search.UsePVS)
}

func useAsp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseAspiration = v
	log.Debugf("Set Use Aspiration Search to %v", Settings.Search.UseAspiration)
}

func useTTMove(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseTTMove = v
	log.Debugf("Set Use Hash Move to %v", Settings.Search.UseTTMove)
}

func useTTValue(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseTTValue = v
	log.Debugf("Set Use Hash Value to %v", Settings.Search.UseTTValue)
}

func useMdp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseMDP = v
	log.Debugf("Set Use MDP to %v", Settings.Search.UseMDP)
}

func useKiller(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseKiller = v
	log.Debugf("Set Use Killer Moves to %v", Settings.Search.UseKiller)
}

func useHC(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseHistoryCounter = v
	log.Debugf("Set Use History Counter to %v", Settings.Search.UseHistoryCounter)
}

func useCM(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseCounterMoves = v
	log.Debugf("Set Use Counter Moves to %v", Settings.Search.UseCounterMoves)
}

func useRazor(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseRazoring = v
	log.Debugf("Set Use Razoring to %v", Settings.Search.UseRazoring)
}

func razorMargin(_ *UciHandler, o *uciOption) {
	v, _ := strconv.Atoi(o.CurrentValue)
	Settings.Search.RazorMargin = v
	log.Debugf("Set Razor Margin to %v", Settings.Search.RazorMargin)
}

func useNullMove(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseNullMove = v
	log.Debugf("Set Use Null Move Pruning to %v", Settings.Search.UseNullMove)
}

func useIID(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseIID = v
	log.Debugf("Set Use IID to %v", Settings.Search.UseIID)
}

func useLmr(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseLmr = v
	log.Debugf("Set use Late Move Reduction to %v", Settings.Search.UseLmr)
}

func useLmp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseLmp = v
	log.Debugf("Set use Late Move Pruning to %v", Settings.Search.UseLmp)
}

func useQSStandpat(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseQSStandpat = v
	log.Debugf("Set use QS Standpat to %v", Settings.Search.UseQSStandpat)
}

func useSee(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseSEE = v
	log.Debugf("Set use SEE to %v", Settings.Search.UseSEE)
}

func useExt(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseExt = v
	log.Debugf("Set use Extensions to %v", Settings.Search.UseExt)
}

func useExtAddDepth(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseExtAddDepth = v
	log.Debugf("Set use Extensions Add to Depth to %v", Settings.Search.UseExtAddDepth)
}

func useCheckExt(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseCheckExt = v
	log.Debugf("Set use Check Extension to %v", Settings.Search.UseCheckExt)
}

func useThreatExt(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseThreatExt = v
	log.Debugf("Set use Threat Extension to %v", Settings.Search.UseThreatExt)
}

func useRfp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseRFP = v
	log.Debugf("Set use Reverse Futility Pruning (RFP) to %v", Settings.Search.UseRFP)
}

func useFp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseFP = v
	log.Debugf("Set use Futility Pruning (FP) to %v", Settings.Search.UseFP)
}

func useQfp(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseQFP = v
	log.Debugf("Set use Quiescence Futility Pruning (QFP) to %v", Settings.Search.UseQFP)
}

func useEvalHash(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Search.UseEvalTT = v
	log.Debugf("Set use Eval Hash to %v", Settings.Search.UseEvalTT)
}

func evalLazy(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Eval.UseLazyEval = v
	log.Debugf("Set use Lazy Eval to %v", Settings.Eval.UseLazyEval)
}

func evalPawns(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Eval.UsePawnEval = v
	log.Debugf("Set use Eval Pawns to %v", Settings.Eval.UsePawnEval)
}

func evalPawnCache(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Eval.UsePawnCache = v
	log.Debugf("Set use Eval Pawn Cache to %v", Settings.Eval.UsePawnCache)
}

func pawnCacheSize(_ *UciHandler, o *uciOption) {
	v, _ := strconv.Atoi(o.CurrentValue)
	Settings.Eval.PawnCacheSize = v
	log.Debugf("Set use Pawn Cache Size to %v", Settings.Eval.PawnCacheSize)
}

func evalMob(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Eval.UseMobility = v
	log.Debugf("Set use Eval Mobility to %v", Settings.Eval.UseMobility)
}

func evalAdv(_ *UciHandler, o *uciOption) {
	v, _ := strconv.ParseBool(o.CurrentValue)
	Settings.Eval.UseAdvancedPieceEval = v
	log.Debugf("Set use Adv Piece Eval to %v", Settings.Eval.UseAdvancedPieceEval)
}
