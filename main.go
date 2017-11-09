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
	http.HandleFunc("/pdd", savePddCall)
	log.Fatalf("ERROR!", http.ListenAndServe(":8080", nil))
}

const (
	fPath = "pdd.db"
)

func savePddCall(w http.ResponseWriter, r *http.Request) {
	timeStr := []byte(fmt.Sprintf("%v", time.Now().Format(time.RFC3339)))
	ioutil.WriteFile(fPath, timeStr, 0600)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created!"))
}

func checkLastCall() {
	fContent, _ := ioutil.ReadFile(fPath)
	t, err := time.Parse(time.RFC3339, string(fContent))
	if err == nil {
		if time.Since(t) > 30*time.Second {
			fmt.Println("Boum! Last time was: ", t)
		}
	}
	checkLastCall()
}
