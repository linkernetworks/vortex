VAGRANT :=$(shell which vagrant)

ifneq ($(VAGRAMT), "")
VAGRANT_SCP := $(shell vagrant plugin list | grep vagrant-scp)
endif

vagrant:
	@if [ "$(VAGRANT)" == "" ] ; then echo "You need to install the vagrant first"; exit 1; fi

vagrant-scp: vagrant
	@if [ "$(VAGRANT_SCP)" == "" ]; then vagrant plugin install vagrant-scp ;fi

up: clean-tmp vagrant-scp
	vagrant up
	mkdir tmp
	vagrant scp vortex:/home/vagrant/.kube/config tmp/admin.conf
clean-tmp:
	rm -rf tmp

clean: clean-tmp
	vagrant destroy -f
