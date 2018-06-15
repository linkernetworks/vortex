This directory contains a friendly environment for developers to develop the kubernetes with OpenvSwitch in
local host.

It provides a vagrant file to boot a ubuntu-based VM and it will install the kubernetes and OpenvSwitch in that ubuntu host.

Before use the vagrant, you need to install the vagrant in your host and you can refer to [here](https://www.vagrantup.com/docs/installation/) to learn more about vagrant installation.

Usage about Vagrant
===================
- Boot
```
make up
```

- Clean VM
```
make clean
```

Development Environment
=======================
There're at least two way to use this vagrant to develop the vortex project.
## Method-one
Since the vagrant alreay install the proto-buf and golang into the ubuntu, you can use
`vagrant ssh` to login into the ubuntu and develop the project in that VM.
It also contains the `kubectl` and `docker` command for your development.

## Method-two
If you want to develop the vortex in your local host, you can read the following instruction
to setup your `docker` and `kubectl`.

### Docker
```
export DOCKER_HOST="tcp://172.17.8.100:2376"
export DOCKER_TLS_VERIFY=`
```
Now, type the `docker images` and you will see the docker images in that ubuntu VM.

### Kubectl
After the `make up`, the script will copy the kubenetes config from the VM into the tmp/config.
You can use the following to merge that config with your own config and use the `kubectl config use-context`
to switch the kubernetes cluster for your kubectl.

```
cp ~/.kube/config ~/.kube/config.bk
KUBECONFIG=~/.kube/config:`pwd`/tmp/admin.conf kubectl config view --flatten > tmpconfig
cp tmpconfig ~/.kube/config
```
If there're any error in the step(2), please copy the `config.bk` to restor your kubernets config.

Now, you can use `kubectl config get-contexts` to see the `kubernetes-admin@kubernetes` in the list and then use the `kuubectl config use-context kubernetes-admin@kubernetes` to manipulate the VM's kubernetes.

