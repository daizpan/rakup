package rakup

import (
	"context"
	"reflect"
	"testing"

	_ "github.com/dikmit/rakup/statik"
)

func TestBrowser_GetTrend(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "GetTrend",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := NewBrowser()
			if err != nil {
				t.Errorf("NewBrowser error: %s", err)
			}
			defer b.Close()
			ctx := context.Background()
			got, err := b.GetTrend(ctx)
			if err != nil {
				t.Errorf("Browser.GetTrend() error = %v", err)
			}
			if len(got) == 0 {
				t.Errorf("Browser.GetTrend() error = %v", err)
			}
		})
	}
}

func TestReadWords(t *testing.T) {
	type args struct {
		file string
	}
	tests := []struct {
		name    string
		args    args
		want    []string
		wantErr bool
	}{
		{
			name:    "FileExits",
			args:    args{file: "testdata/words.txt"},
			want:    []string{"aaa", "bbb"},
			wantErr: false,
		},
		{
			name:    "FileNotFound",
			args:    args{file: "testdata/nothing.txt"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ReadWords(tt.args.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("ReadWords() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ReadWords() = %v, want %v", got, tt.want)
			}
		})
	}
}
