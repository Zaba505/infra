package models

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
			err := ValidateMACAddress(tt.mac)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateMACAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMachineRequest_Validate(t *testing.T) {
	tests := []struct {
		name              string
		req               *MachineRequest
		wantInvalidFields int
	}{
		{
			name: "valid request",
			req: &MachineRequest{
				NICs: []NIC{
					{MAC: "52:54:00:12:34:56"},
				},
			},
			wantInvalidFields: 0,
		},
		{
			name: "missing NICs",
			req: &MachineRequest{
				NICs: []NIC{},
			},
			wantInvalidFields: 1,
		},
		{
			name:              "nil NICs",
			req:               &MachineRequest{},
			wantInvalidFields: 1,
		},
		{
			name: "invalid MAC address format",
			req: &MachineRequest{
				NICs: []NIC{
					{MAC: "invalid-mac"},
				},
			},
			wantInvalidFields: 1,
		},
		{
			name: "empty MAC address",
			req: &MachineRequest{
				NICs: []NIC{
					{MAC: ""},
				},
			},
			wantInvalidFields: 1,
		},
		{
			name: "multiple NICs with one invalid",
			req: &MachineRequest{
				NICs: []NIC{
					{MAC: "52:54:00:12:34:56"},
					{MAC: "invalid"},
					{MAC: "aa:bb:cc:dd:ee:ff"},
				},
			},
			wantInvalidFields: 1,
		},
		{
			name: "multiple invalid MACs",
			req: &MachineRequest{
				NICs: []NIC{
					{MAC: "invalid1"},
					{MAC: "invalid2"},
				},
			},
			wantInvalidFields: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			invalidFields := tt.req.Validate()
			if len(invalidFields) != tt.wantInvalidFields {
				t.Errorf("MachineRequest.Validate() returned %d invalid fields, want %d: %+v",
					len(invalidFields), tt.wantInvalidFields, invalidFields)
			}
		})
	}
}
