# Changelog
All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog (https://keepachangelog.com), and this project adheres to Semantic Versioning (https://semver.org).

## [Unreleased] - 2026-02-23
### Added
- Add `VAL()` support in Apple II Basic. Add relevant unit tests.

## [Unreleased] - 2026-02-22
### Added
- Add `STR$()` support in Apple II Basic. Add relevant unit tests.

## [Unreleased] - 2026-02-21
### Fixed
- Update `architecture.md`.

## [Unreleased] - 2026-02-20
### Added
- Add `PR#0` support in Apple II Basic to switch to 40 columns mode.
- Add dynamic text mode switching.

## [Unreleased] - 2026-02-19
### Added
- Add `PR#3` support in Apple II Basic to switch to 80 columns mode.

### Fixed
- Done some cleansing.
- Fix a bug in magefile as `Clean` target missed some debug files for deletion.

## [Unreleased] - 2026-02-18
### Fixed
- Fix a bug when testing strings for equality.
- Add a small basic game `nicoma.bas`.
- Add operator unit tests.
- `REM` statement will no longer stripped during lexing
- Fix a bug in magefile as `Clean` target doesn't clean anything.

## [Unreleased] - 2026-02-17
### Added
- Add `RIGHT$()` support in Apple II Basic. Add relevant unit tests.
- Add `MID$()` support in Apple II Basic. Add relevant unit tests.
- Add `LEN()` support in Apple II Basic. Add relevant unit tests.

### Fixed
- Add example programs for handling strings length over 255 characters. Update README.

## [Unreleased] - 2026-02-16
### Added
- Add `LEFT$()` support in Apple II Basic. Add relevant unit tests.

### Fixed
- Add missing unit test for `Abs` in `EvalExpr`.
- Fix error in unit test for `Abs` in `EvalExpr`.
- Add missing unit test for `Sgn` in `EvalExpr`.
- Add missing unit test for `Int` in `EvalExpr`.
- Fix error in `Int` when negative value is already an int.
- Add missing unit test for `Sqr` in `EvalExpr`.

## [Unreleased] - 2026-02-14
### Changed
- Some refactorization for future proof parser.

## [Unreleased] - 2026-02-12
### Added
- Add `SQR()` support in Apple II Basic. Add relevant unit tests.
- Add `TAB()` support in Apple II Basic. Add relevant unit tests.
- Add a small basic game `literature-quizz.bas`.

## Changed
- Add more documentation details regarding `HTAB` and `VTAB` statements, and `TAB` instruction.

### Fixed
- Fix of accented letter input with Ebiten.
- Add missing keyword unit tests.
- Add missing expression dump unit tests.
- Add unit test for string expression when using `TAB`.

## [Unreleased] - 2026-02-11
### Fixed
- Fix mirorred accentued letters.
- Complete redesign of the 7x8 font.
- Display "=== PROGRAM RESULTS ===" in TTY mode only.
- Wording changes in README.

## [Unreleased] - 2026-02-10
### Added
- SHA256 checksum generation for archive releases.

## [v0.1.0.5] - 2026-02-10
### Added
- Add release packaging.

### Fixed
- Add icon in 24x24 for interfaces with 125% scaling or some modern controls in Windows 10/11.

## [Unreleased] - 2026-02-09
### Added
- Add `build` target for `Mage`.
- Add `tools` target for `Mage`.
- Add application icons.
- Generate Windows binary.
- Add `stats` target for `Mage`. Compute project code statistics using [cloc](https://github.com/AlDanial/cloc).
- Update build version number before each build.

## [Unreleased] - 2026-02-06
### Added
- Use `Mage` as development toolchain.

### Changed
- Renamed `./cmd/basic` directory to `./cmd/basics`

## [Unreleased] - 2026-02-05
### Fixed
- Fix misuse of assertions in `StripANSI` unit tests.
- Add unit tests for `VarType` and `VarTypeAsInt`.
- Add unit tests for `<`, `<=`, `<>`, `>` and `>=` operators.

## [Unreleased] - 2026-02-04
### Added
- Set and Get values to/from arrays. Add relevant unit tests.

### Fixed
- Add missing unit tests for parser helper (PrintStmt).
- Add missing unit tests for parser debug (DimStmt).
- Add missing unit tests for parser helper (DimStmt).
- Add new test example for flow control (`FOR ... TO ... STEP ... NEXT`).
- Add unit tests for `Flatten(...)` in `common` package.

## [Unreleased] - 2026-02-02
### Added
- Add `DIM` support in Apple II Basic. Add relevant unit tests.
- `CLEAR` support for arrays. Add relevant unit tests.

### Changed
- Update README.

## [Unreleased] - 2026-01-30
### Added
- Add `FLASH` support in Apple II Basic. Add relevant unit tests.
- Add `CLEAR` support in Apple II Basic. Add relevant unit tests.

### Fixed
- Add missing unit tests for parser helper.
- Add missing unit tests for parser debug functions.
- Add missing unit tests for interpreter when running programs.
- Add renderer's interface unit tests for video package.
- Add provider's interface unit tests for video package.
- Add Mode interface and ModeInfo struct unit tests for video package.
- Add EbitenDevice interface unit tests for video package.
- Add Device interface unit tests for video package.
- Add unit tests for font package (DefaultFont).
- Add unit tests for font package (BitmapFont struct and Glyph function).

## [Unreleased] - 2026-01-29
### Added
- Add `NORMAL` support in Apple II Basic. Add relevant unit tests.
- Add `INVERSE` support in Apple II Basic. Add relevant unit tests.

## [Unreleased] - 2026-01-28
### Added
- Add `GET` support in Apple II Basic. Add relevant unit tests.
- Add new `input` package to properly handle inputs in Ebiten, TTY and unit test modes.
- Add support for real and integer variables in GET instruction
- Distribute example files and programs into subdirectories.
- Add some examples in unit test.

### Changed
- Refactor input/output devices implementation.

### Fixed
- Disable keyboard input after program ended
- Fix blinking cursor when hitting backspace or enter keys

## [Unreleased] - 2026-01-27
### Added
- Add blinking cursor support using Ebiten in Apple II Basic.
- Add CHANGELOG. Changelog follows "Keep a Changelog" recommendations (cf https://keepachangelog.com)

### Changed
- Paradigm change in screen display, use Ebiten library even for text modes.

## [Unreleased] - 2026-01-23
### Added
- Add blinking cursor support in Apple II Basic.

## [Unreleased] - 2026-01-21
### Added
- Add `INPUT` support in Apple II Basic. Add relevant unit tests.

## [Unreleased] - 2026-01-20
### Added
- Support `;` and `,` as terminal separator in `PRINT` statement.
- Support standalone `PRINT` statement.

### Changed
- Update README.

## [Unreleased] - 2026-01-19
### Added
- Add `ABS()` support in Apple II Basic. Add relevant unit tests.
- Add `SGN()` support in Apple II Basic. Add relevant unit tests.
- Add missing AST unit tests.

### Changed
- Update environment unit test.

### Fixed
- Fix and upgrade program execution trace.
- Cosmetic test runner fixes.
- Fix boolean operator comparaison.

## [Unreleased] - 2026-01-18
### Added
- Add `INT()` support in Apple II Basic. Add relevant unit tests.

### Changed
- Parser debug and logger refactoring.

### Fixed
- Fix and extend dump and log functions.

## [Unreleased] - 2026-01-17
### Added
- Add `INTEGER` and `STRING` variables support in Apple II Basic. Add relevant unit tests.

## [Unreleased] - 2026-01-16
### Added
- Add more logging stuff.
- Add missing statement names in parser helper.
- Add interpreter tests using the BASIC source code from the example files.

### Changed
- Update README with descriptions for implemented BASIC instructions.
- Updated README.

## [Unreleased] - 2026-01-15
### Fixed
- Fix `IF ... THEN ... ELSE ...` with `GOSUB` statement. Add more logging stuff.

## [Unreleased] - 2026-01-13
### Added
- Add `IF ... THEN ... ELSE ...` support in Apple II Basic. Add missing `ELSE` keyword. Add relevant unit tests. Update runtime environment to handle boolean values. Add `<=`, `>=` and `<>` support in lexer and parser. Add lexer and parser dump to log file.
- Add tests for parser dump: `GOTO` statement.
- Add tests for parser dump: `END` statement.
- Add tests for parser dump: `HTAB` and `VTAB` statements.
- Add tests for parser dump: `GOTO` statement with complex expressions.
- Add support for math power (`^` operator).
- Add README.
- Add `GOSUB ... RETURN` support in Apple II Basic. Add relevant unit tests.
- Add `HOME` support in Apple II Basic. Add relevant unit tests.

### Fixed
- Fix README.
- Fix filename mispelling

## [Unreleased] - 2026-01-13
### Added
- Add `GOTO` support in Apple II Basic. Add relevant unit tests. Also flatten intricate loops in interpreter.

### Fixed
- Fix parser dump, add `END` statement support.
- Fix parser dump, add `HTAB` and `VTAB` statement support.
- Fix parser dump, add `GOTO` statement support.

## [Unreleased] - 2026-01-12
### Added
- Add log mechanism
- Add lexer unit test for `HTAB` and `VTAB` (Apple II Basic)
- Add parser unit test for `HTAB` and `VTAB` (Apple II Basic)
- Add parser error unit test for `HTAB` and `VTAB` (Apple II Basic)

### Fixed
- Fix missing `main_test.go` in `internal/runtime`
- Fix title in unit tests report
- Fix parser error when line is only linenumber and add relevant unit test
- Fix that parsing `10 LET A 3` and `10 A 3` raise same error (EXPECTED '=' IN ...) and add relevant unit test

## [Unreleased] - 2026-01-11
### Added
- Add `HTAB` and `VTAB` support in Apple II Basic
- Add `END` support in Apple II Basic
- Add parser unit test for `END` (Apple II Basic)
- Add lexer unit test for `END` (Apple II Basic)

## [Unreleased] - 2026-01-09
### Added
- Added machine screen display abstractions (Apple II and TTY) using a runtime environment and rendering engine
- Add MIT Licence

### Fixed
- Basic type constants unit tests after adding TTY

## [Unreleased] - 2026-01-08
### Added
- Project initialization
- Add .gitignore
- Lexer, parser, interpreter implemntation for a very little subset of Applesoft BASIC
