package http

import (
	"fmt"
	"testing"
)

func TestHttpGet(t *testing.T) {
	url := "https://cn.bing.com/search?q=hello&form=QBLH&sp=-1&lq=0&pq=hello&sc=10-5&qs=n&sk=&cvid=0983077C589B417A868EBC59B12BE1E4&ghsh=0&ghacc=0&ghpl="
	var payload interface{}
	err := HttpGet(url, &payload)
	if err != nil {
		t.Error("get error")
		return
	}
	fmt.Println(payload)

}

func TestHttpPost(t *testing.T) {

}
