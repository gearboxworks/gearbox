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
	Short: "Connect to the terminal of the running GearboxOS VM",
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

	vmCmd.PersistentFlags().StringVarP(&sshArgs.SshUsername,"user", "u", gearbox.SshDefaultUsername, "Alternate Gearbox SSH username.")
	vmCmd.PersistentFlags().StringVarP(&sshArgs.SshPassword,"password", "p", gearbox.SshDefaultPassword, "Alternate Gearbox SSH password.")
	vmCmd.PersistentFlags().StringVarP(&sshArgs.SshPublicKey,"key-file", "k", gearbox.SshDefaultKeyFile, "Gearbox SSH public key file.")

}
