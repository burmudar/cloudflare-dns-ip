package dns

import (
	"errors"
	"fmt"
	"github.com/burmudar/cloudflare-dns/dns/cloudflare/model"
	"net/http"
	"os"
	"strings"
)

var ErrZoneNotFound = errors.New("Zone not found")
var ErrRecordNotFound = errors.New("Record not found")

type ZoneType string

const (
	AType ZoneType = "A"
)

type DNSClient interface {
	ExternalIP() (string, error)
	UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error)
	NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error)
	DeleteRecord(r *model.DNSDeleteRequest) (string, error)

	ListZones() ([]*model.Zone, error)
	ListRecords(zoneID string) ([]*model.DNSRecord, error)
}

type Credentials interface {
	Apply(req *http.Request) error
}

type Record struct {
	ID       string
	ZoneName string
	Type     ZoneType
	Name     string
	IP       string
	TTL      int
}

func filterByName(records []*model.DNSRecord, name string) *model.DNSRecord {
	for _, r := range records {
		if r.Name == name {
			return r
		}
	}

	return nil
}

func filterZoneByName(zones []*model.Zone, name string) *model.Zone {
	for _, z := range zones {
		if z.Name == name {
			return z
		}
	}

	return nil
}

func UpdateRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	defer func() { fmt.Fprintln(os.Stderr, "") }()
	fmt.Fprintf(os.Stderr, "Locating DNS record: %s ...", record.Name)
	remoteRecord, err := FindRecord(client, record)
	if err != nil {
		fmt.Fprintln(os.Stderr, "NOT FOUND")
		return CreateRecord(client, record)
	}
	fmt.Fprintln(os.Stderr, "FOUND")

	var ip = record.IP
	if record.IP == "" {
		fmt.Fprintln(os.Stderr, "Fetching external ip ...")
		ip, err = client.ExternalIP()
		if err != nil {
			return nil, fmt.Errorf("error getting external ip: %w", err)
		}
	}

	fmt.Fprintf(os.Stderr, "Using IP: %s\n", ip)

	fmt.Fprintf(os.Stderr, "--- Current DNS Record ---\n%s\n", remoteRecord.String())

	if ip == remoteRecord.Content {
		fmt.Fprintf(os.Stderr, "DNS Record [%s %s] content already contains: %s", remoteRecord.Type, record.Name, ip)
		return remoteRecord, nil

	}
	req := model.DNSRecordRequest{
		ID:      remoteRecord.ID,
		ZoneID:  remoteRecord.ZoneID,
		Name:    remoteRecord.Name,
		Type:    remoteRecord.Type,
		Content: ip,
		Proxied: remoteRecord.Proxied,
		TTL:     record.TTL,
	}

	if err := req.Sanitize(); err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stdout, "--- Updating DNS Record ---\n%s\n", req.String())

	return client.UpdateRecord(&req)
}

func CreateRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	zone, err := FindZone(client, record.ZoneName)
	if err != nil {
		return nil, err
	}

	var ip = record.IP
	if record.IP == "" {
		ip, err = client.ExternalIP()
		if err != nil {
			return nil, fmt.Errorf("error getting external ip: %w", err)
		}
	}

	req := model.DNSRecordRequest{
		ZoneID:   zone.ID,
		Name:     record.Name,
		Content:  ip,
		Type:     string(record.Type),
		Proxied:  false,
		TTL:      record.TTL,
		Priority: 10,
	}

	if err := req.Sanitize(); err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "--- Creating DNS Record ---\n%s", req.String())

	return client.NewRecord(&req)
}

func FindZone(client DNSClient, zoneName string) (*model.Zone, error) {
	zones, err := client.ListZones()
	if err != nil {
		return nil, fmt.Errorf("Error while listing zones: %v\n", err)
	}

	zone := filterZoneByName(zones, zoneName)
	if zone == nil {
		return nil, ErrZoneNotFound
	}

	return zone, nil
}

func ListRecords(client DNSClient, zoneID string) ([]*model.DNSRecord, error) {
	zone, err := FindZone(client, zoneID)
	if err != nil {
		return nil, err
	}

	return client.ListRecords(zone.ID)
}

func FindRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	records, err := ListRecords(client, record.ZoneName)
	if err != nil {
		return nil, err
	}

	remoteRecord := filterByName(records, record.Name)
	if remoteRecord == nil {
		return nil, ErrZoneNotFound
	}

	return remoteRecord, nil
}

func DeleteRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	dnsRecord, err := FindRecord(client, record)
	if err != nil {
		return nil, err
	}

	//Ignoring the ID that gets sent back, since it's essentially dnsRecord.ID
	_, err = client.DeleteRecord(&model.DNSDeleteRequest{
		ID:     dnsRecord.ID,
		ZoneID: dnsRecord.ZoneID,
	})

	if err != nil {
		return nil, err
	}
	return dnsRecord, err
}

func NormaliseRecordName(zoneName string, name string) string {

	if strings.HasSuffix(name, zoneName) {
		return name
	}

	return name + "." + zoneName
}
