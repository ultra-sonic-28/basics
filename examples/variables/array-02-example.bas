10 REM Set / get values to / from arrays
15 SIZE = 10
20 DIM A%(SIZE)
30 FOR I=1 TO SIZE:A%(I) = I*I:NEXT I
40 FOR I=1 TO 10:PRINT A%(I):NEXT I
