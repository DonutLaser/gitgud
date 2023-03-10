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
	"github.com/skratchdot/open-golang/open"
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
	Stash         []git.GitStashEntry
}

type App struct {
	Statusbar    Statusbar
	Staging      Staging
	DiffView     DiffView
	Search       QuickSearch
	CommandInput CommandInput
	NoRepos      NoRepos
	NoChanges    NoChanges

	Mode     AppMode
	Repo     Repo
	Settings settings.Settings
	RepoList []string

	Fonts map[string]font.Font
	Icons map[string]image.Image

	Quit        bool
	Initialized bool
}

func NewApp(windowWidth int32, windowHeight int32, renderer *sdl.Renderer) (result App) {
	result.Statusbar = NewStatusbar(windowWidth, windowHeight)
	result.Staging = NewStaging(windowHeight)
	result.DiffView = NewDiffView(windowWidth, windowHeight)
	result.Search = NewQuickSearch(windowWidth, windowHeight)
	result.CommandInput = NewCommandInput(windowWidth, windowHeight)
	result.NoRepos = NewNoRepos(windowWidth, windowHeight)
	result.NoChanges = NewNoChanges(windowWidth, windowHeight)

	result.Mode = MODE_NORMAL

	result.Fonts = make(map[string]font.Font)
	result.Fonts["12"] = font.LoadFont("./assets/fonts/consola.ttf", 12)
	result.Fonts["14"] = font.LoadFont("./assets/fonts/consola.ttf", 14)
	result.Fonts["16"] = font.LoadFont("./assets/fonts/consola.ttf", 16)
	result.Fonts["24"] = font.LoadFont("./assets/fonts/consola.ttf", 24)

	result.Icons = make(map[string]image.Image)
	result.Icons["repo"] = image.LoadImage("./assets/icons/icon_repo.png", renderer)
	result.Icons["branch"] = image.LoadImage("./assets/icons/icon_branch.png", renderer)
	result.Icons["entry_off"] = image.LoadImage("./assets/icons/icon_entry_off.png", renderer)
	result.Icons["entry_on"] = image.LoadImage("./assets/icons/icon_entry_on.png", renderer)
	result.Icons["stash"] = image.LoadImage("./assets/icons/icon_stash.png", renderer)

	result.Settings = settings.LoadSettings()

	result.Quit = false

	return
}

func (app *App) Close() {
	font := app.Fonts["12"]
	font.Unload()
	font = app.Fonts["14"]
	font.Unload()
	font = app.Fonts["16"]
	font.Unload()
	font = app.Fonts["24"]
	font.Unload()

	icon := app.Icons["repo"]
	icon.Unload()
	icon = app.Icons["branch"]
	icon.Unload()
	icon = app.Icons["entry_off"]
	icon.Unload()
	icon = app.Icons["entry_on"]
	icon.Unload()
	icon = app.Icons["stash"]
	icon.Unload()
}

func (app *App) Resize(windowWidth int32, windowHeight int32) {
	app.Statusbar.Resize(windowWidth, windowHeight)
	app.Staging.Resize(windowHeight)
	app.DiffView.Resize(windowWidth, windowHeight)
	app.Search.Resize(windowWidth, windowHeight)
	app.CommandInput.Resize(windowWidth, windowHeight)
	app.NoRepos.Resize(windowWidth, windowHeight)
	app.NoChanges.Resize(windowWidth, windowHeight)
}

func (app *App) Refresh() {
	if app.Initialized {
		if app.Settings.ActiveRepo == "" {
			return
		}

		app.Repo.Branches = git.ListBranches(app.Repo.Path)
		app.Repo.Changes = git.Status(app.Repo.Path)

		app.Staging.ShowEntries(app.Repo.Changes)

		if len(app.Repo.Changes) > 0 {
			activeEntry := app.Staging.GetActiveEntry()
			app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
		}

		return
	}

	if app.Settings.ActiveRepo != "" {
		app.setRepository(app.Settings.ActiveRepo)
	} else if len(app.Settings.RepoList) > 0 {
		app.setRepository(app.Settings.RepoList[0])
	}

	app.Initialized = true
}

