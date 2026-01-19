# Basics - BASIC Interpreter for old computers

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

## Version control system
* **Jujutsu (jj) + Git**

## Supported computers

### APPLE II
### Real, Integer and String variables
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
##### Editing and format related
* `REM`
    * This serves to allow text of any sort to be inserted in a program. A1l characters, including statement separators and blanks may be included. Their usual meanings are ignored. A REM is terminated only by return.
* `HOME`
    * Moves cursor to upper left screen position within the scrolling window and clears all text within the window.
* `HTAB`
    * Moves the cursor to the position that is `aexpr` positions from the left edge of the current screen line.
* `VTAB`
    * Moves the cursor to the line that is `aexpr` lines down on the screen. The top line is line l; the bottom line is line 24. This statement may involve moving the cursor either up or down, but never to the right or left.

##### Input / Output
* `PRINT`
    * Print a string, a float, an integer, variable or an expression.
    * Multiple arguments may be separated by commas (`,`) and/or semicolons (`;`).
    * If an item on the list is followed by a semicolon, then the first character of the next item to be printed will appear immediatly after the current item.
    * If an item on the list is followed by a comma, then the first character of the next item to be printed will appear in the first position of the next available tab field.
    * Tab fields are 14 positions wide
* `LET`
    * Assign a value to a variable, creating it if necessary. Optionnal.

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

### Supported functions
#### Maths functions
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

## Author

Project developed by **ultra-sonic-28**