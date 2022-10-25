package main

import (
	"bytes"
	"encoding/csv"
	"time"
)

type formatFn func([]day) [][]string

// (day, month)
func getVendrediSaint(d time.Time) (int, time.Month) {
	switch d.Year() {
	case 2022:
		return 15, time.April
	case 2023:
		return 7, time.April
	case 2024:
		return 29, time.March
	case 2025:
		return 18, time.April
	case 2026:
		return 3, time.April
	default:
		return -1, 0
	}
}

// TODO: day => add subject etc
func formatGoogleCalendar(data []day) [][]string {
	header := []string{
		"Subject", "Start Date", "Start Time", "End Date", "End Time",
		"All Day Event", "Description", "Location", "Private",
	}
	var columns = 0

	for _, d := range data {
		if d.DayType == OTHER {
			continue
		}
		columns += 1
	}

	result := make([][]string, len(data)+1) // size + header
	result[0] = header
	j := 0

	for i := 1; i < len(result); i++ {
		r := make([]string, columns)
		d := data[j]

		add := true

		switch d.DayType {
		case ETU:
			r[7] = "Mind7"
			break
		case WORK:
			r[7] = "Université"
			break
		case HOLIDAY:
			day := d.Date.Day()
			month := d.Date.Month()

			vDay, vMonth := getVendrediSaint(d.Date)

			if (vDay == day && vMonth == month) || (day == 26 && month == time.December) {
				r[7] = "Mind7"
			} else {
				add = false
			}

		}

		if !add {
			continue
		}

		dateLayout := "01/02/2006"
		timeLayout := "03:04 PM"

		date := d.Date.Format(dateLayout)

		r[0] = "Subject"
		r[1] = date
		r[2] = d.Date.Format(timeLayout)
		r[3] = date
		r[4] = d.Date.Format(timeLayout)
		r[5] = "False"
		r[6] = "This is a description"
		r[8] = "False"

		result[j] = r
		j++
	}

	return result
}

func convertToCsv(data []day, f formatFn) []byte {
	var result []byte

	buf := bytes.NewBuffer(result)
	w := csv.NewWriter(buf)
	converted := f(data)

	for _, d := range converted {
		if err := w.Write(d); err != nil {
			// process errors
		}
	}

	w.Flush()

	if err := w.Error(); err != nil {
		// process errors
	}

	return buf.Bytes()
}
