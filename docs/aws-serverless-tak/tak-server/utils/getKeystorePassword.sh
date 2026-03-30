#!/bin/bash

# Get keystore password from AWS Secrets Manager
get_keystore_password() {
    if [ -n "$KEYSTORE_SECRET_NAME" ]; then
        aws secretsmanager get-secret-value --secret-id "$KEYSTORE_SECRET_NAME" --query SecretString --output text 2>/dev/null | jq -r .password 2>/dev/null || echo "atakatak"
    else
        echo "atakatak"
    fi
}

# Export the password for use in other scripts
export KEYSTORE_PASSWORD=$(get_keystore_password)