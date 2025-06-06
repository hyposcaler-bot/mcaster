package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "mcaster",
		Short: "A multicast connectivity testing tool",
		Long: `A simple tool for testing multicast connectivity by sending and receiving
UDP multicast packets with timestamps and sequence numbers.

Examples:
  mcaster send                           # Send to default group
  mcaster receive                        # Receive from default group
  mcaster send -g 224.0.1.1:8080        # Send to specific group
  mcaster receive -i eth0                # Receive via specific interface
  MULTICAST_GROUP=239.23.23.23:2323 mcaster send  # Use environment variable`,
	}
)

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	// Global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.mcaster.yaml)")
	rootCmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group address:port")
	rootCmd.PersistentFlags().StringP("interface", "i", "", "network interface name")
	rootCmd.PersistentFlags().IntP("dport", "d", 0, "destination port (overrides port in group address)")

	// Bind flags to viper
	viper.BindPFlag("group", rootCmd.PersistentFlags().Lookup("group"))
	viper.BindPFlag("interface", rootCmd.PersistentFlags().Lookup("interface"))
	viper.BindPFlag("dport", rootCmd.PersistentFlags().Lookup("dport"))

	// Environment variable bindings
	viper.SetEnvPrefix("MULTICAST")
	viper.BindEnv("group", "MULTICAST_GROUP")
	viper.BindEnv("interface", "MULTICAST_INTERFACE")
	viper.BindEnv("interval", "MULTICAST_INTERVAL")
	viper.BindEnv("ttl", "MULTICAST_TTL")
	viper.BindEnv("sport", "MULTICAST_SPORT")
	viper.BindEnv("dport", "MULTICAST_DPORT")

	// Add subcommands
	rootCmd.AddCommand(newSendCmd())
	rootCmd.AddCommand(newReceiveCmd())
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mcaster")
	}

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
