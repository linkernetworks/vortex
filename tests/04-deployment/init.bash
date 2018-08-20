if [ -z "$deploymentName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export deploymentName="test-deployment-$name"

    rm -rf deployment.json
    cp deployment.info deployment.json
    sed -i  "s/@DEPLOYMENTNAME@/${deploymentName}/" deployment.json
fi
