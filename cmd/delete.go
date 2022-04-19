package cmd

import (
	"cloudflare-dns/dns"
	"cloudflare-dns/dns/cloudflare"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
        deleteCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")
	deleteCmd.PersistentFlags().StringSliceP("dns-record-name", "r", recordNames, "Name of the DNS record")

	deleteCmd.MarkPersistentFlagRequired("zone-name")
	deleteCmd.MarkPersistentFlagRequired("dns-record-name")
	rootCmd.AddCommand(deleteCmd)
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "delete the DNS record with <dns-record-name>",
	Long:  `Delete the DNS record with <dns-record-name> that is inside zone with <zone-name>`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := cloudflare.NewTokenClient(cloudflare.API_CLOUDFLARE_V4, tokenPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to create cloudflare client: %v", err)
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
				fmt.Fprintf(os.Stderr, "'%s' record delete\n%v\n", recordNames, result)
			}
		}

		if hasErrs {
			return fmt.Errorf("One or more records failed to delete")
		}

		return nil
	},
}
