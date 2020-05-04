:: UCI engines
::============

:: Run to get sts rating using --getrating option
:: when --getrating is used, num threads will be set to 1, and movetime will also be
:: set by the tool depending on the speed of your machine. On my machine without other loads
:: the movetime per pos used by the tool is 200ms. The tool will run a short benchmark
:: to measure your machine speed to get the movetime, before starting the test suite.

STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "FrankyGo.exe" --proto uci -h 256 --getrating



:: Normal run to get score, movetime is in millisec
::STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "FrankyGo.exe" -t 1 -h 256 --movetime 2000





:: WINBOARD engines
:: ================

:: For winboard engines, sts rating is not applicable at the moment, but engine can be tested
:: and still get the score percentage


:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "LambChop_1099.exe -hash 128" --proto wb --mps 300 --tc 1 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "Gerbil_02_x64_ja.exe" --proto wb --mps 40 --tc 0:8

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "Averno081.exe" --proto wb --mps 300 --tc 1 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "Myrddin_0.87-64.exe" --proto wb --mps 40 --tc 0:8 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "Natwarlal_v0.14.exe" --proto wb --mps 300 --tc 1 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "scorpio-276-64-ja.exe" --proto wb --mps 40 --tc 0:8


:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\Knightx1.92\Knightx192.exe" --proto wb --mps 40 --tc 0:16 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\RomiChess_P3L\RomiChessp3L64.exe" --proto wb --mps 40 --tc 0:8 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\Thinker54D\ThinkerInert 5.4D x64 SP.exe" --proto wb --mps 40 --tc 0:8 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\Satana.2.0.7\Satana.2.0.7.w64bit.exe" --proto wb --mps 40 --tc 0:8 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\shallow-rev688-win-ja\shallow-rev688-64-ja.exe" --proto wb --mps 40 --tc 0:8 --log

:: STS_Rating_v3 -f "STS1-STS15_LAN.epd" -e "C:\Chess\engines\nobook\Nemeton\Nemeton\Nemeton_1.exe" --proto wb --mps 40 --tc 0:8 --log





