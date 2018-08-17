#!/usr/bin/env bats

load init

@test "Create Deployment" {
    http -v --check-status 127.0.0.1:7890/v1/deployments < deployment.json
    [ $? = 0 ]
    #Wait the Deployment
    #jsonpath="{.status.phase}"
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    until kubectl get deployments ${deploymentName} -o jsonpath="{.status.readyReplicas}" | grep "1" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       kubectl get deployments ${deploymentName}
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "List Deployment" {
   run bash -c "http http://127.0.0.1:7890/v1/deployments/ 2>/dev/null | jq -r '.[] | select(.name == \"${deploymentName}\").name'"
   [ "$output" = "${deploymentName}" ]
   [ $status = 0 ]
}

@test "Check Deployment Attribute" {
    run kubectl get deployments ${deploymentName} -o jsonpath='{.spec.template.spec.hostNetwork}'
    [ $status = 0 ]
    [ "$output" = "true" ]
}

@test "Check Deployment Env" {
    kubectl get deployments ${deploymentName} -o jsonpath='{.spec.template.spec.containers[0].env}' | grep "myip"
    [ $? = 0 ]
}
@test "Delete Deployment" {
    run bash -c 'http http://127.0.0.1:7890/v1/deployments/ 2>/dev/null | jq -r ".[0].id"'
    run http DELETE http://127.0.0.1:7890/v1/deployments/${output} 2>/dev/null
    [ $status = 0 ]
}
