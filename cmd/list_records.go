package cmd

import (
	"cloudflare-dns/dns"
	"cloudflare-dns/dns/cloudflare"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	listRecordCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")

	rootCmd.AddCommand(listRecordCmd)
}

var listRecordCmd = &cobra.Command{
	Use:   "list-records",
	Short: "list DNS records present in zone <zoneName>",
	Long:  `Using the <zoneName> all the DNS records registered for the zone are fetched`,
	Run: func(cmd *cobra.Command, args []string) {
		client, err := cloudflare.NewTokenClient(cloudflare.API_CLOUDFLARE_V4, token)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create cloudflare client: %v", err)
		}
		fmt.Fprintf(os.Stderr, "Listing records found in zone '%s':\n", zoneName)
		records, err := dns.ListRecords(client, zoneName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error locating records in zone: %v\n", err)
		}
		for _, record := range records {
			fmt.Fprintf(os.Stdout, "%v\n", record)
		}
	},
}
