package cmd

import (
	"cloudflare-dns/dns"
	"cloudflare-dns/dns/cloudflare"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
)

var tokenPath string
var zoneName string
var recordNames []string = make([]string, 0)
var manualIP string
var ttlInSeconds int

var rootCmd = &cobra.Command{
	Use:   "cloudfare-dns",
	Short: "Cloudfare DNS updates specific dns records with public ips",
	Long: `A Personal utility used by @burmudar to update various machines he has in his apartment
                Code at github.com/burmudar/cloudflare-dns-ip`,
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&tokenPath, "token", "t", "", "Cloudflare API token file")
	rootCmd.MarkPersistentFlagRequired("token")
}

func createClient() (dns.DNSClient, error) {
	token, err := readTokenFile(tokenPath)
	if err != nil {
		return nil, err
	}

	return cloudflare.NewTokenClient(cloudflare.API_CLOUDFLARE_V4, string(token))
}

func readTokenFile(path string) ([]byte, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	validPerm := info.Mode() == 0o644 || info.Mode() == 0o600 || info.Mode() == 0o444 || info.Mode() == 0o400
	if !validPerm {
		return nil, fmt.Errorf("invalid permissions %s", info.Mode())
	}

	return ioutil.ReadFile(path)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
