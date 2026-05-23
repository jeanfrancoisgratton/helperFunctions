# helperFunctions
---
<img src="./_images/helperFunctions_banner.png" alt="helperFunctions lib logo" height="400" width="512" /><br>

A personal Go utility library — a single source of truth for functions that keep appearing across my tools.

Current version: **v5**  
Module path: `github.com/jeanfrancoisgratton/helperFunctions/v5`

---

## Installation

```bash
go get github.com/jeanfrancoisgratton/helperFunctions/v5
```

---

## Package overview

The library is organized as a root package plus several subpackages, each with a focused area of responsibility.

| Package | Import path | Responsibility |
|---|---|---|
| `helperFunctions` | `.../v5` | Miscellaneous utilities (number formatting, string reversal, changelog) |
| `terminalfx` | `.../v5/terminalfx` | Terminal size, clearing, text alignment, ANSI colours, glyphs |
| `logging` | `.../v5/logging` | Levelled, structured logging to stdout or a file |
| `prettyjson` | `.../v5/prettyjson` | Colourized JSON pretty-printer (jq-style) |
| `networking` | `.../v5/networking` | DNS/reverse-lookup helpers |
| `repomanagement` | `.../v5/repomanagement` | Parse Git remote URLs (HTTPS and SCP-style SSH) |
| `pager` | `.../v5/pager` | Terminal pager with sticky banner/footer (à la `more`) |

---

## Root package — `helperFunctions`

```go
import hf "github.com/jeanfrancoisgratton/helperFunctions/v5"
```

### Numeric formatting

```go
func SI(nombre interface{}) string
```

Formats any integer or float as a comma-separated SI string.

```go
hf.SI(1234567)   // → "1,234,567"
hf.SI(-9876543)  // → "-9,876,543"
hf.SI(3.14159)   // → "3"  (rounded to nearest integer)
```

Accepts `int`, `int8`–`int64`, `uint`–`uint64`, `float32`, `float64`.
Returns `"Invalid input"` for anything else.

### String utilities

```go
func ReverseString(s string) string
```

Returns the Unicode-safe reversal of `s`.

```go
hf.ReverseString("abcdef")  // → "fedcba"
```

### Changelog display

```go
func ChangeLog(cl string, clear bool)
```

Prints the string `cl` to stdout. If `clear` is `true` the terminal is cleared first via `terminalfx.ClearTTY`.

### Prompts

```go
func GetStringValFromPrompt(prompt string) string
func GetIntValFromPrompt(prompt string) int
func GetBoolValFromPrompt(prompt string) bool
func GetStringSliceFromPrompt(prompt string) []string
func GetValueFromPrompt(prompt string) interface{}
```

Read a typed value from stdin after displaying `prompt`. `GetStringSliceFromPrompt` collects lines until the user submits an empty line. `GetValueFromPrompt` attempts to parse the input as `uint`, `int`, or `bool` in that order, falling back to `string`.

```go
name := hf.GetStringValFromPrompt("Enter your name: ")
age  := hf.GetIntValFromPrompt("Enter your age: ")
ok   := hf.GetBoolValFromPrompt("Proceed? [true/false]: ")
tags := hf.GetStringSliceFromPrompt("Enter tags (blank line to finish):")
```

### Passwords and AES encryption

```go
func GetPassword(prompt string, debugmode bool) string
func EncodeString(string2encode string, privateKey string) string
func DecodeString(encodedstring string, privateKey string) string
```

`GetPassword` reads a password without echoing it. If `debugmode` is `true` it falls back to a plain `GetStringValFromPrompt` (useful during development).

`EncodeString`/`DecodeString` use AES-256-CFB with a SHA-256-derived key.
`privateKey` must be exactly 32 bytes; a built-in default key is used if it is not.

```go
pass    := hf.GetPassword("Password: ", false)
encoded := hf.EncodeString("s3cr3t", myKey)
decoded := hf.DecodeString(encoded, myKey)
```

> **Note:** The default key is intentionally weak. Always supply your own 32-byte key in production.

---

