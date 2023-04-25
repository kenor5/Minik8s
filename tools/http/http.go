package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"minik8s/tools/log"
	"net/http"
	"strings"
)

/*
	REFERENCE TO
	https://cloud.tencent.com/developer/article/1849807
*/

func HttpGet(url string, v interface{}) error {
	res, err := http.Get(url)

	if err != nil {
		log.LOG("http get error")
		return err
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, v)
	return err
}

type contentType int

const (
	JSON contentType = iota
	STRING
	FORM
)

func HttpPost(url string, v interface{}, t contentType) error {
	switch t {
	case JSON:
		marshal, err := json.Marshal(v)
		if err != nil {
			return err
		}
		res, err := http.Post(url, "application/json", bytes.NewReader(marshal))
		if err != nil {
			return err
		}
		log.LOG("post to url")
		fmt.Println(res.Body)
	case STRING:
		value := v.(string)
		_, err := http.Post(url, "text/plain", strings.NewReader(value))
		if err != nil {
			return err
		}
	}

	return nil
}

func ApplyToServer(url string, v interface{}) error {
	err := HttpPost(url, v, JSON)
	return err
}
