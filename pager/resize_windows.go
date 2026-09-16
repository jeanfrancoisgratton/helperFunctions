//go:build windows

// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: pager/resize_windows.go

package pager

// watchResize is a no-op on Windows: there is no SIGWINCH equivalent
// exposed by the standard library, so resize while paging is picked up
// only on the next keystroke (via a fresh call to Page).
func watchResize() (<-chan struct{}, func()) {
	return make(chan struct{}), func() {}
}
