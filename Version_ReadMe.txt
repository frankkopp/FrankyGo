v 0.x (planned)
================================================================================

v 0.6 (planned)
================================================================================
- TODO
    - SearchTreeSize
    - TestSuite Tests
    - Use TT
    - Better Evaluation
    - Pawn Structure Cache

v 0.5 (in progress)
================================================================================
- TODO
    - Testing
    - Known Issues:

-DONE
    - Quiescence search
    - Score as string()
    - Evaluation (simple)

v 0.4 (done)
================================================================================
- DONE
    - Pondering
    - Testing for correct play with Arena against Stockfish
    - log files for standard and search log
    - Implement simple Evaluator
    - Complete simple search
    - Implement simple Search

v 0.3 (done)
================================================================================
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

v 0.2 (done)
================================================================================
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

v 0.1 (done)
================================================================================
- DONE:
    - Perft works (1.353.761 nps / Java 3.5M, C++ 4.5M)
    - MoveGenerator (all required for perft)
    - Position (all required for perft)
    - MoveArray and MoveList - both are for list of moves - MoveArray is faster for sorting
        might be slower when inserting at the front - needs testing
    - Most types (all required for perft)
