# Main

deploy-%: 
	ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		deploy.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-deploy.log

clean: 
	rm -rf *.log **/*.vmdk 

# Vagrant

vagrant-up:
	vagrant up

vagrant-destroy:
	vagrant destroy -f

vagrant: vagrant-up deploy-vagrant

# GCE

gce-up:
	ansible-playbook \
		--inventory=inventory/gce/hosts.ini \
		gce-up.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-gce-up.log

gce: gce-up deploy-gce

