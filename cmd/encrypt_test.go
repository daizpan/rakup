package cmd

import (
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncryptDecrypt(t *testing.T) {
	tests := []struct {
		name string
		text string
		key  string
		want string
	}{
		{
			name: "Encrypt",
			text: "testtest",
			key:  "8af6c23d6540ae2853459a4c2cf1012287ca0dcaf1a47d6e01be14dddfa9e2f7",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := []byte(tt.text)
			key, err := hex.DecodeString(tt.key)
			if err != nil {
				t.Error(err)
			}
			ciphertext, err := Encrypt(key, data)
			if err != nil {
				t.Errorf("encrypt error: %s\n", err)
			}

			plaintext, err := Decrypt(key, ciphertext)
			if err != nil {
				t.Errorf("decrypt error: %s\n", err)
			}
			if string(plaintext) != tt.text {
				t.Errorf("got: %+v want: %+v\n", string(plaintext), tt.text)
			}
		})
	}
}

func TestGenerateKey(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "Generate Key",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateKey()
			if err != nil {
				t.Fatal(err)
			}
			key := hex.EncodeToString(got)
			fmt.Println(key)
		})
	}
}
