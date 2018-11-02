package tasks

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/DECODEproject/iotpolicystore/pkg/logger"
	"github.com/DECODEproject/iotpolicystore/pkg/postgres"
	"github.com/DECODEproject/iotpolicystore/pkg/version"
)

var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Manage Postgres migrations",
	Long: `This task provides a set of subcommands for working with Postgres migrations.

Up migrations are run automatically when the application boots, but here we
also add some commands to create properly named migration files, and a
command to roll back migrations.`,
}

var migrateNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new Postgres migration",
	Long: fmt.Sprintf(`This command allows the caller to create new pairs of migration files, named
in the way the migration library requires.

The desired migration name should be passed via a positional argument after
the new subcommand. For example:

	$ %s migrate new AddPolicyTable`, version.BinaryName),
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := viper.GetString("dir")
		if dir == "" {
			return errors.New("Must supply a directory")
		}

		logger := logger.NewLogger()

		return postgres.NewMigration(dir, args[0], logger)
	},
}

var migrateDownCmd = &cobra.Command{
	Use:   "down",
	Short: "Run down migrations against Postgres",
	Long: `This command can be used to rollback previously executed migrations against
Postgres.

It takes as parameters the number of steps to rollback (default of 1), or a
boolean flag (--all) indicating we should attempt to rollback all
migrations.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		connStr := viper.GetString("database-url")
		if connStr == "" {
			return errors.New("Must provide a database connection string")
		}

		steps, err := cmd.Flags().GetInt("steps")
		if err != nil {
			return err
		}

		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			return err
		}

		logger := logger.NewLogger()

		db, err := postgres.Open(connStr)
		if err != nil {
			return err
		}

		if all {
			return postgres.MigrateDownAll(db, logger)
		}

		return postgres.MigrateDown(db, steps, logger)
	},
}

func init() {
	rootCmd.AddCommand(migrateCmd)
	migrateCmd.AddCommand(migrateNewCmd)
	migrateCmd.AddCommand(migrateDownCmd)

	migrateNewCmd.Flags().String("dir", "pkg/migrations/sql", "The directory into which new migrations will be created")
	migrateDownCmd.Flags().IntP("steps", "s", 1, "Number of down migrations to run")
	migrateDownCmd.Flags().Bool("all", false, "Boolean flag that if true runs all down migrations")

	viper.BindPFlag("dir", migrateNewCmd.Flags().Lookup("dir"))
}
