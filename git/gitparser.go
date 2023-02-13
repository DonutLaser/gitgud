package git

import (
	"strconv"
	"strings"
)

func ParseStatus(text string) (result []GitStatusEntry) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		entryType, entryName, _ := strings.Cut(trimmed, " ")

		result = append(result, GitStatusEntry{
			Filename: strings.TrimSpace(entryName),
			Type:     stringToChangeType(entryType),
			Selected: true,
		})
	}

	return
}

func ParseBranches(text string) (result []string) {
	if text == "" {
		return
	}

	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmed := strings.Trim(line, "'")
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	return
}

func ParseDiff(text string) (result GitDiff) {
	if text == "" {
		return
	}

	// First 4 lines are not useful
	lines := strings.Split(text, "\n")

	chunkStart := -1
	for index, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "@@") {
			chunkStart = index
			break
		}
	}

	if chunkStart == -1 {
		return
	}

	lines = lines[chunkStart:]

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "@@") {
			nextChunk := GitDiffChunk{}
			nextChunk.Old = GitDiffFile{}
			nextChunk.New = GitDiffFile{}

			nextChunk.Old.StartLine, nextChunk.Old.EndLine, nextChunk.New.StartLine, nextChunk.New.EndLine = parseChunkRange(trimmed)

			result.Chunks = append(result.Chunks, nextChunk)
		} else if strings.HasPrefix(trimmed, "+") {
			lastChunk := len(result.Chunks) - 1
			result.Chunks[lastChunk].New.Lines = append(result.Chunks[lastChunk].New.Lines, GitDiffLine{
				Text: strings.TrimPrefix(trimmed, "+"),
				Type: GIT_LINE_NEW,
			})
		} else if strings.HasPrefix(trimmed, "-") {
			lastChunk := len(result.Chunks) - 1
			result.Chunks[lastChunk].Old.Lines = append(result.Chunks[lastChunk].Old.Lines, GitDiffLine{
				Text: strings.TrimPrefix(trimmed, "-"),
				Type: GIT_LINE_REMOVED,
			})
		} else if strings.HasPrefix(trimmed, "\\") {
			// Ignore these kinds of lines
		} else {
			lastChunk := len(result.Chunks) - 1
			result.Chunks[lastChunk].Old.Lines = append(result.Chunks[lastChunk].Old.Lines, GitDiffLine{
				Text: trimmed,
				Type: GIT_LINE_UNMODIFIED,
			})
			result.Chunks[lastChunk].New.Lines = append(result.Chunks[lastChunk].New.Lines, GitDiffLine{
				Text: trimmed,
				Type: GIT_LINE_UNMODIFIED,
			})
		}
	}

	return
}

func parseChunkRange(line string) (uint32, uint32, uint32, uint32) {
	rangesOnly := strings.TrimSpace(strings.Trim(line, "@"))

	split := strings.Split(rangesOnly, " ")
	old := split[0]
	new := split[1]

	oldSplit := strings.Split(strings.TrimPrefix(old, "-"), ",")
	newSplit := strings.Split(strings.TrimPrefix(new, "+"), ",")

	var oldStart int
	var oldEnd int
	var newStart int
	var newEnd int
	if len(oldSplit) > 1 {
		startNumber, _ := strconv.Atoi(oldSplit[0])
		oldStart = startNumber
		endNumber, _ := strconv.Atoi(oldSplit[1])
		oldEnd = endNumber
	} else {
		startNumber, _ := strconv.Atoi(oldSplit[0])
		oldStart = startNumber
		oldEnd = oldStart
	}

	if len(newSplit) > 1 {
		startNumber, _ := strconv.Atoi(newSplit[0])
		newStart = startNumber
		endNumber, _ := strconv.Atoi(newSplit[1])
		newEnd = endNumber
	} else {
		startNumber, _ := strconv.Atoi(newSplit[0])
		newStart = startNumber
		newEnd = newStart
	}

	return uint32(oldStart), uint32(oldEnd), uint32(newStart), uint32(newEnd)
}

func stringToChangeType(str string) GitStatusEntryType {
	switch str {
	case "M":
		return GIT_ENTRY_MODIFIED
	case "??":
		return GIT_ENTRY_NEW_UNSTAGED
	case "A":
		return GIT_ENTRY_NEW
	case "D":
		return GIT_ENTRY_DELETED
	default:
		panic("Unreachable")
	}
}
