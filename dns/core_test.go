package dns

import (
	"cloudflare-dns/dns/cloudflare/model"
	"testing"
)

type DummyDNSClient struct {
	lastReq  *model.DNSRecordRequest
    responses map[string]interface{}
}

func (c *DummyDNSClient) DeleteRecord(r *model.DNSRecordRequest) (string, error) {
	c.lastReq = r

    response := c.responses["DeleteRecord"]

	if v, ok := response.(string); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return "", v
	}

	return "", nil
}

func (c *DummyDNSClient) NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.lastReq = r

    response := c.responses["NewRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
        return nil,v
	}

    return nil, nil
}

func (c *DummyDNSClient) UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.lastReq = r

    response := c.responses["NewRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
        return nil,v
	}

    return nil, nil
}

func (c *DummyDNSClient) ListZones() ([]model.Zone, error) {
	c.lastReq = nil

    response := c.responses["ListZones"]
	if v, ok := response.([]model.Zone); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
        return nil,v
	}

    return nil, nil
}

func (c *DummyDNSClient) ListRecords() ([]model.DNSRecord, error) {
	c.lastReq = nil

    response := c.responses["ListRecords"]
	if v, ok := response.([]model.Zone); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
        return nil,v
	}

    return nil, nil
}


var dummy DNSClient = &DummyDNSClient{}

func TestUpdateRecord(t *testing.T) {
}
