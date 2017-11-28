package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type Data struct {
	LastCall, NotifiedToday string
}

func getDataFromDisk() (*Data, error) {
	var d Data
	f, err := os.Open(fPath)
	if err != nil {
		return nil, fmt.Errorf("could not read file %v: %v", fPath, err)
	}
	content, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("could not read file %v: %v", fPath, err)
	}
	err = json.Unmarshal(content, &d)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal %v: %v", content, err)
	}
	return &d, nil
}

func writeToDisk(d *Data) error {
	f, err := os.Create(fPath)
	defer f.Close()
	if err != nil {
		return err
	}
	j, _ := json.Marshal(d)
	f.Write([]byte(j))
	return nil
}
