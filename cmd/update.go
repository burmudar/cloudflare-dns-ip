package cmd

import (
	"github.com/burmudar/cloudflare-dns/dns"

	"github.com/spf13/cobra"
)

func init() {
	updateCmd.PersistentFlags().StringVarP(&zoneName, "zone-name", "z", "", "Name of the Zone the DNS record resides in")
	updateCmd.PersistentFlags().IntVarP(&ttlInSeconds, "ttl", "", 3600, "TTL (in seconds) to set on the DNS record")
	updateCmd.PersistentFlags().StringVarP(&manualIP, "ip", "", "", "Set the content of the dns record to this ip")
	updateCmd.PersistentFlags().StringSliceVarP(&recordNames, "dns-record-names", "r", recordNames, "Name of one or more DNS records. If more than one record is specified separated them with a comma")

	updateCmd.MarkPersistentFlagRequired("zone-name")
	updateCmd.MarkPersistentFlagRequired("dns-record-names")
	rootCmd.AddCommand(updateCmd)
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a type A DNS record found in the given <zoneId> with the public IP",
	Long:  `Using the zone id the DNS record is retrieved and the content is updated to the latest public ip`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := createClient()
		if err != nil {
			return err
		}

		for _, name := range recordNames {
			dns.UpdateRecord(client, dns.Record{
				ZoneName: zoneName,
				Name:     dns.NormaliseRecordName(zoneName, name),
				TTL:      ttlInSeconds,
				IP:       manualIP,
			})
		}

		return nil
	},
}
