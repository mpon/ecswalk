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
$ ecsctl get service -c default
```

* get ECS tasks for specified ECS cluster and service

```console
$ ecsctl get tasks -c default -s web-service
```

### Get Information from ECS Interactively

* get ECS services by selecting cluster interactively

```console
$ ecsctl walk services
```

* get ECS tasks by selecting cluster and service interactively

```console
$ ecsctl walk tasks
```

### TODO: Run ECS task

* [ ] run a task with running serivce task definition
* [ ] polling cloudwatch logs and task status

```console
$ ecsctl run --c default --s web-service echo hello
```
