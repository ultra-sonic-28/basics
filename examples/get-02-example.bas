1 REM ***** Exemple de sous-routine *****
10 HOME:PRINT "Ce programme est une demo de l'utilisation de sous-routine"
20 GOSUB 100
30 PRINT "La sous-routine peut être appelée autant de fois que nécessaire."
40 GOSUB 100
50 PRINT "Et elle évite de la redondance de code."
60 GOSUB 100
70 PRINT "Elle est très utilisée en BASIC."
80 GOSUB 100
85 PRINT "Elle permet aussi de montrer l'utilisation de l'instruction GET"
87 GOSUB 100
89 PRINT "Merci !!"
90 END
99 REM ***** Sous-routine *****
100 PRINT "Pressez une touche pour continuer. ";
110 GET X$
120 HOME
130 RETURN