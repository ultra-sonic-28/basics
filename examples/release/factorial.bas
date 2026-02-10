10 REM Calculate factorial of 10
20 N=10
30 GOSUB 100
40 PRINT "Factorial of ";N;" is ";F
50 END
100 REM Factorial calc using simple loop
110 F = 1
120 FOR I=1 TO N
130   F = F*I
140 NEXT I
150 RETURN