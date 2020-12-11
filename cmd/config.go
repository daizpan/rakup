package cmd

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type configOptions struct {
	User     string `validate:"required"`
	Password string `validate:"required"`
}

// NewConfigCmd creates a new `rakup config` command
func NewConfigCmd() *cobra.Command {
	options := configOptions{}
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Set config",
		// SilenceUsage: true,
		PreRunE: func(cmd *cobra.Command, args []string) error {
			validate := validator.New()
			err := validate.Struct(options)
			if err != nil {
				return err
			}
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runConfig(cmd, options)
		},
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.User, "user", "u", "", "user")
	flags.StringVarP(&options.Password, "password", "p", "", "password")

	return cmd
}

func runConfig(cmd *cobra.Command, options configOptions) error {
	user, err := EncryptValue([]byte(options.User))
	if err != nil {
		return err
	}
	password, err := EncryptValue([]byte(options.Password))
	if err != nil {
		return err
	}
	viper.Set("user", user)
	viper.Set("password", password)
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}
