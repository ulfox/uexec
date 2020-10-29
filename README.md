# uexec

Golang lib for handling function return codes

## Usage

    package main

    import (
        "encoding/json"
        "io/ioutil"

        "github.com/ulfox/uexec"
    )

    type Config struct {
        Version string
    }

    func main() {
        state := Config{
            Version: "Version",
        }

        try := uexec.NewErrorHandler()
        json.Unmarshal(
            try.Exec(ioutil.ReadFile("config.json")).ByteS(0), &state,
        )
    }
