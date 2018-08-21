#!/usr/bin/env bats

load init

@test "Signup" {
    http -v --check-status 127.0.0.1:7890/v1/user/signup < user.json
    [ $? = 0 ]
}

@test "Signin" {
   token = `http --check-status http://127.0.0.1:7890/v1/users/signin < credential.json 2>/dev/null | jq -r ".message"`
   [ $status = 0 ]
}
