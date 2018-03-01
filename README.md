# ecsctl

A convert tool from AWS Elastic Container Service(ECS) to kubernetes manifest.

## Usage

* ecsctl ecs clusters
* ecsctl ecs services --cluster ${cluster}

```console
$ ecsctl ecs services --cluster $(ecsctl ecs clusters | peco)
```
