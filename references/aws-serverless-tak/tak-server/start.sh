#!/bin/bash

set -euo pipefail

# this is a new deployment if the /opt/tak directory is empty
new_deployment=$([ -z "$(ls -A /opt/tak)" ] && echo true || echo false)
# this is a version upgrade if the tak versions are different
version_upgrade=$([ -f "/opt/tak/version.txt" ] && [ "$(cat /opt/tak/version.txt)" != "$(cat ./tak/version.txt)" ] && echo true || echo false)

if [ $new_deployment == true ] || [ $version_upgrade == true ] ; then
    # Copy takserver files to EFS
    echo "Copying Takserver files to EFS"
    cp -r ./tak/* /opt/tak/
fi

# initialize certificates
./utils/initializeCerts.sh
# Generate CoreConfig.xml
./utils/generateCoreConfig.sh
# initialize database
./utils/initializeDatabase.sh
# generate client certs
./utils/makeClientCerts.sh

# start TAK Server
/opt/tak/configureInDocker.sh init
