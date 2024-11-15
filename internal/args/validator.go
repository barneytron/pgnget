package args

import "strconv"

func IsYearValid(year string) bool {
	if year == "all" {
		return true
	}
	_, err := strconv.Atoi(year)
	if err != nil {
		return false
	}

	if len(year) != 4 {
		return false
	}

	return true
}

func IsMonthValid(month string) bool {
	if month == "all" {
		return true
	}
	m, err := strconv.Atoi(month)
	if err != nil {
		return false
	}

	if len(month) != 2 {
		return false
	}

	return m >= 1 && m <= 12
}

func IsUsernameValid(username string) bool {
	// TODO: think of more validations
	return len(username) >= 1
}
