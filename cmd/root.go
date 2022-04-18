package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var token string
var zoneName string
var dnsRecordName string
var manualIP string
var ttlInSeconds int

var rootCmd = &cobra.Command{
	Use:   "cloudfare-dns",
	Short: "Cloudfare DNS updates specific dns records with public ips",
	Long: `A Personal utility used by @burmudar to update various machines he has in his apartment
                Code at github.com/burmudar/cloudflare-dns-ip`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
