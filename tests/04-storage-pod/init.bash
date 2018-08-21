nodeName=`kubectl get nodes | grep "Ready" | awk '{print $1}'`
if [ -z "$podName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"
    export nfsstorageName="test-nfs-storage-$name"
    export volumeName="test-nfs-volume-$name"
    nfsIP=`kubectl get nodes -o jsonpath="{.items[0].status.addresses[0].address}"`
    #nfsIP="172.17.8.100"

    for i in storage volume pod; do
        rm -f $i.json
        cp $i.info $i.json
    done

    sed -i  "s/@NFSIP@/${nfsIP}/" storage.json
    sed -i  "s/@NFSSTORAGENAME@/${nfsstorageName}/" storage.json
    sed -i  "s/@STORAGENAME@/${nfsstorageName}/" volume.json
    sed -i  "s/@VOLUMENAME@/${volumeName}/" volume.json
    sed -i  "s/@PODNAME@/${podName}/" pod.json
    sed -i  "s/@VOLUMENAME@/${volumeName}/" pod.json
fi
