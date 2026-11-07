# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [5.9.0] - 2026-11-07

### Added
- `GenerateUUIDorPassword()` function, to generate UUIDs, passwords or plain random strings

### Changed
- Updated the build dependencies (go.mod)

## [5.8.0] - 2026-10-18

### Added
- `EncryptString()`/`DecryptString()`/`EncryptFile()`/`DecryptFile()` — authenticated AES-256-GCM encryption with Argon2id key derivation, as a separate scheme alongside (not a replacement for) `EncodeString`/`DecodeString`/`EncodeFile`/`DecodeFile`

### Changed
- Re-added the `golang.org/x/crypto` dependency (dropped in 5.4.0) for Argon2id

## [5.7.0] - 2026-09-21

### Added
- `terminalfx`: `IsTerminal()` and `Interactive()` to centralize TTY detection
- `prettyjson`: `WithPaging()`

### Changed
- `prettyjson`: `Print` now auto-routes overflowing output through `pager`, forcing colour on internally, when stdout is the default writer and the process is fully interactive; falls back to a plain print otherwise or on any pager failure

## [5.6.0] - 2026-09-15

### Fixed
- `pager`: `Ctrl+G` swallowing the next keystroke instead of acting on it (e.g. `Ctrl+G`, `Q` now quits right away)
- `pager`: terminal resize (`SIGWINCH`) is now picked up live instead of requiring a restart; on Windows the new size is picked up on the next keystroke
- `pager`: lines wider than the terminal are now hard-wrapped so page size/status bar counts stay accurate on narrow terminals

### Removed
- `pager`: a redundant terminal-restore call on quit

## [5.5.0] - 2026-09-10

### Added
- Dim/darker-shade variants for all colours, plus new colours Magenta, Cyan, Brown, Copper, Bronze, Gold, Silver, Gray, Orange
- `Bold()`, `Italic()`, `Underline()`, `Reverse()`, `Blink()`, `Strikethrough()`
- `Left()`
- `HorizontalLine()`

### Changed
- Swapped `gchalk` for `lipgloss/v2` in `terminalfx` (unmaintained since 2022; lipgloss gives a larger palette and is actively developed)
- `Center()`/`Right()` now measure width ANSI-escape- and grapheme-aware instead of by byte length
- Go version bump: 1.27.1

## [5.4.0] - 2026-08-25

### Added
- Package now builds on Windows

### Changed
- Replaced deprecated `golang.org/x/crypto/ssh/terminal` with `golang.org/x/term` in `GetPassword()` and `GetTerminalSize()`

### Removed
- The `golang.org/x/crypto` dependency

## [5.3.2] - 2026-07-26

### Fixed
- The InfoSign glyph that returned a monochrome glyph

## [5.3.1] - 2026-07-16

### Added
- Test files

### Changed
- Doc tweaks

## [5.3.0] - 2026-06-28

### Added
- `EncodeFile()` and `DecodeFile()` — AES-256-CFB file encryption/decryption (mirrors `EncodeString`/`DecodeString`)

## [5.2.2] - 2026-06-18

### Fixed
- Better glyph rendering for some "faulty" glyphs in terminalfx/

## [5.2.1] - 2026-05-31

### Fixed
- SI() so malformed numbers with many decimal places display properly (i.e.: no 1,234.1,150)

## [5.2.0] - 2026-05-30

### Added
- A function to convert bytes to MB, GB or TB

### Removed
- Changelog() from misc.go

## [5.1.0] - 2026-05-26

### Added
- The `pager` subpackage to (...kind of...) emulate Linux's `more` pager

## [5.0.2] - 2026-03-27

### Added
- Mathematical symbols in glyphs.go

## [5.0.1] - 2026-03-15

### Added
- A new package, networking

## [4.1.0] - 2025-12-31

### Added
- A prettyjson subpackage, to prettify JSON output, "à la jq"

## [4.0.0] - 2025-11-13

### Changed
- **BREAKING:** Glyph function calls are cleaner
- **BREAKING:** All log init options are now wrapped in a single struct

