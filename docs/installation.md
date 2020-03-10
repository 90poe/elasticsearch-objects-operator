# elasticsearch-objects-operator installation

The contains details of how to install and uninstall elasticsearch-objects-operator

## Requirements

elasticsearch-objects-operator runs on K8S cluster 1.14 and up. To install it you would need:
1. Admin access to cluster
2. `kubectl` which is configured to access your cluster and is in your execution path
3. GNU or *NIX Make which is in your execution path

## Install to K8S cluster

The elasticsearch-objects-operator docker image is located at [DockerHub](https://hub.docker.com/repository/docker/90poe/elasticsearch-objects-operators).

To install it to your K8S cluster:
1. edit `deploy/secret.yaml` and add URL to your ES cluster
2. Install using `make` and `kubectl`:
```
cd deploy
make install
```
