GIT clone operator
==================

Mutation webhook handler, injects initContainer with a properly configured `git` container that does a clone or checkout, taking care about permissions.


Use cases
--------

### Mounting static user-data in official images

With this approach you can still use official application images like Wordpress image without having to maintain your own images
just to inject a theme, plugins or some other custom files.

### Simple administrative jobs without CI system

Simply clone your scripts repository in your pod workspace, execute script and exit.

### Git clone inside CI job

`git-clone-controller checkout` is a CLI command that could be a replacement of `git clone` and `git checkout`.
It's advantage is that it is designed to be running automatic: When repository does not exists, it gets cloned, when exists, then updated with remote.


Setting up
----------

Use helm to install git-clone-controller. For helm values please take a look at [values reference](https://github.com/riotkit-org/git-clone-controller/blob/main/helm/git-clone-controller/values.yaml).

```bash
helm repo add riotkit-org https://riotkit-org.github.io/helm-of-revolution/
helm install my-git-clone-controller riotkit-org/git-clone-controller                
```

Documentation
-------------

### [For more documentation please take a look at Github page](https://github.com/riotkit-org/git-clone-controller)

