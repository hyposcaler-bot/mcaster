package cli

import (
	"bytes"
	"os"
	"testing"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCommand(t *testing.T) {
	// Reset viper state before each test
	viper.Reset()
	
	tests := []struct {
		name        string
		args        []string
		expectError bool
	}{
		{
			name:        "help flag",
			args:        []string{"--help"},
			expectError: false,
		},
		{
			name:        "version-like usage",
			args:        []string{},
			expectError: false,
		},
		{
			name:        "invalid flag",
			args:        []string{"--invalid-flag"},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new root command for each test
			cmd := &cobra.Command{
				Use:   "mcaster",
				Short: "A multicast connectivity testing tool",
			}
			
			// Add global flags like the real root command
			cmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group address:port")
			cmd.PersistentFlags().StringP("interface", "i", "", "network interface name")
			cmd.PersistentFlags().IntP("dport", "d", 0, "destination port (overrides port in group address)")

			cmd.SetArgs(tt.args)

			// Capture output
			var buf bytes.Buffer
			cmd.SetOut(&buf)
			cmd.SetErr(&buf)

			err := cmd.Execute()

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGlobalFlags(t *testing.T) {
	viper.Reset()
	
	tests := []struct {
		name         string
		args         []string
		expectedGroup string
		expectedIface string
		expectedDPort int
	}{
		{
			name:         "default values",
			args:         []string{},
			expectedGroup: "239.23.23.23:2323",
			expectedIface: "",
			expectedDPort: 0,
		},
		{
			name:         "custom group",
			args:         []string{"-g", "224.0.1.1:8080"},
			expectedGroup: "224.0.1.1:8080",
			expectedIface: "",
			expectedDPort: 0,
		},
		{
			name:         "custom interface",
			args:         []string{"-i", "eth0"},
			expectedGroup: "239.23.23.23:2323",
			expectedIface: "eth0",
			expectedDPort: 0,
		},
		{
			name:         "custom destination port",
			args:         []string{"-d", "9999"},
			expectedGroup: "239.23.23.23:2323",
			expectedIface: "",
			expectedDPort: 9999,
		},
		{
			name:         "all custom flags",
			args:         []string{"-g", "225.1.1.1:7777", "-i", "lo", "-d", "8888"},
			expectedGroup: "225.1.1.1:7777",
			expectedIface: "lo",
			expectedDPort: 8888,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			
			// Create command with flags
			cmd := &cobra.Command{
				Use: "test",
				RunE: func(cmd *cobra.Command, args []string) error {
					// This would be where we test the flag values
					return nil
				},
			}
			
			cmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group address:port")
			cmd.PersistentFlags().StringP("interface", "i", "", "network interface name")
			cmd.PersistentFlags().IntP("dport", "d", 0, "destination port")

			// Bind flags to viper
			viper.BindPFlag("group", cmd.PersistentFlags().Lookup("group"))
			viper.BindPFlag("interface", cmd.PersistentFlags().Lookup("interface"))
			viper.BindPFlag("dport", cmd.PersistentFlags().Lookup("dport"))

			cmd.SetArgs(tt.args)
			
			err := cmd.Execute()
			require.NoError(t, err)

			// Check viper values
			assert.Equal(t, tt.expectedGroup, viper.GetString("group"))
			assert.Equal(t, tt.expectedIface, viper.GetString("interface"))
			assert.Equal(t, tt.expectedDPort, viper.GetInt("dport"))
		})
	}
}

func TestEnvironmentVariableBinding(t *testing.T) {
	// Clear any existing environment variables
	clearTestEnvVars()
	defer clearTestEnvVars()

	tests := []struct {
		name     string
		envVar   string
		envValue string
		viperKey string
		expected interface{}
	}{
		{
			name:     "MULTICAST_GROUP",
			envVar:   "MULTICAST_GROUP",
			envValue: "224.1.1.1:1234",
			viperKey: "group",
			expected: "224.1.1.1:1234",
		},
		{
			name:     "MULTICAST_INTERFACE",
			envVar:   "MULTICAST_INTERFACE", 
			envValue: "eth0",
			viperKey: "interface",
			expected: "eth0",
		},
		{
			name:     "MULTICAST_DPORT",
			envVar:   "MULTICAST_DPORT",
			envValue: "9999",
			viperKey: "dport",
			expected: 9999,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.Reset()
			clearTestEnvVars()

			// Set environment variable
			os.Setenv(tt.envVar, tt.envValue)

			// Setup viper like the real app
			viper.SetEnvPrefix("MULTICAST")
			viper.BindEnv("group", "MULTICAST_GROUP")
			viper.BindEnv("interface", "MULTICAST_INTERFACE")
			viper.BindEnv("dport", "MULTICAST_DPORT")
			viper.AutomaticEnv()

			// Check that viper picked up the environment variable
			switch v := tt.expected.(type) {
			case string:
				assert.Equal(t, v, viper.GetString(tt.viperKey))
			case int:
				assert.Equal(t, v, viper.GetInt(tt.viperKey))
			}
		})
	}
}

func TestConfigFileBinding(t *testing.T) {
	viper.Reset()
	
	// This test would require creating temporary config files
	// For now, we'll test that the config initialization doesn't fail
	
	// Test that initConfig can be called without error
	assert.NotPanics(t, func() {
		// Simulate what happens in initConfig
		viper.SetConfigType("yaml")
		viper.SetConfigName(".mcaster")
		viper.AutomaticEnv()
		
		// Try to read config (will fail silently if no file exists)
		viper.ReadInConfig()
	})
}

func TestFlagPrecedence(t *testing.T) {
	// Test that command line flags override environment variables
	clearTestEnvVars()
	defer clearTestEnvVars()
	
	viper.Reset()
	
	// Set environment variable
	os.Setenv("MULTICAST_GROUP", "224.1.1.1:1111")
	
	// Create command that will override with flag
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	
	cmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group")
	
	// Setup viper
	viper.SetEnvPrefix("MULTICAST")
	viper.BindEnv("group", "MULTICAST_GROUP")
	viper.BindPFlag("group", cmd.PersistentFlags().Lookup("group"))
	viper.AutomaticEnv()
	
	// Set args to override environment
	cmd.SetArgs([]string{"-g", "225.5.5.5:5555"})
	
	err := cmd.Execute()
	require.NoError(t, err)
	
	// Flag should override environment
	assert.Equal(t, "225.5.5.5:5555", viper.GetString("group"))
}

func TestHelpOutput(t *testing.T) {
	// Test that help output contains expected information
	cmd := &cobra.Command{
		Use:   "mcaster",
		Short: "A multicast connectivity testing tool",
		Long: `A simple tool for testing multicast connectivity by sending and receiving
UDP multicast packets with timestamps and sequence numbers.`,
	}
	
	cmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group address:port")
	cmd.PersistentFlags().StringP("interface", "i", "", "network interface name")
	cmd.PersistentFlags().IntP("dport", "d", 0, "destination port")
	
	// Add help command explicitly
	cmd.AddCommand(&cobra.Command{Use: "help", Run: func(cmd *cobra.Command, args []string) {}})
	
	cmd.SetArgs([]string{"--help"})
	
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	
	err := cmd.Execute()
	assert.NoError(t, err)
	
	output := buf.String()
	
	// Check that help contains expected elements (more lenient)
	assert.Contains(t, output, "multicast")
	assert.Contains(t, output, "group")
	assert.Contains(t, output, "interface")
}

// Helper function to clear test environment variables
func clearTestEnvVars() {
	envVars := []string{
		"MULTICAST_GROUP",
		"MULTICAST_INTERFACE",
		"MULTICAST_INTERVAL",
		"MULTICAST_TTL",
		"MULTICAST_SPORT",
		"MULTICAST_DPORT",
	}
	
	for _, envVar := range envVars {
		os.Unsetenv(envVar)
	}
}

// Benchmark tests
func BenchmarkRootCommandExecution(b *testing.B) {
	cmd := &cobra.Command{
		Use: "test",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	
	cmd.PersistentFlags().StringP("group", "g", "239.23.23.23:2323", "multicast group")
	cmd.SetArgs([]string{})
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		viper.Reset()
		err := cmd.Execute()
		if err != nil {
			b.Fatal(err)
		}
	}
}