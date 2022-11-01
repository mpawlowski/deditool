package steam

import (
	"encoding/json"
	"fmt"

	"github.com/rumblefrog/go-a2s"
	"github.com/spf13/cobra"
)

var queryPlayersCmd = &cobra.Command{
	Use:   "query-players",
	Short: "Query the players in a gameserver",
	Long:  `Query the players in a gameserver using the Steam A2S_PLAYER protocol.`,
	RunE: func(cmd *cobra.Command, args []string) error {

		client, err := a2s.NewClient(fmt.Sprintf("%s:%d", queryHost, queryPort))
		if err != nil {
			return err
		}
		defer client.Close()

		info, err := client.QueryPlayer()
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

func init() {
	queryPlayersCmd.PersistentFlags().StringVar(&queryHost, "host", "127.0.0.1", "Hostname to query")
	queryPlayersCmd.PersistentFlags().IntVar(&queryPort, "port", 27015, "Port to query")
}
