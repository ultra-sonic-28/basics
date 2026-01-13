10 LET count = 0
20 PRINT "Count: ", count
30 LET count = count + 1
40 IF count < 10 THEN PRINT "Go to line 20" : GOTO 20 ELSE PRINT "Go to line 60" : GOTO 60
50 END
60 PRINT "All done!"