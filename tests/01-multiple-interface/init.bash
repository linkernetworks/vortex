nodeName=`kubectl get nodes | grep "Ready" | awk '{print $1}'`
if [ -z "$networkName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"
    export networkName="test-network-$name"
fi
