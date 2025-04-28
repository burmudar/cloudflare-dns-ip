package cloudflare

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/burmudar/cloudflare-dns/dns"
	"github.com/burmudar/cloudflare-dns/dns/cloudflare/model"
	"github.com/burmudar/cloudflare-dns/retrievers"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const API_CLOUDFLARE_V4 = "https://api.cloudflare.com/client/v4/"

var ErrFailedToCreateRequest = errors.New("failed to create request")

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

func NewHeaderCredentials(headers ...http.Header) dns.Credentials {
	return &HeaderCredentials{
		headers,
	}
}

type Client struct {
	http        *http.Client
	Credentials dns.Credentials
	ipRetriever retrievers.StringRetriever
	api         string
}

func NewTokenClient(apiURL, token string) (dns.DNSClient, error) {
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
		retrievers.DefaultIPRetriever,
		url.String(),
	}, nil
}

func (c *Client) urlJoin(p string) string {
	url := ""
	if strings.HasSuffix(c.api, "/") {
		url = c.api
	} else {
		url = c.api + "/"
	}

	if strings.HasPrefix(p, "/") {
		return url + p[1:]
	} else {
		return url + p
	}
}

func (c *Client) ExternalIP() (string, error) {
	return c.ipRetriever.Get()
}

func (c *Client) NewRequest(method string, url string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequest(method, url, body)

	if err != nil {
		return nil, fmt.Errorf("error creating request. %w", err)
	}

	c.Credentials.Apply(req)

	return req, nil
}

func (c *Client) doRequest(req *http.Request) (*http.Response, error) {
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request. %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResponse(resp)
	}

	return resp, nil
}

func (c *Client) ListZones() ([]*model.Zone, error) {
	var url = c.urlJoin("zones")

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	result := struct {
		Zones []*model.Zone `json:"result"`
	}{
		Zones: []*model.Zone{},
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to Zones: %v", err)
	}

	return result.Zones, nil
}

func (c *Client) ListRecords(zoneId string) ([]*model.DNSRecord, error) {
	var url = c.urlJoin(fmt.Sprintf("/zones/%s/dns_records", zoneId))

	req, err := c.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read response body: %v", err)
	}

	result := struct {
		Records []*model.DNSRecord `json:"result"`
	}{
		Records: []*model.DNSRecord{},
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("Failed to unmarshall response to DNSRecords: %v", err)
	}

	// results returned do not have the zone id, so we add it here
	for _, r := range result.Records {
		r.ZoneID = zoneId
	}

	return result.Records, nil
}

func (c *Client) UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	var url = c.urlJoin(fmt.Sprintf("zones/%s/dns_records/%s", r.ZoneID, r.ID))
	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("PUT", url, bytes.NewBuffer(data))

	if err != nil {
		return nil, ErrFailedToCreateRequest
	}
	resp, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body. %w", err)
	}

	return nil, err
}

func (c *Client) NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	var url = c.urlJoin(fmt.Sprintf("zones/%s/dns_records", r.ZoneID))

	data, err := json.Marshal(r)
	if err != nil {
		return nil, err
	}

	req, err := c.NewRequest("POST", url, bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errorFromResponse(resp)
	}

	data, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body. %w", err)
	}

	var record model.DNSRecord

	if err = json.Unmarshal(data, &record); err != nil {
		return nil, err
	}

	return &record, nil

}

func (c *Client) DeleteRecord(r *model.DNSDeleteRequest) (string, error) {
	var url = c.urlJoin(fmt.Sprintf("zones/%s/dns_records/%s", r.ZoneID, r.ID))

	req, err := c.NewRequest("DELETE", url, nil)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to do request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", errorFromResponse(resp)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body. %w", err)
	}

	var result = struct {
		id string
	}{}

	if err := json.Unmarshal(data, &result); err != nil {
		return "", fmt.Errorf("failed to unmarshall result: %w", err)
	}

	return result.id, nil
}

func errorFromResponse(resp *http.Response) error {
	msg := fmt.Sprintf("Response code <%d>\n", resp.StatusCode)
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		msg = msg + fmt.Sprintf("Body: <error: %s>\n", err)
	} else {
		msg = msg + fmt.Sprintf("Body: %s", string(data))
	}

	return errors.New(msg)

}
