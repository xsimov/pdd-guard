package main

import (
	"encoding/json"
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
	timeLayout   = time.RFC3339
	timeToNotify = 30 * time.Minute
)

type Data struct {
	LastCall, NotifiedToday string
}

var notifiedAt time.Time

func main() {
	go checkLastCall()
	http.HandleFunc("/pdd", writeTimestampToDisk)
	log.Fatalf("error: %v", http.ListenAndServe(":8080", nil))
}

func writeTimestampToDisk(w http.ResponseWriter, r *http.Request) {
	timeStr := fmt.Sprintf("%v", time.Now().Format(timeLayout))
	ioutil.WriteFile(fPath, timeStr, 0600)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Created!"))
}

func checkLastCall() {
	fContent, err := ioutil.ReadFile(fPath)
	if err == nil {
		t, _ := time.Parse(timeLayout, string(fContent))
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
	formData["from"] = []string{"pdd@xsimov.com"}
	formData["to"] = []string{os.Getenv("EMAIL1"), os.Getenv("EMAIL2")}
	formData["subject"] = []string{"**NO** tinc electricitat!"}
	formData["text"] = []string{fmt.Sprintf("L'última trucada registrada fou: %v", t)}

	url := fmt.Sprintf("https://api:%v@api.mailgun.net/v3/xsimov.com/messages", os.Getenv("MAILGUN_API_KEY"))
	resp, _ := http.PostForm(url, formData)
	fmt.Println(resp.Status)
}

func markTodayAsNotified() {
	db, err := getDataFromDisk()
	if err != nil {
		log.Fatalf("could not mark today as not notified: %v", err)
	}
	db.NotifiedToday = fmt.Sprintf("%v", time.Now().Format(timeLayout))
	err := writeToDisk(db)
	if err != nil {
		log.Fatalf("could not mark today as not notified: %v", err)
	}
}

func notNotifiedToday() bool {
	d, err := getDataFromDisk()
	if err != nil {
		log.Fatalf(err)
	}
	lastNotifiedAt, err := time.Parse(timeLayout, d.NotifiedToday)
	if err != nil {
		log.Fatalf(err)
	}
	return lastNotifiedAt.After(dayStart) && check.Before(dayEnd)
}

func getDataFromDisk() (*Data, err) {
	var d Data
	content, err := ioutil.ReadAll(fPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %v: %v", fPath, err)
	}
	err := json.Unmarshal(content, *d)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal %v: %v", content, err)
	}
	return d, nil
}

func writeToDisk(*Data) err {
	f, err := os.Create(fPath)
	defer f.Close()
	if err != nil {
		return err
	}
	f.Write([]byte(json.Marshal(d)))
	return nil
}
