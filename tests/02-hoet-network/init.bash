if [ -z "$podName" ]; then
    export name=$(date | md5sum | cut -b 1-19)
    export podName="test-pod-$name"

    rm -rf pod.json
    cp pod.info pod.json
    sed -i  "s/@PODNAME@/${podName}/" pod.json
fi
