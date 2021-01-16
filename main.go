package main

import (
	"cloudflare-dns/cloudflare"
	"cloudflare-dns/retrievers"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	token, err := ioutil.ReadFile("TOKEN")
	dns := cloudflare.NewTokenClient(string(token))

	zones, err := dns.ListZones()
	if err != nil {
		fmt.Printf("Error while listing zones: %v\n", err)
	}

	zone := cloudflare.FindZoneByName(zones, "burmudar.dev")

	records, err := dns.ListDnsRecords(zone.ID)
	if err != nil {
		fmt.Printf("Error while listing dns records: %v\n", err)
	}
	fmt.Printf("Records: %v\n", records)

	record := cloudflare.FindDNSRecordByName(records, "media.burmudar.dev")
	if record == nil {
		panic("failed to find required record: media.burmduar.dev")
	}

	retriever := retrievers.NewIPRetriever(http.DefaultClient, "https://ifconfig.co", 30*time.Second)

	if ip, err := retriever.Get(); err != nil {
		fmt.Printf("failed to get ip: %v\n", err)
	} else if ip != record.Content {
		fmt.Printf("\n'%s' == '%s'\n", ip, record.Content)
		record.TTL = int((1 * time.Hour).Seconds())
		record.Content = ip
		dns.UpdateDNSRecord(record)
	} else {
		fmt.Printf("DNS Already updated: %v\n", record)
	}

}
