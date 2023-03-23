package checkout_test

import (
	"github.com/riotkit-org/git-clone-controller/cmd/checkout"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

// TestCommand_Run makes a simple functional test - clone, then make a checkout to other reference
func TestCommand_Run(t *testing.T) {
	dir, err := os.MkdirTemp("../../.build/", "test-command-run")
	if err != nil {
		logrus.Fatal(err)
	}
	defer os.RemoveAll(dir)

	c := checkout.Command{
		LogLevel:         "info",
		Path:             dir,
		Url:              "https://github.com/riotkit-org/git-clone-controller",
		Username:         "__token__",
		Token:            "",
		Revision:         "main",
		IsBare:           false,
		CleanUpRemotes:   true,
		CleanUpWorkspace: false,
	}

	// Step 1: Git clone
	runErr := c.Run()

	assert.Nil(t, runErr)
	assert.FileExists(t, dir+"/Makefile")

	// Step 2: Checkout to commit
	c2 := checkout.Command{
		LogLevel:         "info",
		Path:             dir,
		Url:              "https://github.com/riotkit-org/git-clone-controller",
		Username:         "__token__",
		Token:            "",
		Revision:         "69d09e37b8791d106d6c5a62f47e9db0359452ec", // initial commit, only LICENSE and README.md in repository
		IsBare:           false,
		CleanUpRemotes:   true,
		CleanUpWorkspace: false,
	}
	checkoutErr := c2.Run()

	assert.Nil(t, checkoutErr)
	assert.NoFileExists(t, dir+"/Makefile")
	assert.FileExists(t, dir+"/README.md")

	// Step 3: Checkout back to branch (will do a `git pull` and hit `already up-to-date`)
	c3 := checkout.Command{
		LogLevel:         "info",
		Path:             dir,
		Url:              "https://github.com/riotkit-org/git-clone-controller",
		Username:         "__token__",
		Token:            "",
		Revision:         "main", // initial commit, only LICENSE and README.md in repository
		IsBare:           false,
		CleanUpRemotes:   true,
		CleanUpWorkspace: false,
	}
	checkoutToBranchErr := c3.Run()

	assert.Nil(t, checkoutToBranchErr)
	assert.FileExists(t, dir+"/Makefile")
}

func TestCommandRunWithDirtyWorkspaceOnSecondCheckout(t *testing.T) {
	dir, err := os.MkdirTemp("../../.build/", "test-command-run")
	if err != nil {
		logrus.Fatal(err)
	}
	defer os.RemoveAll(dir)

	c := checkout.Command{
		LogLevel:         "info",
		Path:             dir,
		Url:              "https://github.com/riotkit-org/git-clone-controller",
		Username:         "__token__",
		Token:            "",
		Revision:         "main",
		IsBare:           false,
		CleanUpRemotes:   true,
		CleanUpWorkspace: true,
	}

	// Step 1: Git clone
	runErr := c.Run()

	assert.Nil(t, runErr)

	// Assert: Basic file exists
	assert.FileExists(t, dir+"/Makefile")

	// ACTION: Now let's make a dirty change
	assert.Nil(t, os.WriteFile(dir+"/somefile.txt", []byte(""), 0777)) // NEW FILE
	assert.Nil(t, os.Remove(dir+"/Makefile"))                          // EXISTING FILE

	//
	// Stage 2: Git pull
	//
	c2 := checkout.Command{
		LogLevel:         "info",
		Path:             dir,
		Url:              "https://github.com/riotkit-org/git-clone-controller",
		Username:         "__token__",
		Token:            "",
		Revision:         "main",
		IsBare:           false,
		CleanUpRemotes:   true,
		CleanUpWorkspace: true, // NOTICE: This must be enabled
	}

	// Step 1: Git clone
	run2Err := c2.Run()
	assert.Nil(t, run2Err)

	// ASSERT
	assert.FileExists(t, dir+"/Makefile", "This file should be back")
	assert.NoFileExists(t, dir+"/somefile.txt", "This file should be removed by clean up & pull")
}
