package dns

import (
	"errors"
	"fmt"
	"github.com/burmudar/cloudflare-dns/dns/cloudflare/model"
	"testing"
)

var ErrEmptyResponse error = errors.New("Empty Response. Did you forget to add a response for this method ?")

type DummyDNSClient struct {
	IP        string
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
	return "", fmt.Errorf("Error DeleteRecord: %w", ErrEmptyResponse)
}

func (c *DummyDNSClient) NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["NewRecord"] = r

	response := c.Responses["NewRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, fmt.Errorf("Error NewRecord: %w", ErrEmptyResponse)
}

func (c *DummyDNSClient) UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error) {
	c.Requests["UpdateRecord"] = r

	response := c.Responses["UpdateRecord"]
	if v, ok := response.(*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, fmt.Errorf("Error UpdateRecord: %w", ErrEmptyResponse)
}

func (c *DummyDNSClient) ListZones() ([]*model.Zone, error) {
	c.Requests["ListZones"] = &model.DNSRecordRequest{}

	response := c.Responses["ListZones"]
	if v, ok := response.([]*model.Zone); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, fmt.Errorf("Error ListZones: %w", ErrEmptyResponse)
}

func (c *DummyDNSClient) ListRecords(zoneID string) ([]*model.DNSRecord, error) {
	c.Requests["ListRecords"] = nil

	response := c.Responses["ListRecords"]
	if v, ok := response.([]*model.DNSRecord); ok {
		return v, nil
	} else if v, ok := response.(error); ok {
		return nil, v
	}
	return nil, fmt.Errorf("Error ListRecords: %w", ErrEmptyResponse)
}

func (c *DummyDNSClient) ExternalIP() (string, error) {
	return c.IP, nil
}

func NewDummyClient() *DummyDNSClient {
	return &DummyDNSClient{
		Requests:  make(map[string]interface{}),
		Responses: make(map[string]interface{}),
	}
}

var dummy DNSClient = &DummyDNSClient{}

func validateRequest(t *testing.T, request *model.DNSRecordRequest, record *model.DNSRecord) {
	if request.Content != record.Content {
		t.Errorf("Got %+s. Wanted %+s. Incorrect content values", request.Content, record.Content)
	}
	if request.ZoneID != record.ZoneID {
		t.Errorf("Got %s. Wanted %s. Incorrect ZoneID values", request.ZoneID, record.ZoneID)
	}
	if request.Name != record.Name {
		t.Errorf("Got %s. Wanted %s. Incorrect Name values", request.Name, record.Name)
	}
	if request.TTL != record.TTL {
		t.Errorf("Got %d. Wanted %d. Incorrect TTL values", request.TTL, record.TTL)
	}
	if request.Type != record.Type {
		t.Errorf("Got %s. Wanted %s. Incorrect Type values", request.Type, record.Type)
	}
}

func eq(t *testing.T, left, right *model.DNSRecord) bool {
	if left.ID != right.ID {
		t.Errorf("Left %v. Right %v. Incorrect value for ID", left.ID, right.ID)
	}
	if left.Name != right.Name {
		t.Errorf("Left %v. Right %v. Incorrect value for Name", left.Name, right.Name)
	}
	if left.TTL != right.TTL {
		t.Errorf("Left %v. Right %v. Incorrect value for TTL", left.TTL, right.TTL)
	}
	if left.ZoneName != right.ZoneName {
		t.Errorf("Left %v. Right %v. Incorrect value for ZoneName", left.ZoneName, right.ZoneName)
	}
	if left.Content != right.Content {
		t.Errorf("Left %v. Right %v. Incorrect value for Content", left.Content, right.Content)
	}
	if left.ZoneID != right.ZoneID {
		t.Errorf("Left %v. Right %v. Incorrect value for ZoneID", left.ZoneID, right.ZoneID)
	}
	if left.Created != right.Created {
		t.Errorf("Left %v. Right %v. Incorrect value for Created", left.Created, right.Created)
	}
	if left.Type != right.Type {
		t.Errorf("Left %v. Right %v. Incorrect value for Type", left.Type, right.Type)
	}
	if left.Locked != right.Locked {
		t.Errorf("Left %v. Right %v. Incorrect value for Locked", left.Locked, right.Locked)
	}
	if left.Modified != right.Modified {
		t.Errorf("Left %v. Right %v. Incorrect value for Modified", left.Modified, right.Modified)
	}

	return true
}

func TestUpdateRecord(t *testing.T) {
	t.Run("Record not found, Create Record is called", func(t *testing.T) {
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-123",
			ZoneName: "Test Zone",
			Name:     "Thingy",
			Content:  "127.0.0.1",
			Type:     "A",
			TTL:      200,
		}
		var dummy = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListRecords": []*model.DNSRecord{},
				"ListZones": []*model.Zone{
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

		req := dummy.Requests["NewRecord"].(*model.DNSRecordRequest)
		validateRequest(t, req, &wanted)

		eq(t, &wanted, result)

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

		var dummy = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"ListRecords": []*model.DNSRecord{
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
				"UpdateRecord": &wanted,
			},
		}

		result, err := UpdateRecord(dummy, Record{
			ZoneName: "fake-zone-name-222",
			Type:     "A",
			Name:     "fake-record-name-222",
			IP:       "255.255.255.255",
			TTL:      200,
		})

		if err != nil {
			t.Fatalf("failed during update record: %v", err)
		}

		req := dummy.Requests["UpdateRecord"].(*model.DNSRecordRequest)
		validateRequest(t, req, &wanted)
		eq(t, &wanted, result)
	})
	t.Run("Record with missing Type - Type 'A' is added automatically", func(t *testing.T) {
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

		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"ListRecords": []*model.DNSRecord{
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
				"UpdateRecord": &wanted,
			},
		}

		result, err := UpdateRecord(dummy, Record{
			ZoneName: "fake-zone-name-222",
			Type:     "",
			Name:     "fake-record-name-222",
			IP:       "255.255.255.255",
			TTL:      200,
		})

		if err != nil {
			t.Fatalf("failed during update record: %v", err)
		}

		req := dummy.Requests["UpdateRecord"].(*model.DNSRecordRequest)
		validateRequest(t, req, &wanted)
		eq(t, &wanted, result)
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
			ZoneID:   "fake-zone-id-222",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Content:  record.IP,
			Type:     "A",
			TTL:      200,
		}

		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
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

		req := dummy.Requests["NewRecord"].(*model.DNSRecordRequest)
		validateRequest(t, req, &wanted)
		eq(t, result, &wanted)

	})
	t.Run("Record with missing Type - Type 'A' automatically added", func(t *testing.T) {
		var record = Record{
			ZoneName: "fake-zone-name-222",
			Type:     " ",
			Name:     "fake-record-name-222",
			IP:       "255.255.255.255",
			TTL:      200,
		}
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-zone-id-222",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Content:  record.IP,
			Type:     "A",
			TTL:      200,
		}

		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
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

		req := dummy.Requests["NewRecord"].(*model.DNSRecordRequest)
		validateRequest(t, req, &wanted)
		eq(t, result, &wanted)

	})
	t.Run("Matching zone and set ip, creates Record with provided ip", func(t *testing.T) {
		record = Record{
			ZoneName: "fake-zone-name-222",
			Type:     "A",
			Name:     "fake-record-name-222",
			IP:       "999.999.999.999",
			TTL:      200,
		}
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-zone-id",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Content:  "999.999.999.999",
			Type:     "A",
			TTL:      200,
		}

		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
					{
						ID:     wanted.ZoneID,
						Name:   wanted.ZoneName,
						Status: "ACTIVE",
						Type:   wanted.Type,
					},
				},
				"NewRecord": &wanted,
			},
		}

		//for this test, this should not be the discovered IP and the IP in the record should be used!
		dummy.IP = "should not be the IP"

		result, err := CreateRecord(dummy, record)

		if err != nil {
			t.Fatalf("Unexpected error during CreateRecord: %v", err)
		}

		request := dummy.Requests["NewRecord"].(*model.DNSRecordRequest)

		validateRequest(t, request, &wanted)
		eq(t, result, &wanted)
	})
	t.Run("Matching zone and no ip, creates Record in zone and discovers ip", func(t *testing.T) {
		record = Record{
			ZoneName: "fake-zone-name-222",
			Type:     "A",
			Name:     "fake-record-name-222",
			IP:       "",
			TTL:      200,
		}
		wanted := model.DNSRecord{
			ID:       "fake-record-123",
			ZoneID:   "fake-zone-id",
			ZoneName: "fake-zone-name-222",
			Name:     "fake-record-name-222",
			Content:  "169.0.253.37",
			Type:     "A",
			TTL:      200,
		}

		var dummy *DummyDNSClient = &DummyDNSClient{
			Requests: make(map[string]interface{}),
			Responses: map[string]interface{}{
				"ListZones": []*model.Zone{
					{
						ID:     wanted.ZoneID,
						Name:   wanted.ZoneName,
						Status: "ACTIVE",
						Type:   wanted.Type,
					},
				},
				"NewRecord": &wanted,
			},
		}

		dummy.IP = wanted.Content

		result, err := CreateRecord(dummy, record)

		if err != nil {
			t.Fatalf("Unexpected error during CreateRecord: %v", err)
		}

		request := dummy.Requests["NewRecord"].(*model.DNSRecordRequest)

		validateRequest(t, request, &wanted)
		eq(t, result, &wanted)
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
				"ListZones": []*model.Zone{
					{
						ID:     "fake-zone-id-222",
						Name:   "fake-zone-name-222",
						Status: "ACTIVE",
						Type:   "A",
					},
				},
				"ListRecords":  []*model.DNSRecord{&wanted},
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
