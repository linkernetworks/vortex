#!/usr/bin/env bats

load init

# Signup a testuser for next test cases. must run at first
@test "Signup" {
    http -v --check-status 127.0.0.1:7890/v1/users/signup < user.json
    [ $? = 0 ]
}

@test "Test Signin" {
    http -v --check-status 127.0.0.1:7890/v1/users/signin < credential
    [ $? = 0 ]
}
