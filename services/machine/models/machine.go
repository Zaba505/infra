package models

import (
	"fmt"
	"regexp"
)

var macAddressRegex = regexp.MustCompile(`^([0-9a-fA-F]{2}:){5}[0-9a-fA-F]{2}$`)

type MachineRequest struct {
	CPUs          []CPU          `json:"cpus"`
	MemoryModules []MemoryModule `json:"memory_modules"`
	Accelerators  []Accelerator  `json:"accelerators"`
	NICs          []NIC          `json:"nics"`
	Drives        []Drive        `json:"drives"`
}

type CPU struct {
	Manufacturer   string `json:"manufacturer"`
	ClockFrequency int64  `json:"clock_frequency"`
	Cores          int64  `json:"cores"`
}

type MemoryModule struct {
	Size int64 `json:"size"`
}

type Accelerator struct {
	Manufacturer string `json:"manufacturer"`
}

type NIC struct {
	MAC string `json:"mac"`
}

type Drive struct {
	Capacity int64 `json:"capacity"`
}

type MachineResponse struct {
	ID string `json:"id"`
}

type InvalidField struct {
	Field  string `json:"field"`
	Reason string `json:"reason"`
}

func (r *MachineRequest) Validate() []InvalidField {
	var invalidFields []InvalidField

	if len(r.NICs) == 0 {
		invalidFields = append(invalidFields, InvalidField{
			Field:  "nics",
			Reason: "at least one NIC is required",
		})
		return invalidFields
	}

	for i, nic := range r.NICs {
		if err := ValidateMACAddress(nic.MAC); err != nil {
			invalidFields = append(invalidFields, InvalidField{
				Field:  fmt.Sprintf("nics[%d].mac", i),
				Reason: err.Error(),
			})
		}
	}

	return invalidFields
}

func ValidateMACAddress(mac string) error {
	if mac == "" {
		return fmt.Errorf("MAC address cannot be empty")
	}
	if !macAddressRegex.MatchString(mac) {
		return fmt.Errorf("invalid MAC address format, expected format: aa:bb:cc:dd:ee:ff")
	}
	return nil
}
