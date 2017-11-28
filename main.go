package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

const (
	fPath        = "pdd.db"
	tLayout      = time.RFC3339
	timeToNotify = 30 * time.Minute
)

var (
	mailgunApiKey       string
	toEmails, fromEmail []string
	notifiedAt          time.Time
)

func main() {
	setEnvVars()
	go checkLastCall()
	http.HandleFunc("/pdd", writeTimestampToDisk)
	log.Fatalf("error: %v", http.ListenAndServe(":8080", nil))
}

func writeTimestampToDisk(w http.ResponseWriter, r *http.Request) {
	t := time.Now().Format(tLayout)
	tStr := fmt.Sprintf("%v", t)
	ioutil.WriteFile(fPath, []byte(tStr), 0600)
	w.WriteHeader(http.StatusCreated)
}

func checkLastCall() {
	fContent, err := ioutil.ReadFile(fPath)
	if err == nil {
		t, _ := time.Parse(tLayout, string(fContent))
		if time.Since(t) > timeToNotify && notNotifiedToday() {
			sendMail(t)
			markTodayAsNotified()
		}
	}
	time.Sleep(5 * time.Second)
	checkLastCall()
}

func sendMail(t time.Time) {
	formData := make(url.Values)
	formData["from"] = fromEmail
	formData["to"] = toEmails
	formData["subject"] = []string{"**NO** tinc electricitat!"}
	formData["text"] = []string{fmt.Sprintf("Última trucada registrada: %v", t)}

	url := fmt.Sprintf("https://api:%v@api.mailgun.net/v3/xsimov.com/messages", mailgunApiKey)
	resp, _ := http.PostForm(url, formData)
	fmt.Println(resp.Status)
}

func markTodayAsNotified() {
	d, err := getDataFromDisk()
	if err != nil {
		log.Fatalf("could not mark today as not notified: %v", err)
	}
	d.NotifiedToday = fmt.Sprintf("%v", time.Now().Format(tLayout))
	err = writeToDisk(d)
	if err != nil {
		log.Fatalf("could not mark today as not notified: %v", err)
	}
}

func notNotifiedToday() bool {
	d, err := getDataFromDisk()
	if err != nil {
		log.Fatal(err)
	}
	lastNotifiedAt, err := time.Parse(tLayout, d.NotifiedToday)
	if err != nil {
		log.Fatal(err)
	}
	dayStart, dayEnd := dayBoundaries()
	return lastNotifiedAt.After(dayStart) && lastNotifiedAt.Before(dayEnd)
}

func dayBoundaries() (start time.Time, end time.Time) {
	t := time.Now()
	start, _ = time.Parse(tLayout, fmt.Sprintf("%d%d%dT00:00:00Z", t.Year(), t.Month(), t.Day()))
	end = start.Add(24 * time.Hour)
	return
}

func setEnvVars() {
	toEmails = []string{os.Getenv("TO_EMAIL_1"), os.Getenv("TO_EMAIL_2")}
	fromEmail = []string{os.Getenv("FROM_EMAIL")}
	mailgunApiKey = os.Getenv("MAILGUN_API_KEY")
}
