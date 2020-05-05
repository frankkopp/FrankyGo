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

package testsuite

import (
	"testing"
	"time"

	"github.com/frankkopp/FrankyGo/internal/config"
)

func TestFeatureTests(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}

	// setup tests
	searchTime := 2000 * time.Millisecond
	searchDepth := 0

	// Feature Settings
	{
		config.Settings.Search.UseQuiescence = true
		config.Settings.Search.UseQSStandpat = true
		config.Settings.Search.UseSEE = true
		config.Settings.Search.UsePromNonQuiet = true

		config.Settings.Search.UseTT = true
		config.Settings.Search.TTSize = 256
		config.Settings.Search.UseTTValue = true
		config.Settings.Search.UseQSTT = true

		config.Settings.Search.UsePVS = false
		config.Settings.Search.UseAspiration = false
		config.Settings.Search.UseMTDf = true

		config.Settings.Search.UseTTMove = true
		config.Settings.Search.UseIID = true
		config.Settings.Search.IIDDepth = 6
		config.Settings.Search.IIDReduction = 2
		config.Settings.Search.UseKiller = true
		config.Settings.Search.UseHistoryCounter = true
		config.Settings.Search.UseCounterMoves = true

		config.Settings.Search.UseMDP = true
		config.Settings.Search.UseRazoring = true
		config.Settings.Search.RazorMargin = 531
		config.Settings.Search.UseNullMove = true
		config.Settings.Search.NmpDepth = 3
		config.Settings.Search.NmpReduction = 2

		config.Settings.Search.UseExt = true
		config.Settings.Search.UseExtAddDepth = true
		config.Settings.Search.UseCheckExt = true
		config.Settings.Search.UseThreatExt = false

		config.Settings.Search.UseRFP = true
		config.Settings.Search.UseFP = true
		config.Settings.Search.UseQFP = true
		config.Settings.Search.UseLmr = true
		config.Settings.Search.LmrDepth = 3
		config.Settings.Search.LmrMovesSearched = 3
		config.Settings.Search.UseLmp = true

		config.Settings.Eval.Tempo = 34
		config.Settings.Eval.UseLazyEval = true
		config.Settings.Eval.LazyEvalThreshold = 700

		config.Settings.Eval.UsePawnCache = false
		config.Settings.Eval.PawnCacheSize = 64
		config.Settings.Eval.UseAttacksInEval = false
		config.Settings.Eval.UseMobility = false
		config.Settings.Eval.MobilityBonus = 5
		config.Settings.Eval.UseAdvancedPieceEval = false
		config.Settings.Eval.BishopPairBonus = 20
		config.Settings.Eval.MinorBehindPawnBonus = 15
		config.Settings.Eval.BishopPawnMalus = 5
		config.Settings.Eval.BishopCenterAimBonus = 20
		config.Settings.Eval.BishopBlockedMalus = 40
		config.Settings.Eval.RookOnQueenFileBonus = 6
		config.Settings.Eval.RookOnOpenFileBonus = 25
		config.Settings.Eval.RookTrappedMalus = 40
		config.Settings.Eval.KingRingAttacksBonus = 10
		config.Settings.Eval.UseKingEval = false
		config.Settings.Eval.KingDangerMalus = 50
		config.Settings.Eval.KingDefenderBonus = 10
	}

	folder := "test/testdata/featuretests/"

	out.Println(FeatureTests(folder, searchTime, searchDepth)) // nolint:errcheck
}

// ///////////////////////////////
// PVS
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 11:06:38.5668289 +0200 CEST m=+259.965905801
// Test took            : 4m19.2729095s
// Test setup           : search time: 200ms max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       40,9 % |     143.077.661 |        141 |        204 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       36,5 % |     306.780.548 |        281 |        488 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |      100,0 % |       6.749.903 |         13 |          0 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       35,0 % |       9.722.408 |          7 |         13 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |       20,0 % |       2.155.713 |          1 |          4 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       83,6 % |      82.423.542 |        168 |         33 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       45,2 % |     550.909.775 |        611 |        742 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 4m6.1370595s
// Total NPS : 2.238.223
//
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 12:17:59.4372915 +0200 CEST m=+2657.453594601
// Test took            : 44m16.7956001s
// Test setup           : search time: 2s max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       45,8 % |   1.579.318.667 |        158 |        187 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       57,1 % |   3.338.053.303 |        439 |        330 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |      100,0 % |      65.628.211 |         13 |          0 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       65,0 % |     119.331.186 |         13 |          7 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |       20,0 % |      21.012.160 |          1 |          4 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       94,0 % |     879.759.139 |        189 |         12 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       60,1 % |   6.003.102.666 |        813 |        540 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 44m3.4443158s
// Total NPS : 2.270.939

