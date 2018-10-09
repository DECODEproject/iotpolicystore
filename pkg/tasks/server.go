package tasks

import (
	"errors"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thingful/iotpolicystore/pkg/http"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long: `
Starts an implementation of the DECODE policy store RPC server.

This server listens on the specified port for incoming Protocol Buffer or
JSON RPC messages, with all data being persisted to the attached PostgreSQL
database.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := viper.GetString("addr")
		if addr == "" {
			return errors.New("Must provide a bind address for the")
		}

		databaseURL := viper.GetString("database-url")
		if databaseURL == "" {
			return errors.New("Must provide a database connection url")
		}

		verbose := viper.GetBool("verbose")

		config := &http.Config{
			Addr:        addr,
			DatabaseURL: databaseURL,
			Verbose:     verbose,
		}

		s := http.NewServer(config)

		s.Start()

		return nil
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("addr", "a", "0.0.0.0:8082", "Specify the address to which the server binds. May also be passed via the environment as $POLICY_STORE_ADDR")
	serverCmd.Flags().StringP("database-url", "d", "", "The database connection url. May also be passed via the environment as $POLICYSTORE_DATABASE_URL")

	viper.BindPFlag("addr", serverCmd.Flags().Lookup("addr"))
	viper.BindPFlag("database-url", serverCmd.Flags().Lookup("database-url"))
}
