package helpers

import "time"


func ComputeAge(bd time.Time) int {
	today := time.Now()
	age := today.Year() - bd.Year()

	if today.Month() < bd.Month() {
		age--
	}
	return age
}
