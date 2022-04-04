package cloudflare

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

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
}

func NewTokenClient(token string) *Client {
	var headers http.Header = make(http.Header)
	headers.Add("Authorization", "Bearer "+strings.TrimSpace(token))
	headers.Add("Content-Type", "application/json")

	return &Client{
		http.DefaultClient,
		NewHeaderCredentials(headers),
	}
}

type DNSRecordMeta struct {
	AutoAdd       bool   `json:"auto_added"`
	ManagedByApps bool   `json:"managed_by_apps"`
	ManagedByArgo bool   `json:"managed_by_argo_tunnel"`
	Source        string `json:"Source"`
}

type DNSRecord struct {
	ID        string         `json:"id"`
	ZoneID    string         `json:"zone_id"`
	ZoneName  string         `json:"zone_name"`
	Name      string         `json:"name"`
	Type      string         `json:"type"`
	Content   string         `json:"content"`
	Proxiable bool           `json:"proxiable"`
	Proxied   bool           `json:"proxied"`
	TTL       int            `json:"ttl"`
	Locked    bool           `json:"locked"`
	Created   *time.Time     `json:"Created"`
	Modified  *time.Time     `json:"Modified"`
	Meta      *DNSRecordMeta `json:"meta"`
}

type Zone struct {
	ID              string     `json:"id"`
	Name            string     `json:"name"`
	Status          string     `json:"status"`
	Paused          bool       `json:"paused"`
	Type            string     `json:"type"`
	DevelopmentMode int        `json:"development_mode"`
	NameServers     []string   `json:"name_servers"`
	OrigNameServers []string   `json:"original_name_servers"`
	OrigRegistrar   string     `json:"original_registrar"`
	Created         *time.Time `json:"created_on"`
	Modified        *time.Time `json:"modified_on"`
	Activated       *time.Time `json:"activated_on"`
}

func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, fmt.Errorf("Failed to create request for %s: %v\n", url, err)
	}

	c.Credentials.Apply(req)

	return req, nil
}

func (c *Client) ListZones() ([]Zone, error) {
	var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones")

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
		Zones []Zone `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to Zones: %v", err)
	}

	return result.Zones, nil
}

func (c *Client) ListDnsRecords(zoneId string) ([]DNSRecord, error) {
	var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records", zoneId)

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
		Records []DNSRecord `json:"result"`
	}{}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to DNSRecords: %v", err)
	}

	return result.Records, nil
}

func (c *Client) UpdateDNSRecord(r *DNSRecord) error {
	var url = fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", r.ZoneID, r.ID)
	data, err := json.Marshal(r)

	req, err := c.NewRequest("PUT", url, bytes.NewBuffer(data))

	if err != nil {
		return err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("Failed to do request: %v", err)
	}

	data, err = ioutil.ReadAll(resp.Body)
	fmt.Println(string(data))

	return err

}
