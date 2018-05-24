# -*- mode: ruby -*-
# vi:set ft=ruby sw=2 ts=2 sts=2:
# Create vm
# master: 192.168.26.10
# node-1: 192.168.26.11
# node-2: 192.168.26.12
# ...

NUM_NODE = 2
MASTER_IP = "192.168.26.10"
NODE_IP_NW = "192.168.26."
DISK_SIZE = "100" # 100MB
PRIVATE_KEY_PATH = "inventory/vagrant"
PRIVATE_KEY = "#{PRIVATE_KEY_PATH}/id_rsa"

Vagrant.configure("2") do |config|

  # linkernetworks/aurora-base is a virtualbox with pre-install pacakges:
  #   docker-ce=17.12.1~ce-0~ubuntu 
  #   pip, glusterfs
  #config.vm.box = "linkernetworks/aurora-base"
  #config.vm.box_version = "0.0.6"

  config.vm.box = "ubuntu/xenial64"
  config.vm.box_version = "20180522.0.0"

  config.vm.box_check_update = false
  config.vbguest.auto_update = false

  config.vm.provider "virtualbox" do |vb|
    vb.memory = 2048  # 1536 at least for kubespray
    vb.cpus = 2
  end #vb

  # Generate ssh key at .ssh
  unless File.exist?("#{PRIVATE_KEY}")
    `mkdir -p #{PRIVATE_KEY_PATH} && ssh-keygen -f #{PRIVATE_KEY} -t rsa -N ''`
  end
  config.vm.provision "file", source: "#{PRIVATE_KEY}.pub", destination: "id_rsa.pub"
  config.vm.provision "append-public-key", :type => "shell", inline: "cat id_rsa.pub >> ~/.ssh/authorized_keys"

  #config.vm.provision "setup-hosts", :type => "shell", :path => "../scripts/vagrant/setup-hosts" do |s|
  #  s.args = ["enp0s8"]
  #end
  
  config.vm.provision "shell", privileged: true, inline: <<-SHELL
    apt-get update \
      && apt-get install -y python
  SHELL

  config.vm.define "master" do |node|

    node.vm.hostname = "master"
    node.vm.network :private_network, ip: MASTER_IP

    node.vm.provider "virtualbox" do |vb|
      unless File.exist?("disk-master.vmdk")
        vb.customize ["createhd", "--filename", "disk-master.vmdk", "--size", DISK_SIZE]
      end
      vb.customize ['storageattach', :id,  '--storagectl', 'SCSI', '--port', 2, '--device', 0, '--type', 'hdd', '--medium', "disk-master.vmdk"]
    end
  end #master

  (1..NUM_NODE).each do |i|
    config.vm.define "node-#{i}" do |node|

      node.vm.hostname = "node-#{i}"
      node_ip = NODE_IP_NW + "#{10 + i}"
      node.vm.network :private_network, ip: node_ip

      # Add disk for gluster client
      node.vm.provider "virtualbox" do |vb|
        unless File.exist?("disk-#{i}.vmdk")
          vb.customize ["createhd", "--filename", "disk-#{i}.vmdk", "--size", DISK_SIZE]
        end
        vb.customize ['storageattach', :id,  '--storagectl', 'SCSI', '--port', 2, '--device', 0, '--type', 'hdd', '--medium', "disk-#{i}.vmdk"]
      end #vb

    end #node-i
  end #each node

end
