package tasks

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/lestrrat-go/backoff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DECODEproject/iotpolicystore/pkg/config"
	"github.com/DECODEproject/iotpolicystore/pkg/http"
	"github.com/DECODEproject/iotpolicystore/pkg/logger"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long: `Starts an implementation of the DECODE policy store RPC server.

This server listens on the specified port for incoming Protocol Buffer or
JSON RPC messages, with all data being persisted to the attached PostgreSQL
database.

The server requires the following environment variables to be set in order to
operate:

* $POLICYSTORE_DATABASE_URL - connection string for the database
* $POLICYSTORE_ENCRYPTION_PASSWORD - password used to encrypt tokens in the db
* $POLICYSTORE_HASHID_SALT - salt value used when hashing ids
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// has a default so will always be set
		addr := viper.GetString("addr")

		connStr := viper.GetString("database-url")
		if connStr == "" {
			return missingConfig("database connection url", "$POLICYSTORE_DATABASE_URL")
		}

		encryptionPassword := viper.GetString("encryption-password")
		if encryptionPassword == "" {
			return missingConfig("database encryption password", "$POLICYSTORE_ENCRYPTION_PASSWORD")
		}

		hashidLength := viper.GetInt("hashid-length")
		if hashidLength == 0 {
			return errors.New("Must specify a minimum hashid length greater than 0")
		}

		hashidSalt := viper.GetString("hashid-salt")
		if hashidSalt == "" {
			return missingConfig("hashid salt", "$POLICYSTORE_HASHID_SALT")
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

	serverCmd.Flags().StringP("addr", "a", "0.0.0.0:8082", "Specify the address to which the server binds. May also be passed via the environment as $POLICYSTORE_ADDR")
	serverCmd.Flags().Int("hashid-length", 8, "Minimum length of generated id strings. May also be passed via the environment as $POLICYSTORE_HASHID_LENGTH")
	serverCmd.Flags().StringP("cert-file", "c", "", "The path to a TLS certificate file to enable TLS on the server (or $POLICYSTORE_CERT_FILE env variable)")
	serverCmd.Flags().StringP("key-file", "k", "", "The path to a TLS private key file to enable TLS on the server (or $POLICYSTORE_KEY_FILE env variable)")

	viper.BindPFlag("addr", serverCmd.Flags().Lookup("addr"))
	viper.BindPFlag("hashid-length", serverCmd.Flags().Lookup("hashid-length"))
	viper.BindPFlag("cert-file", serverCmd.Flags().Lookup("cert-file"))
	viper.BindPFlag("key-file", serverCmd.Flags().Lookup("key-file"))
}

// missingConfig builds an error for a missing required config value.
func missingConfig(name, envVarName string) error {
	return fmt.Errorf("Missing required config value: %s. Must be passed via the `%s` environment variable", name, envVarName)
}
