#!/usr/bin/env bats

@test "Check the httpie" {
    run which http
    [ $status = 0 ]
}

@test "Check the jq" {
    run which jq
    [ $status = 0 ]
}

@test "Check the OpenvSwitch" {
    run which ovs-vsctl
    [ $status = 0 ]
}

@test "Check kubernetes tools" {
    run which kubectl
    [ $status = 0 ]
}

@test "Check kubernetes cluster" {
    run kubectl get nodes
    [ $status = 0 ]
    [[ ${lines[0]} != "No resources found." ]]
}

@test "Check NFS server setting" {
    showmount -e | grep /nfsshare ; echo $?
    [ $? = 0 ]
}

@test "Waiting vortex server startup" {
    NEXT_WAIT_COUNT=0
    WAIT_LIMIT=5
    until http -hd 127.0.0.1:7890/v1/version 2>&1 | grep HTTP/ || [ $NEXT_WAIT_COUNT -eq $WAIT_LIMIT ]; do
       echo "Waiting server startup"
       sleep 0.4
       NEXT_WAIT_COUNT=$((NEXT_WAIT_COUNT+ 1))
    done
    [ $NEXT_WAIT_COUNT != $WAIT_LIMIT ]
}
