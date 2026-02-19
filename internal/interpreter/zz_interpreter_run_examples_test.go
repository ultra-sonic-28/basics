package interpreter

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"basics/internal/common"
	"basics/internal/constants"
	"basics/internal/input"
	"basics/internal/lexer"
	"basics/internal/machines"
	"basics/internal/parser"
	"basics/testutils"
)

func TestExamplesExecution(t *testing.T) {
	tests := []struct {
		name     string
		file     string
		errors   int
		expected string
	}{
		{
			name:   "Abs-01",
			file:   "maths/abs-01-example.bas",
			errors: 0,
			expected: `1.75
1.75
2.8746841
2.8746841
10.751318534000001
10.751318534000001
5
5
14.3734205
14.3734205
`,
		},
		{
			name:   "Abs-02",
			file:   "maths/abs-02-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 2 (ABS)
`,
		},
		{
			name:   "Abs-03",
			file:   "maths/abs-03-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 3 (ABS)
`,
		},
		{
			name:   "Array-01",
			file:   "variables/array-01-example.bas",
			errors: 0,
			expected: `1
0.9
0.8
0.7
0.6
0.5
0.4
0.3
0.2
0.1
`,
		},
		{
			name:   "Array-02",
			file:   "variables/array-02-example.bas",
			errors: 0,
			expected: `1
4
9
16
25
36
49
64
81
100
`,
		},
		{
			name:   "Array-03",
			file:   "variables/array-03-example.bas",
			errors: 0,
			expected: `A
AA
AAA
AAAA
AAAAA
AAAAAA
AAAAAAA
AAAAAAAA
AAAAAAAAA
AAAAAAAAAA
`,
		},
		{
			name:   "Array-04",
			file:   "variables/array-04-example.bas",
			errors: 0,
			expected: `111           55555
`,
		},
		{
			name:   "Array-05",
			file:   "variables/array-05-example.bas",
			errors: 0,
			expected: `999
`,
		},
		{
			name:   "CubeExample",
			file:   "programs/maths/cube-example.bas",
			errors: 0,
			expected: `Affichage des cubes de 1 a 10
1 ^ 3 = 1
2 ^ 3 = 8
3 ^ 3 = 27
4 ^ 3 = 64
5 ^ 3 = 125
6 ^ 3 = 216
7 ^ 3 = 343
8 ^ 3 = 512
9 ^ 3 = 729
10 ^ 3 = 1000
`,
		},
		{
			name:   "BooleanOperator-01",
			file:   "operators/boolean-operator-01-example.bas",
			errors: 0,
			expected: `0
1
0
1
1
0
1
0
0
1
1
0
1
0
`,
		},
		{
			name:   "BooleanOperator-02",
			file:   "operators/boolean-operator-02-example.bas",
			errors: 0,
			expected: `0
0
1
1
0
1
0
0
1
1
1
1
0
0
1
1
0
0
`,
		},
		{
			name:   "Clear-01",
			file:   "variables/clear-01-example.bas",
			errors: 0,
			expected: `A=1.5
A%=1
A$=A String
A=0
A%=0
A$=
`,
		},
		{
			name:   "Clear-02",
			file:   "variables/clear-02-example.bas",
			errors: 0,
			expected: `1
0.9
0.8
0.7
0.6
0.5
0.4
0.3
0.2
0.1
0
0
0
0
0
0
0
0
0
0
`,
		},
		{
			name:   "Clear-03",
			file:   "variables/clear-03-example.bas",
			errors: 0,
			expected: `1
4
9
16
25
36
49
64
81
100
0
0
0
0
0
0
0
0
0
0
`,
		},
		{
			name:   "Clear-04",
			file:   "variables/clear-04-example.bas",
			errors: 0,
			expected: `A$(1)=A
A$(2)=AA
A$(3)=AAA
A$(4)=AAAA
A$(5)=AAAAA
A$(6)=AAAAAA
A$(7)=AAAAAAA
A$(8)=AAAAAAAA
A$(9)=AAAAAAAAA
A$(10)=AAAAAAAAAA
A$(10)=
A$(9)=
A$(8)=
A$(7)=
A$(6)=
A$(5)=
A$(4)=
A$(3)=
A$(2)=
A$(1)=
`,
		},
		{
			name:   "End-01",
			file:   "flow_control/end-01-example.bas",
			errors: 0,
			expected: `Hello
`,
		},
		{
			name:   "Factorial",
			file:   "programs/maths/factorial.bas",
			errors: 0,
			expected: `6! = 720
`,
		},
		{
			name:   "Fibonacci",
			file:   "programs/maths/fibonacci.bas",
			errors: 0,
			expected: `
Here is(are) your 20 Fibonacci number(s):
0
1
1
2
3
5
8
13
21
34
55
89
144
233
377
610
987
1597
2584
4181
All done!
`,
		},
		{
			name:   "Flash-01",
			file:   "display/flash-01-example.bas",
			errors: 0,
			expected: `HELLO
WORLD
`,
		},
		{
			name:   "Flash-02",
			file:   "display/flash-02-example.bas",
			errors: 0,
			expected: `HELLO WORLD!
`,
		},
		{
			name:   "Flash-03",
			file:   "display/flash-03-example.bas",
			errors: 0,
			expected: `HELLO WORLD!
This is the End..
`,
		},
		{
			name:   "For-01",
			file:   "flow_control/for-01-example.bas",
			errors: 0,
			expected: `0
1
2
3
4
5
6
7
8
9
10
`,
		},
		{
			name:   "For-02",
			file:   "flow_control/for-02-example.bas",
			errors: 0,
			expected: `0
1
2
3
4
5
6
7
8
9
10
`,
		},
		{
			name:   "For-03",
			file:   "flow_control/for-03-example.bas",
			errors: 0,
			expected: `0
2
4
6
8
10
`,
		},
		{
			name:   "For-04",
			file:   "flow_control/for-04-example.bas",
			errors: 0,
			expected: `0
2.5
5
7.5
10
`,
		},
		{
			name:   "For-05",
			file:   "flow_control/for-05-example.bas",
			errors: 0,
			expected: `10
8
6
4
2
0
`,
		},
		{
			name:   "For-06",
			file:   "flow_control/for-06-example.bas",
			errors: 0,
			expected: `A=0, B=0, A*B=0
A=0, B=2, A*B=0
A=0, B=4, A*B=0
A=0, B=6, A*B=0
A=0, B=8, A*B=0
A=0, B=10, A*B=0
A=2, B=0, A*B=0
A=2, B=2, A*B=4
A=2, B=4, A*B=8
A=2, B=6, A*B=12
A=2, B=8, A*B=16
A=2, B=10, A*B=20
A=4, B=0, A*B=0
A=4, B=2, A*B=8
A=4, B=4, A*B=16
A=4, B=6, A*B=24
A=4, B=8, A*B=32
A=4, B=10, A*B=40
A=6, B=0, A*B=0
A=6, B=2, A*B=12
A=6, B=4, A*B=24
A=6, B=6, A*B=36
A=6, B=8, A*B=48
A=6, B=10, A*B=60
A=8, B=0, A*B=0
A=8, B=2, A*B=16
A=8, B=4, A*B=32
A=8, B=6, A*B=48
A=8, B=8, A*B=64
A=8, B=10, A*B=80
A=10, B=0, A*B=0
A=10, B=2, A*B=20
A=10, B=4, A*B=40
A=10, B=6, A*B=60
A=10, B=8, A*B=80
A=10, B=10, A*B=100
`,
		},
		/* {
			name:     "For-07",
			file:     "for-07-example.bas",
			errors:   1,
			expected: ``,
		}, */
		{
			name:   "For-08",
			file:   "flow_control/for-08-example.bas",
			errors: 0,
			expected: `⚠️ STEP CANNOT BE ZERO IN 10 ()
`,
		},
		{
			name:   "For-09",
			file:   "flow_control/for-09-example.bas",
			errors: 0,
			expected: `1
3
5
7
9
`,
		},
		{
			name:   "Gosub-01",
			file:   "flow_control/gosub-01-example.bas",
			errors: 0,
			expected: `Hello
World
!!!
`,
		},
		{
			name:   "Gosub-02",
			file:   "flow_control/gosub-02-example.bas",
			errors: 0,
			expected: `Hello
World
!!!
`,
		},
		{
			name:   "Gosub-03",
			file:   "flow_control/gosub-03-example.bas",
			errors: 0,
			expected: `TABLE DE 4 :
1             4
2             8
3             12
4             16
5             20
6             24
7             28
8             32
9             36
10            40
`,
		},
		{
			name:   "Gosub-04",
			file:   "flow_control/gosub-04-example.bas",
			errors: 0,
			expected: `Hello
World
!!!
`,
		},
		{
			name:   "Goto-01",
			file:   "flow_control/goto-01-example.bas",
			errors: 0,
			expected: `First line
Second line
Third line
Last line
`,
		},
		{
			name:   "HelloWorld-01",
			file:   "others/hello-world-01-example.bas",
			errors: 0,
			expected: `Hello World
`,
		},
		{
			name:     "HelloWorld-02",
			file:     "others/hello-world-02-example.bas",
			errors:   2,
			expected: ``,
		},
		{
			name:   "Home-01",
			file:   "display/home-01-example.bas",
			errors: 0,
			expected: `HELLO
`,
		},
		{
			name:   "Home-02",
			file:   "display/home-02-example.bas",
			errors: 0,
			expected: `A
B
`,
		},
		{
			name:   "HtabVtab-01",
			file:   "tabs/htab-vtab-01-example.bas",
			errors: 0,
			expected: `1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
15
12
9
6
3
0
1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
`,
		},
		{
			name:   "HtabVtab-02",
			file:   "tabs/htab-vtab-02-example.bas",
			errors: 0,
			expected: `1
2
3
4
5
6
7
8
9
10
11
12
13
14
15
`,
		},
		{
			name:   "If-01",
			file:   "flow_control/if-01-example.bas",
			errors: 0,
			expected: `Count:        0
Count:        1
Count:        2
Count:        3
Count:        4
Count:        5
Count:        6
Count:        7
Count:        8
Count:        9
All done!
`,
		},
		{
			name:   "If-02",
			file:   "flow_control/if-02-example.bas",
			errors: 0,
			expected: `Count:        0
Count:        1
Count:        2
Count:        3
Count:        4
Count:        5
Count:        6
Count:        7
Count:        8
Count:        9
All done!
`,
		},
		{
			name:   "If-03",
			file:   "flow_control/if-03-example.bas",
			errors: 0,
			expected: `Count:        0
Go to line 20
Count:        1
Go to line 20
Count:        2
Go to line 20
Count:        3
Go to line 20
Count:        4
Go to line 20
Count:        5
Go to line 20
Count:        6
Go to line 20
Count:        7
Go to line 20
Count:        8
Go to line 20
Count:        9
Go to line 60
All done!
`,
		},
		{
			name:   "If-04",
			file:   "flow_control/if-04-example.bas",
			errors: 0,
			expected: `Count:        0
Count:        1
Count:        2
Count:        3
Count:        4
Count:        5
Count:        6
Count:        7
Count:        8
Count:        9
All done!
`,
		},
		{
			name:   "If-05",
			file:   "flow_control/if-05-example.bas",
			errors: 0,
			expected: `Let's count...
Count:        0
Count:        1
Count:        2
Count:        3
Count:        4
Count:        5
Count:        6
Count:        7
Count:        8
Count:        9
Count:        10
All done!
`,
		},
		{
			name:   "If-06",
			file:   "flow_control/if-06-example.bas",
			errors: 0,
			expected: `Let's count...
Count:        0
Count:        1
Count:        2
Count:        3
Count:        4
Count:        5
Count:        6
Count:        7
Count:        8
Count:        9
Count:        10
And finally...
All done!
`,
		},
		{
			name:   "Int-01",
			file:   "maths/int-01-example.bas",
			errors: 0,
			expected: `1
1
-2
-2
2
10
5
14
`,
		},
		{
			name:   "Int-02",
			file:   "maths/int-02-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 2 (INT)
`,
		},
		{
			name:   "Int-03",
			file:   "maths/int-03-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 3 (INT)
`,
		},
		{
			name:   "Int-04",
			file:   "maths/int-04-example.bas",
			errors: 0,
			expected: `-21
`,
		},
		{
			name:   "Len-01",
			file:   "strings/len-01-example.bas",
			errors: 0,
			expected: `9
`,
		},
		{
			name:   "Len-02",
			file:   "strings/len-02-example.bas",
			errors: 0,
			expected: `9
`,
		},
		{
			name:   "Len-03",
			file:   "strings/len-03-example.bas",
			errors: 0,
			expected: `9
`,
		},
		{
			name:   "Len-04",
			file:   "strings/len-04-example.bas",
			errors: 0,
			expected: `13
`,
		},
		{
			name:   "Len-05",
			file:   "strings/len-05-example.bas",
			errors: 0,
			expected: `13
`,
		},
		{
			name:   "Len-06",
			file:   "strings/len-06-example.bas",
			errors: 0,
			expected: `0
`,
		},
		{
			name:   "Len-07",
			file:   "strings/len-07-example.bas",
			errors: 0,
			expected: `260
`,
		},
		{
			name:   "Len-08",
			file:   "strings/len-08-example.bas",
			errors: 0,
			expected: `550
`,
		},
		{
			name:     "LinenumWithoutInst-01",
			file:     "others/linenum-without-inst-01-example.bas",
			errors:   0,
			expected: ``,
		},
		{
			name:     "LinenumWithoutInst-02",
			file:     "others/linenum-without-inst-02-example.bas",
			errors:   0,
			expected: ``,
		},
		{
			name:     "LinenumWithoutInst-03",
			file:     "others/linenum-without-inst-03-example.bas",
			errors:   0,
			expected: ``,
		},
		{
			name:   "Left-01",
			file:   "strings/left-01-example.bas",
			errors: 0,
			expected: `APPLE
`,
		},
		{
			name:   "Left-02",
			file:   "strings/left-02-example.bas",
			errors: 0,
			expected: `APPLE
`,
		},
		{
			name:   "Left-03",
			file:   "strings/left-03-example.bas",
			errors: 0,
			expected: `APPLESOFT
`,
		},
		{
			name:   "Left-04",
			file:   "strings/left-04-example.bas",
			errors: 0,
			expected: `APPLE
`,
		},
		{
			name:   "Left-05",
			file:   "strings/left-05-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED NUMBER IN 2 (LEFT$)
`,
		},
		{
			name:   "Left-06",
			file:   "strings/left-06-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED STRING IN 3 (LEFT$)
`,
		},
		{
			name:   "Left-07",
			file:   "strings/left-07-example.bas",
			errors: 1,
			expected: `
`,
		},
		{
			name:   "Left-08",
			file:   "strings/left-08-example.bas",
			errors: 0,
			expected: `⚠️ ILLEGAL QUANTITY ERROR IN 2 (LEFT$)
`,
		},
		{
			name:   "Mid-01",
			file:   "strings/mid-01-example.bas",
			errors: 0,
			expected: `ESOFT
`,
		},
		{
			name:   "Mid-02",
			file:   "strings/mid-02-example.bas",
			errors: 0,
			expected: `ESOFT
`,
		},
		{
			name:   "Mid-03",
			file:   "strings/mid-03-example.bas",
			errors: 0,
			expected: `
`,
		},
		{
			name:   "Mid-04",
			file:   "strings/mid-04-example.bas",
			errors: 0,
			expected: `ESOFT
`,
		},
		{
			name:   "Mid-05",
			file:   "strings/mid-05-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED NUMBER IN 2 (MID$)
`,
		},
		{
			name:   "Mid-06",
			file:   "strings/mid-06-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED STRING IN 3 (MID$)
`,
		},
		{
			name:   "Mid-07",
			file:   "strings/mid-07-example.bas",
			errors: 1,
			expected: `
`,
		},
		{
			name:   "Mid-08",
			file:   "strings/mid-08-example.bas",
			errors: 0,
			expected: `⚠️ ILLEGAL QUANTITY ERROR IN 2 (MID$)
`,
		},
		{
			name:   "Mid-09",
			file:   "strings/mid-09-example.bas",
			errors: 0,
			expected: `LE
`,
		},
		{
			name:   "Mid-10",
			file:   "strings/mid-10-example.bas",
			errors: 0,
			expected: `O
`,
		},
		{
			name:   "Mid-11",
			file:   "strings/mid-11-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED NUMBER IN 2 (MID$)
`,
		},
		{
			name:   "Mid-12",
			file:   "strings/mid-12-example.bas",
			errors: 1,
			expected: `
`,
		},
		{
			name:   "Mid-13",
			file:   "strings/mid-13-example.bas",
			errors: 0,
			expected: `⚠️ ILLEGAL QUANTITY ERROR IN 2 (MID$)
`,
		},
		{
			name:   "MultipleOf4",
			file:   "programs/maths/multpile-of-4-example.bas",
			errors: 0,
			expected: `TABLE DE 4 :
1             4
2             8
3             12
4             16
5             20
6             24
7             28
8             32
9             36
10            40
`,
		},
		{
			name:   "NormalInverse-01",
			file:   "display/normal-inverse-01-example.bas",
			errors: 0,
			expected: `HELLO
WORLD
`,
		},
		{
			name:   "NormalInverse-02",
			file:   "display/normal-inverse-02-example.bas",
			errors: 0,
			expected: `HELLO WORLD!
`,
		},
		{
			name:   "NormalInverse-03",
			file:   "display/normal-inverse-03-example.bas",
			errors: 0,
			expected: `HELLO WORLD!
`,
		},
		{
			name:   "Primes-01",
			file:   "programs/maths/primes-01-example.bas",
			errors: 0,
			expected: `NOMBRES PREMIERS JUSQU'A 50
3
5
7
11
13
17
19
23
29
31
37
41
43
47
All done!
`,
		},
		{
			name:   "Primes-02",
			file:   "programs/maths/primes-02-example.bas",
			errors: 0,
			expected: `NOMBRES PREMIERS JUSQU'A 50
3
5
7
11
13
17
19
23
29
31
37
41
43
47
`,
		},
		{
			name:   "Print-01",
			file:   "display/print-01-example.bas",
			errors: 0,
			expected: `A=7, A+1=8
`,
		},
		{
			name:   "Print-02",
			file:   "display/print-02-example.bas",
			errors: 0,
			expected: `7             7             8
`,
		},
		{
			name:   "Print-03",
			file:   "display/print-03-example.bas",
			errors: 0,
			expected: `⚠️ DIVISION BY ZERO IN 1 (/)
`,
		},
		{
			name:   "Print-04",
			file:   "display/print-04-example.bas",
			errors: 0,
			expected: `7
`,
		},
		{
			name:   "Print-05",
			file:   "display/print-05-example.bas",
			errors: 0,
			expected: `7
`,
		},
		{
			name:   "Print-06",
			file:   "display/print-06-example.bas",
			errors: 2,
			expected: `⚠️ UNDEFINED VARIABLE A IN 3 ()
`,
		},
		{
			name:   "Print-07",
			file:   "display/print-07-example.bas",
			errors: 2,
			expected: `⚠️ UNDEFINED VARIABLE A IN 3 ()
`,
		},
		{
			name:   "Print-08",
			file:   "display/print-08-example.bas",
			errors: 1,
			expected: `Hello World
`,
		},
		{
			name:     "Print-09",
			file:     "display/print-09-example.bas",
			errors:   0,
			expected: `0 1 2 3 4 5 6 7 8 9 10 `,
		},
		{
			name:     "Print-10",
			file:     "display/print-10-example.bas",
			errors:   0,
			expected: `2`,
		},
		{
			name:   "Print-11",
			file:   "display/print-11-example.bas",
			errors: 0,
			expected: `
`,
		},
		{
			name:   "Print-12",
			file:   "display/print-12-example.bas",
			errors: 0,
			expected: `Line 1

Line 3
`,
		},
		{
			name:   "Right-01",
			file:   "strings/right-01-example.bas",
			errors: 0,
			expected: `SOFT
`,
		},
		{
			name:   "Right-02",
			file:   "strings/right-02-example.bas",
			errors: 0,
			expected: `ESOFT
`,
		},
		{
			name:   "Right-03",
			file:   "strings/right-03-example.bas",
			errors: 0,
			expected: `APPLESOFT
`,
		},
		{
			name:   "Right-04",
			file:   "strings/right-04-example.bas",
			errors: 0,
			expected: `ESOFT
`,
		},
		{
			name:   "Right-05",
			file:   "strings/right-05-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED NUMBER IN 2 (RIGHT$)
`,
		},
		{
			name:   "Right-06",
			file:   "strings/right-06-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED STRING IN 3 (RIGHT$)
`,
		},
		{
			name:   "Right-07",
			file:   "strings/right-07-example.bas",
			errors: 1,
			expected: `
`,
		},
		{
			name:   "Right-08",
			file:   "strings/right-08-example.bas",
			errors: 0,
			expected: `⚠️ ILLEGAL QUANTITY ERROR IN 2 (RIGHT$)
`,
		},
		{
			name:   "Sgn-01",
			file:   "maths/sgn-01-example.bas",
			errors: 0,
			expected: `1
-1
1
1
0
0
1
-1
1
-1
`,
		},
		{
			name:   "Sgn-02",
			file:   "maths/sgn-02-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 2 (SGN)
`,
		},
		{
			name:   "Sgn-03",
			file:   "maths/sgn-03-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 3 (SGN)
`,
		},
		{
			name:   "Sqr-01",
			file:   "maths/sqr-01-example.bas",
			errors: 0,
			expected: `1.3228756555322954
2
3.867815921162743
7
14
`,
		},
		{
			name:   "Sqr-02",
			file:   "maths/sqr-02-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 2 (SQR)
`,
		},
		{
			name:   "Sqr-03",
			file:   "maths/sqr-03-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 3 (SQR)
`,
		},
		{
			name:   "Sqr-04",
			file:   "maths/sqr-04-example.bas",
			errors: 0,
			expected: `⚠️ EXPRESSION VALUE MUST BE POSITIVE OR NULL IN 2 (SQR)
`,
		},
		{
			name:   "Sqr-05",
			file:   "maths/sqr-05-example.bas",
			errors: 0,
			expected: `⚠️ EXPRESSION VALUE MUST BE POSITIVE OR NULL IN 3 (SQR)
`,
		},
		{
			name:   "Sqr-06",
			file:   "maths/sqr-06-example.bas",
			errors: 0,
			expected: `⚠️ EXPRESSION VALUE MUST BE POSITIVE OR NULL IN 3 (SQR)
`,
		},
		{
			name:   "Sqr-07",
			file:   "maths/sqr-07-example.bas",
			errors: 0,
			expected: `⚠️ EXPRESSION VALUE MUST BE POSITIVE OR NULL IN 3 (SQR)
`,
		},
		{
			name:   "Sqr-08",
			file:   "maths/sqr-08-example.bas",
			errors: 0,
			expected: `⚠️ EXPRESSION VALUE MUST BE POSITIVE OR NULL IN 4 (SQR)
`,
		},
		{
			name:   "Square",
			file:   "programs/maths/square-example.bas",
			errors: 0,
			expected: `Affichage des carres de 1 a 10
1 x 1 = 1
2 x 2 = 4
3 x 3 = 9
4 x 4 = 16
5 x 5 = 25
6 x 6 = 36
7 x 7 = 49
8 x 8 = 64
9 x 9 = 81
10 x 10 = 100
`,
		},
		{
			name:   "Tab-01",
			file:   "tabs/tab-01-example.bas",
			errors: 0,
			expected: `     TAB(5)
`,
		},
		{
			name:   "Tab-02",
			file:   "tabs/tab-02-example.bas",
			errors: 0,
			expected: `                         TAB Testing


   TAB(3) mais curseur en 10!!

               TAB(15)


                    TAB(20)

          TAB(10)   TAB(20)   TAB(30)

        TAB(8)  TAB(16)
`,
		},
		{
			name:   "Tab-03",
			file:   "tabs/tab-03-example.bas",
			errors: 0,
			expected: `⚠️ EXPECTED NUMBER IN 2 ()
`,
		},
		{
			name:   "Vars-01",
			file:   "variables/vars-01-example.bas",
			errors: 0,
			expected: `A=1.5
A%=1
A$=A String
`,
		},
		{
			name:   "Vars-02",
			file:   "variables/vars-02-example.bas",
			errors: 0,
			expected: `A=3
A%=2
A$=A String Another one
`,
		},
		{
			name:   "Vars-03",
			file:   "variables/vars-03-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH: STRING EXPECTED IN 20 ()
`,
		},
		{
			name:   "Vars-04",
			file:   "variables/vars-04-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH: STRING EXPECTED IN 20 ()
`,
		},
		{
			name:   "Vars-05",
			file:   "variables/vars-05-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH: INTEGER EXPECTED IN 20 ()
`,
		},
		{
			name:   "Vars-06",
			file:   "variables/vars-06-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH: FLOAT EXPECTED IN 20 ()
`,
		},
		{
			name:   "Vars-07",
			file:   "variables/vars-07-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 4 (*)
`,
		},
		{
			name:   "Vars-08",
			file:   "variables/vars-08-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 4 (/)
`,
		},
		{
			name:   "Vars-09",
			file:   "variables/vars-09-example.bas",
			errors: 0,
			expected: `⚠️ TYPE MISMATCH IN 4 (^)
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// --- Lecture fichier ---
			rootDir, _ := os.Getwd()
			path := filepath.Join(rootDir, "..", "..", "examples", tt.file)
			data, err := os.ReadFile(path)
			testutils.True(t, fmt.Sprintf("file read ok, attempting to read '%s'", path), err == nil)

			source := string(data)

			// --- Lexer ---
			tokens := lexer.Lex(source)
			testutils.True(t, "tokenization ok", len(tokens) > 0)

			// --- Parser ---
			p := parser.New(tokens)
			prog, errs := p.ParseProgram()
			testutils.Equal(t, "no parser errors", len(errs), tt.errors)

			// --- Capture stdout ---
			/* var buf bytes.Buffer
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w */

			// --- Interpreter ---
			rt, err := machines.NewRuntime(constants.BASIC_TTY, false)
			testutils.True(t, "runtime ok", err == nil)

			out := &bytes.Buffer{}
			rt.Input = input.NewTTYInput(os.Stdin, out)
			rt.Video.SetOutput(out)

			interp := New(rt)
			interp.Run(prog)

			// --- Restore stdout ---
			/* _ = w.Close()
			os.Stdout = oldStdout
			_, _ = buf.ReadFrom(r)

			output := buf.String() */
			output := out.String()

			// --- Assertion ---
			testutils.Equal(t, "program output", common.StripANSI(output), tt.expected)
		})
	}
}
