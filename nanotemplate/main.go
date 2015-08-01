//go:generate nanoinstall -M result
//go:generate nanotemplate -T string --input=result.go.t
package main

import (
	"flag"
	"fmt"
	"github.com/nanoservice/monad.go/nanotemplate/result_string"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

const (
	dirPermission  = 0775
	filePermission = 0664
)

var (
	typeName      = flag.String("T", "", "TYPE to substitute in the template")
	inputFilename = flag.String("input", "", "INPUT template filename")
	importName    = flag.String("I", "", "IMPORT string for importing a type")
)

func main() {
	flag.Parse()

	if *typeName == "" {
		fmt.Println("TYPE is not provided")
		flag.Usage()
		os.Exit(1)
	}

	if *inputFilename == "" {
		fmt.Println("INPUT template filename is not provided")
		flag.Usage()
		os.Exit(1)
	}

	if *importName != "" {
		*importName = "\"" + *importName + "\""
	}

	lowercaseTypeName := strings.ToLower(*typeName)
	packageName := "result_" + lowercaseTypeName
	packageFile := packageName + ".t.go"
	outputFile := path.Join(packageName, packageFile)

	os.Mkdir(packageName, dirPermission)

	readTemplate().
		Bind(replace("{{I}}", *importName)).
		Bind(replace("{{T}}", *typeName)).
		Bind(replace("{{t}}", lowercaseTypeName)).
		Bind(saveTo(outputFile)).
		OnErrorFn(reportGenerationError)
}

func reportGenerationError(err error) {
	fmt.Printf("Unable to generate from template: %v\n", err)
	os.Exit(1)
}

func readTemplate() result_string.Result {
	return result_string.Success("").Bind(_readTemplate)
}

func _readTemplate(_ string) result_string.Result {
	rawTemplate, err := ioutil.ReadFile(*inputFilename)
	return result_string.NewResult(string(rawTemplate), err)
}

func replace(target, value string) func(string) result_string.Result {
	return func(string body) result_string.Result {
		return result_string.Success(
			strings.Replace(body, target, value, -1),
		)
	}
}

func saveTo(outputFile string) func(string) result_string.Result {
	return func(string body) result_string.Result {
		err = ioutil.WriteFile(outputFile, []byte(result), filePermission)
		return result_string.NewResult("", err)
	}
}
