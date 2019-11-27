# Artifactory Token Controller

## Description

Controller to be deployed in cluster to generate access tokens for artifactory and store them in kubernetes secrets

## Options

```sh
-artifactoryCredentialsSecret string
    artifactory admin credentials secret name (default "artifactory-credentials")
-artifactoryNamespace string
    namespace to look for artifactory instance (default "default")
-artifactoryTokenScope string
    comma separated groups for artifactory token
-artifactoryTokenUserPrefix string
    user prefix for artifactory token (default "gitlab-")
-buildNamespaces value
    comma separated ci build namespaces to monitor (default [build])
-createDockerRegistrySecret
    if you want to create a registry credential secret, instead of a normal access-token
-dockerServer string
    url of the docker server
-secretKey string
    key in the secret containing the token if not docker (default "artifactory-access-token")
-secretName string
    name of the secret containing the token or the docker credentials (default "artifactory-access-token")
-help 
    show this help
```