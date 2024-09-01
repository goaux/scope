package scope_test

import (
	"bytes"
	"compress/gzip"
	"encoding/hex"
	"fmt"
	"io"
	"os"

	"github.com/goaux/scope"
)

func ExampleUse() {
	buf := &bytes.Buffer{}
	for gz := range scope.Use(gzip.NewWriter(buf)) {
		fmt.Fprintln(gz, "hello world")
	}
	fmt.Println(hex.EncodeToString(buf.Bytes()))
	// Output:
	// 1f8b08000000000000ffca48cdc9c95728cf2fca49e102040000ffff2d3b08af0c000000
}

func ExampleUse2() {
	defer os.Remove("test.gz")
	for file, err := range scope.Use2(os.Create("test.gz")) {
		if err != nil {
			fmt.Println(err)
			break
		}
		for gz := range scope.Use(gzip.NewWriter(file)) {
			if _, err := fmt.Fprintln(gz, "hello world"); err != nil {
				fmt.Println(err)
				break
			}
			// gz will be closed at the end of the loop body.
		}
		// file will be closed at the end of the loop body.
	}

	for file, err := range scope.Use2(os.Open("test.gz")) {
		if err != nil {
			fmt.Println(err)
			break
		}
		if gz, err := gzip.NewReader(file); err != nil {
			fmt.Println(err)
			break
		} else if b, err := io.ReadAll(gz); err != nil {
			fmt.Println(err)
		} else {
			fmt.Print(string(b))
		}
		// file will be closed at the end of the loop body.
	}
	// Output:
	// hello world
}
