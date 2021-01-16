package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cloudfare-dns",
	Short: "Cloudfare DNS updates specific dns records with public ips",
	Long: `A Personal utility used by @Burmudar to update various machines he has in his apartment
                Code at github.com/Burmudar/cloudflare-dns-ip`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
