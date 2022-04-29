package cmd

import (
	"cloudflare-dns/dns"
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
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create cloudflare client: %v", err)
		}

        fmt.Fprintf(os.Stderr, "--- Listing records in zone '%s' ---\n", zoneName)
		records, err := dns.ListRecords(client, zoneName)
		if err != nil {
			return err
		}
		for _, record := range records {
			fmt.Fprintf(os.Stdout, "--- %s ---\n%s\n", record.Name, record.String())
		}

		return nil
	},
}
