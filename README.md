GIT clone operator
==================

**WARNING: This project is unreleased, the first release is a WORK IN PROGRESS**

Mutation webhook handler, injects initContainer with a properly configured `git` container that does a clone or checkout, taking care about permissions.


Use cases
--------

### Mounting static user-data in official images

With this approach you can still use official application images like Wordpress image without having to maintain your own images
just to inject a theme, plugins or some other custom files.

### Simple administrative jobs without CI system

Simply clone your scripts repository in your pod workspace, execute script and exit.

### Git clone inside CI job

`git-clone-operator checkout` is a CLI command that could be a replacement of `git clone` and `git checkout`. 
It's advantage is that it is designed to be running automatic: When repository does not exists, it gets cloned, when exists, then updated with remote.

Roadmap
-------

### v1

- [ ] Injecting git-clone initContainers into labelled pods
- [ ] Support for Git over HTTPS
- [ ] Specifying user id (owner) of files in workspace
- [ ] CLI command `git-clone-operator clone ...` and single Dockerfile for both initContainer and operator
- [ ] Helm
- [ ] Add configurable security context - runAs and filesystem permissions

### v2

- [ ] Namespaced CRD `GitClonePermissions` to specify which GIT repositories are allowed, where are the clone keys
- [ ] Reacting on webhooks from Gitea and GitHub to update revision on existing pods (using `kind: Job` with cloned volume definitions from `kind: Pod`, using the same configuration as initContainer)

Thanks
------

Special thanks to [simple-kubernetes-webhook](https://github.com/slackhq/simple-kubernetes-webhook) fantastic project, without it such quick development of this project would not be possible.
