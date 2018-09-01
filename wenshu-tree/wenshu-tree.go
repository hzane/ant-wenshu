package main

import (
	"fmt"

	"github.com/otiai10/gosseract"
)

func main() {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetImage("./ValidateCode.jpeg")
	text, err := client.Text()
	fmt.Println(text, err)
}
