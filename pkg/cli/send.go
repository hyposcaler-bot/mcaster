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
  mcaster send --interface eth0`,
		RunE: func(cmd *cobra.Command, args []string) error {
			group := viper.GetString("group")
			iface := viper.GetString("interface")
			interval := viper.GetDuration("interval")

			sender, err := multicast.NewSender(group, iface, interval)
			if err != nil {
				return err
			}

			return sender.Start()
		},
	}

	cmd.Flags().DurationP("interval", "t", time.Second, "send interval")
	viper.BindPFlag("interval", cmd.Flags().Lookup("interval"))

	return cmd
}
