package checkout

import (
	"fmt"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/moby/sys/mountinfo"
	"github.com/pkg/errors"
	"github.com/riotkit-org/git-clone-operator/pkg/context"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/url"
	"os"
	"strings"
	"syscall"
)

type Command struct {
	LogLevel       string
	Path           string
	Url            string
	Username       string
	Token          string
	Revision       string
	IsBare         bool
	CleanUpRemotes bool
}

func (c *Command) Run() error {
	if err := c.checkAndPrepareInputs(); err != nil {
		return errors.Wrap(err, "Validation failed")
	}

	c.inspectEnvironment()

	urlWithCredentials, parseUrlErr := c.getUrlWithCredentials()
	if parseUrlErr != nil {
		return errors.Wrap(parseUrlErr, "Cannot parse GIT url")
	}
	repository, checkoutErr := c.checkout(urlWithCredentials)
	if checkoutErr != nil {
		return errors.Wrap(checkoutErr, "Cannot clone/checkout repository")
	}
	if c.CleanUpRemotes {
		if err := c.cleanUpRemotes(repository); err != nil {
			return errors.Wrap(err, "Clean up error - cannot remove remotes from local repository")
		}
	}

	head, _ := repository.Head()
	logrus.Infof("The local repository is now on '%s', at commit '%s'", head.Name().String(), head.Hash().String())

	return nil
}

// inspectEnvironment is displaying helpful information about the execution environment to help adjust the parameters in case, when the initContainer would fail
func (c *Command) inspectEnvironment() {
	// Permissions - running as user
	logrus.Infof("Running as uid=%v (to adjust set annotation: %s)", os.Getuid(), context.AnnotationFilesOwner)

	// Current working directory and parent directory
	pwd, _ := os.Getwd()
	c.listDirectory(pwd)
	c.listDirectory(pwd + "/../")

	// Mounted volumes in Linux
	logrus.Info("Inspecting volume mount points")
	mounts, err := mountinfo.GetMounts(nil)
	if err != nil {
		logrus.Errorln(err)
	}
	for _, mount := range mounts {
		logrus.Infof("Found volume mounted at '%s' (%s)", mount.Mountpoint, mount.FSType)
	}
}

// listDirectory lists files and directories in given path, the listing includes permissions
func (c *Command) listDirectory(dirPath string) {
	logrus.Infof("Looking around in '%s' (annotation: %s)", dirPath, context.AnnotationGitPath)
	paths, err := ioutil.ReadDir(dirPath)
	if err != nil {
		logrus.Errorln(err)
	}
	for _, path := range paths {
		stat := path.Sys().(*syscall.Stat_t)
		logrus.Infof(" > [%v %v:%v] %s", path.Mode().String(), stat.Uid, stat.Gid, path.Name())
	}
}

// checkout is performing actually a fresh clone of repository, or update of existing repository
func (c *Command) checkout(url string) (*git.Repository, error) {
	if c.isExistingRepository() {
		logrus.Info("Opening existing repository")
		repository, err := git.PlainOpen(c.Path)
		if err != nil {
			return repository, errors.Wrap(err, "Cannot open git repository")
		}

		if err := c.fetch(repository, "origin", url); err != nil {
			return repository, errors.Wrap(err, "Cannot fetch repository (`git fetch`)")
		}

		w, worktreeErr := repository.Worktree()
		if worktreeErr != nil {
			return repository, errors.Wrap(worktreeErr, "Cannot retrieve a work tree for a `git checkout`")
		}

		var branch plumbing.ReferenceName
		var hash plumbing.Hash
		refName, isBranch := c.createReferenceName(repository, c.Revision)
		if isBranch {
			branch = plumbing.ReferenceName(refName)
		} else {
			hash = plumbing.NewHash(refName)
		}

		logrus.Infof("Doing checkout: hash=%v, branch=%v", hash, branch)
		checkoutErr := w.Checkout(&git.CheckoutOptions{
			Hash:   hash,
			Branch: branch,
			Keep:   false,
			Create: false,
		})
		if checkoutErr != nil {
			return repository, errors.Wrap(checkoutErr, "Cannot perform a `git checkout`")
		}

		if isBranch {
			pullErr := w.Pull(&git.PullOptions{
				RemoteName:    "origin",
				ReferenceName: branch,
			})
			if pullErr != nil {
				if !strings.Contains(pullErr.Error(), "up-to-date") {
					return repository, errors.Wrap(pullErr, "Cannot perform `git pull`")
				}
			}
		}

		return repository, nil
	} else {
		logrus.Info("No local repository found, doing clone")

		if _, err := os.Stat(c.Path); errors.Is(err, os.ErrNotExist) {
			logrus.Info("Directory does not exist, creating")
			if err := os.MkdirAll(c.Path, 0755); err != nil {
				return &git.Repository{}, errors.Wrap(err, "Cannot create target directory before doing `git clone`")
			}
		}

		repository, err := git.PlainClone(c.Path, c.IsBare, &git.CloneOptions{
			URL:      url,
			Progress: os.Stdout,
		})
		if err != nil {
			return repository, errors.Wrapf(err, "Cannot clone '%s' into '%s'", c.Url, c.Path)
		}
		return repository, nil
	}
}

