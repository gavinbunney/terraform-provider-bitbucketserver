package bitbucket

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
)

// Error represents a error from the bitbucket api.
type Error struct {
	Errors []struct {
		Context   string `json:"context,omitempty"`
		Message   string `json:"message,omitempty"`
		Exception string `json:"exceptionName,omitempty"`
	} `json:"errors,omitempty"`
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {

	var errorMessages = ""
	if e.Errors != nil {
		for _, err := range e.Errors {
			errorMessages += err.Message + "\n"
		}
	}

	return fmt.Sprintf("API Error: %d %s %s", e.StatusCode, e.Endpoint, errorMessages)
}

type BitbucketClient struct {
	Server     string
	Username   string
	Password   string
	HTTPClient *http.Client
}

func (c *BitbucketClient) Do(method, endpoint string, payload *bytes.Buffer, contentType string) (*http.Response, error) {

	absoluteendpoint := c.Server + endpoint
	log.Printf("[DEBUG] Sending request to %s %s", method, absoluteendpoint)

	var bodyreader io.Reader

	if payload != nil {
		log.Printf("[DEBUG] With payload %s", payload.String())
		bodyreader = payload
	}

	req, err := http.NewRequest(method, absoluteendpoint, bodyreader)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Add("X-Atlassian-Token", "no-check")

	if payload != nil {
		if contentType != "" {
			req.Header.Add("Content-Type", contentType)
		} else {
			req.Header.Add("Content-Type", "application/json")
		}
	}

	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp != nil && (resp.StatusCode >= 400 || resp.StatusCode < 200) {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		_ = json.Unmarshal(body, &apiError)
		return resp, error(apiError)

	}
	return resp, err
}

// Creates a new file upload http request with optional extra params
func (c *BitbucketClient) PostFileUpload(endpoint string, params map[string]string, paramName, path string) (*http.Response, error) {
	absoluteendpoint := c.Server + endpoint
	log.Printf("[DEBUG] Sending request to POST %s", absoluteendpoint)

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, err := writer.CreateFormFile(paramName, filepath.Base(path))
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(part, file)

	for key, val := range params {
		_ = writer.WriteField(key, val)
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", absoluteendpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Add("X-Atlassian-Token", "no-check")
	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp != nil && (resp.StatusCode >= 400 || resp.StatusCode < 200) {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   endpoint,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		_ = json.Unmarshal(body, &apiError)
		return resp, error(apiError)

	}

	return resp, err
}

func (c *BitbucketClient) Get(endpoint string) (*http.Response, error) {
	return c.Do("GET", endpoint, nil, "application/json")
}

func (c *BitbucketClient) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, jsonpayload, "application/json")
}

func (c *BitbucketClient) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("PUT", endpoint, jsonpayload, "application/json")
}

func (c *BitbucketClient) PutOnly(endpoint string) (*http.Response, error) {
	return c.Do("PUT", endpoint, nil, "application/json")
}

func (c *BitbucketClient) Delete(endpoint string) (*http.Response, error) {
	return c.Do("DELETE", endpoint, nil, "application/json")
}
