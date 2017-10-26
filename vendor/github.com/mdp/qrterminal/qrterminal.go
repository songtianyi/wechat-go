package qrterminal

import (
	"io"

	"github.com/mdp/rsc/qr"
)

const BLACK = "\033[40m  \033[0m"
const WHITE = "\033[47m  \033[0m"

// Level - the QR Code's redundancy level
const H = qr.H
const M = qr.M
const L = qr.L

//Config for generating a barcode
type Config struct {
	Level     qr.Level
	Writer    io.Writer
	BlackChar string
	WhiteChar string
}

// GenerateWithConfig expects a string to encode and a config
func GenerateWithConfig(text string, config Config) {
	w := config.Writer
	white := config.WhiteChar
	black := config.BlackChar
	code, _ := qr.Encode(text, config.Level)
	// Frame the barcode in a 1 pixel border
	w.Write([]byte(white))
	for i := 0; i <= code.Size; i++ {
		w.Write([]byte(white))
	}
	w.Write([]byte("\n"))
	for i := 0; i <= code.Size; i++ {
		w.Write([]byte(white))
		for j := 0; j <= code.Size; j++ {
			if code.Black(i, j) {
				w.Write([]byte(black))
			} else {
				w.Write([]byte(white))
			}
		}
		w.Write([]byte("\n"))
	}
}

// Generate a QR Code and write it out to io.Writer
func Generate(text string, l qr.Level, w io.Writer) {
	config := Config{
		Level:     qr.L,
		Writer:    w,
		BlackChar: BLACK,
		WhiteChar: WHITE,
	}
	GenerateWithConfig(text, config)
}
