if [ -z "$deploymentName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export deploymentName="test-deployment-$name"

    rm -rf deployment.json
    cp deployment.info deployment.json
    sed -i  "s/@DEPLOYMENTNAME@/${deploymentName}/" deployment.json
    export JWT_AUTH_TOKEN=$(http --check-status http://127.0.0.1:7890/v1/users/signin < credential 2>/dev/null | jq -r ".message")
fi
