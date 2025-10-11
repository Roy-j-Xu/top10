package tests

import "time"

func wait10msAnd(f func()) {
	time.Sleep(10 * time.Millisecond)
	f()
}
