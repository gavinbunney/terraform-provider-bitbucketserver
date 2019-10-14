package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Error represents a error from the marketplace api.
type Error struct {
	StatusCode int
	Endpoint   string
}

func (e Error) Error() string {
	var errorMessages = ""
	return fmt.Sprintf("Marketplace Error: %d %s %s", e.StatusCode, e.Endpoint, errorMessages)
}

const marketplaceServer = "https://marketplace.atlassian.com"

type Client struct {
	HTTPClient *http.Client
}

func (c *Client) Do(method, endpoint string, payload *bytes.Buffer) (*http.Response, error) {

	absoluteendpoint := marketplaceServer + endpoint
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

	req.Header.Add("Content-Type", "application/json")
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

func (c *Client) DownloadArtifact(url string, dest *os.File) error {

	log.Printf("[DEBUG] Downloading file from %s", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Close = true

	resp, err := c.HTTPClient.Do(req)
	log.Printf("[DEBUG] Resp: %v Err: %v", resp, err)
	if resp != nil && (resp.StatusCode >= 400 || resp.StatusCode < 200) {
		apiError := Error{
			StatusCode: resp.StatusCode,
			Endpoint:   url,
		}

		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}

		log.Printf("[DEBUG] Resp Body: %s", string(body))

		_ = json.Unmarshal(body, &apiError)
		return error(apiError)

	}

	_, err = io.Copy(dest, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Get(endpoint string) (*http.Response, error) {
	return c.Do("GET", endpoint, nil)
}

func (c *Client) Post(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("POST", endpoint, jsonpayload)
}

func (c *Client) Put(endpoint string, jsonpayload *bytes.Buffer) (*http.Response, error) {
	return c.Do("PUT", endpoint, jsonpayload)
}

func (c *Client) PutOnly(endpoint string) (*http.Response, error) {
	return c.Do("PUT", endpoint, nil)
}

func (c *Client) Delete(endpoint string) (*http.Response, error) {
	return c.Do("DELETE", endpoint, nil)
}
