10 REM A Simple Body Mass Index Calculator.
80 PRINT "***************************"
90 PRINT "*                         *"
100 PRINT "*  Simple BMI Calculator  *"
110 PRINT "*                         *"
120 PRINT "***************************"
130 PRINT ""
140 PRINT ""
150 INPUT "Input your height (inches): "; height
160 INPUT "Input your weight (lbs): ";weight
170 bmiCalc = (weight/(height*height))*703
180 PRINT ""
190 PRINT "Your BMI is ";bmiCalc
200 PRINT:PRINT "Thanks for using BMI Calculator"
210 END