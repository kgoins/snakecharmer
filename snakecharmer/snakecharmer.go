package snakecharmer

import (
	"fmt"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

/* Example:
var rootCmd = &cobra.Command{
	Use:   "cli",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		snakeCharmer := snakecharmer.NewSnakeCharmer(".mycli", "MYCLI")
		confPath, _ := cmd.Flags().GetString("conf")
		return cliutils.InitConfig(cmd, confPath)
	},
}
*/

// SnakeCharmer binds a viper instance to a cobra.Command instance
type SnakeCharmer struct {
	envPrefix       string
	defaultConfName string
}

// NewSnakeCharmer constructs a new SnakeCharmer
func NewSnakeCharmer(envPrefix string, confName string) SnakeCharmer {
	return SnakeCharmer{
		envPrefix:       envPrefix,
		defaultConfName: confName,
	}
}

// InitConfig imports values from viper into the input cmd object
// to form a single consistent view of config information.
// Passing an empty confPath will cause viper to look in the current
// and home directories for a config file.
func (sc SnakeCharmer) InitConfig(cmd *cobra.Command, confPath string) error {
	v := viper.New()

	if confPath != "" {
		v.SetConfigFile(confPath)
	} else {
		// Set the base name of the config file, without the file extension.
		v.SetConfigName(sc.defaultConfName)

		home, err := homedir.Dir()
		if err != nil {
			return err
		}

		v.AddConfigPath(".")
		v.AddConfigPath(home)
	}

	// Attempt to read the config file, gracefully ignoring errors
	// caused by a config file not being found. Return an error
	// if we cannot parse the config file.
	if err := v.ReadInConfig(); err != nil {
		// It's okay if there isn't a config file
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	// When we bind flags to environment variables expect that the
	// environment variables are prefixed, e.g. a flag like --number
	// binds to an environment variable STING_NUMBER. This helps
	// avoid conflicts.
	v.SetEnvPrefix(sc.envPrefix)

	// Bind to environment variables
	// Works great for simple config names, but needs help for names
	// like --favorite-color which we fix in the bindFlags function
	v.AutomaticEnv()

	// Bind the current command's flags to viper
	sc.bindFlags(cmd, v)

	return nil
}

func (sc SnakeCharmer) bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", sc.envPrefix, envVarSuffix))
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
