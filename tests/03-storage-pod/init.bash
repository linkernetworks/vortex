nodeName=`kubectl get nodes | grep "Ready" | awk '{print $1}'`
if [ -z "$podName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"
    export nfsstoragename="test-nfs-storage-$name"
    #nfsIP=`kubectl get nodes -o jsonpath="{.items[0].status.addresses[0].address}"`
    nfsIP="172.17.8.100"

    for i in storage volume pod; do
        rm -f $i.json
        cp $i.info $i.json
    done

    sed -i  "s/@NFSIP@/${nfsIP}/" storage.json
    sed -i  "s/@NFSSTORAGENAME@/${nfsstoragename}/" storage.json
    sed -i  "s/@PODNAME@/${podName}/" pod.json
fi
