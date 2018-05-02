package cmd

import (
	"testing"
	"time"
)

func TestAggregateDurationsInString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name    string
		args    args
		want    time.Duration
		wantErr bool
	}{
		{
			name: "5 minute single duration",
			args: args{
				str: `t=5m`,
			},
			want:    time.Minute * 5,
			wantErr: false,
		},
		{
			name: "5 h 2 m 3 s single duration",
			args: args{
				str: `t=5h2m3s`,
			},
			want:    time.Hour*5 + time.Minute*2 + time.Second*3,
			wantErr: false,
		},
		{
			name: "no duration",
			args: args{
				str: `adfasf asdfadfa sdf ajf = 4 2 fam`,
			},
			want:    0,
			wantErr: false,
		},
		{
			name: "multi line multi durations",
			args: args{
				str: `adfasf asdfadfa sdf ajf = 4 2 fam
				: t=3m
				: t=12h
				asdfkajf akdj f
				t=2h
				`,
			},
			want:    14*time.Hour + 3*time.Minute,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := aggregateDurationsInString(tt.args.str)
			if (err != nil) != tt.wantErr {
				t.Errorf("aggregateDurationsInString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if got != tt.want {
				t.Errorf("aggregateDurationsInString() = %v, want %v", got, tt.want)
			}
		})
	}
}
