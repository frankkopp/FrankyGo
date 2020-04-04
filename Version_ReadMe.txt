v 0.x (planned)
================================================================================
- TODO:
    - MultiCut Pruning
       - https://hci.iwr.uni-heidelberg.de/system/files/private/downloads/1935772097/report_qingyang-cao_enhanced-forward-pruning.pdf
    - IID??
    - Aspiration

v 0.7 (in progress)
================================================================================
- TODO
    - Other Prunings
    - Better Evaluation
    - Pawn Structure Cache
    - Use TestSuites, TreeSize and Arena to test features
    - Performance/Profiling/Testing

- DONE
    - Null Move
    - Remove MPP if not worth

v 0.6 (done)
================================================================================
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

- Measurements
    Test	Standpat	Nodes	4.837.838			        Nps	3.507.125   		        Time	1.450
    Test	TT	        Nodes	1.763.431	63,55%	        Nps	3.137.947	10,53%	    	Time	544	62,48%
    Test	QSTT	    Nodes	1.022.318	42,03%	        Nps	2.748.463	12,41%		    Time	379	30,33%
    Test	MDP/MPP	    Nodes	1.010.160	1,19%	        Nps	2.769.100	-0,75%		    Time	367	3,17%
    Test	PVS     	Nodes	893.718	    11,53%	        Nps	2.777.407	-0,30%		    Time	322	12,26%
    Test	Killer	    Nodes	857.379	    4,07%   82,28%	Nps	2.802.207	-0,89%	20,10%	Time	309	4,04%	78,69%

v 0.5 (done)
================================================================================
-DONE
    - Use TT
    - SearchTreeSize
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
