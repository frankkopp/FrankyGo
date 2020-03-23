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

// //////////////////////////////////////////////////////
// Statistics
// //////////////////////////////////////////////////////

// Statistics extra data and stats not essential for a functioning search
type Statistics struct {

}

// // counter for cut off to measure quality of move ordering
//  std::array<uint64_t, MAX_MOVES> betaCutOffs{};
//  std::array<uint64_t, MAX_MOVES> alphaImprovements{};
//
//  // Search info values
//  Ply currentSearchDepth = PLY_ROOT;
//  Ply currentExtraSearchDepth = PLY_ROOT;
//  Move currentRootMove = MOVE_NONE;
//  int bestMoveChanges = 0;
//  int bestMoveDepth = 0;
//  MilliSec lastSearchTime = 0;
//
//  // performance statistics
//  uint64_t movesGenerated = 0;
//  uint64_t nodesVisited = 0; // legal nodes visited
//
//  // PERFT Values
//  uint64_t leafPositionsEvaluated = 0;
//  uint64_t nonLeafPositionsEvaluated = 0;
//  uint64_t checkCounter = 0;
//  uint64_t checkMateCounter = 0;
//  uint64_t captureCounter = 0;
//  uint64_t enPassantCounter = 0;
//
//  // TT Statistics
//  uint64_t tt_Cuts = 0;
//  uint64_t tt_NoCuts = 0;
//
//  // Optimization Values
//  uint64_t aspirationResearches = 0;
//  uint64_t prunings = 0;
//  uint64_t pvs_root_researches = 0;
//  uint64_t pvs_root_cutoffs = 0;
//  uint64_t pvs_researches = 0;
//  uint64_t pvs_cutoffs = 0;
//  uint64_t pv_sortings = 0;
//  uint64_t no_moveForPVsorting = 0;
//  uint64_t positionsNonQuiet = 0;
//  uint64_t qStandpatCuts = 0;
//  uint64_t minorPromotionPrunings = 0;
//  uint64_t mateDistancePrunings = 0;
//  uint64_t nullMovePrunings = 0;
//  uint64_t nullMoveVerifications = 0;
//  uint64_t extensions = 0;
//  uint64_t rfpPrunings = 0;
//  uint64_t razorReductions = 0;
//  uint64_t iidSearches = 0;
//  uint64_t lrReductions = 0;
//  uint64_t efpPrunings = 0;
//  uint64_t fpPrunings = 0;
//  uint64_t qfpPrunings = 0;
//  uint64_t lmpPrunings = 0;
//  uint64_t lmrReductions = 0;
//
//  uint64_t deltaPrunings = 0;
