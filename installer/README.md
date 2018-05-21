Votex installer
===

Vortex Installer is an ansible playbook to install Vortex platform to server(s).

[TOC]

# Install Life Cycle

0. Obtain GCR key
1. Bring up VMs / prepare bare metal servers
2. Create ssh key for VMs/servers
3. Edit config files `inventory/*/host.ini`
4. Run installer
5. Reset cluster after testing

### Prerequsites

ubuntu
```
apt-get upgrade && apt-get update

apt-get install -y python3 python-pip
pip install ansible netaddr
# Don't use apt ansible package
```

Mac
```
pip install ansible netaddr
```

### Obtain a GCR key (needed by Aurora)

Google Cloud Registry key is the access key which is required to access essential data from Google Cloud. Make sure *5g-vortex-gcr-20180508* is provided by Linkernetwork before deployment.

```
mv 5g-vortex-gcr-20180508 aurora/
```

### Edit inventory file

The `inventory/*/hosts.ini` store information about which server(s) and how to access server(s). For example, edit vagrant host configuration.
```
vim inventory/vagrant/hosts.ini
```

NOTE:

1. Add all hostname of server(s). Ansible will modify hostname of server to match hosts.ini
2. All bare metal server(s) should be listed above with <hostname> and <ip>.
3. All Master server(s) should be listed under [kube-master] and [etcd]. Requires at least 1 Master.
4. All Node servers(s) should be listed under [kube-node]. Requires at least 1 Node.

# Deploy with Vagrant (For test purpose only)

### Edit hosts.ini
```
vim inventory/vagrant/hosts.ini
```

### Bring up vms
```
vagrant up
```

Check vms status
```
vagrant status
```

### Deploy
```
make deploy-vagrant
```

ssh with key
```
vagrant ssh master
```

Destroy vms after test
```
vagrant destroy
```

# Deploy to bare metal servers

1. Prepare bare metal servers
2. Create ssh key for VMs/servers
3. Edit config files `inventory/*/host.ini`
4. Run installer

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
# Generate a ssh key pair under inventory/5g directory
ssh-keygen -f inventory/5g/id_rsa -t rsa -N ''

# Scp the ssh public key to each server. You will need a valid user and password
scp inventory/5g/id_rsa.pub <user>@<server-ip>:.

# Login server, and append the ssh public key to authorized_keys under root account
ssh <user>@<server-ip>

# become root. Sudoer password required.
sudo su

# Make dir and append public key to authorized_keys under root account
mkdir /root/.ssh
cat id_rsa.pub >> /root/.ssh/authorized_keys
exit

#Test your key. You should be able to login server without password promp.
ssh -i inventory/5g/id_rsa root@<server-ip>

# Move on to next server-ip
```

### Deploy

```
make deploy-5g
```

### Test

test kubernetes with kubectl on master

```
kubectl get nodes
```

access aurora with node ip and node port `http://ip:32363`

# Deploy to GCE 

### Create engine instances

Instance specification:

- CPU: 2 vCPU
- Memory: 5 GB
- Zone: asia-east1- (for better GCR access)
- Storage: 
 - boot disk: 100 GB
 - Addition disk: 100 GB (Testing) 400+ GB (Production)
- Network:
 - Tags: open firewall ports for services

### Edit inventory file
```
vim inventory/gce/host.ini
```

NOTE:

1. Add hostname (instance name), ansible_ssh_host (public/private ip), ip (private ip) to node list
2. The installer use `ansible_ssh_host` to connect to instance from provisioner (the server execute scripts). If the provisioner is inside private network, use private ip for better network performance.
3. The `ansible_ssh_host` ip is peripheral. Note that it  will change when restarting server.

### Deploy

```
make deploy-gce
```

If the gce already has a running cluster, consider reset the cluster before installation.

### Reset

```
make reset-gce
```

# Architecture

### Kubespray

Forked repository: https://github.com/linkernetworks/kubespray

version: v2.5.0

Modification:

- change docker edge version from `docker-ce=17.12.1~ce-0~ubuntu-{{ ansible_distribution_release|lower }}` to `docker-ce=17.12.1~ce-0~ubuntu`


