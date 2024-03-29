# git-clone-controller

![Version: 0.0-latest-main](https://img.shields.io/badge/Version-0.0--latest--main-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.0-latest-main](https://img.shields.io/badge/AppVersion-0.0--latest--main-informational?style=flat-square)

Simple Pod provisioner using GIT as source. Just label your Pods to get an additional initContainer that will clone your repo before Pod will start up.

## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| env.LOG_JSON | bool | `false` |  |
| env.LOG_LEVEL | string | `"debug"` |  |
| health.liveness.attributes.failureThreshold | int | `1` |  |
| health.liveness.enabled | bool | `true` |  |
| health.readiness.attributes | object | `{}` |  |
| health.readiness.enabled | bool | `true` |  |
| image.repository | string | `"ghcr.io/riotkit-org/git-clone-controller"` |  |
| image.tag | string | `""` |  |
| onlyLabelledNamespaces | bool | `false` |  |
| podAnnotations | object | `{}` |  |
| podSecurityContext.fsGroup | int | `65161` |  |
| podSecurityContext.runAsGroup | int | `65161` |  |
| podSecurityContext.runAsNonRoot | bool | `true` |  |
| podSecurityContext.runAsUser | int | `65161` |  |
| replicas | int | `1` |  |
| resources.limits.cpu | int | `1` |  |
| resources.limits.memory | string | `"128Mi"` |  |
| resources.requests.cpu | int | `0` |  |
| resources.requests.memory | string | `"16Mi"` |  |
| service.port | int | `8080` |  |
| service.type | string | `"ClusterIP"` |  |
| serviceAccount.annotations | object | `{}` |  |
| serviceAccount.create | bool | `true` |  |
| serviceAccount.name | string | `"git-clone-controller-sa"` |  |
| tls.createSecret | bool | `true` |  |
| tls.enabled | bool | `true` |  |
| tls.secretName | string | `"git-clone-controller-tls"` |  |
| webhook.failurePolicy | string | `"Fail"` |  |

----------------------------------------------
Autogenerated from chart metadata using [helm-docs v1.11.0](https://github.com/norwoodj/helm-docs/releases/v1.11.0)
