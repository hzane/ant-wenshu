package main

import (
	"flag"
	"fmt"

	"github.com/otiai10/gosseract"
)

var args struct {
	file string
	uri  string
}

func init() {
	flag.StringVar(&args.file, "file", "ValidateCode.jpeg", "image file")

	flag.Parse()
}
func main() {
	client := gosseract.NewClient()
	defer client.Close()
	client.SetWhitelist("0123456789")
	//client.SetImage("./ValidateCode.jpeg")

	client.SetImage(args.file)
	text, err := client.Text()
	fmt.Println(text, err)
}
