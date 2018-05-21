# Main

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

clean: 
	rm -rf *.log **/*.vmdk **/*.retry

submodule:
	git submodule init && git submodule update

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

