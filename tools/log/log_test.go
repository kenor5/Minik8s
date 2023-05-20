package log

import (
	"testing"
)

func TestLOG(t *testing.T) {
	Print("hello")
	PrintE("hello")
	PrintS("hello")
	PrintW("hello")

}
