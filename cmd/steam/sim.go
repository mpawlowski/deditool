package steam

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var serverPort int
var serverName string
var serverMap string
var serverFolder string
var serverGame string
var serverPlayersCurrent int
var serverPlayersMax int
var serverPlayersBots int
var serverLatency time.Duration

func init() {

	hostname, err := os.Hostname()
	if err != nil {
		panic(err)
	}

	simCmd.PersistentFlags().IntVar(&serverPort, "port", 27015, "Port bound to incoming A2S requests.")
	simCmd.PersistentFlags().StringVar(&serverName, "name", "deditool A2S simulator", "Name of the server.")
	simCmd.PersistentFlags().StringVar(&serverMap, "map", fmt.Sprintf("%s (%s)", hostname, runtime.GOOS), "Map the server has currently loaded.")
	simCmd.PersistentFlags().StringVar(&serverFolder, "folder", "golang", "Name of the folder containing the game files. ")
	simCmd.PersistentFlags().StringVar(&serverGame, "game", fmt.Sprintf("Golang (%s)", runtime.GOARCH), "Full name of the game.")
	simCmd.PersistentFlags().IntVar(&serverPlayersCurrent, "players-current", 0, "Current number of players.")
	simCmd.PersistentFlags().IntVar(&serverPlayersMax, "players-max", 10, "Max number of players.")
	simCmd.PersistentFlags().IntVar(&serverPlayersBots, "players-bots", 0, "Current number of bots.")
	simCmd.PersistentFlags().DurationVar(&serverLatency, "latency", 25*time.Millisecond, "Latency to respond to requests.")
}

var simCmd = &cobra.Command{
	Use:   "sim",
	Short: "Simulate a gameserver",
	Long:  `Simulate a gameserver using the Steam query protocol`,
	Run: func(cmd *cobra.Command, args []string) {
		addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf(":%d", serverPort))
		if err != nil {
			fmt.Printf("udp resolve error: %v\n", err)
			return
		}

		conn, err := net.ListenUDP("udp", addr)
		if err != nil {
			fmt.Printf("udp listen error: %v\n", err)
			return
		}

		serverResponder := func(conn *net.UDPConn) {
			fmt.Printf("now listening on %s\n", conn.LocalAddr().String())

			for {
				p := make([]byte, 2048)
				oob := make([]byte, 2048)
				packetLength, _, _, addr, err := conn.ReadMsgUDP(p, oob)
				if err != nil {
					fmt.Printf("error reading packet: %v", err)
					continue
				}

				if packetLength < 6 {
					continue
				}

				a2sHeader := p[4]
				switch a2sHeader {
				case A2S_INFO_REQUEST_HEADER:
					fmt.Printf("a2s-info packet (%s) request 0x%x\n", addr.IP, a2sHeader)
					demoserver := BuildServerInfo(
						serverName,
						serverMap,
						serverFolder,
						serverGame,
						serverPlayersCurrent,
						serverPlayersMax,
						serverPlayersBots,
					)
					go func() {
						time.Sleep(serverLatency)
						_, err := conn.WriteToUDP(demoserver, addr)
						if err != nil {
							fmt.Printf("a2s-info response error: %v", err)
						}
					}()
				case A2S_PLAYER_REQUEST_HEADER:
					fmt.Printf("a2s-player packet (%s) request 0x%x\n", addr.IP, a2sHeader)
					demoPlayerInfo := BuildPlayerInfo(serverPlayersCurrent)
					go func() {
						time.Sleep(serverLatency)
						_, err := conn.WriteToUDP(demoPlayerInfo, addr)
						if err != nil {
							fmt.Printf("a2s-player response error: %v", err)
						}
					}()
				case A2S_RULES_REQUEST_HEADER:
					fmt.Printf("a2s-rules packet (%s) request 0x%x\n", addr.IP, a2sHeader)
					demoRules := BuildRules()
					go func() {
						time.Sleep(serverLatency)
						_, err := conn.WriteToUDP(demoRules, addr)
						if err != nil {
							fmt.Printf("a2s-rules response error: %v", err)
						}
					}()
				default:
					fmt.Printf("unknown packet (%s) request 0x%x\n", addr.IP, a2sHeader)
					continue
				}
			}
		}

		go serverResponder(conn)

		sigs := make(chan os.Signal, 1)
		done := make(chan bool, 1)

		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-sigs
			fmt.Println()
			fmt.Printf("caught signal: %s\n", sig)
			done <- true
		}()

		<-done
	},
}

const A2S_INFO_REQUEST_HEADER = 0x54
const A2S_PLAYER_REQUEST_HEADER = 0x55
const A2S_RULES_REQUEST_HEADER = 0x56

func BuildServerInfo(serverName string, serverMap string, serverFolder string, serverGame string, playersCurrent int, playersMax int, playersBots int) (b []byte) {
	buf := bytes.Buffer{}

	// info
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0x49)           // a2s-info-response header
	buf.WriteByte(0x02)           // version
	buf.WriteString(serverName)   // server name
	buf.WriteByte(0x00)           // null byte
	buf.WriteString(serverMap)    // server map
	buf.WriteByte(0x00)           // null byte
	buf.WriteString(serverFolder) // server folder
	buf.WriteByte(0x00)           // null byte
	buf.WriteString(serverGame)   // server map
	buf.WriteByte(0x00)           // null byte
	buf.WriteByte(0xF0)           // app id
	buf.WriteByte(0x00)           // null byte

	// players
	buf.WriteByte(byte(playersCurrent))
	buf.WriteByte(byte(playersMax))
	buf.WriteByte(byte(playersBots))

	// settings
	buf.WriteByte(0x64)       // d - dedicated
	buf.WriteByte(0x6c)       // l - linux
	buf.WriteByte(0x01)       // private
	buf.WriteByte(0x00)       // vac disabled
	buf.WriteString("v0.0.1") // server version
	buf.WriteByte(0x00)       // null byte

	b = buf.Bytes()

	return
}

func BuildPlayerInfo(playersCurrent int) (b []byte) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0x44) // a2s-player-response header
	buf.WriteByte(0x00) //TODO implement players
	b = buf.Bytes()
	return
}

func BuildRules() (b []byte) {
	buf := bytes.Buffer{}
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0xFF)
	buf.WriteByte(0x45) // a2s-rules-response header

	var numRules int16 = 0
	numRulesBuf := make([]byte, 2)
	binary.LittleEndian.PutUint16(numRulesBuf, uint16(numRules))
	buf.WriteByte(numRulesBuf[0])
	buf.WriteByte(numRulesBuf[1])

	b = buf.Bytes()
	return
}
