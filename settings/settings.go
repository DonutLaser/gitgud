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

func (settings *Settings) AddRepo(repoPath string, save bool) {
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

	if save {
		saveSettings(*settings)
	}
}

func (settings *Settings) SetActiveRepo(repoPath string, save bool) {
	if settings.ActiveRepo == repoPath {
		// Avoid write to disk when nothing has changed
		return
	}

	settings.ActiveRepo = repoPath

	if save {
		saveSettings(*settings)
	}
}

func (settings *Settings) SetActiveBranch(branchName string, save bool) {
	if settings.ActiveBranch == branchName {
		// Avoid write to disk when nothing has changed
		return
	}

	settings.ActiveBranch = branchName

	if save {
		saveSettings(*settings)
	}
}

func OpenSettingsInExternalProgram() {
	settingsPath := getSettingsPath()
	open.Start(settingsPath)
}

func LoadSettings() (result Settings) {
	settingsPath := getSettingsPath()

	if !filesystem.DoesFileExist(settingsPath) {
		saveSettings(result)
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

func saveSettings(settings Settings) {
	var sb strings.Builder
	for _, repo := range settings.RepoList {
		sb.WriteString(fmt.Sprintf("repo=%s\n", repo))
	}

	sb.WriteString(fmt.Sprintf("active_repo=%s\n", settings.ActiveRepo))
	sb.WriteString(fmt.Sprintf("active_branch=%s\n", settings.ActiveBranch))

	filesystem.WriteFile(getSettingsPath(), sb.String())
}

func getKeyValuePair(text string) (string, string) {
	split := strings.Split(text, "=")
	return split[0], split[1]
}

func getSettingsPath() string {
	cacheDir, _ := os.UserCacheDir()
	return fmt.Sprintf("%s/gitgud.conf", cacheDir)
}
