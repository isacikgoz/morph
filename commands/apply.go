package commands

import (
	"context"

	"github.com/isacikgoz/morph"
	"github.com/isacikgoz/morph/apply"
	"github.com/spf13/cobra"
)

func ApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "apply",
		Short: "Applies migrations",
	}

	cmd.PersistentFlags().StringP("driver", "d", "", "the database driver of the migrations")
	_ = cmd.MarkPersistentFlagRequired("driver")
	cmd.PersistentFlags().String("dsn", "", "the dsn of the database")
	_ = cmd.MarkPersistentFlagRequired("dsn")

	cmd.PersistentFlags().StringP("path", "p", "", "the source path of the migrations")
	_ = cmd.MarkPersistentFlagRequired("path")

	cmd.PersistentFlags().IntP("timeout", "t", 60, "the timeout in seconds for each migration file to run")
	cmd.PersistentFlags().StringP("migrations-table", "m", "db_migrations", "the name of the migrations table")
	cmd.PersistentFlags().StringP("lock-key", "l", "mutex_migrations", "the name of the mutex key")

	cmd.AddCommand(
		UpApplyCmd(),
		DownApplyCmd(),
		MigrateApplyCmd(),
	)

	return cmd
}

func UpApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "up",
		Short:         "Apply migrations forward a number of steps",
		RunE:          upApplyCmdF,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
	cmd.Flags().Int("number", 0, "apply N up migrations")
	return cmd
}

func DownApplyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "down",
		Short:         "Apply migrations backwards a number of steps",
		RunE:          downApplyCmdF,
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.Flags().Int("number", 0, "apply N down migrations")
	return cmd
}

func MigrateApplyCmd() *cobra.Command {
	return &cobra.Command{
		Use:           "migrate",
		Short:         "Apply all migrations",
		RunE:          migrateApplyCmdF,
		SilenceUsage:  true,
		SilenceErrors: false,
	}
}

func upApplyCmdF(cmd *cobra.Command, _ []string) error {
	dsn, _ := cmd.Flags().GetString("dsn")
	driverName, _ := cmd.Flags().GetString("driver")
	path, _ := cmd.Flags().GetString("path")
	steps, _ := cmd.Flags().GetInt("number")
	timeout, _ := cmd.Flags().GetInt("timeout")
	tableName, _ := cmd.Flags().GetString("migrations-table")
	mutexKey, _ := cmd.Flags().GetString("lock-key")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	morph.InfoLogger.Printf("Attempting to apply %d migrations...\n", steps)
	n, err := apply.Up(ctx, steps, dsn, driverName, path, morph.SetMigrationTableName(tableName), morph.SetSatementTimeoutInSeconds(timeout), morph.WithLock(mutexKey))
	if n > 0 {
		morph.SuccessLogger.Printf("%d migrations applied.\n", n)
	} else if n == 0 {
		morph.InfoLogger.Println("no migrations applied.")
	}
	return err
}

func downApplyCmdF(cmd *cobra.Command, _ []string) error {
	dsn, _ := cmd.Flags().GetString("dsn")
	driverName, _ := cmd.Flags().GetString("driver")
	path, _ := cmd.Flags().GetString("path")
	steps, _ := cmd.Flags().GetInt("number")
	timeout, _ := cmd.Flags().GetInt("timeout")
	tableName, _ := cmd.Flags().GetString("migrations-table")
	mutexKey, _ := cmd.Flags().GetString("lock-key")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	morph.InfoLogger.Printf("Attempting to apply  %d migrations...\n", steps)
	n, err := apply.Down(ctx, steps, dsn, driverName, path, morph.SetMigrationTableName(tableName), morph.SetSatementTimeoutInSeconds(timeout), morph.WithLock(mutexKey))
	if n > 0 {
		morph.SuccessLogger.Printf("%d migrations applied.\n", n)
	} else if n == 0 {
		morph.InfoLogger.Println("no migrations applied.")
	}
	return err
}

func migrateApplyCmdF(cmd *cobra.Command, _ []string) error {
	dsn, _ := cmd.Flags().GetString("dsn")
	driverName, _ := cmd.Flags().GetString("driver")
	path, _ := cmd.Flags().GetString("path")
	timeout, _ := cmd.Flags().GetInt("timeout")
	tableName, _ := cmd.Flags().GetString("migrations-table")
	mutexKey, _ := cmd.Flags().GetString("lock-key")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	morph.InfoLogger.Println("Applying all pending migrations...")
	if err := apply.Migrate(ctx, dsn, driverName, path, morph.SetMigrationTableName(tableName), morph.SetSatementTimeoutInSeconds(timeout), morph.WithLock(mutexKey)); err != nil {
		return err
	}
	morph.SuccessLogger.Println("Pending migrations applied.")

	return nil
}