// ASP
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 11:12:36.7468252 +0200 CEST m=+260.006716801
// Test took            : 4m19.2627156s
// Test setup           : search time: 200ms max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       45,2 % |     144.357.628 |        156 |        189 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       37,2 % |     306.910.130 |        286 |        483 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |      100,0 % |       6.600.551 |         13 |          0 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       30,0 % |       9.467.562 |          6 |         14 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |       20,0 % |       2.027.031 |          1 |          4 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       85,1 % |      82.089.578 |        171 |         30 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       46,8 % |     551.452.480 |        633 |        720 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 4m6.1698369s
// Total NPS : 2.240.130
//
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 13:36:04.3200566 +0200 CEST m=+2658.662874701
// Test took            : 44m17.8708741s
// Test setup           : search time: 2s max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       49,0 % |   1.505.239.986 |        169 |        176 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       57,9 % |   3.291.551.493 |        445 |        324 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |      100,0 % |      61.163.719 |         13 |          0 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       55,0 % |     110.665.561 |         11 |          9 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |       20,0 % |      21.848.521 |          1 |          4 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       93,0 % |     890.215.834 |        187 |         14 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       61,0 % |   5.880.685.114 |        826 |        527 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 44m4.1699958s
// Total NPS : 2.224.019

// MTDf
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 11:17:41.8421966 +0200 CEST m=+259.666185201
// Test took            : 4m19.0949783s
// Test setup           : search time: 200ms max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       38,0 % |     105.888.952 |        131 |        214 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       49,3 % |     236.228.283 |        379 |        390 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |       92,3 % |       5.326.075 |         12 |          1 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       30,0 % |       7.741.903 |          6 |         14 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |        0,0 % |       1.787.203 |          0 |          5 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       83,1 % |      70.826.847 |        167 |         34 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       51,4 % |     427.799.263 |        695 |        658 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 4m6.0958508s
// Total NPS : 1.738.344
//
// Feature Test Result Report
// ==============================================================================
// Date                 : 2020-05-05 14:23:50.6109994 +0200 CEST m=+2656.129210801
// Test took            : 44m15.5142827s
// Test setup           : search time: 2s max depth: 0
// Number of testsuites : 6
// Number of tests      : 1.353
//
// ===============================================================================================================================================
// Test Suite                | Success Rate |           Nodes | Successful |     Failed |    Skipped |        N/A |   Tests | File
// ===============================================================================================================================================
// crafty_test.epd           |       43,8 % |   1.117.957.463 |        151 |        194 |          0 |          0 |     345 | test/testdata/featuretests/crafty_test.epd
// ecm98.epd                 |       67,2 % |   2.727.672.876 |        517 |        252 |          0 |          0 |     769 | test/testdata/featuretests/ecm98.epd
// franky_tests.epd          |      100,0 % |      53.577.083 |         13 |          0 |          0 |          0 |      13 | test/testdata/featuretests/franky_tests.epd
// mate_test_suite.epd       |       60,0 % |     106.642.447 |         12 |          8 |          0 |          0 |      20 | test/testdata/featuretests/mate_test_suite.epd
// nullMoveZugZwangTest.epd  |        0,0 % |      19.460.635 |          0 |          5 |          0 |          0 |       5 | test/testdata/featuretests/nullMoveZugZwangTest.epd
// wac.epd                   |       93,5 % |     737.292.489 |        188 |         13 |          0 |          0 |     201 | test/testdata/featuretests/wac.epd
// -----------------------------------------------------------------------------------------------------------------------------------------------
// TOTAL                     |       65,1 % |   4.762.602.993 |        881 |        472 |          0 |          0 |   1.353 |
// ===============================================================================================================================================
//
// Total Time: 44m2.4647145s
// Total NPS : 1.802.333
