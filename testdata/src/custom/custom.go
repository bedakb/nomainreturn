package custom

import (
	"errors"
	"fmt"
)

func ok() bool {
	return true
}

func main() {
	fmt.Println("start")

	myFunc := func() {
		ok := true
		if !ok {
			return
		}
		fmt.Println("ok")
	}
	myFunc()

	if e := err(); e != nil {
		fmt.Println(err)
		return // want "return found in main"
	}

	fmt.Println("done")
}

func err() error {
	return errors.New("the error")
}
