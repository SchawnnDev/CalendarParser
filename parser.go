package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
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
	Date    time.Time `json:"Date"`
	DayType dayType   `json:"day_type"`
}

func parse(data []byte) []day {
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))

	if err != nil {
		return nil
	}

	// Date regex
	r, err := regexp.Compile("^(0?[1-9]|[12][0-9]|3[01])-(0?[1-9]|1[012])-\\d{4}$")

	if err != nil {
		return nil
	}

	var result []day
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
			date, err := time.Parse(layout, id)

			// ignore this day
			if err != nil {
				//fmt.Printf("ERROR: Date %s could not be parsed.\n", id)
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

			result = append(result, day{Date: date, DayType: dType})
		})

	})

	return result
}
