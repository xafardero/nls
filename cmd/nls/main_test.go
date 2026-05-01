package main

import "testing"

func TestParseArgs(t *testing.T) {
	tests := []struct {
		name            string
		args            []string
		wantShowVersion bool
		wantCIDR        string
	}{
		{name: "--version flag", args: []string{"--version"}, wantShowVersion: true, wantCIDR: ""},
		{name: "-v flag", args: []string{"-v"}, wantShowVersion: true, wantCIDR: ""},
		{name: "CIDR arg", args: []string{"10.0.0.0/24"}, wantShowVersion: false, wantCIDR: "10.0.0.0/24"},
		{name: "no args", args: []string{}, wantShowVersion: false, wantCIDR: ""},
		{name: "CIDR with version flag", args: []string{"--version", "10.0.0.0/24"}, wantShowVersion: true, wantCIDR: "10.0.0.0/24"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotVersion, gotCIDR := parseArgs(tt.args)
			if gotVersion != tt.wantShowVersion {
				t.Errorf("showVersion = %v, want %v", gotVersion, tt.wantShowVersion)
			}
			if gotCIDR != tt.wantCIDR {
				t.Errorf("cidr = %q, want %q", gotCIDR, tt.wantCIDR)
			}
		})
	}
}
