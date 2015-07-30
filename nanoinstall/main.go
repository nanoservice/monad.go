package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
)

const (
	rootUrl = "https://github.com/nanoservice/monad.go/raw/"
	ext     = ".go.t"
)

var (
	monad   = flag.String("M", "", "MONAD template to install")
	version = flag.String("v", "master", "VERSION of template")
)

func main() {
	flag.Parse()

	if *monad == "" {
		fmt.Println("MONAD was not provided")
		flag.Usage()
		os.Exit(1)
	}

	url := rootUrl + *version + "/" + *monad + "/" + *monad + ext

	out, err := os.Create(*monad + ext)
	if err != nil {
		reportError(err)
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		reportError(err)
	}
	if resp.StatusCode != 200 {
		reportError(
			errors.New("Expected status code to be 200, but got: " + resp.Status),
		)
	}
	defer resp.Body.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		reportError(err)
	}
}

func reportError(err error) {
	fmt.Printf("Unable to install template: %v\n", err)
	os.Exit(1)
}
