package main

import (
	"fmt"
	"ezzygo/pkg/auth"
)

func main() {
	fmt.Println(auth.GenerateRandomKey())
}
