package main

import (
	"bytes"
	"net/http"
)

func main() {
	posturl := "http://127.0.0.1:8080/task"

	body := []byte(`{
		"Title": "Task title",
		"Description": "Task description"
	}`)

	for i := 0; i < 10; i++ {

		resp, err := http.Post(posturl, "application/json", bytes.NewBuffer(body))

		if err != nil {
			panic(err)
		}

		defer resp.Body.Close()
	}

}
