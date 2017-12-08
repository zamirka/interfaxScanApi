package utils

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

// PrepareRequest a function that prepares a http.Request object with proper headers
func PrepareRequest(httpMethod string, url string, payload io.Reader, accessToken string) (request *http.Request, err error) {
	var req *http.Request
	if req, err = http.NewRequest(httpMethod, url, payload); err != nil {
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
func MakeRequest(request *http.Request, dataObject interface{}) (err error) {
	var res *http.Response
	if res, err = http.DefaultClient.Do(request); err != nil {
		return err
	}
	defer res.Body.Close()
	var body []byte
	if body, err = ioutil.ReadAll(res.Body); err != nil {
		return err
	}
	if err = json.Unmarshal(body, &dataObject); err != nil {
		return err
	}
	return nil
}
