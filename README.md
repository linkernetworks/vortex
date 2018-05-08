votex
===

Vortex is a networking development platform based on Linker Aurora.

Check [Our dropbox paper](https://paper.dropbox.com/doc/5G-NFV-Information-Architecture-WIGrgN3OhdjGwEHkTRFmZ) for more details.

# Installer

- [x] Add submodule kubesrpay: tag: v2.5.0 ( change version if needed )
- [x] Add Vagrantfile to run a local vm
- [ ] Add baremetal ansible playbook to install aurora
- [ ] Add install guide to README.

# Vagrant

```
# bring up vms with kubespray
vagrant up

# Check status
vagrant status

# ssh with key
vagrant ssh k8s-01

# remember to destroy vms
vagrant destroy -f
```
