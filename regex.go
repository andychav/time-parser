package main

import (
	"log"
	"regexp"
	"strconv"
)

func getMinAfterMidnight(t string) int {
	//TODO: Use capturing groups here
	minAfter := 0
	hPat := "(2[0-3]|[01]\\d|\\d){1}"
	mPat := "(:[0-5]\\d)"
	meridianPat := "(am|pm)"

	tPat := hPat + mPat + "?" + meridianPat + "?"

	if getMatch(tPat, t) == "" {
		log.Fatal(t, "does not follow the correct pattern")
	}

	hI, _ := strconv.Atoi(getMatch(hPat, t))
	if hI != 12 {
		minAfter += hI * 60
	}

	mI, _ := strconv.Atoi(getMatch(mPat, t)[1:])
	minAfter += mI

	if getMatch(meridianPat, t) == "pm" {
		minAfter += 12 * 60
	}

	return minAfter
}

func getMatch(p string, t string) string {
	r, _ := regexp.Compile(p)
	m := r.Find([]byte(t))
	if m == nil {
		return ""
	}
	return string(m)
}
