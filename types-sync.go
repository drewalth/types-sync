package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/urfave/cli"
)

func main() {
	var sourcePath = "src"
	var outputPath = "types.ts"
	var excludedTypes = "Fastify"

	app := &cli.App{
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "src",
				Value:       "src",
				Usage:       "Source folder to walk",
				Destination: &sourcePath,
			},
			&cli.StringFlag{
				Name:        "output",
				Value:       "types.ts",
				Usage:       "The destination output",
				Destination: &outputPath,
			},
			&cli.StringFlag{
				Name:        "excludedTypes",
				Value:       "Fastify",
				Usage:       "Exclude type declarations that contain value. Separated by comma.",
				Destination: &excludedTypes,
			},
		},
		Action: func(c *cli.Context) error {
			finalResults := WalkFiles(sourcePath, excludedTypes)

			WriteTsFile(finalResults, outputPath)

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func WalkFiles(sourcePath string, excludedTypesInput string) []string {

	results := make([]string, 0)

	err := filepath.Walk(sourcePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if IsTypeScriptFile(info.Name()) {
				fmt.Println("> ", path)
				fileResults, _ := readFileWithReadString(path, excludedTypesInput)
				results = append(results, fileResults...)
			}

			return nil
		})
	check(err)

	return results
}

func IsTypeScriptFile(fileName string) bool {

	bucket := strings.Split(fileName, ".")

	if len(bucket) == 2 {

		// could add TSX or "fileTypes" as a flag
		return bucket[1] == "ts"
	}

	return false

}

func HasExcludedType(line string, excludedTypes []string) bool {

	for _, t := range excludedTypes {
		if strings.Contains(line, t) {
			return true
		}
	}

	return false
}

func readFileWithReadString(fn string, excludedTypes string) (typesList []string, err error) {

	listOfTypes := make([]string, 0)
	listOfExcludedTypes := strings.Split(excludedTypes, ",")

	file, err := os.Open(fn)
	defer func() {
		closeErr := file.Close()
		check(closeErr)
	}()

	check(err)

	reader := bufio.NewReader(file)

	var line string
	var readingDeclaration = false
	var level = 1
	var tsType = ""

	for {
		line, err = reader.ReadString('\n')

		if len(line) == 0 {
			break
		}

		if HasExport(line) && !HasExcludedType(line, listOfExcludedTypes) {
			readingDeclaration = true
		}

		if readingDeclaration && strings.Contains(line, "{") && !(HasExport(line)) {
			level++
		}

		if readingDeclaration && strings.Contains(line, "}") && level >= 1 {
			level--
		}

		if readingDeclaration && strings.Contains(line, "}") && level == 0 {
			tsType += RemoveSemiColons(line)
			listOfTypes = append(listOfTypes, tsType)
			tsType = ""
			readingDeclaration = false
		}

		if readingDeclaration {
			tsType += RemoveSemiColons(line)
		}

		if err != nil {
			break
		}
	}

	if err != io.EOF {
		fmt.Printf("Failed!: %v\n", err)
	}

	return listOfTypes, err
}

func RemoveSemiColons(input string) string {
	return strings.Replace(input, ";", "", -1)
}

func WriteTsFile(typesArr []string, outputPath string) {

	var fileDataString = `
	
	/*/////////////////////////////////\
	\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\/

	This file was automatically generated
	by types-sync. Do not manually edit.

	https://github.com/drewalth/types-sync

	///////////////////////////////// \
	\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\\*/`

	for _, str := range typesArr {

		fileDataString += "\n\n" + str

	}

	fileData := []byte(fileDataString)

	writeErr := os.WriteFile(outputPath, fileData, 0644)
	check(writeErr)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func HasExport(line string) bool {
	exports := []string{"export type", "export enum", "export interface"}

	for _, t := range exports {
		if strings.Contains(line, t) {
			return true
		}
	}

	return false
}
