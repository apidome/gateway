package main

import "fmt"

func main() {
	fmt.Println("Hello world!")

	var msg string

	_, err := fmt.Scan(&msg)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(msg)
	}
}
