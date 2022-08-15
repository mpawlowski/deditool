package steam

import (
	"encoding/json"
	"fmt"

	"github.com/rumblefrog/go-a2s"
	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Query a gameserver",
	Long:  `Query a gameserver using the Steam A2S protocol`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := a2s.NewClient(fmt.Sprintf("%s:%d", queryHost, queryPort))
		if err != nil {
			return err
		}
		defer client.Close()
	
		info, err := client.QueryInfo() // QueryInfo, QueryPlayer, QueryRules
		if err != nil {
			return err
		}
		infoBytes, err := json.Marshal(info)
		if err != nil {
			return err
		}
		fmt.Println(string(infoBytes))
		return nil
	},
}

var queryHost string
var queryPort int

func init() {
	queryCmd.PersistentFlags().StringVar(&queryHost, "host", "127.0.0.1", "Hostname to query")
	queryCmd.PersistentFlags().IntVar(&queryPort, "port", 27015, "Port to query")
}
