# nomainreturn
A Go linter that checks for return statements in the main.

## Install
```bash
go get github.com/bedakb/nomainreturn
```

## Usage
```bash
nomainreturn ./...
```

## Example
The purpose of this linter is to prevent accidental exits with 0 status code by using return.

```go
package main

import (
    "fmt"

    log "github.com/sirupsen/logrus"
)

func main() {
    fmt.Println("running the program")

    if err := errnousFunc(); err != nil {
        log.WithError(err).Error("error occurred")
        return // exits with code 0
    }

    fmt.Println("done")
}
```