func (app *App) Tick(input *Input) {
	if app.Search.Active {
		app.Search.Tick(input)

		return
	}

	if app.CommandInput.Active {
		app.CommandInput.Tick(input)

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

	// ctrl + alt + o to clone repo
	// ctrl + shift + o to open repo folder
	// ctrl + r to open pull request
	// ctrl + shift + n to new branch
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
		app.CommandInput.Render(renderer, app)
	}

	renderer.Present()
}

func (app *App) setMode(mode AppMode) {
	app.Mode = mode
}

func (app *App) handleNormalInput(input *Input) {
	if input.TypedCharacter == 'j' {
		if len(app.Repo.Changes) > 0 {
			app.Staging.GoToNextEntry()
			activeEntry := app.Staging.GetActiveEntry()

			app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
		}
	} else if input.TypedCharacter == 'k' {
		if len(app.Repo.Changes) > 0 {
			app.Staging.GoToPrevEntry()
			activeEntry := app.Staging.GetActiveEntry()
			app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
		}
	} else if input.TypedCharacter == 'L' {
		app.DiffView.ScrollDown()
	} else if input.TypedCharacter == 'H' {
		app.DiffView.ScrollUp()
	} else if input.TypedCharacter == 'v' {
		if len(app.Repo.Changes) > 0 {
			app.Staging.ToggleEntrySelected()
		}
	} else if input.TypedCharacter == 'V' {
		if len(app.Repo.Changes) > 0 {
			app.Staging.ToggleAllEntriesSelected()
		}
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
					app.Settings.SetActiveRepo(repo)
					app.Settings.SetActiveBranch(app.Repo.CurrentBranch)
					app.Settings.Save()
				})
			} else {
				app.Search.Open("Branch name", app.Repo.Branches, SEARCH_INCLUDES, func(branchName string) {
					git.SwitchToBranch(branchName, app.Repo.Path)
					app.Repo.CurrentBranch = branchName

					app.Settings.SetActiveBranch(branchName)
					app.Settings.Save()

					app.Statusbar.ShowBranchName(app.Repo.CurrentBranch)
					app.Repo.Changes = git.Status(app.Repo.Path)
					app.Repo.Stash = git.ListStash(app.Repo.Path)

					app.Statusbar.ShowStashExists(git.DoesBranchHaveStash(app.Repo.CurrentBranch, app.Repo.Stash))

					app.Staging.ShowEntries(app.Repo.Changes)

					if len(app.Repo.Changes) > 0 {
						activeEntry := app.Staging.GetActiveEntry()
						app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
					}
				})
			}
		}
	} else if input.TypedCharacter == 'o' {
		if input.Ctrl {
			app.CommandInput.Open("Path to repository folder", func(folderPath string) {
				app.setRepository(folderPath)
				app.Settings.AddRepo(folderPath)
				app.Settings.SetActiveRepo(folderPath)
				app.Settings.SetActiveBranch(app.Repo.CurrentBranch)
				app.Settings.Save()
			})
		}
	} else if input.TypedCharacter == 'O' {
		if input.Ctrl {
			open.Start(app.Settings.ActiveRepo)
		}
	} else if input.TypedCharacter == ':' {
		app.CommandInput.Open("Command", func(string) {
		})
	} else if input.TypedCharacter == '<' {
		if input.Ctrl {
			settings.OpenSettingsInExternalProgram()
		}
	} else if input.TypedCharacter == 'I' {
		app.CommandInput.Open("Commit message", func(message string) {
			app.Repo.Changes = git.Commit(app.Repo.Changes, message, app.Repo.Path)
			app.Staging.ShowEntries(app.Repo.Changes)

			if len(app.Repo.Changes) > 0 {
				activeEntry := app.Staging.GetActiveEntry()
				app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
			}
		})
	} else if input.TypedCharacter == 'u' {
		app.Repo.Changes = git.UndoLastCommit(app.Repo.Path)
		app.Staging.ShowEntries(app.Repo.Changes)

		activeEntry := app.Staging.GetActiveEntry()
		app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
	} else if input.TypedCharacter == 'n' {
		if input.Ctrl {
			app.CommandInput.Open("Path to new repository folder", func(folderPath string) {
				if !filesystem.DoesPathExist(folderPath) {
					success := filesystem.CreateDirectory(folderPath)
					if !success {
						return
					}
				}

				git.CreateRepository(folderPath)
				app.setRepository(folderPath)
				app.Settings.AddRepo(folderPath)
				app.Settings.SetActiveRepo(folderPath)
				app.Settings.SetActiveBranch(app.Repo.CurrentBranch)
				app.Settings.Save()
			})
		}
	} else if input.TypedCharacter == 'N' {
		if input.Ctrl {
			app.CommandInput.Open("New branch name", func(branchName string) {
				git.CreateBranch(branchName, app.Repo.Path)
				app.Repo.CurrentBranch = git.GetCurrentBranch(app.Repo.Path)
				app.Repo.Branches = git.ListBranches(app.Repo.Path)
				app.Repo.Changes = git.Status(app.Repo.Path)

				app.Statusbar.ShowRepoName(app.Repo.Name)
				app.Statusbar.ShowBranchName(app.Repo.CurrentBranch)

				app.Staging.ShowEntries(app.Repo.Changes)

				app.Settings.SetActiveBranch(app.Repo.CurrentBranch)
				app.Settings.Save()
			})
		}
	} else if input.TypedCharacter == 'w' {
		if input.Ctrl {
			app.Quit = true
		}
	}
}

