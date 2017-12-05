package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type data struct {
	LastCall, NotifiedToday string
}

func getDataFromDisk() (*data, error) {
	var d data
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

func writeToDisk(d *data) error {
	f, err := os.Create(fPath)
	defer f.Close()
	if err != nil {
		return err
	}
	j, _ := json.Marshal(d)
	f.Write([]byte(j))
	return nil
}
