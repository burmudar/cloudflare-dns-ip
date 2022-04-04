package dns

import (
	"cloudflare-dns/cloudflare"
	"cloudflare-dns/retrievers"
	"fmt"
	"net/http"
	"os"
	"time"
)

type Record struct {
    ZoneName string
    Name string
    IP string
    TTL int
}


func FindRecordByName(records []cloudflare.DNSRecord, name string) *cloudflare.DNSRecord {
	for _, r := range records {
		if r.Name == name {
			return &r
		}
	}

	return nil
}

func FindZoneByName(zones []cloudflare.Zone, name string) *cloudflare.Zone {
	for _, z := range zones {
		if z.Name == name {
			return &z
		}
	}

	return nil
}


func UpdateRecord(token string, record Record) error {
	client := cloudflare.NewTokenClient(token)

	fmt.Fprintf(os.Stderr, "Listing zones ...")
	zones, err := client.ListZones()
	if err != nil {
		return fmt.Errorf("Error while listing zones: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed zones\n", len(zones))

	fmt.Fprintf(os.Stderr, "Locating zone: %s ...", record.ZoneName)
	zone := FindZoneByName(zones, record.ZoneName)
	if zone == nil {
		return fmt.Errorf("No zone with name '%s' found\n", record.ZoneName)
	}
	fmt.Fprintf(os.Stderr, "FOUND!\n")

	fmt.Fprintf(os.Stderr, "Listing DNS Records for zone '%s' using id '%s' ...", zone.Name, zone.ID)
	records, err := client.ListDnsRecords(zone.ID)
	if err != nil {
		return fmt.Errorf("Error while listing dns records: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed dns records\n", len(records))

	fmt.Fprintf(os.Stderr, "Locatin DNS Record: %s ...", record.Name)
	remoteRecord := FindRecordByName(records, record.Name)
	if remoteRecord == nil {
		return fmt.Errorf("No dns record with name '%s' found in zone '%s'\n", record.Name, record.ZoneName)
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
