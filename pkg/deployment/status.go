package deployment

import (
	"fmt"
	"os"
	"path/filepath"
)

const (
	InstallerVersionFilePath = "agent/installer.version"
	ActiveLinkPath           = "oneagent/active"
)

type Status int

const (
	NotDeployed Status = iota // The versioned OneAgent folder does not exist in the target directory
	LinkMissing               // The versioned OneAgent folder exists, but the `active` symlink is missing or points elsewhere
	Deployed                  // The versioned OneAgent folder exists and the `active` symlink points to the deployed OneAgent folder
	Unknown                   // An error occurred during the deployment status check
)

type AgentDeploymentInfo struct {
	Status       Status
	AgentVersion string
	Error        error
}

func NewAgentDeploymentInfo(status Status, version string, err error) AgentDeploymentInfo {
	return AgentDeploymentInfo{Status: status, AgentVersion: version, Error: err}
}

func (ds Status) String() string {
	switch ds {
	case NotDeployed:
		return "Not deployed"
	case LinkMissing:
		return "Deployment is not complete"
	case Deployed:
		return "Deployed"
	default:
		return "Unknown"
	}
}

func CheckAgentDeploymentStatus(sourceBaseDir string, targetBaseDir string) AgentDeploymentInfo {
	agentVersion, err := getAgentVersion(sourceBaseDir)
	if err != nil {
		return NewAgentDeploymentInfo(Unknown, "", fmt.Errorf("failed to determine OneAgent version to deploy: %w", err))
	}

	// check whether the agent directory exists
	agentDirPath := GetAgentFolder(targetBaseDir, agentVersion)

	info, err := os.Stat(agentDirPath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewAgentDeploymentInfo(NotDeployed, agentVersion, nil)
		}

		return NewAgentDeploymentInfo(Unknown, agentVersion, fmt.Errorf("cannot obtain OneAgent directory info: %w", err))
	}

	if !info.IsDir() {
		return NewAgentDeploymentInfo(Unknown, agentVersion, fmt.Errorf("OneAgent deployment target is not a directory: %s", agentDirPath))
	}

	// check whether the oneagent active symlink exists
	activeLink := filepath.Join(targetBaseDir, ActiveLinkPath)

	info, err = os.Lstat(activeLink)
	if err != nil {
		if os.IsNotExist(err) {
			return NewAgentDeploymentInfo(LinkMissing, agentVersion, nil)
		}

		return NewAgentDeploymentInfo(Unknown, agentVersion, fmt.Errorf("cannot obtain OneAgent `active` symlink info: %w", err))
	}

	if info.Mode()&os.ModeSymlink == 0 {
		return NewAgentDeploymentInfo(Unknown, agentVersion, fmt.Errorf("OneAgent `active` is not a symlink: %s", info.Mode().String()))
	}

	activeLinkTarget, err := os.Readlink(activeLink)
	if err != nil {
		return NewAgentDeploymentInfo(Unknown, agentVersion, fmt.Errorf("cannot read OneAgent `active` symlink: %w", err))
	}

	if activeLinkTarget != agentVersion {
		return NewAgentDeploymentInfo(LinkMissing, agentVersion, nil)
	}

	return NewAgentDeploymentInfo(Deployed, agentVersion, nil)
}

// GetAgentFolder returns the absolute path to the specified version of the OneAgent directory
func GetAgentFolder(targetBaseDir string, agentVersion string) string {
	agentFolder := filepath.Join(targetBaseDir, "oneagent", agentVersion)

	return agentFolder
}

// getAgentVersion reads the OneAgent version from the installer.version file
func getAgentVersion(sourceBasePath string) (string, error) {
	versionFilePath := filepath.Join(sourceBasePath, InstallerVersionFilePath)

	version, err := os.ReadFile(versionFilePath)
	if err != nil {
		return "", err
	}

	return string(version), nil
}
