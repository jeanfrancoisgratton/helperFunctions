//go:build !windows

// helperFunctions
// Written by J.F. Gratton <jean-francois@famillegratton.net>
// Original filename: pager/resize_unix.go

package pager

import (
	"os"
	"os/signal"
	"syscall"
)

// watchResize reports terminal resize events (SIGWINCH) on the returned
// channel. The stop function must be called to release the signal
// subscription and the backing goroutine.
func watchResize() (<-chan struct{}, func()) {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGWINCH)

	out := make(chan struct{}, 1)
	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-sigCh:
				select {
				case out <- struct{}{}:
				default: // a redraw is already pending
				}
			case <-done:
				return
			}
		}
	}()

	return out, func() {
		signal.Stop(sigCh)
		close(done)
	}
}
