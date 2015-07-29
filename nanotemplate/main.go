package main

import (
	"flag"
	"fmt"
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

	os.Mkdir(packageName, dirPermission)

	rawTemplate, err := ioutil.ReadFile(*inputFilename)
	if err != nil {
		reportGenerationError(err)
	}

	result := string(rawTemplate)
	result = strings.Replace(result, "{{I}}", *importName, -1)
	result = strings.Replace(result, "{{T}}", *typeName, -1)
	result = strings.Replace(result, "{{t}}", lowercaseTypeName, -1)

	err = ioutil.WriteFile(
		path.Join(packageName, packageFile),
		[]byte(result),
		filePermission,
	)
	if err != nil {
		reportGenerationError(err)
	}
}

func reportGenerationError(err error) {
	fmt.Printf("Unable to generate from template: %v\n", err)
	os.Exit(1)
}