// fetch is making sure that the REMOTE is properly connected, then does a fetch on such remote
func (c *Command) fetch(repository *git.Repository, remoteName string, url string) error {
	// make sure the remote is configured
	remotes, _ := repository.Remotes()
	found := false
	for _, remote := range remotes {
		if remote.Config().Name == remoteName {
			found = true
		}
	}
	if !found {
		_, remoteCreationErr := repository.CreateRemote(&config.RemoteConfig{
			Name: remoteName,
			URLs: []string{url},
		})
		if remoteCreationErr != nil {
			return errors.Wrap(remoteCreationErr, "Cannot add remote URL to the local repository (`git remote add`)")
		}
	}

	// fetch from configured remote 'origin'
	fetchErr := repository.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Depth:      0,
	})
	if fetchErr != nil {
		if strings.Contains(fetchErr.Error(), "up-to-date") {
			logrus.Info("GIT metadata is up-to-date with remote")
			return nil
		}

		return errors.Wrap(fetchErr, "Cannot perform `git fetch` on repository")
	}
	logrus.Info("GIT metadata fetched from remote")
	return nil
}

// createReferenceName Creates a full reference name e.g. refs/tags/v4.0.2. Second arguments tells if this is a branch/tag=true or commit=false
func (c *Command) createReferenceName(repository *git.Repository, ref string) (string, bool) {
	if strings.Contains(ref, "refs/") {
		logrus.Debugln("Preserving original refs/")
		return ref, true
	} else if _, tagErr := repository.Tag(ref); tagErr == nil {
		logrus.Debugln("Detected tag")
		return fmt.Sprintf("refs/tags/%s", c.Revision), true
	} else if _, branchErr := repository.Branch(ref); branchErr == nil {
		logrus.Debugln("Detected branch")
		return fmt.Sprintf("refs/heads/%s", c.Revision), true
	}
	logrus.Debugln("Detected commit")
	return ref, false
}

// isExistingRepository detects if repository already exists by checking if ".git" directory exists
func (c *Command) isExistingRepository() bool {
	if _, err := os.Stat(c.Path + "/.git"); errors.Is(err, os.ErrNotExist) {
		return false
	}
	return true
}

// cleanUpRemotes removes all remotes from git repository to prevent token leak
func (c *Command) cleanUpRemotes(repository *git.Repository) error {
	remotes, listingErr := repository.Remotes()
	if listingErr != nil {
		return listingErr
	}
	for _, remote := range remotes {
		if err := repository.DeleteRemote(remote.Config().Name); err != nil {
			logrus.Errorf("Error deleting remote '%s'", remote.Config().Name)
		}
	}
	return nil
}

// getUrlWithCredentials makes sure that credentials are in the URL (token, username)
func (c *Command) getUrlWithCredentials() (string, error) {
	if c.Username == "" || c.Token == "" {
		logrus.Info("No credentials configured, will not be using authorization")
		return c.Url, nil
	}

	u, err := url.Parse(c.Url)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s://%s:%s@%s:%s/%s", u.Scheme, c.Username, c.Token, u.Hostname(), u.Port(), u.Path), nil
}

// checkAndPrepareInputs performs a pre-validation and mutation of input parameters
func (c *Command) checkAndPrepareInputs() error {
	if c.Username == "" {
		if os.Getenv("GIT_USER") != "" {
			c.Username = os.Getenv("GIT_USER")
		} else {
			return errors.New("missing username")
		}
	}
	if c.Token == "" {
		if os.Getenv("GIT_TOKEN") != "" {
			c.Token = os.Getenv("GIT_TOKEN")
		}
	}
	if c.Revision == "" {
		if os.Getenv("GIT_REVISION") != "" {
			c.Revision = os.Getenv("GIT_REVISION")
		} else {
			c.Revision = "main"
		}
	}
	return nil
}
