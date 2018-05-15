Votex installer
===

Vortex Installer is an ansible playbook to install Vortex platform to server(s).

# Installer

- [x] Add submodule kubesrpay: tag: v2.5.0 ( change version if needed )
- [x] Add Vagrantfile to run a local vm
- [ ] Add baremetal ansible playbook to install GlusterFS 
- [ ] Add install guide to README.

# Prerequisites

### Obtain a GCR key

Google Cloud Registry key is the access key which is required to access essential data from Google Cloud. Make sure *5g-vortex-gcr-20180508* is provided by Linkernetwork before deployment.

```
mv 5g-vortex-gcr-20180508 aurora/
```

# Deploy with Vagrant (for test purpose only)

```
# bring up vms with kubespray
vagrant up

# Check status
vagrant status

# deploy with inventory/vagrant/*
make deploy-vagrant

# ssh with key
vagrant ssh master

# remember to destroy vms
vagrant destroy
```

# Deploy to bare metal servers

### Edit inventory settings

```
vim inventory/5g/host.ini
```

NOTE:
1. All bare metal server(s) should be listed above with <hostname> and <ip>.
2. All Master server(s) should be listed under [kube-master] and [etcd]. Requires at least 1 Master.
3. All Node servers(s) should be listed under [kube-node]. Requires at least 1 Node.

### SSH Key ###

A root account access to SSH into all servers is required. If you don't have this key, ask your system admin for this key. Or we can also create the SSH key by the following steps:

```
# Generate a ssh key pair under keys/ directory
ssh-keygen -f 5g/id_rsa -t rsa -N ''

# Scp the ssh public key to each server. You will need a valid user and password
scp 5g/id_rsa.pub <user>@<server-ip>:.

# Login server, and append the ssh public key to authorized_keys under root account
ssh <user>@<server-ip>

# become root. Sudoer password required.
sudo su

# Make dir and append public key to authorized_keys under root account
mkdir /root/.ssh
cat id_rsa.pub >> /root/.ssh/authorized_keys
exit

#Test your key. You should be able to login server without password promp.
ssh -i 5g/id_rsa root@<server-ip>

# Move on to next server-ip
```

### Run deploy

```
make deploy-5g
```

### Test

```
# test kubernetes with kubectl 
vagrant ssh master
kubectl get nodes

# access aurora with node ip and node port

```

# Architecture

### Kubespray

version: v2.5.0
