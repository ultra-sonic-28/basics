# BASICS - BASIC Interpreter for old computers

The project is currently primarily focused on APPLE II computers, with an architecture that allows for expansion to other `retro computers`.

## Project Objectives

* BASIC Interpreter for old computers

## Planned Developments

* Support for retro computers:
    * Commodore 64
    * Amstrad
    * MSX flavors
    * Others...

## Technical Stack

* **Go 1.25+**
* **Ebiten 2.9.7**

## Version control system
* **Jujutsu (jj) + Git**

## Get a binary release
Releases are published on Github. You can get the release version you wanted under [Releases](https://github.com/ultra-sonic-28/rpg-companion/releases).
Be sure to verifiy SHA256 cheksum by running (ex: for v0.1.0.12 release):
- Windows
    ```powershell
    Get-FileHash basics-win-x64-v0.1.0.12.zip -Algorithm SHA256
    ```
- Linux / MacOS
    ```bash
    sha256sum -c basics-win-x64-v0.1.0.12.zip.sha256
    ```
## Development toolchain
Basics use `Mage` as its toolchain for developping purpose.
Dependencies are tracked via Go modules.

### `Mage` installation
```bash
go install github.com/magefile/mage@latest
```

No need to run `mage -init` as `Basics` already have its own `magefile.go`.

For resolving dependencies
```bash
go mod tidy
```

### Install tools:
```bash
mage tools
```

The ```tools``` target will install all necessary tools for the building toolchain.

### Running `Mage`
Running `Mage` without any arguments will display all available targets.
```
mage

Targets:
  build      Configure and build BASICS binary to /bin directory
  clean      Delete all .exe files under basics/ except in folders starting with "."
  release    Build release archive for github
  stats      Running sources analysis and statistics generation
  test       Run unit tests with coverage support
  tools      Installing tools : winres
```

## Supported computers

### APPLE II
#### Real, Integer and String variables
There are three different types of variables used in APPLESOFT BASIC.
* Real
* Integer
* String

The table below summarizes the three types of variables used in APPLESOFT

| Description | Symbol to append to variable name | Example |
|:---:|:---:|:---:|
| String | `$` | `A$`<br/>`ALPHA$` |
| Integer | `%` | `A%`<br/>`COUNT1%` |
| Real | none | `A` |

An integer or string variable must be followed by a `%` or `$` at each use of that variable. For example, `X`, `X%` and `X$` are different variables.

#### Supported instructions set
##### Variables management
* `LET`
    * Assign a value to a variable, creating it if necessary. Optionnal.
    * For arrays, they must be defined using `DIM` before any assignments take place
* `DIM`
    * Define and allocate space for an array of reals, integers or strings. Use:
        * DIM A(10) for array of real numbers
        * DIM A%(10) for array of integer numbers
        * DIM A$(10) for array of strings
    * Arrays could have any number of dimensions, dimensions can be a number or an expression:
        * DIM A$(10)
        * DIM B%(5, 5)
        * DIM C(1, 2, 3, 4, 5)
        * SIZE=10 : DIM A(SIZE)
        * All these are correct

##### Editing and format related
* `REM`
    * This serves to allow text of any sort to be inserted in a program. A1l characters, including statement separators and blanks may be included. Their usual meanings are ignored. A REM is terminated only by return.
* `HOME`
    * Moves cursor to upper left screen position within the scrolling window and clears all text within the window.
* `HTAB`
    * Moves the cursor to the position that is `aexpr` positions from the left edge of the current screen line.
    * If `aexpr` is a real, it is converted to an integer.
* `VTAB`
    * Moves the cursor to the line that is `aexpr` lines down on the screen. The top line is line l; the bottom line is line 24. This statement may involve moving the cursor either up or down, but never to the right or left.
    * If `aexpr` is a real, it is converted to an integer.
* `TAB`
    * TAB must be used in a PRINT statement, and `aexpr` must be enclosed in parentheses.
    * TAB moves the cursor to the position that is `aexpr` printing positions from the left margin of the text window if `aexpr` is greater than the value of the current cursor position relative to the left margin.
    * If `aexpr` is less than the value of the current cursor position, then the cursor is not moved.
    * TAB never moves the cursor to the left (use HTAB for this).
    * If `aexpr` is a real, it is converted to an integer.
* `NORMAL`
    * Sets the mode to the usual white letters on a black background.
* `INVERSE`
    * Sets the video mode so that the computer's output prints as black letters on a white background.
* `FLASH`
    * Sets the video mode to "flashing", so the output from the computer is alternately shown on the screen in white on black and then reversed to black on white.
* `CLEAR`
    * Zeroes all float, integer and strings variables. Zeroes all arrays.

##### Input / Output
* `PRINT`
    * Print a string, a float, an integer, variable or an expression.
    * Without any arguments, PRINT causes line feed
    * Multiple arguments may be separated by commas (`,`) and/or semicolons (`;`).
    * If an item on the list is followed by a semicolon, then the first character of the next item to be printed will appear immediatly after the current item.
    * If an item on the list is followed by a comma, then the first character of the next item to be printed will appear in the first position of the next available tab field.
    * Tab fields are 14 positions wide
    * If neither a comma nor a semi-colon ends the list, a line feed and return are executed following the last item printed.
* `INPUT`
    * If the optional string is left out, `INPUT` prints a question mark and waits for the user to type a number (if var is an arithmetic variable) or characters (if var is a string variable). The value of this number or string is put into var.
    * When the string is present, it is printed exactly as specified; no question mark, spaces, or other punctuation are printed after the string. Note that only one optional string may be used. It must appear directly after `INPUT` and be followed by a semi-colon. 
    * `INPUT` will accept only a real or an integer as numeric input, not an arithmetic expression. The characters space, +, -, E, and the period are legitimate parts of numeric input. `INPUT` will accept any of these characters or any concatenation of these characters in acceptable form (e.g. +E- is acceptable, +- is not); such input by itself evaluates as Ø.
    * In numeric input, spaces in any position are ignored. If numeric input which is not a real, an integer, a comma or a colon, the message `?REENTER` is displayed and the `INPUT` instruction re-executed.
    * Similarly, a response assigned to a string variable must be a single string or literal, not a string expression. Spaces preceding the first character are ignored.
* `GET` var$ | var | var%
    * Fetches a single character from the keyboard without displaying it on the screen and without requiring that the RETURN key be pressed.
    * The character corresponding to the pressed key is assigned to the variable `var$` | `var%` | `var`

##### Flow Control
* `FOR ... TO ... STEP ... NEXT`
    * Looping constructs.
* `IF ... THEN ... ELSE ...`
    * Conditional execution. `THEN` and/or `ELSE` branchs can have multiple instructions.
* `GOTO`
    * Jump to the given line.
* `GOSUB ... RETURN`
    * Used to call the subroutines at the specified line.

##### System and Utilities
* `END`
    * Exit the program.

#### Supported operators
* `=`
* `<>`
* `<`
* `>`
* `<=`
* `>=`
* `+`
* `-`
* `*`
* `/`
* `^`

#### Supported functions
##### Maths functions
* `INT`
    * Returns the largest integer less than or equal to `aexpr`.

        > If `aexpr` is a positive number, then the largest whole number can be found by chopping off the decimal part.
        
        > If `aexpr` is a negative number, the largest whole number can be found by moving down to the next lowest whole number (that is, make a negative number more  negative).
* `ABS`
    * Returns the absolute value of `aexpr`.

        > `aexpr` if `aexpr` >= 0

        > -`aexpr` if `aexpr` < 0
* `SGN`
    * Returns the sign of `aexpr`.

        > `1` if `aexpr` > 0

        > `0` if `aexpr` = 0

        > `-1` if `aexpr` < 0
* `SQR`
    * Returns the square root of `aexpr`.
    * `aexpr` must be positive or null.

##### String functions
* `LEFT$(sexpr, aexpr)`
    * This function returns the first (leftmost) `aexpr` characters of `sexpr`.
    * If `aexpr` < 1 the message `?ILLEGAL QUANTITY ERROR` is displayed.
    * If `aexpr` is a real, it is converted to an integer.
    * If `aexpr` > `LEN(sexpr)`, only the characters which constitute the string are returned. Any extra positions are ignored.

#### Differences with Applesoft BASIC
##### Variable names
1. In Applesoft BASIC, a variable name may be up to 238 characters long, but APPLESOFT uses only the first two characters to distinguish one name from another. Thus, the names `GOOD4NOUGHT` and `GOLDRUSH` refer to the same variable.

    > With BASICS, all characters are significant. Thus, the names `GOOD4NOUGHT` and `GOLDRUSH` refer to two different variables.

    > Remember that, with BASICS, variable names can be in uppercase, lowercase or mixed case.

2. Certain words used in APPLESOFT BASIC commands are "reserved" for their specific purpose. You cannot use these words as variable names or as part of any variable name. For instance, `FEND` would be illegal because `END` is a reserved word.

    > With BASICS `END` is illegal as a variable name, as `FEND` is totally legal.

##### Extended instructions set
* GOTO support use of identifier and complex expressions. You can write:
```
10 A=10
20 GOTO (A*3)+10
30 PRINT "Hello"
40 END
```

* GOSUB support use of identifier and complex expressions. You can write:
```
10 PRINT "Hello ":A=50:GOSUB A*2:PRINT "!!!"
30 END
100 PRINT "World":RETURN
```

* HTAB and VTAB support user of identifier and complex expressions. You can write:
```
10 FOR A = 1 TO 15
20 HTAB A * 2: VTAB A: PRINT A
30 NEXT A
40 B = 0
50 FOR A = 15 TO 0 STEP -3
60 HTAB A * 2
70 VTAB A/3+12+B
80 PRINT A
90 B = B + 2
100 NEXT A
110 FOR A = 1 TO 15
120 HTAB A * 2: VTAB 24: PRINT A
130 NEXT A
```

### Differences with old computers
#### Extended charset
* Bitmap fonts are handcoded using `MakeBitmapFont.xlsx`. File is located under `./.bintools`.
* BASICS support extended charset, such as:
    * Uppercase letters
        * `ABCDEFGHIJKLMNOPQRSTUVWXYZ`
    * Lowercase letters
        * `abcdefghijklmnopqrstuvwxyz`
    * Accented letters
        * `aàáâä eèéêë iîï oôö uùûü ÿ nñ cç µ`
    * Symbols
        * `[] () {} \/`
        * `!?¡¿ '#&@ $£¥€ %‰ *+-^÷±= .:,; <>«» _~ ²°§¤©¶`
        * `← → ↑ ↓`
    * Box drawing: some blocs, simple, double and mixed lines
        * `░ ▒ ▓`
        * `─│ ┌┐ └┘ ┼ ├┤ ┬┴`
        * `═║ ╔╗ ╚╝ ╬ ╠╣ ╦╩`
        * `╞╡ ╥╨ ╪ ╟╢ ╙╜ ╓╖ ╫ ╒╕ ╘╛ ╧╤`

#### Supported display device
* Results display in program execution could be in `terminal mode` (usefull for debug or test sessions) or in `graphic mode`
* In `terminal mode`, you cannot have any graphic primitives

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Project Status

✅ In active development — architecture being finalized.

## Contribution

Contributions are welcome:

* Fork
* Branch
* Commit with clear messages
* Pull Request

## Current release statistics
Language|Files|Blank|Comment|Code
:-------|-------:|-------:|-------:|-------:
Go|179|2_847|918|17_031
Basic|112|0|0|594
Markdown|3|109|0|454
JSON|1|0|0|66
Text|1|0|0|2
**TOTAL:**|**296**|**2_956**|**918**|**18_147**

## Author

Project developed by **ultra-sonic-28**