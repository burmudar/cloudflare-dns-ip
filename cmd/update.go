package cmd

import (
	"cloudflare-dns/dns"
	"cloudflare-dns/dns/cloudflare"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	updateCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")
	updateCmd.PersistentFlags().StringVarP(&dnsRecordName, "dns-record-name", "r", "", "Name of the DNS record where the A record needs to be updated")
	updateCmd.PersistentFlags().IntVarP(&ttlInSeconds, "ttl", "", 3600, "TTL (in seconds) to set on the DNS record")
	updateCmd.PersistentFlags().StringVarP(&manualIP, "ip", "", "", "Set the content of the dns record to this ip")

	updateCmd.MarkPersistentFlagRequired("zone-name")
	updateCmd.MarkPersistentFlagRequired("dns-record-name")
	rootCmd.AddCommand(deleteCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a type A DNS record found in the given <zoneId> with the public IP",
	Long:  `Using the zone id the DNS record is retrieved and the content is updated to the latest public ip`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := cloudflare.NewTokenClient(cloudflare.API_CLOUDFLARE_V4, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create cloudflare client: %v", err)
		}
		dns.UpdateRecord(client, dns.Record{
			ZoneName: zoneName,
			Name:     dnsRecordName,
			TTL:      ttlInSeconds,
			IP:       manualIP,
		})
	},
}