func (app *App) handleDeleteInput(input *Input) {
	if input.Escape {
		app.setMode(MODE_NORMAL)
		return
	}

	if input.TypedCharacter == 'd' {
		if len(app.Repo.Changes) > 0 {
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

			if len(app.Repo.Changes) > 0 {
				activeEntry = app.Staging.GetActiveEntry()
				app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
			}
		}

		app.setMode(MODE_NORMAL)
	} else if input.TypedCharacter == 'a' {
		git.DiscardAll(app.Repo.Path)

		app.Staging.ShowEntries([]git.GitStatusEntry{})
		app.setMode(MODE_NORMAL)
	}
}

func (app *App) handleStashInput(input *Input) {
	if input.Escape {
		app.setMode(MODE_NORMAL)
		return
	}

	if input.TypedCharacter == 's' {
		git.Stash(app.Repo.Path)

		app.Repo.Changes = git.Status(app.Repo.Path)
		app.Repo.Stash = git.ListStash(app.Repo.Path)
		app.Staging.ShowEntries(app.Repo.Changes)
		app.Statusbar.ShowStashExists(git.DoesBranchHaveStash(app.Repo.CurrentBranch, app.Repo.Stash))

		app.setMode(MODE_NORMAL)
	} else if input.TypedCharacter == 'a' {
		index := git.GetStashIndex(app.Repo.CurrentBranch, app.Repo.Stash)

		if index != "" {
			app.Repo.Changes = git.ApplyStash(index, app.Repo.Path)
			app.Repo.Stash = git.ListStash(app.Repo.Path)
			app.Staging.ShowEntries(app.Repo.Changes)
			app.Statusbar.ShowStashExists(false)

			if len(app.Repo.Changes) > 0 {
				activeEntry := app.Staging.GetActiveEntry()
				app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
			}
		}

		app.setMode(MODE_NORMAL)
	} else if input.TypedCharacter == 'd' {
		index := git.GetStashIndex(app.Repo.CurrentBranch, app.Repo.Stash)

		if index != "" {
			git.DeleteStash(index, app.Repo.Path)
			app.Statusbar.ShowStashExists(false)

			app.Repo.Stash = git.ListStash(app.Repo.Path)
		}

		app.setMode(MODE_NORMAL)
	}
}

func (app *App) setRepository(repoPath string) {
	app.Repo.Name = filepath.Base(repoPath)
	app.Repo.Path = repoPath
	app.Repo.CurrentBranch = git.GetCurrentBranch(app.Repo.Path)
	app.Repo.Branches = git.ListBranches(app.Repo.Path)
	app.Repo.Changes = git.Status(app.Repo.Path)
	app.Repo.Stash = git.ListStash(app.Repo.Path)

	app.Statusbar.ShowRepoName(app.Repo.Name)
	app.Statusbar.ShowBranchName(app.Repo.CurrentBranch)
	app.Statusbar.ShowStashExists(git.DoesBranchHaveStash(app.Repo.CurrentBranch, app.Repo.Stash))

	app.Staging.ShowEntries(app.Repo.Changes)

	if len(app.Repo.Changes) > 0 {
		activeEntry := app.Staging.GetActiveEntry()
		app.DiffView.ShowDiff(git.DiffEntry(activeEntry, app.Repo.Path), activeEntry)
	}
}
