package tasks

import (
	"context"
	"fmt"
	"time"

	"github.com/lestrrat-go/backoff"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/thingful/iotpolicystore/pkg/config"
	"github.com/thingful/iotpolicystore/pkg/http"
	"github.com/thingful/iotpolicystore/pkg/logger"
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
		// has a default so will always be set
		addr := viper.GetString("addr")

		connStr := viper.GetString("database-url")
		if connStr == "" {
			return missingConfig("database connection url", "--database-url", "$POLICYSTORE_DATABASE_URL")
		}

		encryptionPassword := viper.GetString("encryption-password")
		if encryptionPassword == "" {
			return missingConfig("database encryption password", "--encryption-password", "$POLICYSTORE_ENCRYPTION_PASSWORD")
		}

		hashidLength := viper.GetInt("hashid-length")
		if hashidLength == 0 {
			return missingConfig("hashid min length", "--hashid-length", "$POLICYSTORE_HASHID_LENGTH")
		}

		hashidSalt := viper.GetString("hashid-salt")
		if hashidSalt == "" {
			return missingConfig("hashid salt", "--hashid-salt", "$POLICYSTORE_HASHID_SALT")
		}

		verbose := viper.GetBool("verbose")

		config := &config.Config{
			ServerAddr:         addr,
			ConnStr:            connStr,
			EncryptionPassword: encryptionPassword,
			HashidLength:       hashidLength,
			HashidSalt:         hashidSalt,
			Verbose:            verbose,
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
	serverCmd.Flags().StringP("database-url", "d", "", "The database connection url. May also be passed via the environment as $POLICYSTORE_DATABASE_URL")
	serverCmd.Flags().String("encryption-password", "", "Encryption password used to when encrypting policy tokens in the DB. May also be passed via the environment as $POLICYSTORE_ENCRYPTION_PASSWORD")
	serverCmd.Flags().Int("hashid-length", 8, "Minimum length of generated id strings. May also be passed via the environment as $POLICYSTORE_HASHID_LENGTH")
	serverCmd.Flags().String("hashid-salt", "", "Salt value for generated id strings. May also be passed via the environment as $POLICYSTORE_HASHID_SALT")

	viper.BindPFlag("addr", serverCmd.Flags().Lookup("addr"))
	viper.BindPFlag("database-url", serverCmd.Flags().Lookup("database-url"))
	viper.BindPFlag("encryption-password", serverCmd.Flags().Lookup("encryption-password"))
	viper.BindPFlag("hashid-length", serverCmd.Flags().Lookup("hashid-length"))
	viper.BindPFlag("hashid-salt", serverCmd.Flags().Lookup("hashid-salt"))
}

// missingConfig builds an error for a missing required config value.
func missingConfig(name, flagName, envVarName string) error {
	return fmt.Errorf("Missing required config value: %s. May either be passed as a flag: `%s`, or via the `%s` environment variable", name, flagName, envVarName)
}
