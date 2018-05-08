vagrant-up:
	vagrant up

vagrant-destroy:
	vagrant destroy -f

vagrant: vagrant-up deploy-vagrant

deploy-%: 
	ansible-playbook \
		--inventory=inventory/$*/hosts.ini \
		--become-user=root \
		--private-key=inventory/$*/ssh-key \
		deploy.yml 2>&1 | tee aurora-$(shell date +%F-%H%M%S)-deploy.log

clean: 
	rm *.log **/*.vmdk
