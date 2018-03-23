# ecsctl

Show information for Amazon Elastic Container Service(ECS) like the AWS management console.

## Usage

### Get Information from ECS

* get ECS clusters

```console
$ ecsctl get clusters
```

* get ECS services for specified ECS cluster

```console
$ ecsctl get service --cluster default
```

* get ECS services by selecting cluster interactively

```console
$ ecsctl walk
```

### TODO: Run ECS task

* run a task with running serivce task definition

TODO: polling cloudwatch logs and task status

```console
$ ecsctl run --cluster default --service web-service echo hello
```
