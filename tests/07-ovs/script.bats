#!/usr/bin/env bats

load init

@test "Create network" {
    http -v --check-status --auth-type=jwt 127.0.0.1:7890/v1/networks < networks.json
    [ $? = 0 ]
}

@test "List network" {
    run bash -c "http --auth-type=jwt http://127.0.0.1:7890/v1/networks/ 2>/dev/null | jq -r '.[] | select(.name == \"${networkName}\").name'"
    [ "$output" = "${networkName}" ]
    [ $status = 0 ]
}

@test "Check OVS port" {
    run sudo ovs-vsctl get Interface ${ethName} ofport
    [ $status = 0 ]
}

@test "Create Deployment" {
    http -v --check-status --auth-type=jwt 127.0.0.1:7890/v1/deployments < deployment.json
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
   run bash -c "http --auth-type=jwt http://127.0.0.1:7890/v1/deployments/ 2>/dev/null | jq -r '.[] | select(.name == \"${deploymentName}\").name'"
   [ "$output" = "${deploymentName}" ]
   [ $status = 0 ]
}

@test "Get OVS Port Status" {
    run bash -c 'http --auth-type=jwt http://127.0.0.1:7890/v1/networks/ 2>/dev/null | jq -r ".[0].bridgeName"'
    localID=`http --auth-type=jwt "127.0.0.1:7890/v1/ovs/portinfos?nodeName=$nodeName&bridgeName=$output" | jq -r ".[3].portID"`
    [ $? = 0 ]
    [ "$localID" = "-1" ]

    podName=`http --auth-type=jwt "127.0.0.1:7890/v1/ovs/portinfos?nodeName=$nodeName&bridgeName=$output" | jq -r ".[0].podName"`
    echo $podName | grep ${deploymentName}
    [ $? = 0 ]
}

@test "Delete Deployment" {
    run bash -c 'http --auth-type=jwt http://127.0.0.1:7890/v1/deployments/ 2>/dev/null | jq -r ".[0].id"'
    run http --auth-type=jwt DELETE http://127.0.0.1:7890/v1/deployments/${output} 2>/dev/null
    [ $status = 0 ]
}

@test "Delete Network" {
    run bash -c 'http --auth-type=jwt http://127.0.0.1:7890/v1/networks/ 2>/dev/null | jq -r ".[0].id"'
    run http --auth-type=jwt DELETE http://127.0.0.1:7890/v1/networks/${output} 2>/dev/null
    [ $status = 0 ]
}
