package cmd

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/spf13/viper"
)

func initDummyConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("testdata")
}
func TestConfigCmdErrors(t *testing.T) {
	initDummyConfig()
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "all nothing",
			args:          []string{},
			expectedError: "Key: 'configOptions.User' Error:Field validation for 'User' failed on the 'required' tag\nKey: 'configOptions.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:          "not password",
			args:          []string{"-u", "user"},
			expectedError: "Key: 'configOptions.Password' Error:Field validation for 'Password' failed on the 'required' tag",
		},
		{
			name:          "not user",
			args:          []string{"-p", "password"},
			expectedError: "Key: 'configOptions.User' Error:Field validation for 'User' failed on the 'required' tag",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd := NewConfigCmd()
			cmd.SetOutput(buf)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err == nil || err.Error() != tt.expectedError {
				t.Errorf("Validate error = %v, expected %v", err, tt.expectedError)
				return
			}
		})
	}
}

func TestConfigCmdSuccess(t *testing.T) {
	initDummyConfig()
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "test1",
			args: []string{"-u", "user", "-p", "password"},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println(os.Getwd())
			buf := new(bytes.Buffer)
			cmd := NewConfigCmd()
			cmd.SetOutput(buf)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Failed %s error = %v", tt.name, err)
			}
			got := buf.String()
			if got != tt.want {
				t.Errorf("got: %+v want: %+v\n", got, tt.want)
			}
		})
	}
}
