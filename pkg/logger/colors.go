package logger

import (
	"fmt"
)

type Color uint8

const (
	ColorRed Color = iota + 31
	ColorGreen
	ColorYellow
	ColorBlue
	ColorPurple
	ColorCyan
	ColorGray
)

func (c Color) Fill(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(c), s)
}
