#!/bin/bash

# Certificate renewal script for TAK Server
# This script renews Let's Encrypt certificates and copies them to the application directory

if [ -z "$CERT_DOMAIN" ]; then
    echo "ERROR: CERT_DOMAIN environment variable not set"
    exit 1
fi

# Get keystore password from Secrets Manager
source ./utils/getKeystorePassword.sh

echo "$(date): Starting certificate renewal process for domain: $CERT_DOMAIN"

# Renew certificates
/usr/bin/certbot renew -q

# Check if renewal was successful
if [ $? -eq 0 ]; then
    echo "$(date): Certificate renewal completed successfully"
    
    # Copy renewed certificates to application directory if they exist
    if [ -d "/etc/letsencrypt/live/${CERT_DOMAIN}" ]; then
        echo "$(date): Copying renewed certificates to application directory"
        mkdir -p "/opt/tak/certs/files/${CERT_DOMAIN}/"
        cp "/etc/letsencrypt/live/${CERT_DOMAIN}/"* "/opt/tak/certs/files/${CERT_DOMAIN}/"
        
        # Regenerate PKCS12 format with renewed certs
        if [ -f "/etc/letsencrypt/live/${CERT_DOMAIN}/fullchain.pem" ]; then
            echo "$(date): Regenerating PKCS12 certificate"
            openssl pkcs12 \
                -export \
                -in "/etc/letsencrypt/live/${CERT_DOMAIN}/fullchain.pem" \
                -inkey "/etc/letsencrypt/live/${CERT_DOMAIN}/privkey.pem" \
                -out "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.p12" \
                -name "${CERT_DOMAIN}" \
                -password "pass:${KEYSTORE_PASSWORD}"
            
            # Regenerate JKS format with renewed certs
            echo "$(date): Regenerating JKS certificate"
            keytool \
                -importkeystore \
                -srcstorepass "${KEYSTORE_PASSWORD}" \
                -deststorepass "${KEYSTORE_PASSWORD}" \
                -destkeystore "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.jks" \
                -srckeystore "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.p12" \
                -srcstoretype "pkcs12" \
                -noprompt
            
            echo "$(date): Certificate renewal and conversion completed successfully"
        else
            echo "$(date): WARNING: Renewed certificate files not found"
        fi
    else
        echo "$(date): WARNING: Let's Encrypt live directory not found for domain: $CERT_DOMAIN"
    fi
else
    echo "$(date): Certificate renewal failed"
    exit 1
fi
