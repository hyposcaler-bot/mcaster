package cli

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/yourusername/mcaster/internal/multicast"
)

func newReceiveCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "receive",
		Short: "Receive multicast packets",
		Long: `Listen for multicast packets and display their contents including
timing information and network delay calculations.`,
		Example: `  # Receive from default group
  mcaster receive

  # Receive from specific group via specific interface
  mcaster receive -g 224.0.1.1:8080 -i eth0

  # Receive on specific destination port
  mcaster receive --dport 8080`,
		RunE: func(cmd *cobra.Command, args []string) error {
			group := viper.GetString("group")
			iface := viper.GetString("interface")
			dport := viper.GetInt("dport")

			receiver, err := multicast.NewReceiver(group, iface, dport)
			if err != nil {
				return err
			}

			return receiver.Start()
		},
	}
}
