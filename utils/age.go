package utils

import (
	"strconv"
	"time"
)

func GetAge(birthday time.Time) int {
	now := time.Now()
	age := now.Year() - birthday.Year()
	if now.YearDay() < birthday.YearDay() {
		age--
	}
	return age
}

func GetAgeGroup(age int) string {
	return strconv.Itoa(age/10*10) + "s"
}
