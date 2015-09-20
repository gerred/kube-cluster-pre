package cli

import "github.com/spf13/cobra"

// KubeClusterCmd is the root command. Attach all other commands to this.
var KubeClusterCmd = &cobra.Command{
	Use:   "kube-cluster",
	Short: "kube-cluster provisions, scales, and manages kubernetes environments",
	Long:  "kube-cluster provisions, scales, and manages kubernetes environments",
	Run: func(cmd *cobra.Command, args []string) {
	},
}

// Execute adds all child commands to the root command KubeClusterCmd and sets all flags appropriately.
func Execute() {
	KubeClusterCmd.Execute()
}
