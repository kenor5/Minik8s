package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func main_5() {
	url := "http://10.0.17.2:8070/function/hello_function"
	data := map[string]interface{}{
		"name": "John",
		"age":  18,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("JSON marshal error:", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("NewRequest error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Request error:", err)
		return
	}
	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	fmt.Println(buf.String())
}
