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
