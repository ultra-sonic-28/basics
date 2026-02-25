10 REM GR mode test - Dégradé de couleurs
20 GR
30 FOR Y = 0 TO 39
40 C = INT(Y/2) : IF C>15 THEN C=15
50 COLOR = C
60 FOR X = 0 TO 39
70 PLOT X,Y
80 NEXT X
90 NEXT Y