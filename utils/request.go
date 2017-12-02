package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// PrepareRequest a function that prepares a http.Request object with proper headers
func PrepareRequest(httpMethod string, url string, payload io.Reader, accessToken string) (request *http.Request, err error) {
	req, err := http.NewRequest(httpMethod, url, payload)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Cache-Control", "no-cache")
	if accessToken != "" {
		var buffer bytes.Buffer
		buffer.WriteString("Bearer ")
		buffer.WriteString(accessToken)
		req.Header.Add("Authorization", buffer.String())
	}
	return req, nil
}

// MakeRequest executes http.Request and fills passed object from JSON response
func MakeRequest(request *http.Request, dataObject interface{}) error {
	res, err := http.DefaultClient.Do(request)

	if err != nil {
		return err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	// bodyStr := string(body)
	// fmt.Println(bodyStr)
	err = json.Unmarshal(body, &dataObject)

	if err != nil {
		return err
	}
	return nil
}
