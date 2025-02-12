package cmd

import (
	"github.com/sighup-io/furyagent/pkg/component"
	"github.com/spf13/cobra"
	"log"
)

// restoreCmd represents the `furyctl restore` subcommand
var restoreCmd = &cobra.Command{
	Use:   "restore",
	Short: "Executes restores",
	Long:  ``,
}

// etcdRestoreCmd represents the `furyctl restore etcd` command
var etcdRestoreCmd = &cobra.Command{
	Use:   "etcd",
	Short: "Restores etcd node",
	Long:  `Restores etcd node`,
	Run: func(cmd *cobra.Command, args []string) {
		var etcd component.ClusterComponent = component.Etcd{data}
		err := etcd.Restore()
		if err != nil {
			log.Fatal(err)
		}
	},
}

// masterBackupCmd represents the `furyctl restore master` command
var masterRestoreCmd = &cobra.Command{
	Use:   "master",
	Short: "Restores master node",
	Long:  `Restores master node`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(restoreCmd)
	restoreCmd.AddCommand(etcdRestoreCmd)
	restoreCmd.AddCommand(masterRestoreCmd)
}
