package util

import "fmt"

func Check(e error, message string) {
	if e != nil {
		fmt.Println(message)
		panic(e)
	}
}

func LogError(e error) {
	fmt.Printf("Error: %s\n", e)
}
