package main

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/benbjohnson/clock"
)

func TestExecuteOnTick(t *testing.T) {
	type args struct {
		tickerInterval  time.Duration
		timeoutInterval time.Duration
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "Ticker 1s, Timeout 5s",
			args: args{
				tickerInterval:  1 * time.Second,
				timeoutInterval: 5 * time.Second,
			},
			want: []string{"OK", "OK", "OK", "OK"},
		},
		{
			name: "Ticker 2s, Timeout 5s",
			args: args{
				tickerInterval:  2 * time.Second,
				timeoutInterval: 5 * time.Second,
			},
			want: []string{"OK", "OK"},
		},
		{
			name: "Ticker 2s, Timeout 6s",
			args: args{
				tickerInterval:  2 * time.Second,
				timeoutInterval: 6 * time.Second,
			},
			want: []string{"OK", "OK"},
		},
		{
			name: "Ticker 3s, Timeout 10s",
			args: args{
				tickerInterval:  3 * time.Second,
				timeoutInterval: 10 * time.Second,
			},
			want: []string{"OK", "OK", "OK"},
		},
		{
			name: "Ticker 5s, Timeout 10s",
			args: args{
				tickerInterval:  5 * time.Second,
				timeoutInterval: 10 * time.Second,
			},
			want: []string{"OK"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			// Create a new mock clock
			mockclock := clock.NewMock()

			// Create a context with a timeout
			mockctx, cancel := mockclock.WithTimeout(context.Background(), tt.args.timeoutInterval)
			defer cancel()

			// Start the ExecuteOnTick function with our mock clock
			resultCh := ExecuteOnTick(mockctx, mockclock, tt.args.tickerInterval)

			got := []string{}

			for {
				// Advance the mock clock
				mockclock.Add(tt.args.tickerInterval)

				// Read the result
				v, open := <-resultCh
				if !open {
					t.Log("resultCh closed")
					break
				} else {
					got = append(got, v)
				}
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExecuteOnTick(%v, %v) = %v, want %v", tt.args.tickerInterval, tt.args.timeoutInterval, got, tt.want)
			}
		})
	}
}
