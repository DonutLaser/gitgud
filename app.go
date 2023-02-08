package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DonutLaser/git-client/filesystem"
	"github.com/DonutLaser/git-client/font"
	"github.com/DonutLaser/git-client/git"
	"github.com/DonutLaser/git-client/image"
	"github.com/DonutLaser/git-client/settings"
	"github.com/veandco/go-sdl2/sdl"
)

type AppMode uint8

const (
	MODE_NORMAL AppMode = iota
	MODE_DELETE
	MODE_STASH
)

type Repo struct {
	Name          string
	Path          string
	CurrentBranch string
	Branches      []string
	Changes       []git.GitStatusEntry
}

type App struct {
	Statusbar Statusbar
	Staging   Staging
	DiffView  DiffView
	Search    QuickSearch
	Commit    Commit
	NoRepos   NoRepos
	NoChanges NoChanges

	Mode     AppMode
	Repo     Repo
	Settings settings.Settings
	RepoList []string

	Fonts map[string]font.Font
	Icons map[string]image.Image
}

func NewApp(windowWidth int32, windowHeight int32, renderer *sdl.Renderer) (result App) {
	result.Statusbar = NewStatusbar(windowWidth, windowHeight)
	result.Staging = NewStaging(windowHeight)
	result.DiffView = NewDiffView(windowWidth, windowHeight)
	result.Search = NewQuickSearch(windowWidth, windowHeight)
	result.Commit = NewCommit(windowWidth, windowHeight)
	result.NoRepos = NewNoRepos(windowWidth, windowHeight)
	result.NoChanges = NewNoChanges(windowWidth, windowHeight)

	result.Mode = MODE_NORMAL

	result.Fonts = make(map[string]font.Font)
	result.Fonts["12"] = font.LoadFont("./assets/fonts/consola.ttf", 12)
	result.Fonts["14"] = font.LoadFont("./assets/fonts/consola.ttf", 14)
	result.Fonts["16"] = font.LoadFont("./assets/fonts/consola.ttf", 16)

	result.Icons = make(map[string]image.Image)
	result.Icons["repo"] = image.LoadImage("./assets/icons/icon_repo.png", renderer)
	result.Icons["branch"] = image.LoadImage("./assets/icons/icon_branch.png", renderer)
	result.Icons["entry_off"] = image.LoadImage("./assets/icons/icon_entry_off.png", renderer)
	result.Icons["entry_on"] = image.LoadImage("./assets/icons/icon_entry_on.png", renderer)

	result.Settings = settings.LoadSettings()

	result.Refresh()

	return
}

func (app *App) Resize(windowWidth int32, windowHeight int32) {
	app.Statusbar.Resize(windowWidth, windowHeight)
	app.Staging.Resize(windowHeight)
	app.DiffView.Resize(windowWidth, windowHeight)
	app.Search.Resize(windowWidth, windowHeight)
	app.Commit.Resize(windowWidth, windowHeight)
	app.NoRepos.Resize(windowWidth, windowHeight)
	app.NoChanges.Resize(windowWidth, windowHeight)
}

func (app *App) Refresh() {
	if app.Settings.ActiveRepo != "" && app.Repo.Path != app.Settings.ActiveRepo {
		app.setRepository(app.Settings.ActiveRepo)
	} else if len(app.Settings.RepoList) > 0 && app.Repo.Path != app.Settings.RepoList[0] {
		app.setRepository(app.Settings.RepoList[0])
	}
}

func (app *App) Tick(input *Input) {
	if app.Search.Active {
		app.Search.Tick(input)

		return
	}

	if app.Commit.Active {
		app.Commit.Tick(input)

		return
	}

	if app.Mode == MODE_NORMAL {
		app.handleNormalInput(input)
	} else if app.Mode == MODE_DELETE {
		app.handleDeleteInput(input)
	} else if app.Mode == MODE_STASH {
		app.handleStashInput(input)
	} else {
		panic("Unreachable")
	}

	// ctrl + alt + p to select repo
	// ctrl + p to select branch
	// ctrl + n to new repo
	// ctrl + o to import existing repo
	// ctrl + alt + o to clone repo
	// ctrl + shift + o to open repo folder
	// ctrl + r to open pull request
	// ctrl + w to close app
	// ctrl + shift + n to new branch
	// shift + J to scroll diff down
	// shift + K to scroll diff up
	// i to enter commit message
}

func (app *App) Render(renderer *sdl.Renderer) {
	renderer.SetDrawColor(18, 17, 20, 255)
	renderer.Clear()

	if app.Settings.ActiveRepo == "" {
		app.NoRepos.Render(renderer, app)
	} else {
		app.Statusbar.Render(renderer, app)

		if len(app.Repo.Changes) > 0 {
			app.Staging.Render(renderer, app)
			app.DiffView.Render(renderer, app)
		} else {
			app.NoChanges.Render(renderer, app)
		}

		app.Search.Render(renderer, app)
		app.Commit.Render(renderer, app)
	}

	renderer.Present()
}

func (app *App) setMode(mode AppMode) {
	app.Mode = mode
}

