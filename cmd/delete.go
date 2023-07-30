package cmd

import (
	"fmt"
	"os"

	"github.com/burmudar/cloudflare-dns/dns"

	"github.com/spf13/cobra"
)

func init() {
	deleteCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")
	deleteCmd.PersistentFlags().StringSliceVarP(&recordNames, "dns-record-name", "r", recordNames, "Name of the DNS record")

	deleteCmd.MarkPersistentFlagRequired("zone-name")
	deleteCmd.MarkPersistentFlagRequired("dns-record-name")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete the DNS record with <dns-record-name>",
	Long:  `Delete the DNS record with <dns-record-name> that is inside zone with <zone-name>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return err
		}

		hasErrs := false

		for _, name := range recordNames {
			result, err := dns.DeleteRecord(client, dns.Record{
				ZoneName: zoneName,
				Name:     dns.NormaliseRecordName(zoneName, name),
			})

			if err != nil {
				hasErrs = true
				fmt.Fprintf(os.Stderr, "error deleting dns record %s. %v", name, err)
			} else {
				fmt.Fprintf(os.Stderr, "--- DNS '%s' record deleted ---\n%s\n", name, result)
			}
		}

		if hasErrs {
			return fmt.Errorf("One or more records failed to delete")
		}

		return nil
	},
}
