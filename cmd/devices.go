package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasongerard/remoteit-cli/client"
	"github.com/jasongerard/remoteit-cli/storage"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"text/tabwriter"
)

var deviceCmd = &cobra.Command{
	Use:   "devices",
	Short: "list devices",
	Long: "Retrieves list of devices registered in remote.it",
	Run: func(cmd *cobra.Command, args []string) {

		// big ugly function follows
		printJson, _ := cmd.Flags().GetBool("json")

		config := viper.GetViper()

		rc := client.NewClient(config, nil)

		token := config.GetString("token")

		if token == "" {
			token = getTokenFromFile()
			if token == "" {
				errorAndExit(errors.New("token not provided"), 1)
			}
		}

		apikey := config.GetString("apikey")

		if apikey == "" {
			errorAndExit(errors.New("apikey not provided"), 1)
		}

		req := client.ListDevicesRequest{
			BaseReqeust:client.BaseReqeust{
				APIKey: apikey,
				Token:  token,
			},
		}

		resp, err := rc.ListDevices(req)

		errorAndExit(err, 1)

		b, err := json.Marshal(resp)

		errorAndExit(err, 1)

		if printJson {

			var out bytes.Buffer
			json.Indent(&out, b, "", "  ")
			fmt.Println(out.String())

		} else {

			w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
			defer w.Flush()

			if !config.GetBool("noheader") {
				fields := []string{
					"Alias",
					"Address",
					"Service",
					"Last IP",
				}
				fmt.Fprintln(w, createHeader(fields))
			}

			for _, v := range resp.Devices {
				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t\n", v.Alias, v.Address, v.ServiceTitle, v.LastIP)
			}
		}

		// write cache
		err = storage.WriteFile(storage.DeviceCacheFile, b)

		errorAndExit(err, 1)
	},
}

func init() {

}
