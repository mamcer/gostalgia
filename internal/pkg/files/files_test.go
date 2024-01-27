package files

import (
	"testing"
)

var validateEquivalence = func(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got:%v, want:%v", got, want)
	}
}

func TestSizeString(t *testing.T) {
	unit := int64(1000)

	t.Run("0 size, 0.0 bytes", func(t *testing.T) {
		got := SizeString(0)
		want := "0.0 Bytes"

		validateEquivalence(t, got, want)
	})

	t.Run("1000.0 bytes 1 KB", func(t *testing.T) {
		got := SizeString(unit)
		want := "1.0 KB"

		validateEquivalence(t, got, want)
	})

	t.Run("1000.0 Kbytes 1 MB", func(t *testing.T) {
		got := SizeString(unit * unit)
		want := "1.0 MB"

		validateEquivalence(t, got, want)
	})

	t.Run("1000.0 Mbytes 1 GB", func(t *testing.T) {
		got := SizeString(unit * unit * unit)
		want := "1.0 GB"

		validateEquivalence(t, got, want)
	})

	t.Run("1000.0 Gbytes 1 TB", func(t *testing.T) {
		got := SizeString(unit * unit * unit * unit)
		want := "1.0 TB"

		validateEquivalence(t, got, want)
	})

	t.Run("4300 Gbytes 4.3 TB", func(t *testing.T) {
		got := SizeString(4300 * unit * unit * unit)
		want := "4.3 TB"

		validateEquivalence(t, got, want)
	})
}
