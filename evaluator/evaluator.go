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

// Package evaluator contains structures and functions to calculate
// the value of a chess position to be used in a chess engine search
package evaluator

import (
	"github.com/op/go-logging"

	"github.com/frankkopp/FrankyGo/config"
	myLogging "github.com/frankkopp/FrankyGo/logging"
	"github.com/frankkopp/FrankyGo/position"
	. "github.com/frankkopp/FrankyGo/types"
)

const trace = true

// Evaluator  represents a data structure and functionality fo
// evaluating chess positions by using various evaluation
// heuristics like material, positional values, pawn structure, etc.
//  Create a new instance with NewEvaluator()
type Evaluator struct {
	log *logging.Logger
}

// NewEvaluator creates a new instance of an Evaluator.
func NewEvaluator() *Evaluator {
	return &Evaluator{
		log: myLogging.GetLog(),
	}
}

// Evaluate calculates a value for a chess positions by
// using various evaluation heuristics like material,
// positional values, pawn structure, etc.
func (e *Evaluator) Evaluate(position *position.Position) Value {
	var value Value = 0

	gamePhaseFactor := float64(position.GamePhase() / GamePhaseMax)

	// Each position is evaluated from the view of the white
	// player. Before returning the value this will be adjusted
	// to the next player's color.
	// All heuristic should return a value in centi pawns or
	// have a dedicated configurable weight to adjust and test

	// Material
	value += e.material(position, gamePhaseFactor)

	// Positional values
	value += e.positional(position, gamePhaseFactor)

	// value is always from the view of the next player
	if position.NextPlayer() == Black {
		value *= -1
	}

	// TEMPO Bonus for the side to move (helps with evaluation alternation -
	// less difference between side which makes aspiration search faster
	// (not empirically tested)
	value += Value(float64(config.Settings.Eval.Tempo) * gamePhaseFactor)

	return value
}

func (e *Evaluator) material(position *position.Position, gamePhaseFactor float64) Value {
	return position.Material(White) - position.Material(Black)
}

func (e *Evaluator) positional(position *position.Position, gamePhaseFactor float64) Value {
	return Value(float64(position.PsqMidValue(White)-position.PsqMidValue(Black))*gamePhaseFactor +
		float64(position.PsqEndValue(White)-position.PsqEndValue(Black))*(1-gamePhaseFactor))
}
