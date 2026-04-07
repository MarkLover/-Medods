package task

import (
	"reflect"
	"testing"
	"time"
)

func TestGenerateRecurrenceDates(t *testing.T) {
	parseTime := func(s string) time.Time {
		tm, err := time.Parse(time.RFC3339, s)
		if err != nil {
			t.Fatal(err)
		}
		return tm
	}

	start := parseTime("2023-10-25T10:00:00Z")
	end := parseTime("2023-11-05T10:00:00Z")

	startDatePtr := &start

	tests := []struct {
		name     string
		req      *RecurrenceInput
		expected []string
	}{
		{
			name: "daily, interval 2",
			req: &RecurrenceInput{
				Type:     RecurrenceDaily,
				Interval: 2,
				EndDate:  &end,
			},
			expected: []string{
				"2023-10-25T10:00:00Z",
				"2023-10-27T10:00:00Z",
				"2023-10-29T10:00:00Z",
				"2023-10-31T10:00:00Z",
				"2023-11-02T10:00:00Z",
				"2023-11-04T10:00:00Z",
			},
		},
		{
			name: "monthly dates",
			req: &RecurrenceInput{
				Type:       RecurrenceMonthly,
				MonthDates: []int{28, 1},
				EndDate:    &end,
			},
			expected: []string{
				"2023-10-28T10:00:00Z",
				"2023-11-01T10:00:00Z",
			},
		},
		{
			name: "specific dates",
			req: &RecurrenceInput{
				Type:          RecurrenceSpecificDates,
				SpecificDates: []string{"2023-11-01T10:00:00Z", "2023-11-03T10:00:00Z"},
				EndDate:       &end,
			},
			expected: []string{
				"2023-11-01T10:00:00Z",
				"2023-11-03T10:00:00Z",
			},
		},
		{
			name: "even days",
			req: &RecurrenceInput{
				Type:    RecurrenceEvenOdd,
				Parity:  ParityEven,
				EndDate: &end,
			},
			expected: []string{
				"2023-10-26T10:00:00Z",
				"2023-10-28T10:00:00Z",
				"2023-10-30T10:00:00Z",
				"2023-11-02T10:00:00Z",
				"2023-11-04T10:00:00Z",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GenerateRecurrenceDates(startDatePtr, tt.req)
			
			var gotStrings []string
			for _, tm := range got {
				gotStrings = append(gotStrings, tm.Format(time.RFC3339))
			}

			if !reflect.DeepEqual(gotStrings, tt.expected) {
				t.Errorf("GenerateRecurrenceDates() = %v, want %v", gotStrings, tt.expected)
			}
		})
	}
}
