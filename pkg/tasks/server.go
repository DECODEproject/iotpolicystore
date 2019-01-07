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

		hashidLength := viper.GetInt("hashid-length")
		if hashidLength == 0 {
			return errors.New("Must specify a minimum hashid length greater than 0")
		}

		hashidSalt := viper.GetString("hashid-salt")
		if hashidSalt == "" {
			return errors.New("Must supply a hashid salt value")
		}

		certFile := viper.GetString("cert-file")
		keyFile := viper.GetString("key-file")
		verbose := viper.GetBool("verbose")

		config := &config.Config{
			ServerAddr:         addr,
			ConnStr:            connStr,
			EncryptionPassword: encryptionPassword,
			HashidLength:       hashidLength,
			HashidSalt:         hashidSalt,
			Verbose:            verbose,
			CertFile:           certFile,
			KeyFile:            keyFile,
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

	serverCmd.Flags().StringP("addr", "a", "0.0.0.0:8082", "address to which the server binds")
	serverCmd.Flags().StringP("cert-file", "c", "", "path to a TLS certificate file to enable TLS on the server")
	serverCmd.Flags().StringP("database-url", "d", "", "URL at which Postgres is listening (e.g. postgres://user:password@host:5432/dbname?sslmode=enable)")
	serverCmd.Flags().StringP("key-file", "k", "", "path to a TLS private key file to enable TLS on the server")
	serverCmd.Flags().String("encryption-password", "", "password used to encrypt secret tokens we write to the database")
	serverCmd.Flags().IntP("hashid-length", "l", 8, "minimum length of generated id strings for policies")
	serverCmd.Flags().String("hashid-salt", "", "salt value used when generating IDs for policies")

	viper.BindPFlag("addr", serverCmd.Flags().Lookup("addr"))
	viper.BindPFlag("cert-file", serverCmd.Flags().Lookup("cert-file"))
	viper.BindPFlag("database-url", serverCmd.Flags().Lookup("database-url"))
	viper.BindPFlag("key-file", serverCmd.Flags().Lookup("key-file"))
	viper.BindPFlag("encryption-password", serverCmd.Flags().Lookup("encryption-password"))
	viper.BindPFlag("hashid-length", serverCmd.Flags().Lookup("hashid-length"))
	viper.BindPFlag("hashid-salt", serverCmd.Flags().Lookup("hashid-salt"))

	raven.SetRelease(version.Version)
	raven.SetTagsContext(map[string]string{"component": "policystore"})
}
