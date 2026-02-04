10 REM Set / get values to / from arrays
15 SIZE = 10
20 DIM A(SIZE)
30 FOR I=1 TO SIZE:A(I) = I/SIZE:NEXT I
40 FOR I=10 TO 1 STEP -1:PRINT A(I):NEXT I