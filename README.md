# Artifactory Token Controller

## Description

Controller to be deployed in cluster to generate access tokens for artifactory and store them in kubernetes secrets

## Options

```sh
-artifactoryAdminCredentialsSecret string
    artifactory admin credentials secret name (default "artifactory-admin-credentials")
-artifactoryNamespace string
    namespace to look for artifactory instance (default "default")
-artifactoryTokenScope string
    comma separated groups for artifactory token
-artifactoryTokenUserPrefix string
    user prefix for artifactory token (default "gitlab-")
-buildNamespaces value
    comma separated ci build namespaces to monitor (default [build])
-help 
    show this help
```