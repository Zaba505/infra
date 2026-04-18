package service

type MachineRequest struct {
	CPUs          []CPU          `firestore:"cpus"`
	MemoryModules []MemoryModule `firestore:"memory_modules"`
	Accelerators  []Accelerator  `firestore:"accelerators"`
	NICs          []NIC          `firestore:"nics"`
	Drives        []Drive        `firestore:"drives"`
}

type CPU struct {
	Manufacturer   string `firestore:"manufacturer"`
	ClockFrequency int64  `firestore:"clock_frequency"`
	Cores          int64  `firestore:"cores"`
}

type MemoryModule struct {
	Size int64 `firestore:"size"`
}

type Accelerator struct {
	Manufacturer string `firestore:"manufacturer"`
}

type NIC struct {
	MAC string `firestore:"mac"`
}

type Drive struct {
	Capacity int64 `firestore:"capacity"`
}
