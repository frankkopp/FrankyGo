v 0.x (planned)
================================================================================

v 0.4 (planned)
================================================================================
    - Implement simple Search
    - Implement simple Evaluator

v 0.3 (in progress)
================================================================================
- TODO:
    - search handling
    - limits
    - time control

- DONE:
    - starting / stopping


v 0.2 (in progress)
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
