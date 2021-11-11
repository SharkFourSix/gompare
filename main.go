package main

import (
	"bufio"
	"encoding/csv"
	_ "encoding/csv"
	"errors"
	"fmt"
	"github.com/devfacet/gocmd"
	"io"
	"os"
	"strings"
)

var (
	name       = "gompare"
	version    = "0.0.1"
	repository = "https://github.com/SharkFourSix/gompare"
)

const (
	OutputFileName                  = "gompare-results.txt"
	HeaderFormatString              = "%s\n---------------------------\n"
	MissingColumnHeaderFormatString = "\n%s\n---------------------\n"
)

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

func IsEmptyString(str string) bool {
	return len(str) == 0
}

// ReadCsvColumns Read columns from a csv file
func ReadCsvColumns(filename string) (columns []string, err error) {
	var (
		handle *os.File
		reader *csv.Reader
	)
	handle, err = os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer handle.Close()
	reader = csv.NewReader(handle)
	columns, err = reader.Read()
	if err != nil {
		if err == io.EOF {
			if len(columns) == 0 {
				err = fmt.Errorf("read empty columns from file: '%s'", filename)
			}
		}
	}
	return columns, err
}

func main() {
	flags := struct {
		Version    bool   `short:"v" long:"version" description:"Display version"`
		Help       bool   `short:"h" long:"help" description:"Show help"`
		Template   string `short:"t" long:"template-file" description:"Template file containing comma separated column names"`
		InputFile  string `short:"i" long:"input-file" description:"Target CSV file to be inspected"`
		OutputFile bool   `short:"o" long:"output-to-file" description:"Prints results to file instead of standard output"`
	}{}

	var (
		templateFile    string
		inputFile       string
		templateColumns []string
		targetColumns   []string
		writer          *bufio.Writer
		writeToFile     = false
		outputFile      *os.File
		missingColumns  []string
	)
	_, _ = gocmd.HandleFlag("Help", func(cmd *gocmd.Cmd, args []string) error {
		cmd.PrintUsage()
		return nil
	})

	_, _ = gocmd.HandleFlag("Template", func(cmd *gocmd.Cmd, args []string) error {
		templateFile = args[0]
		if !fileExists(templateFile) {
			return fmt.Errorf("template file does not exist: '%s'", templateFile)
		}
		return nil
	})

	_, _ = gocmd.HandleFlag("InputFile", func(cmd *gocmd.Cmd, args []string) error {
		inputFile = args[0]
		if !fileExists(templateFile) {
			return fmt.Errorf("target csv file does not exist: '%s'", inputFile)
		}
		return nil
	})

	_, _ = gocmd.HandleFlag("OutputFile", func(cmd *gocmd.Cmd, args []string) error {
		writeToFile = true
		return nil
	})

	_, _ = gocmd.New(gocmd.Options{
		Name:        name,
		Description: "Compare and validate CSV files",
		Version:     fmt.Sprintf("%s v%s | %s", name, version, repository),
		Flags:       &flags,
		ConfigType:  gocmd.ConfigTypeAuto,
		AutoHelp:    true,
	})

	if IsEmptyString(inputFile) {
		fmt.Println(errors.New("missing parameter '--input-file'. Use '--help' for more information"))
		return
	}

	if IsEmptyString(templateFile) {
		fmt.Println(errors.New("missing parameter '--template-file'. Use '--help' for more information"))
		return
	}

	templateColumns, err := ReadCsvColumns(templateFile)
	if err != nil {
		fmt.Println(fmt.Errorf("error reading template file: %s", err))
		return
	}

	targetColumns, err = ReadCsvColumns(inputFile)
	if err != nil {
		fmt.Println(fmt.Errorf("error reading target file: %s", err))
		return
	}

	if writeToFile {
		outputFile, err = os.OpenFile(OutputFileName, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
		if err != nil {
			fmt.Println(errors.New("error creating output file"))
			return
		}
		writer = bufio.NewWriter(outputFile)
	} else {
		writer = bufio.NewWriter(os.Stdout)
	}

	targetColumnCount := len(targetColumns)
	templateColumnCount := len(templateColumns)

	writeLine := func(format string, a ...interface{}) {
		_, _ = writer.WriteString(fmt.Sprintf(format, a...))
	}

	writeLine(HeaderFormatString, "Matching Columns")

	for _, templateColumn := range templateColumns {
		inTargetCsv := false
		for _, targetColumn := range targetColumns {
			inTargetCsv = strings.EqualFold(strings.ToLower(templateColumn), strings.ToLower(targetColumn))
			if inTargetCsv {
				writeLine("%s\n", templateColumn)
				break
			}
		}
		if !inTargetCsv {
			missingColumns = append(missingColumns, templateColumn)
		}
	}

	writeLine(MissingColumnHeaderFormatString, "Missing Columns")
	if len(missingColumns) == 0 {
		writeLine("N/A\n")
	} else {
		for _, column := range missingColumns {
			writeLine("%s\n", column)
		}
	}

	writeLine("\n")
	writeLine("Column Count (template/target/missing): %d/%d/%d\n",
		templateColumnCount, targetColumnCount, len(missingColumns))

	_ = writer.Flush()
	if writeToFile {
		_ = outputFile.Close()
	}
}
