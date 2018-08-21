#!/usr/bin/env bats

load init

@test "Signup" {
    http -v --check-status 127.0.0.1:7890/v1/users/signup < user.json
    [ $? = 0 ]
}

@test "Signin" {
   http -v --check-status http://127.0.0.1:7890/v1/users/signin < credential 2>/dev/null | jq -r ".message"
   [ $? = 0 ]
}
