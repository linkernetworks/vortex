if [ -z "$podName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"

    rm -rf pod.json
    cp pod.info pod.json
    sed -i  "s/@PODNAME@/${podName}/" pod.json
    export JWT_AUTH_TOKEN=$(http --check-status http://127.0.0.1:7890/v1/users/signin < credential.json 2>/dev/null | jq -r ".message")
fi
