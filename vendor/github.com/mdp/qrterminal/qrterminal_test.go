package qrterminal

import (
	"os"
	"testing"
)

func TestGenerate(t *testing.T) {
	Generate("https://github.com/mdp/qrterminal", L, os.Stdout)
}

func TestGenerateWithConfig(t *testing.T) {
	config := Config{
		Level:     M,
		Writer:    os.Stdout,
		BlackChar: WHITE,
		WhiteChar: BLACK,
	}
	GenerateWithConfig("https://github.com/mdp/qrterminal", config)
}
