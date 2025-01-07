// Package systemspecs provides functionality to collect and report hardware specifications
// of the compute node. This information is critical for job matching and execution
// capabilities reporting in the Lumino network.
package systemspecs

import (
	"encoding/json"
	"fmt"
	"log"

	nvml "github.com/mindprince/gonvml"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

// SystemSpecs represents the complete hardware specification of a compute node
// including GPU, CPU, and memory information. This data is used during staking
// to inform the network of node capabilities.
type SystemSpecs struct {
	GPU []GPUSpec `json:"gpu,omitempty"`
	CPU *CPUSpec  `json:"cpu,omitempty"`
	Mem string    `json:"mem,omitempty"`
}

type GPUSpec struct {
	Model      string `json:"model"`
	Memory     string `json:"memory"`
	SomeMemory string `json:"somememory"`
}

type CPUSpec struct {
	ModelName      string  `json:"model_name"`
	Cores          int32   `json:"cores"`
	ThreadsPerCore int     `json:"threads_per_core"`
	Mhz            float64 `json:"mhz"`
}

// GetSystemSpecs Collects comprehensive system specifications for the compute node. This function:
// 1. Gathers GPU information including model and memory
// 2. Collects CPU specifications like cores and frequency
// 3. Determines available system memory
// 4. Formats data into JSON for blockchain storage
// Returns error if spec collection fails for any component.
func GetSystemSpecs() (string, error) {
	specs := &SystemSpecs{}

	gpu, err := getGPUSpec()
	if err == nil {
		specs.GPU = gpu
	} else {
		log.Printf("Error getting GPU specs: %v\n", err)
	}

	cpu, err := getCPUSpec()
	if err == nil {
		specs.CPU = cpu
	} else {
		log.Printf("Error getting CPU specs: %v\n", err)
	}

	mem, err := getMemSpec()
	if err == nil {
		specs.Mem = mem
	} else {
		log.Printf("Error getting memory specs: %v\n", err)
	}

	return specs.ToJSON()
}

// getGPUSpec Retrieves detailed GPU specifications using NVML:
// 1. Initializes NVIDIA Management Library
// 2. Enumerates available GPU devices
// 3. Collects memory, model, and capability information
// 4. Handles cleanup of NVML resources
// Returns error if GPU information cannot be retrieved.
func getGPUSpec() ([]GPUSpec, error) {
	err := nvml.Initialize()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize NVML: %v", err)
	}
	defer func() {
		if err := nvml.Shutdown(); err != nil {
			log.Printf("Error shutting down NVML: %v\n", err)
		}
	}()

	deviceCount, err := nvml.DeviceCount()
	if err != nil {
		return nil, fmt.Errorf("failed to get device count: %v", err)
	}

	var gpuSpecs []GPUSpec
	for i := uint(0); i < deviceCount; i++ {
		device, err := nvml.DeviceHandleByIndex(i)
		if err != nil {
			return nil, fmt.Errorf("failed to get device handle for index %d: %v", i, err)
		}

		name, err := device.Name()
		if err != nil {
			return nil, fmt.Errorf("failed to get device name: %v", err)
		}

		memory, someInt, err := device.MemoryInfo()
		if err != nil {
			return nil, fmt.Errorf("failed to get memory info: %v", err)
		}

		// powerLimit, err := device.PowerManagementLimit()
		// if err != nil {
		// 	return nil, fmt.Errorf("failed to get power limit: %v", err)
		// }

		gpuSpecs = append(gpuSpecs, GPUSpec{
			Model:      name,
			Memory:     fmt.Sprintf("%.2f MiB", float64(memory)/(1024*1024)),
			SomeMemory: fmt.Sprintf("%.2f MiB", float64(someInt)/(1024*1024)),
			// PwrLimit: fmt.Sprintf("%.2f W", float64(powerLimit)/1000),
		})
	}

	return gpuSpecs, nil
}

// getCPUSpec Gathers CPU specifications including:
// 1. Model name and architecture
// 2. Core count and threads per core
// 3. CPU frequency information
// Returns error if CPU information is unavailable.
func getCPUSpec() (*CPUSpec, error) {
	cpuInfo, err := cpu.Info()
	if err != nil {
		return nil, err
	}

	cpuCounts, err := cpu.Counts(true)
	if err != nil {
		return nil, err
	}

	// Assuming at least one CPU info is available
	ci := cpuInfo[0]

	threadsPerCore := cpuCounts / int(ci.Cores)

	cpuSpec := &CPUSpec{
		ModelName:      ci.ModelName,
		Cores:          ci.Cores,
		ThreadsPerCore: threadsPerCore,
		Mhz:            ci.Mhz,
	}

	return cpuSpec, nil
}

// getMemSpec Retrieves system memory information:
// 1. Queries total available memory
// 2. Converts to human-readable format (GiB)
// Returns error if memory information cannot be accessed.
func getMemSpec() (string, error) {
	vmStat, err := mem.VirtualMemory()
	if err != nil {
		return "", err
	}

	totalMemGB := float64(vmStat.Total) / (1024 * 1024 * 1024)
	return fmt.Sprintf("%.2f GiB", totalMemGB), nil
}

// Converts system specifications to JSON format for blockchain storage.
// Ensures consistent formatting and handles conversion errors.
func (s *SystemSpecs) ToJSON() (string, error) {
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed to marshal system specs to JSON: %v", err)
	}
	return string(jsonData), nil
}
