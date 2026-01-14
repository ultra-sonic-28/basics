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
#### Supported instructions set
* `REM`
* `PRINT`
* `LET`
* `FOR ... TO ... STEP ... NEXT`
* `IF ... THEN ... ELSE ...`
* `GOTO`
* `GOSUB ... RETURN`
* `HTAB`
* `VTAB`  
* `END`

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

#### Extended instructions set
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