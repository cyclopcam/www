package www

import (
	"encoding/json"
	"io"
	"net/http"
)

// Perform the request, and if any errors occurs (transport or non-200 status code), return an error
// This does the work for you of checking for a non-200 response, reading the response body,
// and turning it into an error.
func Do(req *http.Request) (*http.Response, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		respB, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, Error(resp.StatusCode, string(respB))
	}
	return resp, nil
}

// Returns nil if the status code is 200 and the JSON decodes, or an error in all other cases.
func FetchJSON(req *http.Request, output any) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	return HandleJSONFetch(resp, output)
}

// Returns nil if the status code is 200 and the JSON decodes, or an error in all other cases.
// Closes resp.Body before returning.
func HandleJSONFetch(resp *http.Response, output any) error {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respB, _ := io.ReadAll(resp.Body)
		return Error(resp.StatusCode, string(respB))
	}
	return json.NewDecoder(resp.Body).Decode(output)
}

// Returns (string(body), nil) if the status code is 200, or an error in all other cases.
func FetchText(req *http.Request) (string, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	return HandleTextFetch(resp)
}

// Returns (string(body), nil) if the status code is 200, or an error in all other cases.
// Closes resp.Body before returning.
func HandleTextFetch(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		respB, _ := io.ReadAll(resp.Body)
		return "", Error(resp.StatusCode, string(respB))
	}
	respB, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(respB), nil
}
