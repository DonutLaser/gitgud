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

	lines := strings.Split(text, "\n")

	chunkStart := -1
	for index, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "@@") {
			chunkStart = index
			break
		}
	}

	if chunkStart == -1 {
		binaryFile := false
		for _, line := range lines {
			if strings.HasPrefix(line, "Binary files") {
				binaryFile = true
				break
			}
		}

		if binaryFile {
			result.Chunks = append(result.Chunks, GitDiffChunk{
				Old: GitDiffFile{BinaryFile: true},
				New: GitDiffFile{BinaryFile: true},
			})
		}

		return
	}

	lines = lines[chunkStart:]

	tempChunksOld := make([][]string, 0)
	tempChunksNew := make([][]string, 0)

	// First parse the chunks into their respective sides without turning them into GitDiffLines
	for _, line := range lines {
		trimmed := strings.ReplaceAll(line, "\t", "    ")
		if trimmed == "" {
			continue
		}

		if strings.HasPrefix(trimmed, "@@") {
			nextChunk := GitDiffChunk{}
			nextChunk.Old = GitDiffFile{}
			nextChunk.New = GitDiffFile{}

			nextChunk.Old.StartLine, nextChunk.Old.EndLine, nextChunk.New.StartLine, nextChunk.New.EndLine = parseChunkRange(trimmed)

			result.Chunks = append(result.Chunks, nextChunk)

			tempChunksOld = append(tempChunksOld, []string{})
			tempChunksNew = append(tempChunksNew, []string{})
		} else if strings.HasPrefix(trimmed, "+") {
			lastChunk := len(tempChunksNew) - 1
			tempChunksNew[lastChunk] = append(tempChunksNew[lastChunk], trimmed)
		} else if strings.HasPrefix(trimmed, "-") {
			lastChunk := len(tempChunksOld) - 1
			tempChunksOld[lastChunk] = append(tempChunksOld[lastChunk], trimmed)
		} else if strings.HasPrefix(trimmed, "\\") {
			// Ignore these kinds of lines
		} else {
			lastChunk := len(tempChunksNew) - 1
			tempChunksNew[lastChunk] = append(tempChunksNew[lastChunk], trimmed)
			tempChunksOld[lastChunk] = append(tempChunksOld[lastChunk], trimmed)
		}
	}

	// Now produce proper chunks
	for chIndex := 0; chIndex < len(tempChunksOld); chIndex += 1 {
		oldIndex := 0
		newIndex := 0

		for true {
			if oldIndex < len(tempChunksOld[chIndex]) && newIndex < len(tempChunksNew[chIndex]) && !strings.HasPrefix(tempChunksOld[chIndex][oldIndex], "-") && !strings.HasPrefix(tempChunksNew[chIndex][newIndex], "+") {
				result.Chunks[chIndex].New.Lines = append(result.Chunks[chIndex].New.Lines, GitDiffLine{
					Text: tempChunksNew[chIndex][newIndex],
					Type: GIT_LINE_UNMODIFIED,
				})
				result.Chunks[chIndex].Old.Lines = append(result.Chunks[chIndex].Old.Lines, GitDiffLine{
					Text: tempChunksOld[chIndex][oldIndex],
					Type: GIT_LINE_UNMODIFIED,
				})

				oldIndex += 1
				newIndex += 1
			} else if oldIndex < len(tempChunksOld[chIndex]) && newIndex < len(tempChunksNew[chIndex]) && strings.HasPrefix(tempChunksOld[chIndex][oldIndex], "-") && strings.HasPrefix(tempChunksNew[chIndex][newIndex], "+") {
				result.Chunks[chIndex].New.Lines = append(result.Chunks[chIndex].New.Lines, GitDiffLine{
					Text: tempChunksNew[chIndex][newIndex][1:],
					Type: GIT_LINE_NEW,
				})
				result.Chunks[chIndex].Old.Lines = append(result.Chunks[chIndex].Old.Lines, GitDiffLine{
					Text: tempChunksOld[chIndex][newIndex][1:],
					Type: GIT_LINE_REMOVED,
				})

				oldIndex += 1
				newIndex += 1
			} else if newIndex < len(tempChunksNew[chIndex]) && (oldIndex >= len(tempChunksOld[chIndex]) || !strings.HasPrefix(tempChunksOld[chIndex][oldIndex], "-")) && strings.HasPrefix(tempChunksNew[chIndex][newIndex], "+") {
				result.Chunks[chIndex].New.Lines = append(result.Chunks[chIndex].New.Lines, GitDiffLine{
					Text: tempChunksNew[chIndex][newIndex][1:],
					Type: GIT_LINE_NEW,
				})
				result.Chunks[chIndex].Old.Lines = append(result.Chunks[chIndex].Old.Lines, GitDiffLine{
					Text: "",
					Type: GIT_LINE_EMPTY,
				})

				newIndex += 1
			} else if oldIndex < len(tempChunksOld[chIndex]) && strings.HasPrefix(tempChunksOld[chIndex][oldIndex], "-") && (newIndex >= len(tempChunksNew[chIndex]) || !strings.HasPrefix(tempChunksNew[chIndex][newIndex], "+")) {
				result.Chunks[chIndex].Old.Lines = append(result.Chunks[chIndex].Old.Lines, GitDiffLine{
					Text: tempChunksOld[chIndex][oldIndex][1:],
					Type: GIT_LINE_REMOVED,
				})
				result.Chunks[chIndex].New.Lines = append(result.Chunks[chIndex].New.Lines, GitDiffLine{
					Text: "",
					Type: GIT_LINE_EMPTY,
				})

				oldIndex += 1
			}

			if oldIndex >= len(tempChunksOld[chIndex]) && newIndex >= len(tempChunksNew[chIndex]) {
				break
			}
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
