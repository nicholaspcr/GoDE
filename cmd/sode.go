package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// sodeCmd represents the sode command
var sodeCmd = &cobra.Command{
	Use:   "sode",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("sode called")
	},
}

func init() {
	rootCmd.AddCommand(sodeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// sodeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// sodeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
