package main

import (
	"errors"
	"fmt"
)

func ok() bool {
	return true
}

func main() {
	fmt.Println("start")

	if e := err(); e != nil {
		fmt.Println(err)
		return // want "return found in main"
	}

	fmt.Println("done")
}

func err() error {
	return errors.New("the error")
}
