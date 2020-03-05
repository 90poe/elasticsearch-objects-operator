# elasticsearch-operator installation

The contains details of how to install and uninstall elasticsearch-operator

## Requirements

elasticsearch-operator runs on K8S cluster 1.14 and up.

## Download the binary

The elasticsearch-operator docker image is located at [DockerHub](https://hub.docker.com/repository/docker/90poe/elasticsearch-operators).

To install it to your K8S cluster:
1. edit `deploy/secret.yaml` and add URL to your ES cluster
2. Install using `make` and `kubectl`:
```
cd deploy
make install
```
