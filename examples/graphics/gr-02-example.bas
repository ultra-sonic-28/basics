10 REM GR mode test - Rod's Color Pattern
20  GR
30  FOR W = 3 TO 50
40  FOR I = 1 TO 19
50  FOR J = 0 TO 19
60  K = I + J
80  COLOR= J * 3 / (I + 3) + I * W / 12
90  PLOT I,K: PLOT K,I: PLOT 40 - I,40 - K: PLOT 40 - K,40 - I
100  PLOT K,40 - I: PLOT 40 - I,K: PLOT I,40 - K: PLOT 40 - K,I
110  NEXT J : NEXT I : NEXT W : GOTO 30
120  TEXT : HOME : END