package tormenta

import (
	"encoding/json"
	"log"
	"math/rand"
	"strings"
	"time"
)

func RandomiseRecords(slice []Record) {
	for i := range slice {
		j := rand.Intn(i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
}

func MemberString(valid []string, target string) bool {
	for _, validOption := range valid {
		if target == validOption {
			return true
		}
	}
	return false
}

var nonContentWords = []string{"on", "at", "the", "in", "a"}

func removeNonContentWords(strings []string) (results []string) {
	for _, s := range strings {
		if !MemberString(nonContentWords, s) {
			results = append(results, s)
		}
	}

	return
}

func timerMiliseconds(t time.Time) int {
	t1 := time.Now()
	duration := t1.Sub(t)
	return int(duration.Seconds() * 1000)
}

func toJSON(m interface{}) string {
	js, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(string(js), ",", ", ", -1)
}
