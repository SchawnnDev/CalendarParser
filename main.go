package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	path := flag.String("path", "", "HTML file containing the calendar")
	port := flag.Int("port", 8000, "Port of webserver")
	in := flag.Bool("in", false, "Using std input stream.")
	fileType := flag.String("fileType", "csv", "Export file type")
	fileName := flag.String("fileName", "export.csv", "Export file name")

	flag.Parse()

	// Temp only csv is valid
	if *fileType != "csv" {
		*fileType = "csv"
	}

	// If path is empty, check stdin
	if *path != "" || *in {
		var content []byte
		var err error

		if *in {
			fmt.Println("Reading from Stdin. Make sure you specified an input stream.")
			content, err = io.ReadAll(os.Stdin)
		} else {
			content, err = os.ReadFile(*path)
		}

		check(err)

		if len(content) > 0 {
			fmt.Println("Content successfully loaded.")
		} else {
			panic("No input given")
		}

		data := parse(content)

		fmt.Printf("Successfully parsed %d Date(s).\n", len(data))

		csv := convertToCsv(data, formatGoogleCalendar)

		err = os.WriteFile(*fileName, csv, 0644)
		
		if err != nil {
			fmt.Printf("Could not write to file %s\n", *fileName)
		}

		return
	}

	startHttpServer(":" + strconv.Itoa(*port))
}
