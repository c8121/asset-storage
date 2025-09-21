package util

import "fmt"

func Check(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}
