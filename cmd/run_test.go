package cmd

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

func TestRunCmdErrors(t *testing.T) {
	tests := []struct {
		name          string
		args          []string
		expectedError string
	}{
		{
			name:          "no config",
			args:          []string{"-v"},
			expectedError: "please set config",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd := NewRunCmd()
			cmd.SetOutput(buf)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err == nil || err.Error() != tt.expectedError {
				t.Errorf("Validate error = %v, expected %v", err, tt.expectedError)
			}
		})
	}
}

func TestRunCmdLoginError(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "test1",
			args: []string{"-v", "-n", "1", "-u", "user", "-p", "password"},
			want: "login error",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd := NewRunCmd()
			cmd.SetOutput(buf)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err == nil || err.Error() != tt.want {
				t.Errorf("RunCmd error = %v, expected %v", err, tt.want)
			}
		})
	}
}
func TestRunCmdSuccess(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{
			name: "test1",
			args: []string{
				"-v", "-n", "1", "-u", os.Getenv("RAKUTEN_USER"), "-p", os.Getenv("RAKUTEN_PASSWORD"),
				"--word-file", "testdata/words.txt",
			},
			want: "Searched",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buf := new(bytes.Buffer)
			cmd := NewRunCmd()
			cmd.SetOutput(buf)
			cmd.SetArgs(tt.args)
			err := cmd.Execute()
			if err != nil {
				t.Errorf("Failed %s error = %v", tt.name, err)
			}
			got := buf.String()
			if !strings.Contains(got, tt.want) {
				t.Errorf("got: %+v want: %+v\n", got, tt.want)
			}
		})
	}
}
