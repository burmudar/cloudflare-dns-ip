package dns

import (
	"cloudflare-dns/cloudflare"
	"cloudflare-dns/retrievers"
	"fmt"
	"net/http"
	"os"
	"time"
)

type ZoneType string

const (
    AType ZoneType = "A"
)

type Record struct {
    ZoneName string
    Type ZoneType
    Name string
    IP string
    TTL int
}


func filterByName(records []cloudflare.DNSRecordResponse, name string) *cloudflare.DNSRecordResponse {
	for _, r := range records {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

func filterZoneByName(zones []cloudflare.ZoneResponse, name string) *cloudflare.ZoneResponse {
	for _, z := range zones {
		if z.Name == name {
			return &z
		}
	}

	return nil
}


func UpdateRecord(client *cloudflare.Client, record Record) error {
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
		retriever = retrievers.NewIPRetriever(http.DefaultClient, "https://ifconfig.co", 30*time.Second)
		fmt.Fprintf(os.Stderr, "Discovering public ip ...")
	}

	if ip, err := retriever.Get(); err != nil {
		return fmt.Errorf("Failed to get ip: %v\n", err)
	} else if ip != remoteRecord.Content {
		fmt.Fprintf(os.Stderr, "%s\n", ip)
		remoteRecord.TTL = record.TTL
		remoteRecord.Content = ip
		fmt.Fprintf(os.Stderr, "Updating DNS [%s %s] record content with ip: %s\n", remoteRecord.Type, remoteRecord.Name, ip)
		client.UpdateDNSRecord(remoteRecord)
		fmt.Fprintln(os.Stderr, "Updated!")
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", ip)
		fmt.Fprintf(os.Stderr, "DNS [%s %s] content already contains: %s\n", remoteRecord.Type, record.Name, ip)
	}

	return nil
}

func CreateRecord(client *cloudflare.Client, record Record) error {
    zone, err := FindZone(client, record)
    if err != nil {
        return err
    }

    resp, err := client.NewDNSRecord(&cloudflare.DNSRecordRequest{
        ZoneID: zone.ID,
        Name: record.Name,
        Content: record.IP,
        Type: string(record.Type),
        Proxied: false,
        TTL: record.TTL,
        Priority: 10,
    })

    if data, err := resp.ok() {
    }
    return nil
}

func FindZone(client *cloudflare.Client, record Record) (*cloudflare.ZoneResponse, error) {
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


func FindRecord(client *cloudflare.Client, record Record) (*cloudflare.DNSRecordResponse, error) {

	fmt.Fprintf(os.Stderr, "Listing DNS Records for zone '%s' using id '%s' ...", zone.Name, zone.ID)
	records, err := client.ListDnsRecords(zone.ID)
	if err != nil {
		return nil, fmt.Errorf("Error while listing dns records: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed dns records\n", len(records))

	fmt.Fprintf(os.Stderr, "Locatin DNS Record: %s ...", record.Name)
	remoteRecord := filterByName(records, record.Name)
	if remoteRecord == nil {
		return nil, fmt.Errorf("No dns record with name '%s' found in zone '%s'\n", record.Name, record.ZoneName)
	}

    return remoteRecord, nil
}
