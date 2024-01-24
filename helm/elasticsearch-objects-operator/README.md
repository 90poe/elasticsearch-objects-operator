# elasticsearch-objects-operator

[elasticsearch-objects-operator](https://github.com/90poe/elasticsearch-objects-operator) is ElasticSearch index and template managing operator.

![Version: 0.1.0](https://img.shields.io/badge/Version-0.1.0-informational?style=flat-square) ![Type: application](https://img.shields.io/badge/Type-application-informational?style=flat-square) ![AppVersion: 0.1.1](https://img.shields.io/badge/AppVersion-0.1.1-informational?style=flat-square)

To use, add ElasticSearchIndex or ElasticSearchTemplate CRD object in K8S.

This chart bootstraps an elasticsearch-objects-operator deployment on a [Kubernetes](http://kubernetes.io) cluster using the [Helm](https://helm.sh) package manager.

## Prerequisites

- Chart version 0.1.x: Kubernetes v1.28+

## Get Repo Info

```
git clone git@github.com:90poe/elasticsearch-objects-operator.git
```

## Install Chart

**Important:** only helm3 is supported

```console
cd helm/elasticsearch-objects-operator
helm install [RELEASE_NAME] .
```

The command deploys elasticsearch-objects-operator on the Kubernetes cluster in the default configuration.

_See [configuration](#configuration) below._

_See [helm install](https://helm.sh/docs/helm/helm_install/) for command documentation._

## Uninstall Chart

```console
helm uninstall [RELEASE_NAME]
```

This removes all the Kubernetes components associated with the chart and deletes the release.

_See [helm uninstall](https://helm.sh/docs/helm/helm_uninstall/) for command documentation._

## Upgrading Chart

```console
helm upgrade [RELEASE_NAME] [CHART] --install
```

_See [helm upgrade](https://helm.sh/docs/helm/helm_upgrade/) for command documentation._

