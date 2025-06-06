package cli

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourusername/mcaster/internal/multicast"
)

func newSendCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "send",
		Short: "Send multicast packets",
		Long: `Send multicast packets continuously to test multicast connectivity.
Each packet contains an incrementing ID, timestamp, and source hostname.`,
		Example: `  # Send with default settings
  mcaster send

  # Send to specific group with fast interval
  mcaster send -g 224.0.1.1:8080 -t 100ms

  # Send via specific interface
  mcaster send --interface eth0

  # Send with custom TTL
  mcaster send --ttl 10

  # Send from specific source port
  mcaster send --sport 12345

  # Send to specific destination port
  mcaster send --dport 8080`,
		RunE: func(cmd *cobra.Command, args []string) error {
			group := viper.GetString("group")
			iface := viper.GetString("interface")
			interval := viper.GetDuration("interval")
			ttl := viper.GetInt("ttl")
			sport := viper.GetInt("sport")
			dport := viper.GetInt("dport")

			sender, err := multicast.NewSender(group, iface, interval, ttl, sport, dport)
			if err != nil {
				return err
			}

			return sender.Start()
		},
	}

	cmd.Flags().DurationP("interval", "t", time.Second, "send interval")
	cmd.Flags().IntP("ttl", "", 1, "TTL (Time To Live) for multicast packets (1-255)")
	cmd.Flags().IntP("sport", "s", 0, "source port for sending packets (0 = random)")
	viper.BindPFlag("interval", cmd.Flags().Lookup("interval"))
	viper.BindPFlag("ttl", cmd.Flags().Lookup("ttl"))
	viper.BindPFlag("sport", cmd.Flags().Lookup("sport"))

	return cmd
}
