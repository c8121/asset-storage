package util

import (
	"fmt"
	"io"
)

func Check(e error, message string) {
	if e != nil {
		fmt.Println(message)
		fmt.Printf("Check/Panic: %s\n", e)
		panic(e)
	}
}

func LogError(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e)
	}
}

func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		fmt.Printf("Close error: %s\n", err)
	}
}
