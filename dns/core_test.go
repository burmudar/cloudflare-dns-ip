package dns

import (
	"cloudflare-dns/dns/cloudflare/model"
	"fmt"
	"testing"
)

type DummyDNSClient struct {
	Requests  map[string]*model.DNSRecordRequest
	Responses map[string]interface{}
}

func (c *DummyDNSClient) DeleteRecord(r *model.DNSRecordRequest) (string, error) {
	c.Requests["DeleteRecord"] = r

	response := c.Responses["DeleteRecord"]

	if v, ok := response.(string); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return "", v
	}

	return "", nil
}

func (c *DummyDNSClient) NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["NewRecord"] = r

	response := c.Responses["NewRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}

	return nil, nil
}

func (c *DummyDNSClient) UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["UpdateRecord"] = r

	response := c.Responses["UpdateRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}

	return nil, nil
}

func (c *DummyDNSClient) ListZones() ([]model.Zone, error) {
	c.Requests["ListZones"] = &model.DNSRecordRequest{}

	response := c.Responses["ListZones"]
	if v, ok := response.([]model.Zone); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}

	return nil, nil
}

func (c *DummyDNSClient) ListRecords(zoneID string) ([]model.DNSRecord, error) {
	c.Requests["ListRecords"] = nil

	response := c.Responses["ListRecords"]
	if v, ok := response.([]model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}

	return nil, nil
}

func NewDummyClient() *DummyDNSClient {
	return &DummyDNSClient{
		Requests:  make(map[string]*model.DNSRecordRequest),
		Responses: make(map[string]interface{}),
	}
}

var dummy DNSClient = &DummyDNSClient{}

func eq(left, right *model.DNSRecord) bool {
	if left.ID != right.ID {
		return false
	}
	if left.Name != right.Name {
		return false
	}
	if left.TTL != right.TTL {
		return false
	}
	if left.ZoneName != right.ZoneName {
		return false
	}
	if left.Content != right.Content {
		return false
	}
	if left.ZoneID != right.ZoneID {
		return false
	}
	if left.Created != right.Created {
		return false
	}
	if left.Type != right.Type {
		return false
	}
	if left.Locked != right.Locked {
		return false
	}
	if left.Modified != right.Modified {
		return false
	}

	return true
}

func TestUpdateRecord(t *testing.T) {
	t.Run("Record not found, Create Record is called", func(t *testing.T) {
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-123",
			ZoneName: "Test Zone",
			Content:  "127.0.0.1",
			Type:     "A",
		}
		var dummy DNSClient = &DummyDNSClient{
			Requests: make(map[string]*model.DNSRecordRequest),
			Responses: map[string]interface{}{
				"ListRecords": []model.DNSRecord{},
				"ListZones": []model.Zone{
					{
						ID:     "fake-123",
						Name:   "Test Zone",
						Status: "Test",
						Type:   "A",
					},
				},
				"NewRecord": &wanted,
			},
		}

		result, err := UpdateRecord(dummy, Record{
			ZoneName: "Test Zone",
			Type:     "A",
			Name:     "Thingy",
			IP:       "127.0.0.1",
			TTL:      200,
		})

		if err != nil {
			t.Fatalf("failure during update record: %v", err)
		}

		if !eq(&wanted, result) {
			t.Errorf("Wanted %v got %v", wanted, result)
		}

	})

	t.Run("Record found, Update record with IP", func(t *testing.T) {
		wanted := model.DNSRecord{
			ID:       "fake-record-222",
			ZoneID:   "fake-zone-id-222",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Content:  "255.255.255.255",
			Proxied:  false,
			Type:     "A",
			Locked:   false,
			TTL:      200,
		}

		var dummy DNSClient = &DummyDNSClient{
			Requests: make(map[string]*model.DNSRecordRequest),
			Responses: map[string]interface{}{
				"ListZones": []model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"ListRecords": []model.DNSRecord{
					{
						ID:       "fake-record-222",
						ZoneID:   "fake-zone-id-222",
						ZoneName: "fake-zone-name-222",
						Name:     "fake-record-name-222",
						Type:     "A",
						Content:  "128.127.1.1",
						Locked:   false,
						Proxied:  false,
					},
				},
			},
		}

		result, err := UpdateRecord(dummy, Record{
			ZoneName: "fake-zone-name-222",
			Type:     "A",
			Name:     "fake-record-name-222",
			IP:       "255.255.255.255",
			TTL:      200,
		})

		fmt.Printf("result: %v\n", result)

		if err != nil {
			t.Fatalf("failed during update record: %v", err)
		}

		if !eq(&wanted, result) {
			t.Errorf("wanted:\n %v \ngot %v", wanted, result)
		}
	})
}
