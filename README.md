GIT clone operator
------------------

Mutation webhook handler, injects initContainer with a properly configured `git` container that does a clone or checkout, taking care about permissions.


Use case
--------

With this approach you can still use official application images like Wordpress image without having to maintain your own images
just to inject a theme, plugins or some other custom files.

Roadmap
-------

### v1

- [ ] Injecting git-clone initContainers into labelled pods
- [ ] Support for Git over HTTPS
- [ ] Specifying user id (owner) of files in workspace
- [ ] CLI command `git-clone-operator clone ...` and single Dockerfile for both initContainer and operator
- [ ] Helm

### v2

- [ ] Namespaced CRD `GitClonePermissions` to specify which GIT repositories are allowed, where are the clone keys
- [ ] Reacting on webhooks from Gitea and GitHub to update revision on existing pods (using `kind: Job` with cloned volume definitions from `kind: Pod`, using the same configuration as initContainer)

Thanks
------

Special thanks to [simple-kubernetes-webhook](https://github.com/slackhq/simple-kubernetes-webhook) project, without it such quick development of this project would not be possible.
