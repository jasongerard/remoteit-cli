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
	"strings"
	"text/tabwriter"
)

var connectCmd = &cobra.Command{
	Use:   "connect [devicealias] [hostip]",
	Short: "create proxy for device",
	Long: "Creates and returns a proxy that can be used to connect to a device",
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {

		// big ugly function follows
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

		req := client.ConnectRequest{
			BaseReqeust:client.BaseReqeust{
				APIKey: apikey,
				Token:  token,
			},
		}

		nocache, _ := cmd.Flags().GetBool("nocache")

		var devices []*client.DeviceEntry
		var address, lastIP string

		if nocache || !storage.CacheExists() {
			req := client.ListDevicesRequest{
				BaseReqeust:client.BaseReqeust{
					APIKey: apikey,
					Token:  token,
				},
			}

			resp, err := rc.ListDevices(req)

			errorAndExit(err, 1)

			devices = resp.Devices
		} else {
			var err error
			devices, err = storage.GetDevicesFromCache()

			errorAndExit(err, 1)
		}

		for _, d := range devices {
			if d.Alias == args[0] {
				lastIP = d.LastIP
				address = d.Address
			}
		}

		if len(args) == 2 {
			lastIP = args[1]
		}

		if lastIP == "" {
			errorAndExit(errors.New("no hostip provided and device LastIP not set"), 1)
		}

		req.DeviceAddress = address
		req.HostIP = lastIP

		resp, err := rc.Connect(req)

		errorAndExit(err, 1)

		printJson, _ := cmd.Flags().GetBool("json")

		if printJson {

			b, err := json.Marshal(resp)

			errorAndExit(err, 1)
			var out bytes.Buffer
			json.Indent(&out, b, "", "  ")
			fmt.Println(out.String())

		} else {

			format, _ := cmd.Flags().GetString("format")

			if format != "ssh" {

				w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
				defer w.Flush()

				if !config.GetBool("noheader") {
					fields := []string{
						"Proxy",
						"Device Address",
						"Status",
						"Expiration",
						"Requested",
					}
					fmt.Fprintln(w, createHeader(fields))
				}

				fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\t\n", resp.Connection.Proxy, resp.Connection.DeviceAddress, resp.Status, resp.Connection.ExpirationInSeconds, resp.Connection.RequestedAt)
			} else {
				p := resp.Connection.Proxy
				p = strings.Replace(strings.Replace(p, "http://", "", 1), "https://", "", 1)

				parts := strings.Split(p, ":")

				fmt.Printf("%s -p %s", parts[0], parts[1])
			}
		}
	},
}

func init() {
	connectCmd.PersistentFlags().String("format", "", "format for proxy connection (currently ssh is only valid value)")
	viper.BindPFlag("format", loginCmd.PersistentFlags().Lookup("format"))

	connectCmd.PersistentFlags().Bool("nocache", false, "do not use the local device cache, fetch device list fresh from remote.it")
	viper.BindPFlag("nocache", loginCmd.PersistentFlags().Lookup("nocache"))
}
