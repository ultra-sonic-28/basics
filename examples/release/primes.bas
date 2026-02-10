2 REM Calculate the first 50 prime numbers
5 MAX = 50
10 PRINT "Prime numbers up to "; MAX
20 FOR N = 2 TO MAX
30 P = 1
40 FOR D = 2 TO N/2
50 T = D
60 IF T >= N THEN 90
70 T = T + D
80 GOTO 60
90 IF T = N THEN P = 0
100 NEXT D
110 IF P = 1 THEN PRINT N
120 NEXT N
130 PRINT "All done!"