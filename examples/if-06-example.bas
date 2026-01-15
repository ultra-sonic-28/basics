5 HOME
8 PRINT "Let's count..."
10 LET count = 0
20 PRINT "Count: ", count
30 IF count < 10 THEN GOSUB 40 : GOTO 20 ELSE PRINT "And finally..." : GOSUB 60 : GOTO 70
40 LET count = count + 1
50 RETURN
60 PRINT "All done!" : RETURN
70 END