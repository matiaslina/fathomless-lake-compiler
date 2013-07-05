// Package firebase implements the functions to manipulate Firebase
package firebase

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
)

// Firebase struct
type firebaseRoot struct {
	base_url   string
	auth_token string
}

/*
Class Methods
*/

// Construct a firebase type
func New(url string) *firebaseRoot {
	return &firebaseRoot{base_url: url}
}

/*
Instance Methods
Firebase Package - REST Interface
*/

// Writes and returns the data to the firebase endpoint
// Example: firebase.Set('/users/info', [struct]) => []byte(data), error
func (f *firebaseRoot) Set(path string, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	json := bytes.NewBuffer(b)

	body, err := f.SendRequest("PUT", path, json)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Returns the data at a path
func (f *firebaseRoot) Get(path string) ([]byte, error) {
	body, err := f.SendRequest("GET", path, nil)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// Writes the data, returns the key name of the data added
// Example: firebase.Push('users', [struct]) => {'id':'-INOQPH-aV_psbk3ZXEX'}
func (f *firebaseRoot) Push(path string, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	json := bytes.NewBuffer(b)

	body, err := f.SendRequest("POST", path, json)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Updates the data at path and returns the data. Does not delete omitted children.
func (f *firebaseRoot) Update(path string, v interface{}) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	json := bytes.NewBuffer(b)

	body, err := f.SendRequest("PATCH", path, json)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Deletes the data and returns true or false
func (f *firebaseRoot) Delete(path string) (bool, error) {
	_, err := f.SendRequest("DELETE", path, nil)
	if err != nil {
		return false, err
	}
	return true, nil
}

// Send HTTP Request, return data
func (f *firebaseRoot) SendRequest(method string, path string, body io.Reader) ([]byte, error) {
	url := f.BuildURL(path)
	// Append .json to use Firebase REST API
	url += ".json"

	// Append auth token if one exists
	if len(f.auth_token) != 0 {
		url += "?auth=" + f.auth_token
	}

	// create a request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}

	// send JSON to firebase
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Bad HTTP Response: %v", resp.Status)
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return b, nil
}

// Build URL based on base_url and relative path
func (f *firebaseRoot) BuildURL(path string) string {
	base_url, _ := url.Parse(f.base_url)
	url, _ := toAbsURL(base_url, path)
	return url
}

/*
Helper Methods
*/

func toAbsURL(base_url *url.URL, path string) (string, error) {
	relurl, err := url.Parse(path)
	if err != nil {
		return "", err
	}
	absurl := base_url.ResolveReference(relurl)

	return absurl.String(), nil
}
