package global

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
)

func SaveJson(data interface{}) error {
	b, err := json.Marshal(data)
	if err != nil {
		return err
	}
	var buf bytes.Buffer
	if err = json.Indent(&buf, []byte(b), "", "  "); err != nil {
		return err
	}
	return ioutil.WriteFile(SshhJson(), buf.Bytes(), 0777)
}
