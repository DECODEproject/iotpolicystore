package tasks

import (
	"log"
	"strings"

	raven "github.com/getsentry/raven-go"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DECODEproject/iotpolicystore/pkg/version"
)

var rootCmd = &cobra.Command{
	Use:   version.BinaryName,
	Short: "Policy store service for the DECODE IoT Pilot",
	Long: `This tool is an implementation of the policy store service being
 developed as part of the IoT pilot for the DECODE project
 (https://decodeproject.eu).

This service exposes a simple RPC API implemented using a tool called Twirp
that provides both a JSON/HTTP and a more performant Protocol Buffer/HTTP API
to clients.append

Entitlement policy data is persisted to a PostgreSQL database, and the server
exposes methods to write, delete and read these entitlement policies.
Encryption is provided by other components within the pilot.`,
	Version: version.VersionString(),
}

// Execute is our main entry point - called by main.go in cmd/policystore
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		raven.CaptureErrorAndWait(err, nil)
		log.Fatal(err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().Bool("verbose", false, "enable verbose output")
	viper.BindPFlag("verbose", rootCmd.PersistentFlags().Lookup("verbose"))
}

// initConfig sets up any stuff for viper to do with reading from environment
// variables for example.
func initConfig() {
	viper.SetEnvPrefix("policystore")
	viper.AutomaticEnv()
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
}
