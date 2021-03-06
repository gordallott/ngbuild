package github

import (
	"errors"
	"fmt"
	"os/exec"

	"github.com/watchly/ngbuild/core"
)

func (g *Github) cloneAndMerge(directory string, config *core.BuildConfig) error {

	baseBranch := config.BaseBranch
	if baseBranch == "" {
		baseBranch = "master"
	}
	script := ""
	if config.GetMetadata("github:BuildType") == "pullrequest" {
		if config.HeadRepo == "" || config.HeadHash == "" || config.BaseRepo == "" {
			return errors.New("Config is not filled out properly")
		}

		pullNumber := config.GetMetadata("github:PullNumber")
		if pullNumber == "" {
			return errors.New("Config is missing a pull request number for a pull request type build")
		}

		script += fmt.Sprintf(`git clone -q %s "%s"; `, config.BaseRepo, directory)
		script += fmt.Sprintf(`cd %s ; `, directory)
		script += fmt.Sprintf(`git fetch origin pull/%s/head:pull-requestMerge ; `, pullNumber)
		script += fmt.Sprintf(`git checkout -q -f %s ; `, config.HeadHash)
		script += fmt.Sprintf(`git merge --no-edit %s ; `, config.BaseBranch)

	} else if config.GetMetadata("github:BuildType") == "commit" {
		if config.BaseRepo == "" || config.BaseHash == "" {
			return errors.New("Config is not filled out properly")
		}

		script += fmt.Sprintf(`git clone -q --branch %s %s "%s"; `, baseBranch, config.BaseRepo, directory)
		script += fmt.Sprintf(` cd %s ; `, directory)
		script += fmt.Sprintf(`git checkout -q -f %s ; `, config.BaseHash)
	}

	cmd := exec.Command("/bin/sh", "-c", "-e", script)
	output, err := cmd.Output()
	if err != nil {
		logcritf("Error cloning repo: \nscript: %s\nstdout: %s", script, string(output))
		return err
	}
	return nil
}
