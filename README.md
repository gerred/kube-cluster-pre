# Deprecation Notice

This has been moved to pre-planning at https://github.com/gerred/kube-cluster. It's being refactored for simplicity, and will have a much more `docker-machine` like focus.

# kube-cluster

`kube-cluster` tool for launching and managing Kubernetes environments.

## Usage
```ShellSession
$ kube-cluster create env dev --driver=vbox
Creating environment "dev".
Running magic... Let there be more light.
Done.

$ kube-cluster create env stage1 --driver=aws --interactive
Creating environment "stage1".
Give AWS Credential ID: some-aws-id
Give AWS Credential KEY: some-aws-key
How many nodes [5]: 6
Autoscale [Y/n]: n
...
Running magic... Let there be more light.
Done.

$ kube-cluster get env
environment    driver    nodes
dev            vbox      1
stage1         aws       6
live-aws       aws       10
live-gce       gce       20

$ kube-cluster env dev
Current environment is "dev"

$ kube-cluster describe env stage1
Environment: stage1
Driver: AWS
Nodes: 6
Autoscale: No
...

$ kube-cluster get rc,svc,po,no
 ...kubectl output...

```

## Options

 * *create env*: creates a local environment configuration to communicate with a
 kubernetes deployment.

 * *env [name]*: changes to given environment. It reconfigures the tool so all
 calls are sent to the right deployment.

 * *delete env*: removes an existing local configuration.


## Troubleshooting

1. I get the error `environment not set` for any kube-cluster call I make.

It happens because you have not chosen a working environment yet. You can help
yourself out of this situation by updating your shell profile with the variable
which sets a default environment. In case of Bash:

```ShellSession
$ echo export KUBE_CLUSTER_ENVIRONMENT="environment-name" >> ~/.bash_profile
```


## How to Install

```ShellSession
# go get github.com/gerred/kube-cluster/...
```

## Dependencies

## Contributing

 If you want to contribute, you are going to need Go installed and configured
 in your machina. Please, refer to the [CONTRIBUTING](CONTRIBUTING.md) file.
