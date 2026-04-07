package task

import (
	"time"
)

const defaultRecurrenceDuration = 365 * 24 * time.Hour // 1 year by default

// GenerateRecurrenceDates returns a slice of generated dates based on the recurrence rules.
// If startDate is nil, it uses current time.
func GenerateRecurrenceDates(startDate *time.Time, req *RecurrenceInput) []time.Time {
	if req == nil {
		return nil
	}

	start := time.Now().UTC()
	if startDate != nil {
		start = *startDate
	}

	end := start.Add(defaultRecurrenceDuration)
	if req.EndDate != nil && req.EndDate.After(start) {
		end = *req.EndDate
	}

	var dates []time.Time

	switch req.Type {
	case RecurrenceDaily:
		interval := req.Interval
		if interval <= 0 {
			interval = 1 // default to 1 day
		}
		
		current := start
		for !current.After(end) {
			dates = append(dates, current)
			current = current.AddDate(0, 0, interval)
		}

	case RecurrenceMonthly:
		current := start
		for !current.After(end) {
			// Check if current day is in month dates
			for _, md := range req.MonthDates {
				if current.Day() == md {
					dates = append(dates, current)
					break
				}
			}
			current = current.AddDate(0, 0, 1)
		}

	case RecurrenceSpecificDates:
		for _, ds := range req.SpecificDates {
			if t, err := time.Parse(time.RFC3339, ds); err == nil {
				dates = append(dates, t)
			} else if t, err := time.Parse("2006-01-02", ds); err == nil {
				dates = append(dates, t)
			}
			// if parsing fails we continue, but in reality we might want to validate this earlier
		}

	case RecurrenceEvenOdd:
		current := start
		for !current.After(end) {
			day := current.Day()
			isEven := day%2 == 0
			
			if (req.Parity == ParityEven && isEven) || (req.Parity == ParityOdd && !isEven) {
				dates = append(dates, current)
			}
			current = current.AddDate(0, 0, 1)
		}
	}

	return dates
}
