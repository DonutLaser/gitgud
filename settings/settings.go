package settings

import (
	"fmt"
	"os"
	"strings"

	"github.com/DonutLaser/git-client/filesystem"
	"github.com/skratchdot/open-golang/open"
)

type Settings struct {
	RepoList     []string
	ActiveRepo   string
	ActiveBranch string
}

func (settings *Settings) AddRepo(repoPath string) {
	found := false
	for _, repo := range settings.RepoList {
		if repo == repoPath {
			found = true
			break
		}
	}

	if found {
		return
	}

	settings.RepoList = append(settings.RepoList, repoPath)
}

func (settings *Settings) SetActiveRepo(repoPath string) {
	settings.ActiveRepo = repoPath
}

func (settings *Settings) SetActiveBranch(branchName string) {
	settings.ActiveBranch = branchName
}

func (settings *Settings) Save() {
	var sb strings.Builder
	for _, repo := range settings.RepoList {
		sb.WriteString(fmt.Sprintf("repo=%s\n", repo))
	}

	sb.WriteString(fmt.Sprintf("active_repo=%s\n", settings.ActiveRepo))
	sb.WriteString(fmt.Sprintf("active_branch=%s\n", settings.ActiveBranch))

	filesystem.WriteFile(getSettingsPath(), sb.String())

}

func OpenSettingsInExternalProgram() {
	settingsPath := getSettingsPath()
	open.Start(settingsPath)
}

func LoadSettings() (result Settings) {
	settingsPath := getSettingsPath()

	if !filesystem.DoesPathExist(settingsPath) {
		result.Save()
		return
	}

	contents, success := filesystem.ReadFile(settingsPath)
	if !success {
		return
	}

	lines := strings.Split(contents, "\n")

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		key, value := getKeyValuePair(trimmed)

		if key == "repo" {
			result.RepoList = append(result.RepoList, value)
		} else if key == "active_repo" {
			result.ActiveRepo = value
		} else if key == "active_branch" {
			result.ActiveBranch = value
		}
	}

	return
}

func getKeyValuePair(text string) (string, string) {
	split := strings.Split(text, "=")
	return split[0], split[1]
}

func getSettingsPath() string {
	cacheDir, _ := os.UserCacheDir()
	return fmt.Sprintf("%s/gitgud.conf", cacheDir)
}
