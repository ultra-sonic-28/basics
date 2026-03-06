10 REM Hello World Sine Wave Demo
20 PR#3
30 x = 0: count = 0
40 sp = INT(30 + (30* SIN(x* 3.14159 / 180 ) ))
50 PRINT SPC(sp);"Hello World"
60 x = x + 15 : count = count + 1
70 WAIT 100
80 IF count < 100 THEN GOTO 40
90 END
