package tasks

import (
	"context"
	"errors"
	"time"

	raven "github.com/getsentry/raven-go"
	"github.com/lestrrat-go/backoff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/http"
	"github.com/DECODEproject/iotpolicystore/pkg/logger"
	"github.com/DECODEproject/iotpolicystore/pkg/version"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long: `Starts an implementation of the DECODE policy store RPC server.

This server listens on the specified port for incoming Protocol Buffer or
JSON RPC messages, with all data being persisted to the attached PostgreSQL
database.

Configuration values can be provided either by flags, or generally by
environment variables. If a flag is named: --example-flag, then it will also be
able to be supplied via an environment variable: $POLICYSTORE_EXAMPLE_FLAG`,
	RunE: func(cmd *cobra.Command, args []string) error {
		addr := viper.GetString("addr")
		if addr == "" {
			return errors.New("Must supply an address to bind to")
		}

		connStr := viper.GetString("database-url")
		if connStr == "" {
			return errors.New("Must supply a database connection string")
		}

		encryptionPassword := viper.GetString("encryption-password")
		if encryptionPassword == "" {
			return errors.New("Must supply a database encryption password")
		}

		dashboardURL := viper.GetString("dashboard-url")
		if dashboardURL == "" {
			return errors.New("Must supply a valid URL for the dashboard API")
		}

		clientTimeout := viper.GetInt("client-timeout")
		if clientTimeout == 0 {
			return errors.New("Must supply a non-zero timeout value")
		}

		config := &config.Config{
			ServerAddr:         addr,
			ConnStr:            connStr,
			EncryptionPassword: encryptionPassword,
			Verbose:            viper.GetBool("verbose"),
			Domains:            viper.GetStringSlice("domains"),
			DashboardURL:       dashboardURL,
			ClientTimeout:      clientTimeout,
			Logger:             logger.NewLogger(),
		}

		e := backoff.ExecuteFunc(func(_ context.Context) error {
			s := http.NewServer(config)
			return s.Start()
		})

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()

		policy := backoff.NewExponential()
		return backoff.Retry(ctx, policy, e)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringP("addr", "a", ":8082", "address to which the server binds")
	serverCmd.Flags().StringP("database-url", "d", "", "URL at which Postgres is listening (e.g. postgres://user:password@host:5432/dbname?sslmode=enable)")
	serverCmd.Flags().StringSlice("domains", []string{}, "comma separate list of domains for which we obtain TLS certificates")
	serverCmd.Flags().String("encryption-password", "", "password used to encrypt secret tokens we write to the database")
	serverCmd.Flags().String("dashboard-url", "", "URL at which the dashboard API is listening")
	serverCmd.Flags().Int("client-timeout", 10, "timeout in seconds for the http client")

	viper.BindPFlag("addr", serverCmd.Flags().Lookup("addr"))
	viper.BindPFlag("database-url", serverCmd.Flags().Lookup("database-url"))
	viper.BindPFlag("domains", serverCmd.Flags().Lookup("domains"))
	viper.BindPFlag("encryption-password", serverCmd.Flags().Lookup("encryption-password"))
	viper.BindPFlag("dashboard-url", serverCmd.Flags().Lookup("dashboard-url"))
	viper.BindPFlag("client-timeout", serverCmd.Flags().Lookup("client-timeout"))

	raven.SetRelease(version.Version)
	raven.SetTagsContext(map[string]string{"component": "policystore"})
}
