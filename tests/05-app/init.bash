if [ -z "$appName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export appName="test-app-$name"
    export svcName="test-app-svc-$name"

    rm -rf app.json
    cp app.info app.json
    sed -i  "s/@DEPLOYMENTNAME@/${appName}/" app.json
    sed -i  "s/@SERVICENAME@/${svcName}/" app.json
fi
