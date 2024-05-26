package cmd

import (
	"log"

	"github.com/spf13/cobra"
)

var (
	bootstrapServer string
	saslEnabled     bool
	saslUser        string
	saslPassword    string
	saslMechanism   string
	saslProtocol    string
	file            string
	topic           string
	resouceName     string
)

var rootCmd = &cobra.Command{
	Use:   "kafkactl",
	Short: "kafkactl is a command line tool for interacting with kafka",
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Error executing kafkactl: %v", err)
	}
}

// init flag bootstrap, username password for authen to kafka
func init() {
	rootCmd.PersistentFlags().StringVar(&bootstrapServer, "bootstrap-server", "", "bootstrap server")
	rootCmd.PersistentFlags().BoolVar(&saslEnabled, "sasl-enable", false, "Enable sasl authentication")
	rootCmd.PersistentFlags().StringVar(&saslUser, "sasl-user", "", "sasl username")
	rootCmd.PersistentFlags().StringVar(&saslPassword, "sasl-password", "", "sasl password")
	rootCmd.PersistentFlags().StringVar(&saslMechanism, "sasl-mechanism", "", "sasl mechanism")
	rootCmd.PersistentFlags().StringVar(&saslProtocol, "sasl-protocol", "", "sasl protocol")
	rootCmd.PersistentFlags().StringVar(&file, "file", "", "file permissions")
	rootCmd.PersistentFlags().StringVar(&topic, "topic", "", "name topic")
	rootCmd.PersistentFlags().StringVar(&resouceName, "resourcename", "", "resoucename")
}
