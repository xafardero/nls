package app

import (
	"testing"
	"time"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.CIDR != "192.168.1.0/24" {
		t.Errorf("CIDR = %q; want %q", cfg.CIDR, "192.168.1.0/24")
	}

	if cfg.Timeout != 5*time.Minute {
		t.Errorf("Timeout = %v; want %v", cfg.Timeout, 5*time.Minute)
	}

	if !cfg.ShowProgress {
		t.Error("ShowProgress should be true by default")
	}
}

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name: "valid config",
			config: &Config{
				CIDR:         "192.168.1.0/24",
				Timeout:      5 * time.Minute,
				ShowProgress: true,
			},
			wantErr: false,
		},
		{
			name: "invalid CIDR",
			config: &Config{
				CIDR:         "invalid",
				Timeout:      5 * time.Minute,
				ShowProgress: true,
			},
			wantErr: true,
		},
		{
			name: "zero timeout",
			config: &Config{
				CIDR:         "192.168.1.0/24",
				Timeout:      0,
				ShowProgress: true,
			},
			wantErr: true,
		},
		{
			name: "negative timeout",
			config: &Config{
				CIDR:         "192.168.1.0/24",
				Timeout:      -1 * time.Second,
				ShowProgress: true,
			},
			wantErr: true,
		},
		{
			name: "valid different CIDR",
			config: &Config{
				CIDR:         "10.0.0.0/8",
				Timeout:      1 * time.Minute,
				ShowProgress: false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
