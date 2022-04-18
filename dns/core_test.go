package dns

import (
	"cloudflare-dns/dns/cloudflare/model"
	"fmt"
	"testing"
)

var ErrEmptyResponse error = fmt.Errorf("Empty Response. Did you forget to add a response for this method ?")

type DummyDNSClient struct {
	Requests  map[string]interface{}
	Responses map[string]interface{}
}

func (c *DummyDNSClient) DeleteRecord(r *model.DNSDeleteRequest) (string, error) {
	c.Requests["DeleteRecord"] = r

	response := c.Responses["DeleteRecord"]

	if v, ok := response.(string); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return "", v
	}
	return "", ErrEmptyResponse
}

func (c *DummyDNSClient) NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["NewRecord"] = r

	response := c.Responses["NewRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, ErrEmptyResponse
}

func (c *DummyDNSClient) UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["UpdateRecord"] = r

	response := c.Responses["UpdateRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, ErrEmptyResponse
}

func (c *DummyDNSClient) ListZones() ([]model.Zone, error) {
	c.Requests["ListZones"] = &model.DNSRecordRequest{}

	response := c.Responses["ListZones"]
	if v, ok := response.([]model.Zone); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, ErrEmptyResponse
}

func (c *DummyDNSClient) ListRecords(zoneID string) ([]model.DNSRecord, error) {
	c.Requests["ListRecords"] = nil

	response := c.Responses["ListRecords"]
	if v, ok := response.([]model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, ErrEmptyResponse
}

func NewDummyClient() *DummyDNSClient {
	return &DummyDNSClient{
		Requests:  make(map[string]interface{}),
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
			Requests: make(map[string]interface{}),
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
			Requests: make(map[string]interface{}),
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

func TestCreateRecord(t *testing.T) {
	var record = Record{
		ZoneName: "fake-zone-name-222",
		Type:     "A",
		Name:     "fake-record-name-222",
		IP:       "255.255.255.255",
		TTL:      200,
	}
	t.Run("No matching zone, returns error", func(t *testing.T) {
		var dummy DNSClient = &DummyDNSClient{
			Requests:  make(map[string]interface{}),
			Responses: map[string]interface{}{},
		}

		result, err := CreateRecord(dummy, record)

		if result != nil {
			t.Errorf("With no matching zone, result should be nil")
		}
		if err == nil {
			t.Errorf("With no matching zone, err should not be nil")
		}
	})

	t.Run("Matching zone, creates Record in zone", func(t *testing.T) {
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-123",
			ZoneName: "Test Zone",
			Content:  "127.0.0.1",
			Type:     "A",
			TTL:      300,
		}

		var dummy DNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"NewRecord": &wanted,
			},
		}

		result, err := CreateRecord(dummy, record)

		if err != nil {
			t.Fatalf("Unexpected error during CreateRecord: %v", err)
		}

		if !eq(result, &wanted) {
			t.Errorf("wanted %v got %v", result, wanted)
		}
	})
}

func TestDeleteRecord(t *testing.T) {
	var record = Record{
		ZoneName: "fake-zone-name-222",
		Type:     "A",
		Name:     "fake-record-name-222",
		IP:       "255.255.255.255",
		TTL:      200,
	}
	t.Run("No matching zone, returns error", func(t *testing.T) {
		var dummy DNSClient = &DummyDNSClient{
			Requests:  map[string]interface{}{},
			Responses: map[string]interface{}{},
		}

		_, err := DeleteRecord(dummy, record)

		if err == nil {
			t.Errorf("With no matching zone, err should not be nil")
		}
	})

	t.Run("Matching zone, Deletes Record", func(t *testing.T) {
		wanted := model.DNSRecord{
			ID:       "fake-record-222",
			ZoneID:   "fake-zone-id-222",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Type:     "A",
			Content:  "128.127.1.1",
			Locked:   false,
			Proxied:  false,
		}
		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"ListRecords":  []model.DNSRecord{wanted},
				"DeleteRecord": wanted.ID,
			},
		}

		result, err := DeleteRecord(dummy, record)
		if err != nil {
			t.Fatalf("Unexpected error during DeleteRecord: %v", err)
		}

		//Check the request that was sent
		actualReq := dummy.Requests["DeleteRecord"].(*model.DNSDeleteRequest)

		expectedReq := model.DNSDeleteRequest{
			ID:     wanted.ID,
			ZoneID: "fake-zone-id-222",
		}

		if expectedReq.ID != actualReq.ID {
			t.Errorf("Wanted '%s' Got '%s'. Incorrect Record ID used in Delete Request", expectedReq.ID, actualReq.ID)
		}

		if expectedReq.ZoneID != actualReq.ZoneID {
			t.Errorf("Wanted '%s' Got '%s'. Incorrect Zone ID used in Delete Request", expectedReq.ZoneID, actualReq.ZoneID)
		}

		//Now we can check the result

		if result == nil {
			t.Errorf("Wanted %v Got nil. Should return record deleted", record)
		}

		if result.ZoneName != record.ZoneName || ZoneType(result.Type) != record.Type || result.Name != record.Name {
			t.Errorf("Name, Type, ZoneName mistmatch. Wanted %v Got %v as result of DeleteRecord", record, result)
		}

	})
}
