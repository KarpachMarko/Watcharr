package util

import "time"

// Returns age in years.
func GetAge(birthDate time.Time, deathDate time.Time) int {
	var endDate time.Time
	if !deathDate.IsZero() {
		endDate = deathDate
	} else {
		endDate = time.Now()
	}

	years := endDate.Year() - birthDate.Year()

	// Adjust years if birthday hasn't occurred yet
	md := endDate.Month() - birthDate.Month()
	if md < 0 || (md == 0 && endDate.Day() < birthDate.Day()) {
		years--
	}

	return years
}
