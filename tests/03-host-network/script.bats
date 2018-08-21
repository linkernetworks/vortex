#!/usr/bin/env bats

load init

@test "Create Pod" {
    http -v --check-status 127.0.0.1:7890/v1/pods < pod.json
    [ $? = 0 ]
    #Wait the Pod
    #jsonpath="{.status.phase}"
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    until kubectl get pods ${podName} -o jsonpath="{.status.phase}" | grep "Running" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       kubectl get pods ${podName}
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "List Pod" {
   run bash -c "http http://127.0.0.1:7890/v1/pods/ 2>/dev/null | jq -r '.[] | select(.name == \"${podName}\").name'"
   [ "$output" = "${podName}" ]
   [ $status = 0 ]
}

@test "Check Pod Attribute" {
    run kubectl get pods ${podName} -o jsonpath='{.spec.hostNetwork}'
    [ $status = 0 ]
    [ "$output" = "true" ]
}

@test "Check Pod Env" {
    kubectl get pods ${podName} -o jsonpath='{.spec.containers[0].env}' | grep "myip"
    [ $? = 0 ]
}
@test "Delete Pod" {
    run bash -c 'http http://127.0.0.1:7890/v1/pods/ 2>/dev/null | jq -r ".[0].id"'
    run http DELETE http://127.0.0.1:7890/v1/pods/${output} 2>/dev/null
    [ $status = 0 ]
}
