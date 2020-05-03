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

package testsuite

import (
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/frankkopp/FrankyGo/internal/config"
	"github.com/frankkopp/FrankyGo/internal/util"
)

// FeatureTests runs all epd tests in a folder and prints a report
func FeatureTests(folder string, searchTime time.Duration, searchDepth int) string {

	// get all tests in folder
	files, err := ioutil.ReadDir(folder)
	if err != nil {
		log.Fatal(err)
	}
	var list []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".epd" {
			list = append(list, f.Name())
		}
	}

	// prepare
	config.Settings.Search.UseBook = false
	result := make(map[string]TestSuite, 10)
	executedTests := 0
	totalNodes := uint64(0)
	var totalTime time.Duration
	totalSuccess := 0
	totalFailed := 0
	totalSkipped := 0
	totalNotTested := 0
	totalTests := 0

	// run all test files
	start := time.Now()
	for _, t := range list {

		// Run test
		ts, _ := NewTestSuite(folder+t, searchTime, searchDepth)
		ts.RunTests()

		// save result
		executedTests += len(ts.Tests)
		result[t] = *ts
	}
	duration := time.Since(start)

	// sort order
	keys := make([]string, 0, len(result))
	for name := range result {
		keys = append(keys, name)
	}
	sort.Strings(keys)

	// print report
	os := strings.Builder{}
	os.WriteString(out.Sprintf("Feature Test Result Report\n"))
	os.WriteString(out.Sprintf("==============================================================================\n"))
	os.WriteString(out.Sprintf("Date                 : %s\n", time.Now()))
	os.WriteString(out.Sprintf("Test took            : %s\n", duration))
	os.WriteString(out.Sprintf("Test setup           : search time: %s max depth: %d\n", searchTime, searchDepth))
	os.WriteString(out.Sprintf("Number of testsuites : %d\n", len(result)))
	os.WriteString(out.Sprintf("Number of tests      : %d\n", executedTests))
	os.WriteString(out.Sprintln())
	os.WriteString(out.Sprintf("===============================================================================================================================================\n"))
	os.WriteString(out.Sprintf(" %-25s | %-12s | %-15s | %-10s | %-10s | %-10s | %-10s | %-6s | %s\n", "Test Suite", "Success Rate", "          Nodes", "Successful", "    Failed", "   Skipped", "       N/A", "  Tests", "File"))
	os.WriteString(out.Sprintf("===============================================================================================================================================\n"))
	for _, name := range keys {
		r := result[name]
		successRate := float64(r.LastResult.SuccessCounter) / float64(r.LastResult.Counter) * 100
		totalNodes += r.LastResult.Nodes
		totalTime += r.LastResult.Time
		totalSuccess += r.LastResult.SuccessCounter
		totalFailed += r.LastResult.FailedCounter
		totalSkipped += r.LastResult.SkippedCounter
		totalNotTested += r.LastResult.NotTestedCounter
		totalTests += r.LastResult.Counter
		// os.WriteString(out.Sprintf(" %-25s |      %5.1f %% | %11d |   %8d |   %8d |   %8d |   %8d |  %6d | %s\n", name, 100.0, 999999999, 99999, 99999, 99999, 99999, 9999, folder+name))
		os.WriteString(out.Sprintf(" %-25s |      %5.1f %% | %15d |   %8d |   %8d |   %8d |   %8d |  %6d | %s\n", name, successRate, r.LastResult.Nodes, r.LastResult.SuccessCounter, r.LastResult.FailedCounter, r.LastResult.SkippedCounter, r.LastResult.NotTestedCounter, len(r.Tests), folder+name))
	}
	successRate := float64(totalSuccess) / float64(totalTests) * 100
	os.WriteString(out.Sprintf("-----------------------------------------------------------------------------------------------------------------------------------------------\n"))
	os.WriteString(out.Sprintf(" %-25s |      %5.1f %% | %15d |   %8d |   %8d |   %8d |   %8d |  %6d | %s\n", "TOTAL", successRate, totalNodes, totalSuccess, totalFailed, totalSkipped, totalNotTested, totalTests, ""))
	os.WriteString(out.Sprintf("===============================================================================================================================================\n"))
	os.WriteString(out.Sprintln())
	os.WriteString(out.Sprintf("Total Time: %s\n", totalTime))
	os.WriteString(out.Sprintf("Total NPS : %d\n", util.Nps(totalNodes, totalTime)))
	os.WriteString(out.Sprintln())
	os.WriteString(out.Sprintf("Configuration: %s\n", config.Settings.String()))
	os.WriteString(out.Sprintln())

	return os.String()
}
