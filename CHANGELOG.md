# Changelog
All notable changes to this project will be documented in this file.

The format is based on Keep a Changelog (https://keepachangelog.com), and this project adheres to Semantic Versioning (https://semver.org).

## [Unreleased] - 2026-01-28
### Added
- Add `GET` support in Apple II Basic. Add relevant unit tests.
- Add new `input` package to properly handle inputs in Ebiten, TTY and unit test modes.
- Add support for real and integer variables in GET instruction
- Distribute example files and programs into subdirectories.

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
