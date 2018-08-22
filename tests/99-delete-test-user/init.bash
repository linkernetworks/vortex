export JWT_AUTH_TOKEN=$(http --check-status http://127.0.0.1:7890/v1/users/signin < credential 2>/dev/null | jq -r ".message")
