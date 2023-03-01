package git

import (
	"bytes"
	"os/exec"
	"strings"
)

type GitStatusEntryType uint16
type GitDiffLineType uint8

const (
	GIT_ENTRY_MODIFIED GitStatusEntryType = iota
	GIT_ENTRY_NEW_UNSTAGED
	GIT_ENTRY_NEW
	GIT_ENTRY_DELETED
)

const (
	GIT_LINE_UNMODIFIED GitDiffLineType = iota
	GIT_LINE_NEW
	GIT_LINE_REMOVED
	GIT_LINE_EMPTY
)

type GitStatusEntry struct {
	Filename string
	Type     GitStatusEntryType
	Selected bool
}

type GitDiff struct {
	OldChunks []GitDiffFile
	NewChunks []GitDiffFile
}

type GitDiffFile struct {
	StartLine  uint32
	EndLine    uint32
	Lines      []GitDiffLine
	BinaryFile bool
}

type GitDiffLine struct {
	Text string
	Type GitDiffLineType
}

type GitStashEntry struct {
	BranchName string
	Index      string
}

func Status(pathToRepo string) (result []GitStatusEntry) {
	output := executeGit([]string{"status", "--porcelain", "-u"}, pathToRepo)
	return ParseStatus(output)
}

func Discard(filename string, pathToRepo string) {
	executeGit([]string{"restore", filename}, pathToRepo)
}

func DiscardAll(pathToRepo string) {
	executeGit([]string{"reset", "--hard"}, pathToRepo)
	executeGit([]string{"clean", "-fxd"}, pathToRepo)
}

func Stash(pathToRepo string) {
	executeGit([]string{"stash", "-u"}, pathToRepo)
}

func SwitchToBranch(branchName string, pathToRepo string) {
	executeGit([]string{"checkout", branchName}, pathToRepo)
}

func ListBranches(pathToRepo string) (result []string) {
	output := executeGit([]string{"branch", "-l", "--format='%(refname:short)'"}, pathToRepo)
	return ParseBranches(output)
}

func ListStash(pathToRepo string) (result []GitStashEntry) {
	output := executeGit([]string{"stash", "list"}, pathToRepo)
	return ParseStashList(output)
}

func DoesBranchHaveStash(branchName string, stash []GitStashEntry) bool {
	for _, entry := range stash {
		if entry.BranchName == branchName {
			return true
		}
	}

	return false
}

func GetStashIndex(branchName string, stash []GitStashEntry) string {
	for _, entry := range stash {
		if entry.BranchName == branchName {
			return entry.Index
		}
	}

	return ""
}

func GetCurrentBranch(pathToRepo string) string {
	output := executeGit([]string{"branch", "--show-current"}, pathToRepo)
	return strings.TrimSpace(output)
}

func DiffEntry(entry GitStatusEntry, pathToRepo string) (result GitDiff) {
	switch entry.Type {
	case GIT_ENTRY_NEW_UNSTAGED:
		fallthrough
	case GIT_ENTRY_NEW:
		return diffNew(entry.Filename, pathToRepo)
	case GIT_ENTRY_MODIFIED:
		fallthrough
	case GIT_ENTRY_DELETED:
		return diffModified(entry.Filename, pathToRepo)
	default:
		panic("Unreachable")
	}
}

func Commit(entries []GitStatusEntry, message string, pathToRepo string) (result []GitStatusEntry) {
	fileNames := make([]string, 0)

	for _, entry := range entries {
		if entry.Selected {
			if entry.Type == GIT_ENTRY_NEW_UNSTAGED {
				executeGit([]string{"add", entry.Filename}, pathToRepo)
			}

			fileNames = append(fileNames, entry.Filename)
		}
	}

	executeGit([]string{"commit", "-o", strings.Join(fileNames, " "), "-m", message}, pathToRepo)
	return Status(pathToRepo)
}

func UndoLastCommit(pathToRepo string) (result []GitStatusEntry) {
	executeGit([]string{"reset", "--soft", "HEAD~"}, pathToRepo)
	return Status(pathToRepo)
}

func ApplyStash(index string, pathToRepo string) (result []GitStatusEntry) {
	executeGit([]string{"stash", "pop", index}, pathToRepo)
	return Status(pathToRepo)
}

func DeleteStash(index string, pathToRepo string) {
	executeGit([]string{"stash", "drop", index}, pathToRepo)
}

func CreateRepository(pathToRepo string) {
	executeGit([]string{"init"}, pathToRepo)
}

func CreateBranch(branchName string, pathToRepo string) {
	executeGit([]string{"checkout", "-b", branchName}, pathToRepo)
}

func diffNew(filename string, pathToRepo string) (result GitDiff) {
	output := executeGit([]string{"diff", "--no-index", "/dev/null", filename}, pathToRepo)
	return ParseDiff(output)
}

func diffModified(filename string, pathToRepo string) (result GitDiff) {
	output := executeGit([]string{"diff", "HEAD", "--", filename}, pathToRepo)
	return ParseDiff(output)

}

func executeGit(command []string, cwd string) string {
	var result bytes.Buffer
	var er bytes.Buffer

	cmd := exec.Command("git", command...)
	cmd.Stdout = &result
	cmd.Stderr = &er

	if cwd != "" {
		cmd.Dir = cwd
	}

	cmd.Run()

	return result.String()
}
