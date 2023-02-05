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
	output := executeGit("status --porcelain", pathToRepo)
	return ParseStatus(output)
}

func Discard(filename string, pathToRepo string) {
	executeGit(fmt.Sprintf("restore %s", filename), pathToRepo)
}

func DiscardAll(pathToRepo string) {
	executeGit("reset --hard", pathToRepo)
	executeGit("clean -fxd", pathToRepo)
}

func Stash(pathToRepo string) {
	executeGit("stash -u", pathToRepo)
}

func SwitchToBranch(branchName string, pathToRepo string) {
	executeGit(fmt.Sprintf("checkout %s", branchName), pathToRepo)
}

func ListBranches(pathToRepo string) (result []string) {
	output := executeGit("branch -l --format='%(refname:short)'", pathToRepo)
	return ParseBranches(output)
}

func GetCurrentBranch(pathToRepo string) string {
	output := executeGit("branch --show-current", pathToRepo)
	return strings.TrimSpace(output)
}

func DiffEntry(entry GitStatusEntry, pathToRepo string) (result GitDiff) {
	switch entry.Type {
	case GIT_ENTRY_NEW:
		return diffNew(entry.Filename, pathToRepo)
	case GIT_ENTRY_MODIFIED:
		return diff(entry.Filename, pathToRepo)
	case GIT_ENTRY_DELETED:
		return diffRemoved(entry.Filename, pathToRepo)
	default:
		panic("Unreachable")
	}
}

func diff(filename string, pathToRepo string) (result GitDiff) {
	output := executeGit(fmt.Sprintf("diff %s", filename), pathToRepo)
	return ParseDiff(output)
}

func diffNew(filename string, pathToRepo string) (result GitDiff) {
	output := executeGit(fmt.Sprintf("diff --cached %s", filename), pathToRepo)
	return ParseDiff(output)
}

func diffRemoved(filename string, pathToRepo string) (result GitDiff) {
	output := executeGit(fmt.Sprintf("diff -- %s", filename), pathToRepo)
	return ParseDiff(output)
}

func executeGit(command string, cwd string) string {
	var result bytes.Buffer

	args := strings.Split(command, " ")
	cmd := exec.Command("git", args...)
	cmd.Stdout = &result

	if cwd != "" {
		cmd.Dir = cwd
	}

	err := cmd.Run()
	if err != nil {
		panic(err)
	}

	return result.String()
}
