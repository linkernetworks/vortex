if [ -z "$appName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export appName="test-app-$name"

    rm -rf app.json
    cp app.info app.json
    sed -i  "s/@APPNAME@/${appName}/" app.json
    export JWT_AUTH_TOKEN=$(http --check-status http://127.0.0.1:7890/v1/users/signin < credential 2>/dev/null | jq -r ".message")
fi
