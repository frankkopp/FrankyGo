# FrankyGo
Go implementation of a UCI compatible chess engine.

[![Build Status](https://travis-ci.org/frankkopp/FrankyGo.svg?branch=master)](https://travis-ci.org/frankkopp/FrankyGo)
[![codecov](https://codecov.io/gh/frankkopp/FrankyGo/branch/master/graph/badge.svg)](https://codecov.io/gh/frankkopp/FrankyGo)
[![Go Report Card](https://goreportcard.com/badge/github.com/frankkopp/FrankyGo)](https://goreportcard.com/report/github.com/frankkopp/FrankyGo)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/frankkopp/WorkerPool/blob/master/LICENSE)

## Description
FrankyGo is the fourth chess engine I wrote. The first one was written in Java and had a JavaFX UI 
([Chessly](https://github.com/frankkopp/Chessly)) but did not support UCI. I was very inexperienced then, and the 
engine itself was not good. But I learned a lot and also liked developing the JavaFx part. I then rewrote it in Java 
again but as an UCI engine. This time the structure was much better thanks to many great chess engine developers I 
learned from. I also found out that my interest is in actual programming rather than building the "best" chess engine. 
Therefore, I used what I learned from other chess engines but tried to build mine in a clean and easy to understand way. 
This resulted in Franky 1.0
([Github](https://github.com/frankkopp/FrankyUCIChessEngine),
[CCRL](http://www.computerchess.org.uk/ccrl/404/cgi/engine_details.cgi?print=Details&each_game=1&eng=Franky%201.0%2064-bit#Franky_1_0_64-bit))
 
I then started to re-write it in C++ ([FrankyCPP](https://github.com/frankkopp/FrankyCPP)) and came quite far. But C++ 
(and its infrastructure) forced me to focus a lot on the language itself, and the surrounding tools. E.g. I learned a lot 
(too much?) about different compilers on different platforms, CMake, Boost, Google Test, etc. but it rather distracted 
me a lot from developing the chess engine itself. My time is limited so this was very annoying. 

Then I looked into [Go](https://golang.org/doc/). A completely new language compiled directly into machine code, 
platform independent, with a garbage collector and the promise to be as fast as C++/C. This was the motivation to start 
FrankyGo and Go kept its promises. Easy to learn, great tool set, no distraction from the code by the language and 
indeed fast. See below "Learning Go".

So FrankyGo is now a rather clean chess engine with a lot of potential. As mentioned above I do not claim to have 
invented any of algorithms and in fact I used many ideas from Stockfish, Crafty, Beowulf, etc. See "Credits" below.
 But I had a lot of fun developing it and I do hope to have more fun improving it in the future.

There is a small issue though - when developing a chess engine up to a certain level, testing starts to completely 
overshadow development. Every feature, every evaluation, every change needs lots of testing against KPIs, test suites, 
itself and real other engines. These tests are extremely time-consuming if they are to be meaningful. It needs thousands
 of games to reliably prove that a feature is an improvement, and these games need time. Too short thinking times make 
 the tests unreliable, especially for features which have increased effect in deeper searches. So several computers 
 doing tests at the same time is normal but also rather stressful :)

I will go on improving FrankyGo, and I'm happy to receive any feedback on my Go code but also of course on improving 
the engine itself.

A word on v1.0: I have implemented most major and common chess engine features, and the search itself is ok for now 
but can of course always be improved. There are many great ideas out there to make it even more effective. Evaluation 
in v1.0 is very basic. Only material and positional differences are counted. There are already some other evaluations 
implemented but deactivated as they need a lot of testing and tuning.

## Features (v1.0.0)
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
    * ~53 Million moves per sec on my i7-7700
* Perft tests
    * ~4.8+Million nps 
* Search as PVS AlphaBeta search
    * all UCI search modes 
    * Iterative Deepening
    * Pondering
    * Quiescence
    * SEE for deciding on good qsearch moves
    * Internal Iterative Deepening
    * Killer Moves
    * History Pruning (History Count, Counter Move) 
    * Mate Distance Pruning
    * Null Move Pruning
    * Late Move Reductions (needs further tuning)
    * Late Move Pruning (needs further tuning)
    * Check extension (needs further tuning/testing)
    * Mate threat extension (deactivated - needs further tuning/testing - currently weakens play)
    * Reverse Futility Pruning (Static Null Move)
    * Futility Pruning (also in quiescence search)
    * Razoring (Stockfish like)
    * History Pruning (History Count, Counter Moves)
    * ~2.4 nps on average on i7-7700 and 256 MB hash
    * TODO: Aspiration Windows, MTDf, Multi-cut Pruning,  
* Transposition Table
* Opening Book from PGN, SAN and Simple move list files and persistent cache
* Evaluation
    * Very simple yet: Material, positional piece values
    * Implemented but not yet tested/activated: Attacks, Mobility, Piece specific evaluations
    * TODO: pawn evaluations and pawn hash table  
* Tests:   
    * Test framework to run EPD test suites
    * Test to determine search tree size for features
    * Test tool to run multiple test suites to test features
* General: 
    * logging, command lines options, configuration file, search and eval parameters in config files

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

In addition to this all the tools included and standardized in Go are extremely valuable. I had so many headaches in 
the past to find good free tools for Java and especially C++ for profiling, assembly code review, etc. comfortable on 
all platforms I'm using (Windows, Mac, Linux). 
Just look at this output from the Go pprof profiler which comes for free with Go:

![pprof](https://github.com/frankkopp/FrankyGo/raw/master/docs/pprof_graph.png)

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

To configure the engine (log level, search and evaluation features and parameters, etc.) a config file can be used. 
Default place to look for it is ./config.toml. Use command line option "-config" to change the location.

Also helpful:
* logs files are stored in folder "../logs" or folder defined by cmd line option "-logpath"
* opening books are searched in folder "../assets/book" or folder defined in cmd line option "-bookpath"

Command line options:

```
 Usage of D:\_DEV\go\src\github.com\frankkopp\FrankyGo\bin\FrankyGo.exe:
   -bookFormat string
         format of opening book
         (Simple|San|Pgn)
   -bookfile string
         opening book file
         provide path if file is not in same directory as executable
         Please also provide bookFormat otherwise this will be ignored
   -bookpath string
         path to opening book files (default "../assets/books")
   -config string
         path to configuration settings file (default "./config.toml")
   -fen string
         fen for perft and nps test (default "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
   -loglvl string
         standard log level
         (critical|error|warning|notice|info|debug) (default "info")
   -logpath string
         path where to write log files to (default "../logs")
   -nps int
         starts nodes per second test for given amount of seconds
   -perft int
         starts perft with the given depth
   -searchloglvl string
         search log level
         (critical|error|warning|notice|info|debug)
   -testdepth int
         search depth limit for each test position
   -testsuite string
         path to file containing EPD tests or folder containing EPD files
   -testtime int
         search time for each test position in milliseconds (default 2000)
   -version
         prints version and exits
```

## Roadmap
### vx.x (planned)
- TODO:
    - Aspiration
    - MTDf
    - Better Evaluation and testing
    - Pawn Structure Cache
    - MultiCut Pruning
    - Other Prunings
    - Continuously Performance/Profiling/Testing
    - Tuning and Testing of all search features and parameters

## Versions
### v1.0.1 (done)
- FIX:
    - MoveGen hasLegalMove did not check pawn doubles but necessary when pawn double blocks attacker i an otherwise
    - nearly mate position.
     
### v1.0.0 (done)
- DONE:
    - change default behavior for log files so running the executable without logs folder works smoothly
    - make it runnable without config file / config file optional
    - Razor (Stockfish)
    - QSearch Futility Pruning
    - Additional Feature Test tool
    - History Heuristics (History Counter, Counter Moves)

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

### v0.5 (done)
- DONE
    - Use TT
    - SearchTreeSize
    - Quiescence search
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

## Credits
- https://www.chessprogramming.org
- Stockfish / Tord Romstad, Marco Costalba, Joona Kiiski and Gary Linscott
- Crafty / Robert Hyatt
- TalkChess.org
- Flux, Pulse / Phokham Nonava
- Beowulf / Colin Frayn
- Mediocre / Jonatan Pettersson
- Bruce Moreland & GERBIL (http://www.brucemo.com/compchess/programming/index.htm)
- CPW-Engine / Pawel Koziol, Edmund Moshammer
- DarkThought / Ernst A. Heinz
- Arena Chess GUI / Martin Blume (Germany) et al
- and many more