### Removed
- **BREAKING:** The customError package from this package, to avoid circular dependency hell

## [3.1.0] - 2025-11-10

### Changed
- Swapped parameters in most glyph functions in terminalfx so that the text will stay in the usual white color

## [3.0.8] - 2025-11-10

### Added
- Another option to prepend the message text with the username calling the tool

## [3.0.6] - 2025-11-09

### Added
- An optional prefix option for each logfile entry

## [3.0.5] - 2025-10-08

### Added
- A new "in progress" glyph

## [3.0.4] - 2025-10-07

### Changed
- Simplified glyph handling and printing

## [3.0.3] - 2025-10-06

### Changed
- Go version bump: 1.25.2

### Removed
- Mentions of customError/v2

## [3.0.2] - 2025-09-26

### Changed
- **BREAKING:** terminal functions are now in their own subpackage

## [2.4.1] - 2025-09-16

### Added
- Coloured (terminal) symbols in terminal.go

## [2.3.1] - 2025-08-25

### Added
- Proc information to Init()

## [2.2.4] - 2025-08-14

### Changed
- Go version bump to 1.25.0
- Updated CustomError dependency (version bump)

## [2.2.2] - 2025-08-13

### Added
- A ParseLevel() helper to convert a string loglevel to a LogLevelIota variable

## [2.2.1] - 2025-08-13

### Changed
- Simplified the logging subpackage, which now handles loglevels itself instead of relying on the calling tool

## [2.1.2] - 2025-08-11

### Removed
- Reverted that warning loglevel; some brain freeze where I confused loglevel and errorlevel

## [2.1.1] - 2025-08-11

### Added
- A "Warning" loglevel; how could this have been forgotten??

## [2.1.0] - 2025-08-09

### Changed
- Hybrid approach to logging (WIP)

## [2.0.4] - 2025-08-09

### Added
- A String method to LogLevel

### Changed
- Updated builddeps

## [2.0.3] - 2025-08-09

### Changed
- Minor refactoring in the logging subpackage
- Go version bump

## [2.0.2] - 2025-08-05

### Changed
- Version bump to ensure that we got rid of customError v1.x.y

## [2.0.1] - 2025-07-27

### Added
- Basic logging facilities

## [1.12.2] - 2025-05-27

### Changed
- Modified the paginating terminal size logic

## [1.12.1] - 2025-05-27

### Fixed
- Error-handling in the Pager() functions

## [1.12.0] - 2025-05-26

### Changed
- Reworked GetPassword() a bit in debug mode

## [1.11.1] - 2025-05-15

### Added
- A pager function, like the `more` command

## [1.10.1] - 2025-04-21

### Added
- A repository information extraction function

## [1.9.0] - 2025-04-17

### Added
- A wrapper so that if we call GetPassword in DEBUG mode, we use

## [1.8.0] - 2025-02-18

### Added
- Center() and Right() for more text formatting features

### Changed
- Go version bump (1.23.6)

## [1.7.0] - 2024-11-19

### Changed
- Updated to Go 1.23.3

### Fixed
- Blue() actually printed in yellow; how come I never noticed that?!

## [1.6.x] - 2024

_No entries recorded._

## [1.5.2] - 2024-06-22

### Added
- Function to check mountpoint type (only supports Linux)

## [1.4.0] - 2024-04-29

### Added
- A generic prompt function

## [1.3.2] - 2024-04-14

### Changed
- Reverted v1.3.1

## [1.3.1] - 2024-04-14

### Changed
- Reverted Encode/Decode strings to a default cipherkey

## [1.3.0] - 2024-04-12

### Changed
- SI() now supports more than uint64, and is also sign-aware

## [1.2.1] - 2024-04-10

### Added
- Changelog()

## [1.1.0] - 2024-04-10

### Added
- Doc, CHANGELOG, extra functions in {encodeDecodePassword,terminal}.go

## [1.0.0] - 2024-04-10

### Added
- Initial version
