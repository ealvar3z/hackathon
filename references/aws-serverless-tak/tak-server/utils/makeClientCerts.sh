#!/bin/bash
new_certs_created=false

# Make sure the current user has Java in their PATH
export PATH="/opt/java/openjdk/bin:$PATH"

# Get keystore password from Secrets Manager
source ./utils/getKeystorePassword.sh 

# Create admin cert if it doesn't exist
if [ ! -f "/opt/tak/certs/files/admin.pem" ]; then
    echo "Generating admin certificate"
    cd /opt/tak/certs
    CAPASS=${KEYSTORE_PASSWORD} ./makeCert.sh client admin

    # Elevate admin cert permissions
    echo "Elevating admin cert permissions"
    cd /opt/tak
    java -jar /opt/tak/utils/UserManager.jar certmod -A /opt/tak/certs/files/admin.pem
    new_certs_created=true
fi

# Create user cert if it doesn't exist
if [ ! -f "/opt/tak/certs/files/user.pem" ]; then
    echo "Generating user certificate"
    cd /opt/tak/certs
    CAPASS=${KEYSTORE_PASSWORD} ./makeCert.sh client user
    
    new_certs_created=true
fi

# upload certs to S3 if new certs were created
if [ "$new_certs_created" = true ]; then
    echo "Uploading certs to S3"
    aws s3 cp /opt/tak/certs/files s3://"$BUCKET_NAME"/TAK/certs/ \
        --recursive \
        --exclude "*" \
        --include "admin.p12" \
        --include "user.p12" \
        --include "truststore-root.p12" \
        --include "ca.pem"
fi