package context

const (
	LabelIsEnabled           = "riotkit.org/git-clone-controller"
	AnnotationGitUrl         = "git-clone-controller/url"
	AnnotationGitPath        = "git-clone-controller/path"
	AnnotationCleanUp        = "git-clone-controller/cleanWorkspace"
	AnnotationFilesOwner     = "git-clone-controller/owner"
	AnnotationFilesGroup     = "git-clone-controller/group"
	AnnotationRev            = "git-clone-controller/revision"
	AnnotationSecretName     = "git-clone-controller/secretName"
	AnnotationSecretTokenKey = "git-clone-controller/secretTokenKey"
	AnnotationSecretUserKey  = "git-clone-controller/secretUsernameKey"
)
