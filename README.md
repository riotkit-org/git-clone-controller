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

Example usage
-------------

Every `Pod` labelled with `riotkit.org/git-clone-operator: "true"` will be processed by `git-clone-operator`.

`Pod` annotations are the place, where `git-clone-operator` short specification is kept.

```yaml
apiVersion: v1
kind: Pod
metadata:
    name: "tagged-pod"
    labels:
        # required: only labelled Pods are processed
        riotkit.org/git-clone-operator: "true"
    annotations:
        # required: commit/tag/branch
        git-clone-operator/revision: main
        # required: http/https url
        git-clone-operator/url: "https://github.com/jenkins-x/go-scm"
        # required: target path, where the repository should be cloned, should be placed on a shared Volume mount point with other containers in same Pod
        git-clone-operator/path: /workspace/source
        # optional: user id (will result in adding `securityContext`), in effect: running `git` as selected user and creating files as selected user
        git-clone-operator/owner: "1000"
        # optional: group id (will result in adding `securityContext`), same behavior as in "git-clone-operator/owner"
        git-clone-operator/group: "1000"
        # optional: `kind: Secret` name from same namespace as Pod is (if not specified, then global defaults from operator will be taken, or no authorization would be used)
        git-clone-operator/secretName: git-secrets
        # optional: entry name in `.data` section of selected `kind: Secret`
        git-clone-operator/secretKey: jenkins-x
spec:
    restartPolicy: Never
    automountServiceAccountToken: false
    containers:
        - command:
              - /bin/sh
              - "-c"
              - "find /workspace/source; ls -la /workspace/source"
          image: busybox:latest
          name: test
          volumeMounts:
              - mountPath: /workspace/source
                name: workspace
    volumes:
        - name: workspace
          emptyDir: {}
          
    # PERMISSIONS:
    #  If `git-clone-operator/owner` and `git-clone-operator/group` specified, then `fsGroup` should have same value there
    #  so the mounted volume would have proper permissions
    securityContext:
        fsGroup: 1000
```

**Running this example:**

```bash
kubectl delete -f docs/examples/pod-tagged.yaml; kubectl apply -f docs/examples/pod-tagged.yaml

# observe checkout
kubectl logs -f tagged-pod -c git-checkout

# check the example command that will print list of files and permissions
kubectl logs -f tagged-pod
```

Behavior
--------

| Circumstances                                                     | Behavior                                                            |
|-------------------------------------------------------------------|---------------------------------------------------------------------|
| Pods NOT marked with `riotkit.org/git-clone-operator: "true"`     | Do Nothing                                                          |
| Pods MARKED with `riotkit.org/git-clone-operator: "true"`         | Process                                                             |
| Missing required annotation                                       | Do not schedule that `Pod`                                          |
| `kind: Secret` was specified, but is invalid                      | Do not schedule that `Pod`                                          |
| Unknown error while processing labelled `Pod`                     | Do not schedule that `Pod`                                          |
| GIT credentials are invalid                                       | Fail inside initContainer and don't let Pod's containers to execute |
| Revision is invalid                                               | Fail inside initContainer and don't let Pod's containers to execute |
| Volume permissions are invalid                                    | Fail inside initContainer and don't let Pod's containers to execute |
| Unknown error while trying to checkout/clone inside initContainer | Fail inside initContainer and don't let Pod's containers to execute |

Security and reliability
------------------------

- Using [distroless](https://github.com/GoogleContainerTools/distroless/#why-should-i-use-distroless-images) **static image (No operating system, it does not contain even glibc)**
- Static golang binary, without dynamic libraries, no dependency on libc
- No dependency on `git` binary, thanks to [go-git](https://github.com/go-git/go-git)
- Namespaced `kind: Secret` are used close to `kind: Pod`
- Admission Webhooks are [limited in scope on API level](./helm/git-clone-operator/templates/mutatingwebhookconfiguration.yaml) - **only labelled Pods are touched**
- Default Pod's securityContext runs as non-root, with high uid/gid, should work on OpenShift
- API is using internally mutual TLS to talk with Kubernetes

Roadmap
-------

### v1

- [x] Injecting git-clone initContainers into labelled pods
- [x] Support for Git over HTTPS
- [x] Specifying user id (owner) of files in workspace
- [x] CLI command `git-clone-operator clone ...` and single Dockerfile for both initContainer and operator
- [x] Helm
- [x] Add configurable security context - runAs and filesystem permissions

### v2

- [ ] Namespaced CRD `GitClonePermissions` to specify which GIT repositories are allowed, where are the clone keys
- [ ] `chmod user:group -R` as an alternative to `securityContext`

### v3

- [ ] Possibly: Reacting on webhooks from Gitea and GitHub to update revision on existing pods (using `kind: Job` with cloned volume definitions from `kind: Pod`, using the same configuration as initContainer)

Thanks
------

Special thanks to [simple-kubernetes-webhook](https://github.com/slackhq/simple-kubernetes-webhook) fantastic project, without it such quick development of this project would not be possible.
