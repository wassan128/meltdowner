package cmd

import (
	"fmt"
	"net/http"

	"github.com/spf13/cobra"
	"github.com/wassan128/meltdowner/meltdowner/build"
)

type Opts struct {
	optBool bool
}
var o = &Opts{}

var RootCmd = &cobra.Command{
	Use: "melt",
	Short: "CLI tool for MeltDonwer(blog generator)",
	Run: func(cmd *cobra.Command, args []string) {},
}

var versionCmd = &cobra.Command{
	Use: "version",
	Short: "Print the version number of melt",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("melt v0.1")
	},
}

var generateCmd = &cobra.Command{
	Use: "generate",
	Short: "Generate blog(htmls, static files)",
	Run: func(cmd *cobra.Command, args []string) {
		build.Run()
	},
}

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "Serve public/ on localhost",
	Run: func(cmd *cobra.Command, args []string) {
		if o.optBool {
			fmt.Println("[*] found option generate before serve.")
			build.Run()
		}
		http.Handle("/", http.FileServer(http.Dir("public")))

		fmt.Println("[*] public/ is being served on http://localhost:5000")
		if err := http.ListenAndServe(":5000", nil); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	cobra.OnInitialize()

	RootCmd.AddCommand(versionCmd)
	RootCmd.AddCommand(generateCmd)
	RootCmd.AddCommand(serverCmd)
	serverCmd.Flags().BoolVarP(&o.optBool, "generate", "g", false, "generate before serve")
}

