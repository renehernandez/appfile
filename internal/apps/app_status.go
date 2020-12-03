package apps

type AppStatus struct {
	Name string

	Status       DeploymentStatus
	DeploymentID string
	UpdatedAt    string
	URL          string
}

type DeploymentStatus string

const (
	DeploymentStatusUnknown    DeploymentStatus = "unknown"
	DeploymentStatusDeployed   DeploymentStatus = "deployed"
	DeploymentStatusInProgress DeploymentStatus = "in progress"
)
