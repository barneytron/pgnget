package args_test

import (
	"pgnget/internal/args"
	"testing"
)

func TestIsMonthValid(t *testing.T) {
	var month string

	month = "1x"
	if args.IsMonthValid(month) {
		t.Errorf("%v should be invalid", month)
	}
	month = "1"
	if args.IsMonthValid(month) {
		t.Errorf("%v should be invalid", month)
	}
	month = "100"
	if args.IsMonthValid(month) {
		t.Errorf("%v should be invalid", month)
	}
	month = "13"
	if args.IsMonthValid(month) {
		t.Errorf("%v should be invalid", month)
	}

	month = "all"
	if !args.IsYearValid(month) {
		t.Errorf("%v should be valid", month)
	}
	month = "01"
	if !args.IsMonthValid(month) {
		t.Errorf("%v should be valid", month)
	}
}

func TestIsYeahValid(t *testing.T) {
	var year string
	year = "199"

	year = "202x"
	if args.IsYearValid(year) {
		t.Errorf("%v should be invalid", year)
	}
	if args.IsYearValid(year) {
		t.Errorf("%v should be invalid", year)
	}
	year = "20024"
	if args.IsYearValid(year) {
		t.Errorf("%v should be invalid", year)
	}

	year = "all"
	if !args.IsYearValid(year) {
		t.Errorf("%v should be valid", year)
	}
	year = "2024"
	if !args.IsYearValid(year) {
		t.Errorf("%v should be valid", year)
	}
}

func TestIsUsernameValid(t *testing.T) {
	var username string

	username = ""
	if args.IsUsernameValid(username) {
		t.Errorf("%v should be invalid", username)
	}
	username = "birdmaster3000"
	if !args.IsUsernameValid(username) {
		t.Errorf("%v should be valid", username)
	}
}
