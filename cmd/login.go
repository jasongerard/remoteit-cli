package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasongerard/remoteit-cli/client"
	"github.com/jasongerard/remoteit-cli/storage"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"syscall"
	"text/tabwriter"
	"time"
)

var loginCmd = &cobra.Command{
	Use:   "login [username]",
	Short: "Retrieves login token. Username is passed to login via argument or read from REMOTEIT_USERNAME environment variable if missing.",
	Args: cobra.MaximumNArgs(1),
	ValidArgs: []string{"username"},
	Run: func(cmd *cobra.Command, args []string) {

		// big ugly function follows
		config := viper.GetViper()

		rc := getClient(config)

		apikey := config.GetString("apikey")

		if apikey == "" {
			errorAndExit(errors.New("apikey not provided"), 1)
		}

		username := config.GetString("username")

		if len(args) == 1 {
			username = args[0]
		}

		if username == "" {
			errorAndExit(errors.New("username not provided"), 1)
		}

		password := config.GetString("password")

		if config.GetBool("prompt") || password == "" {
			fmt.Print("Enter password: ")
			pw, err := terminal.ReadPassword(int(syscall.Stdin))

			errorAndExit(err, 1)

			password = string(pw)
			fmt.Println()
		}

		if password == "" {
			errorAndExit(errors.New("password not provided"), 1)
		}

		req := client.LoginRequest{
			APIKey:  apikey,
			Username: username,
			Password: password,
		}

		resp, err := rc.Login(req)

		errorAndExit(err, 1)

		b, err := json.Marshal(resp)

		err = storage.WriteFile(storage.LoginFile, b)

		errorAndExit(err, 1)

		printJson, _ := cmd.Flags().GetBool("json")

		if printJson {
			fmt.Print(string(b))
		} else {
			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()

			t := time.Unix(int64(resp.AuthExpiration), 0)

			if !config.GetBool("noheader") {
				fields := []string{
					"Token",
					"Expiry Unix",
					"Expiry",
				}
				fmt.Fprintln(w, createHeader(fields))
			}
			fmt.Fprintf(w, "%s\t%v\t%s\n", resp.Token, resp.AuthExpiration, t.String())
		}
	},
}
func init() {
	loginCmd.PersistentFlags().Bool("prompt", false, "prompt for password, default is to read REMOTEIT_PASSWORD environment var")
	viper.BindPFlag("prompt", loginCmd.PersistentFlags().Lookup("prompt"))
}