## `terminalfx`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/terminalfx"
```

### Terminal introspection and control

```go
func GetTerminalSize() (cols int, rows int)
func ClearTTY()
```

`GetTerminalSize` queries the kernel via `TIOCGWINSZ`; returns `(0, 0)` on failure.  
`ClearTTY` sends ANSI escape sequences to clear the screen and home the cursor.

### Text alignment

```go
func Center(input string) string
func Right(input string) string
```

Returns `input` padded with leading spaces so it appears centred or right-aligned in the current terminal width (defaults to 80 columns if the size cannot be determined).

### ANSI colour wrappers

Each function wraps `input` with bold ANSI colour codes via [gchalk](https://github.com/jwalton/gchalk) and returns the colourized string.

```go
func Red(input string) string
func Green(input string) string
func White(input string) string
func Yellow(input string) string
func Blue(input string) string
```

These are also used by `prettyjson.DefaultStyle` to colourize JSON tokens.

### Glyphs

Prepend a Unicode glyph to `sentence` and return the result. The sentence itself is always rendered in the default terminal colour.

| Function | Glyph | Semantic |
|---|---|---|
| `EuropeanStopSign(s)` | ⛔ | Hard stop |
| `AmericanStopSign(s)` | 🛑 | Hard stop |
| `BombSign(s, coloured)` | 💥 | Fatal error |
| `SkullBonesSign(s)` | ☠ | Fatal error |
| `EnabledSign(s)` | ✅ | Accepted / success |
| `ErrorSign(s)` | ❌ | Rejected / failure |
| `GreenGoSign(s)` | 🟢 | All clear |
| `InProgressSign(s)` | ⏳ | Task in progress |
| `WarningSign(s)` | ⚠ | Warning |
| `InfoSign(s)` | 🛈 | Informational |
| `NoteSign(s)` | 💬 | Note / comment |
| `ScrollSign(s)` | 📜 | Document |
| `TipSign(s)` / `LightbulbSign(s)` | 💡 | Tip |
| `ThumbsUpSign(s)` | 👍 | Positive |
| `ThumbsDownSign(s)` | 👎 | Negative |
| `NotExistMathSign(s)` | ∄ | Does not exist |
| `ExistMathSign(s)` | ∃ | Exists |
| `NotIncludedInMathSign(s)` | ∉ | Not a member of |
| `IsIncludedInMathSign(s)` | ∈ | Member of |
| `DeltaSymbolMathSign(s)` | ∆ | Delta / change |

---

## `logging`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/logging"
```

A levelled logger that writes to stdout or a file.

### Log levels

```
None < Error < Info < Debug
```

`Userf` is outside the ladder: it fires whenever the level is anything other than `None`.

### Initialization

```go
func Init(path string, level LogLevel, opts LogInitOptions) error
func Close()
```

`path` of `""` or `"-"` routes output to stdout; any other value opens (or creates) a log file with mode `0640`. Calling `Init` again rotates to the new target.

```go
err := logging.Init("/var/log/mytool.log", logging.Info, logging.LogInitOptions{
    EntryPrefix:        "mytool",
    DisplayPID:         true,
    DisplayExecName:    true,
    DisplayCurrentUser: true,
})
defer logging.Close()
```

`LogInitOptions` fields:

| Field | Type | Description |
|---|---|---|
| `EntryPrefix` | `string` | Prepended to every log entry (wrapped in `<>`) |
| `UserHeader` | `string` | Header label used by `Userf` when none is supplied (default `[USER]`) |
| `DisplayCurrentUser` | `bool` | Append the OS username to each entry |
| `DisplayExecName` | `bool` | Append the executable name (`os.Args[0]`) |
| `DisplayPID` | `bool` | Append the process PID |

### Emitting log lines

```go
func Debugf(msg string, args ...any)
func Infof(msg string, args ...any)
func Errorf(msg string, args ...any)
func Userf(msg string, header string, args ...any)
```

All functions accept a `fmt.Sprintf`-style format string. `Userf` takes an optional `header` string; if empty, the `UserHeader` set in `Init` is used. The header is automatically bracketed if it is not already.

### Level management

```go
func SetLevel(l LogLevel)
func GetLevel() LogLevel
func Enabled(l LogLevel) bool
func ParseLevel(s string) LogLevel   // "none"|"error"|"info"|"debug" → LogLevel
```

```go
logging.SetLevel(logging.Debug)
logging.Debugf("cache miss for key %q", key)
logging.Infof("server started on :%d", port)
logging.Errorf("dial failed: %v", err)
logging.Userf("deploying %s", "release/v2", "BUILD")
```

Log line format:

```
2025-11-09 10:58:42 [INFO] <mytool> (alice) mytoolbin (PID 12345) server started on :8080
```

---

## `prettyjson`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/prettyjson"
```

Pretty-prints JSON with optional ANSI colourization, à la `jq`.

### Functions

```go
func Print(payload []byte, opts ...Option) error
func SPrint(payload []byte, opts ...Option) (string, error)
func Format(payload []byte, opts ...Option) ([]byte, error)
```

`Print` writes to stdout by default. `SPrint` returns a string. `Format` returns bytes.

### Options

| Option constructor | Effect |
|---|---|
| `WithWriter(w io.Writer)` | Redirect output |
| `WithIndent(s string)` | Indentation string (default: two spaces) |
| `WithSortKeys(b bool)` | Sort object keys (default: `true`) |
| `WithColorMode(m ColorMode)` | `ColorAuto` (default) / `ColorAlways` / `ColorNever` |
| `WithStyle(s Style)` | Override the colour style |

### Colour modes

`ColorAuto` (the default) enables colour only when the writer is a real TTY. Use `ColorAlways` to force ANSI codes into a captured string, or `ColorNever` to disable them entirely.

### Styles

```go
prettyjson.DefaultStyle()  // Blue keys, Green strings, Yellow numbers/bools, Red null, White punctuation
prettyjson.PlainStyle()    // No colour at all
```

You can also build a custom `Style`:

```go
custom := prettyjson.Style{
    Key:    terminalfx.Yellow,
    String: terminalfx.Green,
}
prettyjson.Print(data, prettyjson.WithStyle(custom))
```

### Example

```go
data := []byte(`{"name":"alice","age":30,"active":true}`)
prettyjson.Print(data)
// with default style and a TTY:
// {
//   "active": true,
//   "age": 30,
//   "name": "alice"
// }

