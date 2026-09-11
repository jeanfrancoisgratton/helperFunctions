// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: terminalfx/colours.go
// Original timestamp: 2026.09.10 22:28:21
package terminalfx

import "charm.land/lipgloss/v2"

// COLOR FUNCTIONS
// ===============
//
// Every colour has a bright variant (e.g. Red) and a dim/darker-shade
// variant (e.g. DimRed). The base eight ANSI colours use the terminal's
// native bright (9-15) and normal (1-7) intensity indices, so they degrade
// correctly on 16-colour terminals. Colours with no ANSI equivalent (Brown,
// Copper, Bronze, Gold, Silver, Gray, Orange) use hand-picked truecolor hex
// pairs instead.

var (
	redBright = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
	redDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))

	greenBright = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	greenDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	yellowBright = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
	yellowDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	blueBright = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	blueDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))

	magentaBright = lipgloss.NewStyle().Foreground(lipgloss.Color("13"))
	magentaDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))

	cyanBright = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
	cyanDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))

	whiteBright = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	whiteDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))

	brownBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#A0522D"))
	brownDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#5C2E15"))

	copperBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#B87333"))
	copperDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#6E4520"))

	bronzeBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#CD7F32"))
	bronzeDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#7A4B1E"))

	goldBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700"))
	goldDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#B8860B"))

	silverBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#C0C0C0"))
	silverDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#808080"))

	grayBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#A9A9A9"))
	grayDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#555555"))

	orangeBright = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFA500"))
	orangeDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("#B36200"))
)

func Red(sentence string) string    { return redBright.Render(sentence) }
func DimRed(sentence string) string { return redDim.Render(sentence) }

func Green(sentence string) string    { return greenBright.Render(sentence) }
func DimGreen(sentence string) string { return greenDim.Render(sentence) }

func Yellow(sentence string) string    { return yellowBright.Render(sentence) }
func DimYellow(sentence string) string { return yellowDim.Render(sentence) }

func Blue(sentence string) string    { return blueBright.Render(sentence) }
func DimBlue(sentence string) string { return blueDim.Render(sentence) }

func Magenta(sentence string) string    { return magentaBright.Render(sentence) }
func DimMagenta(sentence string) string { return magentaDim.Render(sentence) }

func Cyan(sentence string) string    { return cyanBright.Render(sentence) }
func DimCyan(sentence string) string { return cyanDim.Render(sentence) }

func White(sentence string) string    { return whiteBright.Render(sentence) }
func DimWhite(sentence string) string { return whiteDim.Render(sentence) }

func Brown(sentence string) string    { return brownBright.Render(sentence) }
func DimBrown(sentence string) string { return brownDim.Render(sentence) }

func Copper(sentence string) string    { return copperBright.Render(sentence) }
func DimCopper(sentence string) string { return copperDim.Render(sentence) }

func Bronze(sentence string) string    { return bronzeBright.Render(sentence) }
func DimBronze(sentence string) string { return bronzeDim.Render(sentence) }

func Gold(sentence string) string    { return goldBright.Render(sentence) }
func DimGold(sentence string) string { return goldDim.Render(sentence) }

func Silver(sentence string) string    { return silverBright.Render(sentence) }
func DimSilver(sentence string) string { return silverDim.Render(sentence) }

func Gray(sentence string) string    { return grayBright.Render(sentence) }
func DimGray(sentence string) string { return grayDim.Render(sentence) }

func Orange(sentence string) string    { return orangeBright.Render(sentence) }
func DimOrange(sentence string) string { return orangeDim.Render(sentence) }
