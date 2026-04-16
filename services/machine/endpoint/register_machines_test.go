package endpoint

import (
	"testing"
)

func TestValidateMACAddress(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr bool
	}{
		{
			name:    "valid MAC address lowercase",
			mac:     "52:54:00:12:34:56",
			wantErr: false,
		},
		{
			name:    "valid MAC address uppercase",
			mac:     "AA:BB:CC:DD:EE:FF",
			wantErr: false,
		},
		{
			name:    "valid MAC address mixed case",
			mac:     "aA:bB:cC:dD:eE:fF",
			wantErr: false,
		},
		{
			name:    "empty MAC address",
			mac:     "",
			wantErr: true,
		},
		{
			name:    "invalid format - missing colons",
			mac:     "aabbccddeeff",
			wantErr: true,
		},
		{
			name:    "invalid format - wrong separator",
			mac:     "aa-bb-cc-dd-ee-ff",
			wantErr: true,
		},
		{
			name:    "invalid format - too short",
			mac:     "aa:bb:cc:dd:ee",
			wantErr: true,
		},
		{
			name:    "invalid format - too long",
			mac:     "aa:bb:cc:dd:ee:ff:gg",
			wantErr: true,
		},
		{
			name:    "invalid characters",
			mac:     "zz:yy:xx:ww:vv:uu",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateMACAddress(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
