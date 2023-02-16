package git

import (
	"bytes"
	"fmt"
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
)

type GitStatusEntry struct {
	Filename string
	Type     GitStatusEntryType
	Selected bool
}

type GitDiff struct {
	Chunks []GitDiffChunk
}

type GitDiffChunk struct {
	Old GitDiffFile
	New GitDiffFile
}

type GitDiffFile struct {
	StartLine uint32
	EndLine   uint32
	Lines     []GitDiffLine
}

type GitDiffLine struct {
	Text string
	Type GitDiffLineType
}

func (diff *GitDiff) ToString() string {
	var sb strings.Builder

	for _, chunk := range diff.Chunks {
		sb.WriteString(fmt.Sprintf("Chunk: old (%d - %d) | new (%d - %d)\n", chunk.Old.StartLine, chunk.Old.EndLine, chunk.New.StartLine, chunk.New.EndLine))
		sb.WriteString("NEW:\n")
		for _, line := range chunk.New.Lines {
			lineType := " "
			if line.Type == GIT_LINE_NEW {
				lineType = "+"
			} else {
				lineType = "-"
			}

			sb.WriteString(fmt.Sprintf("%s %s\n", lineType, line.Text))
		}

		sb.WriteString("OLD:\n")
		for _, line := range chunk.Old.Lines {
			lineType := " "
			if line.Type == GIT_LINE_NEW {
				lineType = "+"
			} else {
				lineType = "-"
			}

			sb.WriteString(fmt.Sprintf("%s %s\n", lineType, line.Text))
		}
	}

	return sb.String()
}

// How to read diff output
// https://stackoverflow.com/questions/27508982/interpreting-git-diff-output

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

func CreateRepository(pathToRepo string) {
	executeGit([]string{"init"}, pathToRepo)
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
