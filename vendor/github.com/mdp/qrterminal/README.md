# QRCode Terminal

[![Build Status](https://secure.travis-ci.org/mdp/qrterminal.png)](https://travis-ci.org/mdp/qrterminal)

Pretty simple, I stole this from the NodeJS version at https://github.com/gtanner/qrcode-terminal and turned it into a Golang library.

## Install

`go get github.com/mdp/qrterminal`

## Usage

```go
import (
    "github.com/mdp/qrterminal"
    "os"
    )

func main() {
  qrterminal.Generate("https://github.com/mdp/qrterminal", os.Stdout)
}
```

### More complicated

Inverted barcode with medium redundancy
```go
import (
    "github.com/mdp/qrterminal"
    "os"
    )

func main() {
  config := qrterminal.Config{
      Level: qrterminal.L,
      Writer: os.Stdout,
      BlackChar: qrterminal.WHITE,
      WhiteChar: qrterminal.BLACK,
  }
  qrterminal.GenerateWithConfig("https://github.com/mdp/qrterminal", config)
}
```

