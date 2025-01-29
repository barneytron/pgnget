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

// minimal validation of username
func IsUsernameValid(username string) bool {
	// Taken from https://support.chess.com/en/articles/8557649-how-can-i-change-my-username
	// Your username must be at least 3 characters long.
	// It can only include letters, numbers, underscores, and dashes.
	// It must start and end with a letter or number. Spaces are not allowed.
	return len(username) >= 3
}
