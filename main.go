package main

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
	"strings"
)

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	log.Error().Str("error: ", err.Error()).Send()
	os.Exit(1)
}

// CheckArgs should be used to ensure the right command line arguments are
// passed before executing an example.
func CheckArgs(arg ...string) {
	if len(os.Args) < len(arg)+1 {
		os.Exit(1)
	}
}

// Basic example of how to clone a repository using clone options.
func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	CheckArgs("<url>")
	url := os.Args[1]

	repoName := getRepoName(url)
	directory := "./repositories/" + repoName

	if _, err := os.Stat(directory); !os.IsNotExist(err) {
		log.Info().Str("directory", directory).Msg("cats")

		r, err := git.PlainOpen(directory)
		CheckIfError(err)

		worktree, err := r.Worktree()
		CheckIfError(err)

		result := new(strings.Builder)

		log.Info().Msg("made it to the pull")

		err = worktree.Pull(&git.PullOptions{RemoteName: "origin", Depth: 1, Progress: result})
		CheckIfError(err)

		log.Info().Msg("made it to after the pull")

		ref, err := r.Head()
		commit, err := r.CommitObject(ref.Hash())

		log.Info().Stringer("commit", commit).Send()

	} else {
		// checkout
		// Clone the given repository to the given directory
		log.Info().Str("git clone", url).Send()

		r, err := git.PlainClone(directory, false, &git.CloneOptions{
			URL:   url,
			Depth: 1,
		})
		CheckIfError(err)

		// ... retrieving the branch being pointed by HEAD
		ref, err := r.Head()
		CheckIfError(err)
		// ... retrieving the commit object
		commit, err := r.CommitObject(ref.Hash())
		CheckIfError(err)

		fmt.Println(commit)
	}

}

func getRepoName(inputUrl string) string {
	urlPathComponents := strings.Split(inputUrl, "/")
	repoName := strings.Replace(urlPathComponents[len(urlPathComponents)-1], ".git", "", -1)
	return repoName
}
