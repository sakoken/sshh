package json

import (
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
