5 REM **** Ce programme affiche la table de 4 ****
10 PRINT "TABLE DE 4 :"
20 FOR I=1 TO 10
25 GOSUB 100
30 PRINT I, V
40 NEXT I
50 END
100 V = I * 4 : RETURN