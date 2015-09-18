# kube-cluster

`kube-cluster` is a wrapper for `kubectl.sh` that provides environment
configuration switch.

## Usage
```ShellSession
$ kube-cluster create env dev --driver=vbox
Creating environment "dev".
Running magic... Let there be more light.

$ kube-cluster get env
environment     driver
dev             vbox
stage1          aws
live-aws        aws
live-gce        gce

$ kube-cluster env dev
Current environment is "dev"

$ kube-cluster get rc,svc,po,no
 ...kubectl.sh output...

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
# gb get github.com/gerred/kube-cluster
```

## Dependencies

 * [GB](http://getgb.io)
 * Installed and configured `kubectl.sh` (version ???)
