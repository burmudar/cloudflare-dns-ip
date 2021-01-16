package cmd

import (
	"cloudflare-dns/cloudflare"
	"cloudflare-dns/retrievers"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var token string
var zoneName string
var dnsRecordName string
var manualIP string
var ttlInSeconds int

func init() {
	updateCmd.PersistentFlags().StringVarP(&token, "token", "t", "", "Cloudflare API token")
	updateCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")
	updateCmd.PersistentFlags().StringVarP(&dnsRecordName, "dns-record-name", "r", "", "Name of the DNS record where the A record needs to be updated")
	updateCmd.PersistentFlags().IntVarP(&ttlInSeconds, "ttl", "", 3600, "TTL (in seconds) to set on the DNS record")
	updateCmd.PersistentFlags().StringVarP(&manualIP, "ip", "", "", "Set the content of the dns record to this ip")

	updateCmd.MarkPersistentFlagRequired("token")
	updateCmd.MarkPersistentFlagRequired("zone-name")
	updateCmd.MarkPersistentFlagRequired("dns-record-name")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a type A DNS record found in the given <zoneId> with the public IP",
	Long:  `Using the zone id the DNS record is retrieved and the content is updated to the latest public ip`,
	Run: func(cmd *cobra.Command, args []string) {
		updateDNSRecord(token, zoneName, dnsRecordName)
	},
}

func updateDNSRecord(token, zoneName, dnsRecordName string) error {
	dns := cloudflare.NewTokenClient(token)

	fmt.Fprintf(os.Stderr, "Listing zones ...")
	zones, err := dns.ListZones()
	if err != nil {
		return fmt.Errorf("Error while listing zones: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed zones\n", len(zones))

	fmt.Fprintf(os.Stderr, "Locating zone: %s ...", zoneName)
	zone := cloudflare.FindZoneByName(zones, zoneName)
	if zone == nil {
		return fmt.Errorf("No zone with name '%s' found\n", zoneName)
	}
	fmt.Fprintf(os.Stderr, "FOUND!\n")

	fmt.Fprintf(os.Stderr, "Listing DNS Records for zone '%s' using id '%s' ...", zone.Name, zone.ID)
	records, err := dns.ListDnsRecords(zone.ID)
	if err != nil {
		return fmt.Errorf("Error while listing dns records: %v\n", err)
	}
	fmt.Fprintf(os.Stderr, "%d listed dns records\n", len(records))

	fmt.Fprintf(os.Stderr, "Locatin DNS Record: %s ...", dnsRecordName)
	record := cloudflare.FindDNSRecordByName(records, dnsRecordName)
	if record == nil {
		return fmt.Errorf("No dns record with name '%s' found in zone '%s'\n", dnsRecordName, zoneName)
	}
	fmt.Fprintf(os.Stderr, "FOUND!\n")

	var retriever retrievers.StringRetriever
	if manualIP != "" {
		retriever = retrievers.NewStaticStringRetriever(manualIP)
		fmt.Fprintf(os.Stderr, "Manually setting ip ...")
	} else {
		retriever = retrievers.NewIPRetriever(http.DefaultClient, "https://ifconfig.co", 30*time.Second)
		fmt.Fprintf(os.Stderr, "Discovering public ip ...")
	}

	if ip, err := retriever.Get(); err != nil {
		return fmt.Errorf("Failed to get ip: %v\n", err)
	} else if ip != record.Content {
		fmt.Fprintf(os.Stderr, "%s\n", ip)
		record.TTL = ttlInSeconds
		record.Content = ip
		fmt.Fprintf(os.Stderr, "Updating DNS [%s %s] record content with ip: %s\n", record.Type, record.Name, ip)
		dns.UpdateDNSRecord(record)
		fmt.Fprintln(os.Stderr, "Updated!")
	} else {
		fmt.Fprintf(os.Stderr, "%s\n", ip)
		fmt.Fprintf(os.Stderr, "DNS [%s %s] content already contains: %s\n", record.Type, record.Name, ip)
	}

	return nil
}
