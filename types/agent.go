package types

// MesosAgent ..
type MesosAgent struct {
	Slaves []struct {
		ID         string `json:"id"`
		Hostname   string `json:"hostname"`
		Port       int    `json:"port"`
		Attributes struct {
		} `json:"attributes"`
		Pid              string  `json:"pid"`
		RegisteredTime   float64 `json:"registered_time"`
		ReregisteredTime float64 `json:"reregistered_time"`
		Resources        struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"resources"`
		UsedResources struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"used_resources"`
		OfferedResources struct {
			Disk float64 `json:"disk"`
			Mem  float64 `json:"mem"`
			Gpus float64 `json:"gpus"`
			Cpus float64 `json:"cpus"`
		} `json:"offered_resources"`
		ReservedResources struct {
		} `json:"reserved_resources"`
		UnreservedResources struct {
			Disk  float64 `json:"disk"`
			Mem   float64 `json:"mem"`
			Gpus  float64 `json:"gpus"`
			Cpus  float64 `json:"cpus"`
			Ports string  `json:"ports"`
		} `json:"unreserved_resources"`
		Active                bool     `json:"active"`
		Deactivated           bool     `json:"deactivated"`
		Version               string   `json:"version"`
		Capabilities          []string `json:"capabilities"`
		ReservedResourcesFull struct {
		} `json:"reserved_resources_full"`
		UnreservedResourcesFull []struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar,omitempty"`
			Role   string `json:"role"`
			Ranges struct {
				Range []struct {
					Begin int `json:"begin"`
					End   int `json:"end"`
				} `json:"range"`
			} `json:"ranges,omitempty"`
		} `json:"unreserved_resources_full"`
		UsedResourcesFull []struct {
			Name   string `json:"name"`
			Type   string `json:"type"`
			Scalar struct {
				Value float64 `json:"value"`
			} `json:"scalar,omitempty"`
			Role           string `json:"role"`
			AllocationInfo struct {
				Role string `json:"role"`
			} `json:"allocation_info"`
			Ranges struct {
				Range []struct {
					Begin int `json:"begin"`
					End   int `json:"end"`
				} `json:"range"`
			} `json:"ranges,omitempty"`
		} `json:"used_resources_full"`
		OfferedResourcesFull []interface{} `json:"offered_resources_full"`
	} `json:"slaves"`
	RecoveredSlaves []interface{} `json:"recovered_slaves"`
}

// MesosAgentContainers ..
type MesosAgentContainers []struct {
	ContainerID  string `json:"container_id"`
	ExecutorID   string `json:"executor_id"`
	ExecutorName string `json:"executor_name"`
	FrameworkID  string `json:"framework_id"`
	Source       string `json:"source"`
	Status       struct {
		ContainerID struct {
			Value string `json:"value"`
		} `json:"container_id"`
	} `json:"status"`
}
