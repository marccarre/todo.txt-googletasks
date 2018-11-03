package cmd

import (
	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks"
	"github.com/marccarre/todo.txt-googletasks/pkg/gtasks/credentials"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete all tasks from Google Tasks",
	Long:  "",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := gtasks.NewClient(credentials.DefaultPath)
		if err != nil {
			log.Fatal(err)
		}
		if err := client.DeleteAll(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// deleteCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// deleteCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
