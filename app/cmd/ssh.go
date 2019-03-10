package cmd

import (
	"gearbox"
	"github.com/spf13/cobra"
)

var sshCmd = &cobra.Command{
	Use: "ssh",
	SuggestFor: []string{
		"login",
		"logon",
		"access",
	},
	Short: "Connect to the terminal of the Box running GearboxOS",
	Run: func(cmd *cobra.Command, args []string) {
		err := gearbox.Instance.ConnectSSH(sshArgs)
		if err != nil {
			//fmt.Printf("%s", err)
		}
	},
}

var sshArgs gearbox.SshArgs

func init() {

	RootCmd.AddCommand(sshCmd)

	sshCmd.PersistentFlags().StringVarP(&sshArgs.Username, "user", "u", gearbox.SshDefaultUsername, "Alternate Gearbox SSH username.")
	sshCmd.PersistentFlags().StringVarP(&sshArgs.Password, "password", "p", gearbox.SshDefaultPassword, "Alternate Gearbox SSH password.")
	sshCmd.PersistentFlags().StringVarP(&sshArgs.PublicKey, "key-file", "k", gearbox.SshDefaultKeyFile, "Gearbox SSH public key file.")
	sshCmd.PersistentFlags().BoolVarP(&sshArgs.StatusLine.Disable, "no-status", "", false, "Disable Gearbox status line.")

}
