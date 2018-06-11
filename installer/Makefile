# Main

clean: 
	rm -rf *.log **/*.vmdk **/*.retry

submodule:
	git submodule init && git submodule update

.PHONY: ansible
UNAME := $(shell uname)
ifeq ($(UNAME), Linux)

ansible: submodule
  sudo apt-get upgrade && apt-get update
  sudo apt-get install -y python3 python3-pip jq
  sudo pip3 install yq ansible netaddr
  sudo pip3 install -r kubespray/requirements.txt

else ifeq ($(UNAME), Darwin)

ansible: submodule
  sudo port install jq # sudo brew install jq
  sudo pip3 install yq ansible
  sudo pip3 install -r kubespray/requirements.txt

endif

deploy-%: submodule
	ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		deploy.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-deploy.log

reset-%: submodule
	ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		-e reset_confirmation=yes \
		kubespray/reset.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-reset.log

scale-%: submodule
	ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		kubespray/scale.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-scale.log

upgrade-%: submodule
	ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		kubespray/upgrade-cluster.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-upgrade.log

# Vagrant

vagrant-up:
	vagrant up

vagrant-destroy:
	vagrant destroy -f

vagrant: vagrant-up deploy-vagrant

# GCE

gce-up: submodule
	ANSIBLE_HOST_KEY_CHECKING=false ansible-playbook \
		--inventory=inventory/gce/hosts.ini \
		gce-up.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-gce-up.log

gce: gce-up deploy-gce
