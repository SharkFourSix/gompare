package main

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/devfacet/gocmd"
	"io"
	"os"
	"strings"
)

var (
	name       = "gompare"
	version    = "0.0.2"
	repository = "https://github.com/SharkFourSix/gompare"
)

const (
	OutputFileName            = "gompare-results.txt"
	ColumnHeaderFormatString  = "\n%s\n---------------------\n"
	StatisticsHeader          = "\nStatistics\n------------------------\n"
	StatisticsNumberRowFormat = "%-20s%d\n"
	StatisticsTextRowFormat   = "%-20s%s\n"
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

func StringArrayContains(needle string, array []string) bool {
	for _, str := range array {
		if strings.EqualFold(needle, str) {
			return true
		}
	}
	return false
}

func PrintColumns(title string, columns []string, writer *bufio.Writer) {
	_, _ = writer.WriteString(fmt.Sprintf(ColumnHeaderFormatString, title))
	if len(columns) == 0 {
		_, _ = writer.WriteString("N/A")
	} else {
		_, _ = writer.WriteString(strings.Join(columns, "\n"))
	}
	_, _ = writer.WriteString("\n")
}

func main() {
	flags := struct {
		Version       bool   `short:"v" long:"version" description:"Display version"`
		Help          bool   `short:"h" long:"help" description:"Show help"`
		Template      string `short:"t" long:"template-file" description:"Template file containing comma separated column names"`
		InputFile     string `short:"i" long:"input-file" description:"Target CSV file to be inspected"`
		OutputFile    bool   `short:"o" long:"output-to-file" description:"Prints results to file instead of standard output"`
		ShowUnmatched bool   `short:"u" long:"show-unmatched-cols" description:"Show unmatched columns"`
	}{}

	var (
		templateFile      string
		inputFile         string
		templateColumns   []string
		targetColumns     []string
		writer            *bufio.Writer
		writeToFile       = false
		outputFile        *os.File
		missingColumns    []string
		unmatchedColumns  []string
		matchedColumns    []string
		showUnmatchedCols = false
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

	_, _ = gocmd.HandleFlag("ShowUnmatched", func(cmd *gocmd.Cmd, args []string) error {
		showUnmatchedCols = true
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

	writeLine := func(format string, a ...interface{}) {
		_, _ = writer.WriteString(fmt.Sprintf(format, a...))
	}

	for _, templateColumn := range templateColumns {
		if StringArrayContains(templateColumn, targetColumns) {
			matchedColumns = append(matchedColumns, templateColumn)
		} else {
			missingColumns = append(missingColumns, templateColumn)
		}
	}

	// in target and not in template
	if showUnmatchedCols {
		for _, column := range targetColumns {
			if !StringArrayContains(column, templateColumns) {
				unmatchedColumns = append(unmatchedColumns, column)
			}
		}
	}

	PrintColumns("Matching Columns", matchedColumns, writer)
	PrintColumns("Missing Columns", missingColumns, writer)
	if showUnmatchedCols {
		PrintColumns("Unmatched Columns", unmatchedColumns, writer)
	}

	targetColumnCount := len(targetColumns)
	templateColumnCount := len(templateColumns)
	matchedColumnCount := len(matchedColumns)
	unmatchedColumnCount := len(unmatchedColumns)
	missingColumnCount := len(missingColumns)

	status := "Invalid"
	if templateColumnCount == matchedColumnCount {
		status = "Valid"
	}

	writeLine(StatisticsHeader)
	writeLine(StatisticsNumberRowFormat, "Template Columns", templateColumnCount)
	writeLine(StatisticsNumberRowFormat, "Target Columns", targetColumnCount)
	writeLine(StatisticsNumberRowFormat, "Missing Columns", missingColumnCount)
	writeLine(StatisticsNumberRowFormat, "Matching Columns", matchedColumnCount)
	writeLine(StatisticsNumberRowFormat, "Unmatched Columns", unmatchedColumnCount)
	writeLine(StatisticsTextRowFormat, "Status", status)

	_ = writer.Flush()
	if writeToFile {
		_ = outputFile.Close()
	}
}
