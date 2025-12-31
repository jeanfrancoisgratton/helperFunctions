// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: /json/util.go
// Original timestamp: 2025/12/31

package prettyjson

import (
	"io"
	"os"

	"golang.org/x/term"
)

func isTerminalWriter(w io.Writer) bool {
	f, ok := w.(*os.File)
	if !ok {
		return false
	}
	return term.IsTerminal(int(f.Fd()))
}
