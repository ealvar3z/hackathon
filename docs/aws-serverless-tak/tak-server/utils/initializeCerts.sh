#!/bin/bash

# Ensure user-installed Python packages and java are in PATH
export PATH="$HOME/.local/bin:/opt/java/openjdk/bin:$PATH"

# Get keystore password from Secrets Manager
source ./utils/getKeystorePassword.sh

# copy letsencrypt certs from EFS if they exist
if [ -d "/opt/tak/certs/files/${CERT_DOMAIN}" ]; then
    echo "Copying Let's Encrypt certs from EFS"
    mkdir -p "/etc/letsencrypt/live/${CERT_DOMAIN}"
    cp "/opt/tak/certs/files/${CERT_DOMAIN}/"* "/etc/letsencrypt/live/${CERT_DOMAIN}"
else # otherwise generate new ones
    echo "No Let's Encrypt certs found in EFS."
    echo "Generating new Let's Encrypt certs."
    max_attempts=3
    attempt=1

    while ! certbot certonly -v --dns-route53 -d "${CERT_DOMAIN}" -d "*.${CERT_DOMAIN}" --email "${CERT_EMAIL}" --non-interactive --test-cert --agree-tos; do
        if [ $attempt -ge $max_attempts ]; then
            echo "Command failed after $max_attempts attempts. Exiting."
            exit 1
        fi

        echo "Command failed, retrying in 10 seconds... (Attempt $attempt of $max_attempts)"
        sleep 10
        ((attempt++))
    done

    mkdir -p "/opt/tak/certs/files/${CERT_DOMAIN}/"
    cp "/etc/letsencrypt/live/${CERT_DOMAIN}/"* "/opt/tak/certs/files/${CERT_DOMAIN}/"
fi

# setup an automated renewal process for certbot
# (every 2 months on the 1st at 2:00 AM)
echo "Setting up automated renewal process for Let's Encrypt certs."
(crontab -l 2>/dev/null; echo "0 2 1 */2 * cd /home/takserver && CERT_DOMAIN=${CERT_DOMAIN} KEYSTORE_PASSWORD=${KEYSTORE_PASSWORD} ./utils/renewCerts.sh >> /var/log/letsencrypt/renewal.log 2>&1") | crontab -

# convert Let's Encrypt certs to PKCS12 format
if [ ! -f "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.p12" ]; then
    echo "Converting Let's Encrypt certs to PKCS12 format"
    openssl pkcs12 \
        -export \
        -in "/etc/letsencrypt/live/${CERT_DOMAIN}/fullchain.pem" \
        -inkey "/etc/letsencrypt/live/${CERT_DOMAIN}/privkey.pem" \
        -out "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.p12" \
        -name "${CERT_DOMAIN}" \
        -password "pass:${KEYSTORE_PASSWORD}"
fi

# convert Let's Encrypt certs to Java KeyStore format
if [ ! -f "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.jks" ]; then
    echo "Converting Let's Encrypt certs to Java KeyStore format"
    keytool \
        -importkeystore \
        -srcstorepass "${KEYSTORE_PASSWORD}" \
        -deststorepass "${KEYSTORE_PASSWORD}" \
        -destkeystore "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.jks" \
        -srckeystore "/opt/tak/certs/files/${CERT_DOMAIN}/letsencrypt.p12" \
        -srcstoretype "pkcs12"
fi

# create a CA if needed
if [ ! -f "/opt/tak/certs/files/ca.pem" ]; then
    echo "Creating CA and intermediate signing cert"
    cd /opt/tak/certs/
    CAPASS=${KEYSTORE_PASSWORD} ./makeRootCa.sh --ca-name tak
fi

# create an intermediate signing cert if needed
if [ ! -f "/opt/tak/certs/files/tak-signing.jks" ]; then
    echo "Creating intermediate signing cert"
    cd /opt/tak/certs/
    echo y | CAPASS=${KEYSTORE_PASSWORD} ./makeCert.sh ca tak
fi

# create a server cert if needed
if [ ! -f "/opt/tak/certs/files/${CERT_SUBDOMAIN}.${CERT_DOMAIN}.jks" ]; then
    echo "Creating server certificate"
    cd /opt/tak/certs/
    CAPASS=${KEYSTORE_PASSWORD} ./makeCert.sh server "${CERT_SUBDOMAIN}.${CERT_DOMAIN}"
fi
