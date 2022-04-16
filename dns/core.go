package dns

import (
	"cloudflare-dns/dns/cloudflare/model"
	"cloudflare-dns/retrievers"
	"fmt"
	"net/http"
	"os"
)

type ZoneType string

const (
	AType ZoneType = "A"
)

type DNSClient interface {
	UpdateRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error)
	NewRecord(r *model.DNSRecordRequest) (*model.DNSRecord, error)
	DeleteRecord(r *model.DNSDeleteRequest) (string, error)

	ListZones() ([]model.Zone, error)
	ListRecords(zoneID string) ([]model.DNSRecord, error)
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

func filterByName(records []model.DNSRecord, name string) *model.DNSRecord {
	for _, r := range records {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

func filterZoneByName(zones []model.Zone, name string) *model.Zone {
	for _, z := range zones {
		if z.Name == name {
			return &z
		}
	}

	return nil
}

func UpdateRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	remoteRecord, err := FindRecord(client, record)
	if err != nil {
		return CreateRecord(client, record)
	}
	fmt.Fprintf(os.Stderr, "FOUND!\n")

	var retriever retrievers.StringRetriever
	if record.IP != "" {
		retriever = retrievers.NewStaticStringRetriever(record.IP)
		fmt.Fprintf(os.Stderr, "Manually setting ip ...")
	} else {
		retriever = retrievers.DefaultIPRetriever
		fmt.Fprintf(os.Stderr, "Discovering public ip ...")
	}

	ip, err := retriever.Get()

    if err != nil {
		return nil, fmt.Errorf("Failed to get ip: %v\n", err)
	}

    if ip != remoteRecord.Content {
		fmt.Fprintf(os.Stderr, "%s\n", ip)
		remoteRecord.TTL = record.TTL
		remoteRecord.Content = ip
		fmt.Fprintf(os.Stderr, "Updating DNS [%s %s] record content with ip: %s\n", remoteRecord.Type, remoteRecord.Name, ip)
		client.UpdateRecord(&model.DNSRecordRequest{
            ID: remoteRecord.ID,
            ZoneID: remoteRecord.ZoneID,
            Name: remoteRecord.Name,
            Type: remoteRecord.Type,
            Content: ip,
            Proxied: remoteRecord.Proxied,
            TTL: record.TTL,
        })
		fmt.Fprintln(os.Stderr, "Updated!")
	}

    fmt.Fprintf(os.Stderr, "%s\n", ip)
    fmt.Fprintf(os.Stderr, "DNS [%s %s] content already contains: %s\n", remoteRecord.Type, record.Name, ip)

	return remoteRecord, nil
}

func CreateRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	zone, err := FindZone(client, record)
	if err != nil {
		return nil, err
	}

    fmt.Println("\nCreate DNS Record:")
    fmt.Printf("Zone ID: %20s\n", zone.ID)
    fmt.Printf("Name: %20s\n", record.Name)
    fmt.Printf("Content: %20s\n", record.IP)
    fmt.Printf("Type: %20s\n", record.Type)
    fmt.Printf("TTL: %20d\n", record.TTL)

	return client.NewRecord(&model.DNSRecordRequest{
		ZoneID:   zone.ID,
		Name:     record.Name,
		Content:  record.IP,
		Type:     string(record.Type),
		Proxied:  false,
		TTL:      record.TTL,
		Priority: 10,
	})
}

func FindZone(client DNSClient, record Record) (*model.Zone, error) {
	zones, err := client.ListZones()
	if err != nil {
		return nil, fmt.Errorf("Error while listing zones: %v\n", err)
	}

	zone := filterZoneByName(zones, record.ZoneName)
	if zone == nil {
		return nil, fmt.Errorf("No zone with name '%s' found\n", record.ZoneName)
	}

	return zone, nil
}

func FindRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
	zone, err := FindZone(client, record)
	if err != nil {
		return nil, err
	}

	fmt.Fprintf(os.Stderr, "Listing DNS Records for zone '%s' using id '%s' ...", zone.Name, zone.ID)
	records, err := client.ListRecords(zone.ID)
	if err != nil {
		return nil, fmt.Errorf("Error while listing dns records: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed dns records\n", len(records))

	fmt.Fprintf(os.Stderr, "Locating DNS Record: %s ...", record.Name)
	remoteRecord := filterByName(records, record.Name)
	if remoteRecord == nil {
		return nil, fmt.Errorf("No dns record with name '%s' found in zone '%s'\n", record.Name, record.ZoneName)
	}

	return remoteRecord, nil
}

func DeleteRecord(client DNSClient, record Record) (*model.DNSRecord, error) {
    return nil, nil
}
