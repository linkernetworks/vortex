#!/usr/bin/env bats

load init

@test "Create App" {
    http -v --check-status 127.0.0.1:7890/v1/apps < app.json
    [ $? = 0 ]
    #Wait the Deployment
    #jsonpath="{.status.phase}"
    NEXT_WAIT_TIME=0
    WAIT_LIMIT=40
    until kubectl get deployment ${appName} -o jsonpath="{.status.readyReplicas}" | grep "2" || [ $NEXT_WAIT_TIME -eq $WAIT_LIMIT ]; do
       sleep 2
       kubectl get deployment ${appName}
       NEXT_WAIT_TIME=$((NEXT_WAIT_TIME+ 1))
    done
    [ $NEXT_WAIT_TIME != $WAIT_LIMIT ]
}

@test "List Service" {
   run bash -c "http http://127.0.0.1:7890/v1/services 2>/dev/null | jq -r '.[] | select(.name == \"${appName}\").name'"
   [ "$output" = "${appName}" ]
   [ $status = 0 ]
}

@test "List Deployment" {
   run bash -c "http http://127.0.0.1:7890/v1/deployments 2>/dev/null | jq -r '.[] | select(.name == \"${appName}\").name'"
   [ "$output" = "${appName}" ]
   [ $status = 0 ]
}

@test "Check Deployment Label" {
    run kubectl get deployment ${appName} -o jsonpath="{.spec.template.metadata.labels.app}"
    [ "$output" = "${appName}" ]
    [ $status = 0 ]
}

@test "Check Service Selector" {
    run kubectl get svc ${appName} -o jsonpath="{.spec.selector.app}"
    [ "$output" = "${appName}" ]
    [ $status = 0 ]
}

@test "Delete Deployment" {
    run bash -c "http http://127.0.0.1:7890/v1/deployments 2>/dev/null | jq -r '.[] | select(.name == \"${appName}\").id'"
    run http DELETE http://127.0.0.1:7890/v1/deployments/${output} 2>/dev/null
    [ $status = 0 ]
}

@test "Delete Services" {
    run bash -c "http http://127.0.0.1:7890/v1/services 2>/dev/null | jq -r '.[] | select(.name == \"${appName}\").id'"
    run http DELETE http://127.0.0.1:7890/v1/services/${output} 2>/dev/null
    [ $status = 0 ]
}
