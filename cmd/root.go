package cmd

import (
	"fmt"
	"github.com/jasongerard/remoteit-cli/client"
	"github.com/jasongerard/remoteit-cli/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
)

const urlKey = "remoteit_url"

var rootCmd = &cobra.Command{
	Use:   "remoteit",
	Short: "remote.it CLI",
	Long: `Command Line Interface for remote.it allowing login, device list, and connections`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s\n\n%s\n", cmd.Long, cmd.UsageString())
		os.Exit(0)
	},
}

func init() {

	err := storage.Initialize()

	if err != nil {
		panic(err)
	}

	config := viper.GetViper()
	config.SetEnvPrefix("remoteit")
	config.AutomaticEnv()

	config.SetDefault(urlKey, "https://api.remot3.it/apv/v23.5")

	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(deviceCmd)
	rootCmd.AddCommand(connectCmd)

	rootCmd.PersistentFlags().String("apikey", "", "API key for accessing remote.it API. Uses REMOTEIT_APIKEY env var if not set")
	viper.BindPFlag("apikey", rootCmd.PersistentFlags().Lookup("apikey"))

	rootCmd.PersistentFlags().Bool("json", false, "Output the results of the command in JSON format")
	viper.BindPFlag("json", rootCmd.PersistentFlags().Lookup("json"))

	rootCmd.PersistentFlags().Bool("noheader", false, "Disables header in output.")
	viper.BindPFlag("noheader", rootCmd.PersistentFlags().Lookup("noheader"))

	rootCmd.PersistentFlags().Bool("loghttp", false, "Log HTTP request/response from API to ~/.remoteit/http.log")
	viper.BindPFlag("loghttp", rootCmd.PersistentFlags().Lookup("loghttp"))
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		errorAndExit(err, 1)
	}
}

func errorAndExit(err error, code int) {
	if err == nil {
		return
	}

	fmt.Println(err)

	os.Exit(code)
}

func getClient(config *viper.Viper) client.Client {

	out, err := storage.GetHTTPLogWriter()

	if err != nil {
		panic(err)
	}

	logger := log.New(out, "HTTPLOG ", log.LstdFlags)

	return client.NewClient(config, nil, logger)
}

func getTokenFromFile() string {
	token, err := storage.GetToken()

	if os.IsNotExist(err) {
		fmt.Println("cannot find token, did you run `remote login` ?")
	}

	errorAndExit(err, 1)

	return token
}

func createHeader(fields []string) string {

	var underlined []string
	for _, s := range fields {
		underlined = append(underlined, strings.Repeat("-", len(s)))
	}
	n := strings.Join(fields, "\t")
	u := strings.Join(underlined, "\t")
	return n + "\n" + u + "\t"

}