func (app *App) handleNormalInput(input *Input) {
	if input.TypedCharacter == 'j' {
		app.Staging.GoToNextEntry()
		activeEntry := app.Staging.GetActiveEntry()

		app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path))
	} else if input.TypedCharacter == 'k' {
		app.Staging.GoToPrevEntry()
		activeEntry := app.Staging.GetActiveEntry()
		app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path))
	} else if input.TypedCharacter == 'v' {
		app.Staging.ToggleEntrySelected()
	} else if input.TypedCharacter == 'V' {
		app.Staging.ToggleAllEntriesSelected()
	} else if input.TypedCharacter == 'd' {
		app.setMode(MODE_DELETE)
	} else if input.TypedCharacter == 's' {
		app.setMode(MODE_STASH)
	} else if input.TypedCharacter == '`' {

	} else if input.TypedCharacter == 'p' {
		if input.Ctrl {
			if input.Alt {
				repoNames := make([]string, len(app.Settings.RepoList))
				for index, repoPath := range app.Settings.RepoList {
					repoNames[index] = filepath.Base(repoPath)
				}

				app.Search.Open("Repository name", repoNames, SEARCH_BEGINS_WITH, func(repoName string) {
					index := -1
					for i, name := range repoNames {
						if name == repoName {
							index = i
							break
						}
					}

					repo := app.Settings.RepoList[index]

					app.setRepository(repo)
					app.Settings.SetActiveRepo(repo, false)
					app.Settings.SetActiveBranch(app.Repo.CurrentBranch, true)
				})
			} else {
				app.Search.Open("Branch name", app.Repo.Branches, SEARCH_INCLUDES, func(branchName string) {
					app.Repo.CurrentBranch = branchName
					app.Statusbar.ShowBranchName(app.Repo.CurrentBranch)

					git.SwitchToBranch(app.Repo.CurrentBranch, app.Repo.Path)

					app.Settings.SetActiveBranch(branchName, true)
				})
			}
		}
	} else if input.TypedCharacter == 'O' {
		if input.Ctrl {
			path, success := filesystem.OpenDirectory()
			if success {
				app.setRepository(path)
				app.Settings.AddRepo(path, false)
				app.Settings.SetActiveRepo(path, false)
				app.Settings.SetActiveBranch(app.Repo.CurrentBranch, true)
			}
		}
	} else if input.TypedCharacter == ':' {
		// Run arbitrary git command
	} else if input.TypedCharacter == '<' {
		if input.Ctrl {
			settings.OpenSettingsInExternalProgram()
		}
	} else if input.TypedCharacter == 'I' {
		app.Commit.Open(func(message string) {
			filesToCommit := make([]string, 0)
			for _, change := range app.Repo.Changes {
				if change.Selected {
					filesToCommit = append(filesToCommit, change.Filename)
				}
			}

			git.Commit(filesToCommit, message, app.Repo.Path)

			app.Repo.Changes = git.Status(app.Repo.Path)
			app.Staging.ShowEntries(app.Repo.Changes)

			activeEntry := app.Staging.GetActiveEntry()
			app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path))
		})
	}
}

func (app *App) handleDeleteInput(input *Input) {
	if input.Escape {
		app.setMode(MODE_NORMAL)
		return
	}

	if input.TypedCharacter == 'd' {
		activeEntry := app.Staging.GetActiveEntry()
		if activeEntry.Type == git.GIT_ENTRY_NEW {
			// `git restore`` doesn't work on new files so we manually delete them
			// because that's what `git restore` would do anyway
			os.Remove(fmt.Sprintf("%s/%s", app.Repo.Path, activeEntry.Filename))
		} else {
			git.Discard(activeEntry.Filename, app.Repo.Path)
		}

		app.Repo.Changes = git.Status(app.Repo.Path)
		app.Staging.ShowEntries(app.Repo.Changes)

		activeEntry = app.Staging.GetActiveEntry()
		app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path))

		app.setMode(MODE_NORMAL)
	} else if input.TypedCharacter == 'a' {
		git.DiscardAll(app.Repo.Path)

		app.Staging.ShowEntries([]git.GitStatusEntry{})
		app.setMode(MODE_NORMAL)
	} else if input.TypedCharacter == 's' {
		// Delete stash
	}
}

func (app *App) handleStashInput(input *Input) {
	if input.Escape {
		app.setMode(MODE_NORMAL)
		return
	}

	if input.TypedCharacter == 'a' {
		// Stash everything
	}
}

func (app *App) setRepository(repoPath string) {
	app.Repo.Name = filepath.Base(repoPath)
	app.Repo.Path = repoPath
	app.Repo.CurrentBranch = git.GetCurrentBranch(app.Repo.Path)
	app.Repo.Branches = git.ListBranches(app.Repo.Path)
	app.Repo.Changes = git.Status(app.Repo.Path)

	app.Statusbar.ShowRepoName(app.Repo.Name)
	app.Statusbar.ShowBranchName(app.Repo.CurrentBranch)

	app.Staging.ShowEntries(app.Repo.Changes)

	activeEntry := app.Staging.GetActiveEntry()
	app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path))
}
