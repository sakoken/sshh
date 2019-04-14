package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func CreateJson(filePath string) {
	if _, err := os.Stat(filePath); !os.IsNotExist(err) {
		return
	}
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0700)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	fmt.Fprintln(file, "{}")
}

func ReadJson(filePath string, v interface{}) {
	dat, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Fatal(err)
	}
	json.Unmarshal(dat, v)
}

func WriteJson(filePath string, v interface{}) error {
	b, err := json.Marshal(v)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = json.Indent(&buf, []byte(b), "", "  "); err != nil {
		return err
	}
	return ioutil.WriteFile(filePath, buf.Bytes(), 0777)
}
