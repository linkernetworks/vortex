#!/usr/bin/env bats

load init

@test "Delete test User" {
    run bash -c 'http --auth-type=jwt http://127.0.0.1:7890/v1/users/ 2>/dev/null | jq -r ".[0].id"'
    run http --auth-type=jwt DELETE http://127.0.0.1:7890/v1/users/${output} 2>/dev/null
    [ $status = 0 ]
}
