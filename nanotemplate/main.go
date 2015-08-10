//go:generate nanoinstall -M result -v v1.3.0
//go:generate nanotemplate -T string --input=_result.tt.go
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
	typeName          = flag.String("T", "", "TYPE to substitute in the template")
	lowercaseTypeName = flag.String("t", "", "LOWERCASETYPE to substitute in package name")
	inputFilename     = flag.String("input", "", "INPUT template filename")
	importName        = flag.String("I", "", "IMPORT string for importing a type")
)

func main() {
	flag.Parse()

	if *typeName == "" {
		reportIsNotProvided("TYPE")
	}

	if *inputFilename == "" {
		reportIsNotProvided("INPUT")
	}

	if *importName != "" {
		*importName = "\"" + *importName + "\""
	}

	if *lowercaseTypeName == "" {
		*lowercaseTypeName = strings.ToLower(*typeName)
	}

	packageName := "result_" + *lowercaseTypeName
	packageFile := packageName + ".t.go"
	outputFile := path.Join(packageName, packageFile)

	os.Mkdir(packageName, dirPermission)

	readTemplate().Chain(
		replace("{{I}}", *importName),
		replace("{{T}}", *typeName),
		replace("{{t}}", *lowercaseTypeName),
		saveTo(outputFile),
	).OnErrorFn(reportGenerationError)
}

func reportIsNotProvided(what string) {
	fmt.Printf("%s is not provided\n", what)
	flag.Usage()
	os.Exit(1)
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
	return func(body string) result_string.Result {
		return result_string.Success(
			strings.Replace(body, target, value, -1),
		)
	}
}

func saveTo(outputFile string) func(string) result_string.Result {
	return func(content string) result_string.Result {
		err := ioutil.WriteFile(outputFile, []byte(content), filePermission)
		return result_string.NewResult("", err)
	}
}
