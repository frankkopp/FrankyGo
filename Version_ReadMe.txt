v 0.x (planned)
================================================================================
- TODO:
    - SEE
    - Pawn Structure Cache
    - Better Evaluation
    - FP
    - Ext
    - Other Prunings
    - MultiCut Pruning
       - https://hci.iwr.uni-heidelberg.de/system/files/private/downloads/1935772097/report_qingyang-cao_enhanced-forward-pruning.pdf
    - Aspiration
    - NullMove Threat Detection

v 0.8 (planned)
================================================================================
- TODO

v 0.7 (in progress)
================================================================================
- TODO
    - Performance/Profiling/Testing
    - LMR and LMP weaken engine - more testing needed

- DONE
    - Magic Bitboards
    - Use TestSuites, TreeSize and Arena to test features
    - IID
    - Null Move
    - Remove MPP if not worth

Measurements:
    LMR
    -----------------00 FrankyGo Config 1-----------------
    00 FrankyGo Config 1 - 00 FrankyGo Config 2 : 54,0/100 35-27-38 (=01101==11=0111=1====10=1=10001=011=1==0====110=0=01011===110110==1=010==00===10011=01=0=0==101010=1)  54%   +28
    -----------------00 FrankyGo Config 2-----------------
    00 FrankyGo Config 2 - 00 FrankyGo Config 1 : 46,0/100 27-35-38 (=10010==00=1000=0====01=0=01110=100=0==1====001=1=10100===001001==0=101==11===01100=10=1=1==010101=0)  46%   -28

    With Null Move
    Without IID
    -----------------00 FrankyGo-----------------
    00 FrankyGo - 00 FrankyGo -NMP : 71,0/100 59-17-24  71%  +156
    00 FrankyGo - 20 Franky-1.0    : 37,0/100 31-57-12  37%   -92
    00 FrankyGo - Stockfish Weak   : 64,5/100 61-32-7   65%  +108

    With Null Move and IID
    -----------------00 FrankyGo-----------------
    00 FrankyGo - 00 FrankyGo -IID : 51,0/100 27-25-48 51%    +7
    00 FrankyGo - 20 Franky-1.0    : 37,5/100 32-57-11 38%   -85
    00 FrankyGo - Stockfish Weak   : 66,5/100 61-28-11 67%  +123

    Config 1 uses Lmr, Config 2 without
    -----------------00 FrankyGo Config 1-----------------
    00 FrankyGo Config 1 - 00 FrankyGo Config 2 : 50,5/100 32-31-37 51%    +7
    00 FrankyGo Config 1 - 20 Franky-1.0        : 48,0/100 40-44-16 48%   -14
    00 FrankyGo Config 1 - Stockfish Weak       : 51,5/100 46-43-11 52%   +14


    Number of feature tests: 4
    Number of fens         : 30
    Total tests            : 120
    Depth                  : 8
    Test: Killer        Nodes: 336.173.679     Nps: 2.488.845    Time: 139.059    Depth:   9/22  Special: 0
    Test: NMP           Nodes: 340.970.294     Nps: 2.545.551    Time: 141.480    Depth:   9/22  Special: 0
    Test: TTMove        Nodes: 176.388.883     Nps: 2.532.209    Time: 77.622     Depth:   9/22  Special: 0
    Test: IID           Nodes: 174.964.952     Nps: 2.512.220    Time: 77.746     Depth:   9/22  Special: 912

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
