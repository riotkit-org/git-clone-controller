package checkout_test

import (
	"github.com/riotkit-org/git-clone-controller/cmd/checkout"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

// TestCommand_Run makes a simple functional test - clone, then make a checkout to other reference
func TestCommand_Run(t *testing.T) {
	dir, err := ioutil.TempDir("../../.build/", "test-command-run")
	defer os.RemoveAll(dir)
	if err != nil {
		logrus.Fatal(err)
	}

	c := checkout.Command{
		LogLevel:       "info",
		Path:           dir,
		Url:            "https://github.com/riotkit-org/git-clone-controller",
		Username:       "__token__",
		Token:          "",
		Revision:       "main",
		IsBare:         false,
		CleanUpRemotes: true,
	}

	// Step 1: Git clone
	runErr := c.Run()

	assert.Nil(t, runErr)
	assert.FileExists(t, dir+"/Makefile")

	// Step 2: Checkout to commit
	c2 := checkout.Command{
		LogLevel:       "info",
		Path:           dir,
		Url:            "https://github.com/riotkit-org/git-clone-controller",
		Username:       "__token__",
		Token:          "",
		Revision:       "69d09e37b8791d106d6c5a62f47e9db0359452ec", // initial commit, only LICENSE and README.md in repository
		IsBare:         false,
		CleanUpRemotes: true,
	}
	checkoutErr := c2.Run()

	assert.Nil(t, checkoutErr)
	assert.NoFileExists(t, dir+"/Makefile")
	assert.FileExists(t, dir+"/README.md")

	// Step 3: Checkout back to branch (will do a `git pull` and hit `already up-to-date`)
	c3 := checkout.Command{
		LogLevel:       "info",
		Path:           dir,
		Url:            "https://github.com/riotkit-org/git-clone-controller",
		Username:       "__token__",
		Token:          "",
		Revision:       "main", // initial commit, only LICENSE and README.md in repository
		IsBare:         false,
		CleanUpRemotes: true,
	}
	checkoutToBranchErr := c3.Run()

	assert.Nil(t, checkoutToBranchErr)
	assert.FileExists(t, dir+"/Makefile")
}
