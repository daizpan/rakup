package cmd

import (
	"encoding/hex"
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

type config struct {
	User     string
	Password string
}

var conf config

var cryptoKey = "81e28ad8c4784eca72902c6810579fa645f9300b7109424438d69500a017aa27"

func NewRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Version:       "1.0.1",
		SilenceUsage:  true,
		SilenceErrors: true,
		Use:           "rakup",
		Short:         "Get Rakuten points for Rakuten Web Search",
		Long:          "Get Rakuten points for Rakuten Web Search",
	}

	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.rakup.yaml)")

	rootCmd.AddCommand(
		NewConfigCmd(),
		NewRunCmd(),
	)

	return rootCmd
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cmd := NewRootCmd()
	cmd.SetOutput(os.Stdout)
	cmd.SetErr(os.Stderr)
	if err := cmd.Execute(); err != nil {
		cmd.PrintErrln(err)
		os.Exit(1)
	}
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".rakup" (without extension).
		dir := ".rakup"
		viper.SetConfigName("config")
		viper.SetConfigType("yaml")
		viper.AddConfigPath(path.Join(".", dir))
		viper.AddConfigPath(path.Join(home, dir))
	}

	viper.SetEnvPrefix("rakup")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Config file not read:", viper.ConfigFileUsed())
	}

	user := viper.GetString("user")
	if user != "" {
		user, err := DecryptValue(user)
		if err != nil {
			fmt.Println(err)
			return
		}
		conf.User = string(user)
	}
	password := viper.GetString("password")
	if password != "" {
		password, err := DecryptValue(password)
		if err != nil {
			fmt.Println(err)
			return
		}
		conf.Password = string(password)
	}
	// if err := viper.Unmarshal(&conf); err != nil {
	// 	fmt.Println(err)
	// }
}

func DecryptValue(text string) ([]byte, error) {
	key, err := hex.DecodeString(cryptoKey)
	if err != nil {
		return nil, err
	}
	ciphertext, err := hex.DecodeString(text)
	if err != nil {
		return nil, err
	}
	plaintext, err := Decrypt(key, ciphertext)
	if err != nil {
		return nil, err
	}
	return plaintext, nil
}

func EncryptValue(data []byte) (string, error) {
	key, err := hex.DecodeString(cryptoKey)
	if err != nil {
		return "", err
	}
	ciphertext, err := Encrypt(key, data)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(ciphertext), nil
}
