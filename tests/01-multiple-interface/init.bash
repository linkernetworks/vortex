nodeName=`kubectl get nodes | grep "Ready" | awk '{print $1}'`
if [ -z "$networkName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"
    export networkName="test-network-$name"
    export ethName="eth-$name"

    rm -rf networks.json pod.json
    cp networks.info networks.json
    cp pod.info pod.json
    sed -i  "s/@NODENAME@/${nodeName}/" networks.json
    sed -i  "s/@NETWORKNAME@/${networkName}/" networks.json
    sed -i  "s/@ETHNAME@/${ethName}/" networks.json
    sed -i  "s/@NETWORKNAME@/${networkName}/" pod.json
    sed -i  "s/@PODNAME@/${podName}/" pod.json
fi
