package checkout

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetUrlWithCredentials(t *testing.T) {
	c := Command{
		Username: "riotkit",
		Token:    "psst",
	}

	c.Url = "https://git.myexample.org/example/wordpress-theme.git"

	url, err := c.getUrlWithCredentials()
	assert.Nil(t, err)
	assert.Equal(t, "https://riotkit:psst@git.myexample.org:443/example/wordpress-theme.git", url)
}

func TestGetUrlWithCredentials_HTTP(t *testing.T) {
	c := Command{
		Username: "riotkit",
		Token:    "psst",
	}

	c.Url = "http://git.myexample.org/example/wordpress-theme.git"

	url, err := c.getUrlWithCredentials()
	assert.Nil(t, err)
	assert.Equal(t, "http://riotkit:psst@git.myexample.org:80/example/wordpress-theme.git", url)
}

func TestGetUrlWithCredentials_GIT(t *testing.T) {
	c := Command{
		Username: "riotkit",
		Token:    "psst",
	}

	c.Url = "git@github.com:riotkit-org/git-clone-controller.git"

	url, err := c.getUrlWithCredentials()
	assert.Nil(t, err)
	assert.Equal(t, "git@github.com:riotkit-org/git-clone-controller.git", url)
}
