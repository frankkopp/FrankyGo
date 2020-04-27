# FrankyGo
Go implementation of a UCI compatible chess engine.

[![Build Status](https://travis-ci.org/frankkopp/FrankyGo.svg?branch=master)](https://travis-ci.org/frankkopp/FrankyGo)
[![codecov](https://codecov.io/gh/frankkopp/FrankyGo/branch/master/graph/badge.svg)](https://codecov.io/gh/frankkopp/FrankyGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/frankkopp/FrankyGo)](https://goreportcard.com/report/github.com/frankkopp/FrankyGo)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/frankkopp/WorkerPool/blob/master/LICENSE)

## Features (v0.9.0)
* UCI protocol (to use the engine in Arena, xboard, Fritz or other chess user interfaces)
    * UCI Options
    * UCI Search Modes
* FIDE compliant (move repetition, 50-moves rules, draw for low material)
* Board representation bitboards and 8x8 piece array
* Bitboards, Magic Bitboards, pre-computed bitboards and data 
* Move Generation using bitboards 
    * all moves or on demand in phases 
    * mode: all, capture only, non capture, check evasions
    * move sorting for estimated value of move, pv move and killer moves
    * ~50 Million moves per sec on my i7-7700
* Perft tests
    * ~4.8+Million nps 
* Search as PVS AlphaBeta search
    * all UCI search modes 
    * Iterative Deepening
    * Quiescence
    * Pondering
    * Killer Moves
    * Mate Distance Pruning
    * Null Move Pruning
    * Internal Iterative Deepening
    * Late Move Reductions (needs further tuning)
    * Late Move Pruning (needs further tuning)
    * Check extension (needs further tuning/testing)
    * Mate threat extension (deactivated - needs further tuning/testing - currently weakens play)
    * SEE for deciding on good qsearch moves
    * Reverse Futility Pruning (Static Null Move)
    * Futility Pruning
    * History Pruning (History Count, Counter Move) 
    * ~2.5-3.5M nps
    * TODO: Aspiration Windows, Multi-cut Pruning, History Pruning 
* Transposition Table
* Opening Book from PGN, SAN and Simple move list files and persistent cache
* Evaluation
    * Very simple yet: Material, positional piece values
    * Implemented but not yet tested/activated: Attacks, Mobility, Piece specific evaluations
    * TODO: pawn evaluations and pawn hash table  
* Tests:   
    * Test framework to run EPD test suites
    * Test to determine search tree size for features
* General: 
    * logging, command lines options, configuration files, search and eval parameters in config files

* Open topics: 
    * Parallel search
            
     
## Learning Go
Learning a new programming language is always most efficient within a real project.
As I have developed a chess engine in Java in the past and recently migrated it to C++ 
I thought that implementing a chess engine in Go is a great opportunity to learn the language. 

Chess engines offer a lot of different challenges in many aspects of a programming language. 
E.g. efficient data types and data structures, high performance code where even nanoseconds count, reading 
an opening book from files and creating a cache, bit twiddling, caches, recursion, unit testing, 
performance testing and optimization, logging, configurations, communication with other processes over 
pipes (UCI protocol), protocol implementation, etc. 
 
Especially the high performance code is something I'm interested in. The usual advice in programming 
is readability/understandability/maintainability of code is more important than performance and also, don't optimize 
too early, etc. In chess there are some recurring hot spots in the algorithms, which need extra attention 
on how they are implemented from the beginning, or they would be very slow. So a chess engine is a special 
case where efficiency and performance has a higher priority than in a typical application.  
 
E.g. a chess engine evaluates (calculating a numerical estimation of how "good" a position is) millions
of positions per second. For each position there are in average 35 moves to play. So to be able to search 
a million positions per second you need to generate at least the same amount of moves - usually many more 
as the alpha-beta search algorithm and other pruning techniques discard many positions before
they are even evaluated.   

For evaluating 1 million positions per second you have a time budget of 1.000 nanoseconds for creating 
the move that leads to this position, executing the move on your chess board data structure, evaluate 
the position, and the many other things you need to do so that your chess engines has at least some
playing strength. So you see, nanoseconds matter. 
  
As an example many years ago, during my first attempt at this in Java, it was quickly clear, that you can't model 
all of your data types as Java objects. Especially a chess move would be too expensive when implemented 
as a Java object as you need to create hundreds of millions in a very short time just to throw them away 
quickly. Object creation and garbage collection made this approach extremely slow. So one of the first "optimizations"
I did was to use a plain integer to represent a move. With bit twiddling I encoded the necessary data into the 
integer (from square, to square, move type, etc.). This is a very common approach in all the serious chess engines. 

This example also shows a big difference between Java, C++ and Go. In Java, when implementing a move as an int you
loose type checking by the compiler. In C++ you have typedef and in Go you have something similar to typdef. With a 
user defined primitive type the engine can use a plain int internally, but the compiler still can check the type 
correctness for moves. This is one of my most missed features in Java. 

What I already can say about Go is that implementing the chess engine in it was a joy and much easier done in Go than 
in Java or C++. Of course, I have experience building a chess engine now but mostly I noticed that in Go I could focus 
on the program itself and was rarely distracted by the "quirks" of the programming language although it was new to me. 
I'm aware I did not follow many of Go's best practices but rather "transferred" my existing code from C++/Java to it. 
Also, some aspects needed some re-thinking (e.g. no cyclic imports) but overall Go required much less language awareness
than for example C++ where I still find the compile-process unnecessary complicated (.h, cpp - declaration and 
definition, etc.) or Java where I can't have user defined primitive types.

## Installation
Windows build: 

go build -o FrankyGo.exe github.com/frankkopp/FrankyGo/cmd/FrankyGo

Run FrankyGo.exe

Unix/Mac build: 

go build -o FrankyGo github.com/frankkopp/FrankyGo/cmd/FrankyGo

Run FrankyGo

## Usage
Typically, a UCI engine will be used with a UCI compatible chess interface like Arena, xboard, Fritz, etc.
Just configure a new engine in the UCI interface pointing to the executable of FrankyGo. If necessary use 
command line options to find config file, logs folder and opening books folder.    

To run FrankyGo needs to find its config file:

Default place to look for it is ../configs/config.toml

Use command line option "-config" to change location

Also helpful:
* logs files are stored in folder "../logs" or folder defined by cmd line option "-logpath"
* opening books are searched in folder "../assets/book" or folder defined in cmd line option "-bookpath"

Use --help for more command line options

## Roadmap
### vx.x (planned)
- TODO:
    - Better Evaluation and testing
    - Pawn Structure Cache
    - MultiCut Pruning
       - https://hci.iwr.uni-heidelberg.de/system/files/private/downloads/1935772097/report_qingyang-cao_enhanced-forward-pruning.pdf
    - Other Prunings
    - Aspiration
    - Continuously Performance/Profiling/Testing
    - Tuning and Testing of all search features and parameters

## Versions
### v0.9.0 (in progress)
- TODO:
    - History Heuristics 

### v0.8.0 (done)
- DONE
    - SEE
    - Move generation considering check evasion
    - Reverse Futility Pruning
    - Futility Pruning
    - Search extensions:
        - Check extension
        - Mate threat extension (not active - search tree gets too big)
    - Restructuring of packages to better match Go best practices    

### v0.7 (done)
- DONE
    - LMP and LMR 
    - Null Move
    - Use TestSuites, TreeSize and Arena to test features
    - Better Evaluation (not active in config yet - needs testing)
    - Performance/Profiling/Testing
    - Removed MPP (Minor Promotion Pruning) - more harm than benefit

### v0.6 (done)
- DONE
    - Enhance TestSuite / run from command line options
    - TestSuite Tests
    - PVS
    - Killer
    - TT in QS
    - MDP/MPP

- Remarks
    PVS and TT might have some dependencies which I have not fully understood yet.
    Some engines for example only cut with TT values for alpha/beta value in non PV nodes.
    Tests show no drop in Search strength either way and also search tree size shows no
    obvious issues

### v0.5 (done)
- DONE
    - Use TT
    - SearchTreeSize
    - Quiescence search
    - Score as string()
    - Evaluation (simple)

### v0.4 (done)
- DONE
    - Pondering
    - Testing for correct play with Arena against Stockfish
    - log files for standard and search log
    - Implement simple Evaluator
    - Complete simple search
    - Implement simple Search

### v0.3 (done)
- DONE:
    - CleanUp and additional documentation
    - complete uci options
    - add log files to command line options
    - make log files configurable
    - make book configurable
    - added uci options
    - added configuration via file and command line
    - search handling
    - starting / stopping
    - time control
    - limits (except depth limit - needs simple search minimax)

### v0.2 (done)
- DONE
    - CleanUp
    - Added logging
    - TranspositionTable
    - Perft enhanced and more tests
    - UCI Handler enhanced
    - Completed MoveGen
    - OpeningBook base framework (reading and caching)
    - Improve performance of Perft - otherwise not worth continuing
    - Added MoveSlice - little optimization of MoveArray - usable directly as Slice

### v0.1 (done)
- DONE:
    - Perft works (1.353.761 nps / Java 3.5M, C++ 4.5M) - needs improvement
    - MoveGenerator (all required for perft)
    - Position (all required for perft)
    - MoveArray and MoveList - both are for list of moves - MoveArray is faster for sorting
        might be slower when inserting at the front - needs testing
    - Most types (all required for perft)
