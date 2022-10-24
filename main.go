package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"os"
	"regexp"
	"strings"
	"time"
)

type dayType int

const (
	ETU     dayType = iota
	WORK    dayType = iota
	HOLIDAY dayType = iota
	OTHER   dayType = iota
)

type day struct {
	date  time.Time
	dType dayType
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	path := flag.String("path", "", "HTML file containing the calendar")

	var content []byte
	var err error

	// If path is empty, check stdin
	if *path == "" {
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

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))

	check(err)

	// Date regex
	r, _ := regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-\\d{4}$")
	var data []day
	layout := "2-1-2006"

	// Get all months
	doc.Find(".month").Each(func(i int, month *goquery.Selection) {

		// For each month get all dates
		month.Find("td").FilterFunction(func(i int, selection *goquery.Selection) bool {
			id, exists := selection.Attr("id")
			return exists && r.MatchString(id)
		}).Each(func(i int, s *goquery.Selection) {
			id, _ := s.Attr("id")
			class, _ := s.Attr("class")
			fmt.Printf("'%s'\n", id)
			date, err := time.Parse(layout, id)
			fmt.Printf("Day: %d Month: %d Year: %d\n", date.Day(), date.Month(), date.Year())

			// ignore this day
			if err != nil {
				fmt.Printf("ERROR: Date %s could not be parsed.\n", id)
				return
			}

			var dType dayType

			class = strings.ReplaceAll(class, "\t", "")

			switch class {
			case "tdDisabled":
				dType = HOLIDAY
				break
			case "en":
				dType = WORK
				break
			case "fo":
				dType = ETU
				break
			default:
				dType = OTHER
				break
			}

			data = append(data, day{date: date, dType: dType})
		})

	})

}
