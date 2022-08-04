package types

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
