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
func CloseOrLog(c io.Closer) {
	err := c.Close()
	if err != nil {
		fmt.Printf("Close error: %s\n", err)
	}
}
