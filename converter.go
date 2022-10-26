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
	var columns = len(header)

	result := make([][]string, len(data)+1) // size + header
	result[0] = header
	j := 1

	for i := 1; i < len(result); i++ {
		r := make([]string, columns)
		d := data[i-1]

		add := true

		//fmt.Printf("Processing day %s: ", d.Date.Format("01/02/2006"))

		switch d.DayType {
		case ETU:
			r[7] = "Université de Strasbourg"
			break
		case WORK:
			r[7] = "Mind7 Consulting"
			break
		case HOLIDAY:
			day := d.Date.Day()
			month := d.Date.Month()

			vDay, vMonth := getVendrediSaint(d.Date)

			if (vDay == day && vMonth == month) || (day == 26 && month == time.December) {
				r[7] = "Mind7 Consulting"
			} else {
				add = false
			}
			break
		default:
			add = false
			break
		}

		if !add {
			//fmt.Printf("Day not added : %s\n", d.Date.Format("01/02/2006"))
			continue
		}

		dateLayout := "01/02/2006"
		// timeLayout := "03:04 PM"

		date := d.Date.Format(dateLayout)

		if r[7] == "Mind7 Consulting" {
			r[0] = "Présent"
		} else {
			r[0] = "Absent"
		}
		r[1] = date
		r[2] = "09:30 AM"
		r[3] = date
		r[4] = "06:00 PM"
		r[5] = "False"
		r[6] = ""
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
