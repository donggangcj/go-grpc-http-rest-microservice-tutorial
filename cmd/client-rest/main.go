package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

func main() {
	// get configuration
	address := flag.String("server", "http://localhost:8080", "HTTP gateway url, e.g. http://localhost:8080")
	flag.Parse()

	pfx := time.Now().Format(time.RFC3339Nano)
	var body string

	// Call Create
	resp, err := http.Post(*address+"/v1/todo", "application/json", strings.NewReader(fmt.Sprintf(`
		{
			"api":"v1",
			"toDo": {
				"title":"title (%s)",
				"description":"description (%s)",
				"reminder":"%s"
			}
		}
	`, pfx, pfx, pfx)))
	if err != nil {
		log.Fatalf("failed to call Created method: %v ", err)
	}
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Create respones body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Create Response :code = %d,Body= %v", resp.StatusCode, body)

	// Parse Id of created ToDo
	var created struct {
		API string `json:"api"`
		ID  string `json:"id"'`
	}
	err = json.Unmarshal(bodyBytes, &created)
	if err != nil {
		log.Fatalf("failed to unmarshal JSON response of Create method: %v", err)
		fmt.Println("error:", err)
	}

	// Call Read
	resp, err = http.Get(fmt.Sprintf("%s%s/%s", *address, v1/todo", created.ID))
	if err != nil {
		log.Fatalf("failed to call Read method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Read response boyd: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Read Response: code=%d,Body=%s\n\n", resp.StatusCode, body)

	// Call update
	req, err := http.NewRequest("PUT", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID),
		strings.NewReader(fmt.Sprintf(`
		{
					"api":"v1",
					"toDo": {
						"title":"title (%s) + updated",
						"description":"description (%s) + updated",
						"reminder":"%s"
					}
				}
	`, pfx, pfx, pfx)))
	req.Header.Set("content-Type", "application/json")
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Update method:%v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read update response body:%v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Update response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call ReadAll
	resp, err = http.Get(*address + "/v1/todo/all")
	if err != nil {
		log.Fatalf("failed to call ReadAll method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read ReadAll response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("ReadAll response: Code=%d, Body=%s\n\n", resp.StatusCode, body)

	// Call Delete
	req, err = http.NewRequest("DELETE", fmt.Sprintf("%s%s/%s", *address, "/v1/todo", created.ID), nil)
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		log.Fatalf("failed to call Delete method: %v", err)
	}
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		body = fmt.Sprintf("failed read Delete response body: %v", err)
	} else {
		body = string(bodyBytes)
	}
	log.Printf("Delete response: Code=%d, Body=%s\n\n", resp.StatusCode, body)
}
