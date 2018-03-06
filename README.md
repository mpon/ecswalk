# ecsctl

A convert tool from AWS Elastic Container Service(ECS) to kubernetes manifest.

## Usage

### Get Information from ECS

* list ECS services for specified ECS cluster

```console
$ ecsctl get service --cluster default`
```

* describe ECS services for specified ECS cluster

```console
$ ecsctl describe service --cluster default
```


* describe ECS services by selecting cluster interactively

```console
$ ecsctl walk
```

### TODO: Run ECS task

* run a task with running serivce task definition

TODO: polling cloudwatch logs and task status

```console
$ ecsctl run --cluster default --service web-service echo hello
```

### TODO: Convert ECS service to a kubernetes manifest

```console
$ ecsctl convert --cluster default --service --web-service
```
