package main

import (
	"io"
	"os"

	"github.com/rmaylor-base/transform/pkg/primitive"
)

func main() {
	inFile, err := os.Open("31866399.png")
	if err != nil {
		panic(err)
	}
	defer inFile.Close()

	out, err := primitive.Transform(inFile, 100)
	if err != nil {
		panic(err)
	}

	os.Remove("out.png")
	outFile, err := os.Create("out.png")
	if err != nil {
		panic(err)
	}
	io.Copy(outFile, out)
}
