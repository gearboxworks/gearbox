package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
)

var databaseCmd = &cobra.Command{
	Use: "database",
	Aliases: []string{
		"db",
	},
	Short: "Manage project databases.",
}

func init() {
	var importExportHelp string = "" +
		"\n\t1.) The default filename and location specified in `project.json` if no parameters are specified, " +
		"\n\t2.) The default directory is only a filename is specified, " +
		"\n\t3.) The default filename is only a directory is specified, " +
		"\n\t4.) A specific file if a fullpath with filename is only a directory are specified." +
		"\n" +
		"\nDefaults can be found in `project.json`."

	RootCmd.AddCommand(databaseCmd)
	databaseCmd.AddCommand(&cobra.Command{
		Use:        "import",
		SuggestFor: []string{"restore"},
		Args:       cobra.RangeArgs(0, 1),
		Short:      "Import/Restore a database into a Gearbox project.",
		Long: "Import/Restore a database into the current Gearbox project. " +
			"You can import from:\n " + importExportHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Import a import into a Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:        "export",
		SuggestFor: []string{"backup", "dump"},
		Args:       cobra.RangeArgs(0, 1),
		Short:      "Export/Backup a database from a Gearbox project.",
		Long: "Export/Backup a database from the current Gearbox project. " +
			"You can export to:\n " + importExportHelp,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Import a export into a Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:   "chunk",
		Args:  cobra.ExactArgs(1),
		Short: "Unchunk the database for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Unchunk the database for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:   "unchunk",
		Args:  cobra.RangeArgs(0, 1),
		Short: "Unchunk the database for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Unchunk the database for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:   "credentials",
		Args:  cobra.NoArgs,
		Short: "Output all database credentials for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Output database credentials for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:   "name",
		Args:  cobra.NoArgs,
		Short: "Output database name for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Output database name for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:   "host",
		Args:  cobra.NoArgs,
		Short: "Output database host for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Output database host for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:        "user",
		SuggestFor: []string{"username"},
		Args:       cobra.NoArgs,
		Short:      "Output username for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Output username for the current Gearbox project goes here.")
		},
	})
	databaseCmd.AddCommand(&cobra.Command{
		Use:     "password",
		Aliases: []string{"pw"},
		Args:    cobra.NoArgs,
		Short:   "Output database password for the current Gearbox project.",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Output database password for the current Gearbox project goes here.")
		},
	})
}
