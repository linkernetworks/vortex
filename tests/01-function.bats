#!/usr/bin/env bats

load data
setup() {
    sed -i "s/@NODENAME@/${nodeName}/" networks.json
    sed -i "s/@NETWORKNAME@/${networkName}/" networks.json
    sed -i "s/@NETWORKNAME@/${networkName}/" pod.json
    sed -i "s/@PODNAME@/${podName}/" pod.json
}

teardown() {
    git checkout networks.json
    git checkout pod.json
}

@test "Create network" {
    http -v --check-status 127.0.0.1:32326/v1/networks < networks.json
    [ $? = 0 ]
}

@test "List network" {
    run bash -c "http http://127.0.0.1:32326/v1/networks/ 2>/dev/null | jq -r '.[] | select(.name == \"${networkName}\").name'"
    [ "$output" = "${networkName}" ]
    [ $status = 0 ]
}

@test "Create Pod" {
    http -v --check-status 127.0.0.1:32326/v1/pods < pod.json
    [ $? = 0 ]
    #Wait the Pod
    #jsonpath="{.status.phase}"
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    until kubectl get pods ${podName} -o jsonpath="{.status.phase}" | grep "Running" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "List Pod" {
   run bash -c "http http://127.0.0.1:32326/v1/pods/ 2>/dev/null | jq -r '.[] | select(.name == \"${podName}\").name'"
   [ "$output" = "${podName}" ]
   [ $status = 0 ]
}

@test "Delete Pod" {
    run bash -c 'http http://127.0.0.1:32326/v1/pods/ 2>/dev/null | jq -r ".[0].id"'
    run http DELETE http://127.0.0.1:32326/v1/pods/${output} 2>/dev/null
    [ $status = 0 ]
}

@test "Delete Network" {
    run bash -c 'http http://127.0.0.1:32326/v1/networks/ 2>/dev/null | jq -r ".[0].id"'
    run http DELETE http://127.0.0.1:32326/v1/networks/${output} 2>/dev/null
    [ $status = 0 ]
}
