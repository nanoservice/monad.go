//go:generate nanoinstall -M result -v v1.2.1
//go:generate nanotemplate -T *os.File -t file -I os --input=result.go.t
//go:generate nanotemplate -T *http.Response -t response -I net/http --input=result.go.t
package main

import (
	"errors"
	"flag"
	"fmt"
	"github.com/nanoservice/monad.go/nanoinstall/result_file"
	"github.com/nanoservice/monad.go/nanoinstall/result_response"
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

	openOutputFile().Chain(
		writeTemplateToFileFrom(url),
		printSuccess,
	).
		OnErrorFn(reportError).
		Err()
}

func reportError(err error) {
	fmt.Printf("Unable to install template: %v\n", err)
	os.Exit(1)
}

func openOutputFile() result_file.Result {
	out, err := os.Create(*monad + ext)
	return result_file.
		NewResult(out, err).
		Defer(closeOutputFile)
}

func closeOutputFile(out *os.File) {
	out.Close()
}

func printSuccess(out *os.File) result_file.Result {
	fmt.Println("Successfully installed template.")
	return result_file.Success(out)
}

func fetchTemplate(url string) result_response.Result {
	resp, err := http.Get(url)
	return result_response.
		NewResult(resp, err).
		Defer(closeResponseBody)
}

func closeResponseBody(resp *http.Response) {
	resp.Body.Close()
}

func validateResponse(resp *http.Response) result_response.Result {
	if resp.StatusCode != 200 {
		return result_response.Failure(
			errors.New("Expected status code to be 200, but got: " + resp.Status),
		)
	}

	return result_response.Success(resp)
}

func copyResponseBodyToOutput(out *os.File) func(resp *http.Response) result_response.Result {
	return func(resp *http.Response) result_response.Result {
		_, err := io.Copy(out, resp.Body)
		return result_response.NewResult(resp, err)
	}
}

func writeTemplateToFileFrom(url string) func(out *os.File) result_file.Result {
	return func(out *os.File) result_file.Result {
		return result_file.NewResult(
			out,
			fetchTemplate(url).Chain(
				validateResponse,
				copyResponseBodyToOutput(out),
			).Err(),
		)
	}
}
