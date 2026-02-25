10 REM GR mode test - Checkerboard
20 GR
30 FOR Y = 0 TO 39
40 FOR X = 0 TO 39
50 IF (X+Y)/2 = INT((X+Y)/2) THEN C = 0 : GOTO 60
55 C = 15
60 COLOR = C
70 PLOT X,Y
80 NEXT X
90 NEXT Y