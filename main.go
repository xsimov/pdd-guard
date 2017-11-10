package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

func main() {
	go checkLastCall()
	http.HandleFunc("/pdd", writeTimestampToDisk)
	log.Fatalf("error: %v", http.ListenAndServe(":8080", nil))
}

const (
	fPath      = "pdd.db"
	timeLayout = time.RFC3339
)

func writeTimestampToDisk(w http.ResponseWriter, r *http.Request) {
	timeStr := []byte(fmt.Sprintf("%v", time.Now().Format(timeLayout)))
	ioutil.WriteFile(fPath, timeStr, 0600)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created!"))
}

func checkLastCall() {
	fContent, err := ioutil.ReadFile(fPath)
	if err == nil {
		parseTimestamp(fContent)
	}
	time.Sleep(5 * time.Second)
	checkLastCall()
}

func parseTimestamp(content []byte) {
	t, err := time.Parse(timeLayout, string(content))
	if err != nil {
		return
	}
	if time.Since(t) > 30*time.Second {
		fmt.Println("It's been a long time since last call! 30 seconds!", t)
	}
}
