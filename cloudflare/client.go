package cloudflare

import (
	"bytes"
	"cloudflare-dns/cloudflare/model"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const API_CLOUDFLARE_V4 = "https://api.cloudflare.com/client/v4/"

type Credentials interface {
	Apply(req *http.Request) error
}

type HeaderCredentials struct {
	Headers []http.Header
}


func (h *HeaderCredentials) Apply(req *http.Request) error {

	for _, item := range h.Headers {
		for k, v := range item {
			for _, i := range v {
				req.Header.Add(k, i)
			}
		}
	}
	return nil
}

func NewHeaderCredentials(headers ...http.Header) Credentials {
	return &HeaderCredentials{
		headers,
	}
}

type Client struct {
	http        *http.Client
	Credentials Credentials
	apiURL      *url.URL
}

func NewTokenClient(apiURL, token string) (*Client, error) {
	var headers http.Header = make(http.Header)
	headers.Add("Authorization", "Bearer "+strings.TrimSpace(token))
	headers.Add("Content-Type", "application/json")

	url, err := url.Parse(apiURL)
	if err != nil {
		return nil, err
	}

	return &Client{
		http.DefaultClient,
		NewHeaderCredentials(headers),
		url,
	}, nil
}

func (c *Client) formatURL(format string, params ...interface{}) string {
	c.apiURL.Path = fmt.Sprintf(format, params ...)

	return c.apiURL.RequestURI()
}

func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create request for %s: %v\n", url, err)
	}

	c.Credentials.Apply(req)

	return req, nil
}

func (c *Client) ListZones() ([]model.Zone, error) {
	var url = c.formatURL("zones")

	req, err := c.NewRequest("GET", url, nil)
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %v", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	result := struct {
		Zones []model.Zone `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to Zones: %v", err)
	}

	return result.Zones, nil
}

func (c *Client) ListDnsRecords(zoneId string) ([]model.DNSRecord, error) {
	var url = c.formatURL("/zones/%s/dns_records", zoneId)

	req, err := c.NewRequest("GET", url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %v", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	result := struct {
		Records []model.DNSRecord `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to DNSRecords: %v", err)
	}

	return result.Records, nil
}

func (c *Client) UpdateDNSRecord(r *model.DNSRecordRequest) error {
	var url = c.formatURL("zones/%s/dns_records/%s", r.ZoneID, r.ID)
	data, err := json.Marshal(r)
	if err != nil {
        return err
	}

	req, err := c.NewRequest("PUT", url, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to do request: %w", err)
	}

	data, err = ioutil.ReadAll(resp.Body)

	return err
}

func (c *Client) NewDNSRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	var url = c.formatURL("zones/%s/dns_records", r.ZoneID)

	data, err := json.Marshal(r)
	if err != nil {
        return nil, err
	}

    req, err := c.NewRequest("POST", url, bytes.NewBuffer(data))

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %w", err)
	}

	data, err = ioutil.ReadAll(resp.Body)

    var record model.DNSRecord

    if err  = json.Unmarshal(data, &record); err != nil {
        return nil, err
    }

    return &record, nil

}