s, _ := prettyjson.SPrint(data, prettyjson.WithColorMode(prettyjson.ColorNever))
fmt.Println(s)
```

---

## `networking`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/networking"
```

DNS utility functions built on the system resolver and `golang.org/x/net/publicsuffix`.

```go
func ReverseLookupIP(ip string) ([]string, error)
func GetDomainFromIP(ip string) (fqdn string, tld string, err error)
func GetDomainFromHostname(hostname string) (tld string, err error)
```

`ReverseLookupIP` returns all PTR records for an IP address (trailing dots stripped).

`GetDomainFromIP` calls `ReverseLookupIP` and then extracts the effective TLD+1 (i.e. the registrable domain) from the first PTR record. Returns both the full hostname and the registrable domain.

`GetDomainFromHostname` extracts the registrable domain directly from a hostname string, without a DNS lookup.

```go
names, err := networking.ReverseLookupIP("8.8.8.8")
// names → ["dns.google"]

fqdn, tld, err := networking.GetDomainFromIP("8.8.8.8")
// fqdn → "dns.google", tld → "google.com"

tld, err := networking.GetDomainFromHostname("mail.example.co.uk")
// tld → "example.co.uk"
```

---

## `repomanagement`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/repomanagement"
```

Parses Git remote URLs — both standard HTTPS and SCP-style SSH — into a structured form.

```go
func ExtractRepoInfo(raw string) (RepoInfo, error)
```

`RepoInfo` fields:

| Field | Example | Description |
|---|---|---|
| `Scheme` | `"https"`, `"ssh"` | URL scheme |
| `Host` | `"github.com"` | Hostname |
| `Port` | `"2222"` | Port, if explicitly present (otherwise empty) |
| `TopLevelOwner` | `"myorg"` | First path segment (top-level group or username) |
| `FullOwnerPath` | `"myorg/subgroup"` | All path segments except the repository name |
| `Repo` | `"myrepo"` | Repository name (`.git` suffix stripped) |
| `RawURL` | — | The original input string |

```go
// HTTPS
info, _ := repomanagement.ExtractRepoInfo("https://github.com/myorg/subgroup/myrepo.git")
// info.Host          → "github.com"
// info.FullOwnerPath → "myorg/subgroup"
// info.Repo          → "myrepo"

// SCP-style SSH
info, _ = repomanagement.ExtractRepoInfo("git@gitlab.com:myorg/myrepo.git")
// info.Scheme → "ssh"
// info.Host   → "gitlab.com"
// info.Repo   → "myrepo"
```

---

## `pager`

```go
import "github.com/jeanfrancoisgratton/helperFunctions/v5/pager"
```

A terminal pager similar to `more`, with optional sticky banner and footer lines.

### Entry point

```go
func Page(lines []string, bannerSize, footerSize int) error
```

`lines` is the complete content to page through, banner and footer included:

- The first `bannerSize` lines are the **sticky banner** — they are reprinted at the top of every page.
- The last `footerSize` lines are the **sticky footer** — reprinted at the bottom of every page.
- Everything in between is the **scrollable content window**.

If `bannerSize + footerSize >= terminal height`, both values are silently ignored and all lines are treated as scrollable content. Negative values are clamped to zero.

`Page` switches stdin to raw mode for the duration of the call and restores it before returning, whether it exits normally or because of an error.

### Keybindings

| Key | Action |
|---|---|
| `Space` | Advance one full page |
| `Enter` | Advance one line |
| `Ctrl+U` | Go back one full page |
| `Ctrl+G` | Show current line position in the status bar; press any key to continue |
| `Q` / `q` | Quit and clear the screen |

### Status bar

While paging:
```
-- More -- [42/200]  Space=▶page  Enter=▶line  ^U=◀page  ^G=pos  Q=quit
```

At the end of content:
```
(END) [200/200]  Q to quit
```

After `Ctrl+G`:
```
[Line 42/200]  -- press any key --
```

### Example

```go
import (
    "strings"
    "github.com/jeanfrancoisgratton/helperFunctions/v5/pager"
)

// Assume `report` is a multi-line string.
// First 3 lines: title/header block (sticky banner)
// Last 1 line:   legend or key-binding reminder (sticky footer)
lines := strings.Split(report, "\n")

if err := pager.Page(lines, 3, 1); err != nil {
    log.Fatal(err)
}
```

```go
// No sticky regions — plain pager
if err := pager.Page(lines, 0, 0); err != nil {
    log.Fatal(err)
}
```

---

## Dependencies

| Module | Used by |
|---|---|
| `github.com/jwalton/gchalk` | `terminalfx` (ANSI colour output) |
| `golang.org/x/crypto` | root package (SSH terminal / password prompt) |
| `golang.org/x/net` | `networking` (public suffix list) |
| `golang.org/x/term` | `prettyjson` (TTY detection), `pager` (raw mode) |

---

## Changelog

See [CHANGELOG.md](CHANGELOG.md).
