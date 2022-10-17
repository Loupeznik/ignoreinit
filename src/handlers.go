package src

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/devfacet/gocmd/v3"
	"github.com/google/go-github/v47/github"
)

const (
	fncInit    = "Init"
	fncReplace = "Replace"
	gitOwner   = "github"
	gitRepo    = "gitignore"
)

func InitHandlers() {
	gocmd.HandleFlag(fncInit, func(cmd *gocmd.Cmd, args []string) error {
		language, location, err := handleParams(cmd.FlagArgs(fncInit)[1:])

		if err != nil {
			return err
		}

		if _, err := os.Stat(path.Join(location, ".gitignore")); errors.Is(err, nil) {
			fmt.Printf(".gitignore already exists in %s", location)
			return err
		}

		err = getIgnore(language, location, true)

		if err != nil {
			fmt.Print(err.Error())
			return err
		}

		fmt.Printf("Created .gitignore in %s\n", location)

		return nil
	})

	gocmd.HandleFlag(fncReplace, func(cmd *gocmd.Cmd, args []string) error {
		language, location, err := handleParams(cmd.FlagArgs(fncReplace)[1:])

		if err != nil {
			return err
		}

		if _, err := os.Stat(path.Join(location, ".gitignore")); errors.Is(err, os.ErrNotExist) {
			fmt.Printf(".gitignore does not exist in %s", location)
			return err
		}

		err = getIgnore(language, location, false)

		if err != nil {
			fmt.Print(err.Error())
			return err
		}

		fmt.Printf("Replaced .gitignore in %s\n", location)

		return nil
	})
}

func handleParams(params []string) (string, string, error) {
	if len(params) == 0 {
		return "", "", errors.New("no arguments supplied")
	}

	if len(params) == 1 {
		params = append(params, ".")
	}

	return params[0], params[1], nil
}

func getIgnore(language string, location string, isNew bool) error {
	client := github.NewClient(nil)
	ctx := context.Background()
	options := &github.RepositoryContentGetOptions{}

	_, directoryContent, _, err := client.Repositories.GetContents(ctx, gitOwner, gitRepo, "/", options)

	if err != nil {
		return err
	}

	var url string

	for _, file := range directoryContent {
		if !strings.EqualFold(strings.Split(*file.Name, ".")[0], language) {
			continue
		}

		url = file.GetPath()

		break
	}

	if url == "" {
		return fmt.Errorf("could not find .gitignore for %s", language)
	}

	reader, _, err := client.Repositories.DownloadContents(ctx, gitOwner, gitRepo, url, options)

	if err != nil {
		return err
	}

	defer reader.Close()

	bytes, err := io.ReadAll(reader)

	if err != nil {
		return err
	}

	var file *os.File
	fullPath := path.Join(location, ".gitignore")

	if isNew {
		file, err = os.Create(fullPath)
	} else {
		file, err = os.OpenFile(fullPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0755)
	}

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.Write(bytes)

	if err != nil {
		return err
	}

	return nil
}
