if [ -z "$networkName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export deploymentName="test-deploy-$name"
    export networkName="test-network-$name"
    export ethName="eth-$name"
    export nodeName=`kubectl get nodes | grep "Ready" | awk '{print $1}'`

    rm -rf networks.json deployment.json
    cp networks.info networks.json
    cp deployment.info deployment.json
    sed -i  "s/@NODENAME@/${nodeName}/" networks.json
    sed -i  "s/@NETWORKNAME@/${networkName}/" networks.json
    sed -i  "s/@ETHNAME@/${ethName}/" networks.json
    sed -i  "s/@NETWORKNAME@/${networkName}/" deployment.json
    sed -i  "s/@DEPLOYMENTNAME@/${deploymentName}/" deployment.json
    # login
    export JWT_AUTH_TOKEN=$(http --check-status http://127.0.0.1:7890/v1/users/signin < credential 2>/dev/null | jq -r ".message")
fi
