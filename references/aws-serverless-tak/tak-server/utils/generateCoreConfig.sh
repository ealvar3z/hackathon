#!/bin/bash

# Generate the CoreConfig.xml file if it doesn't exist or if forced
if [ ! -f "/opt/tak/CoreConfig.xml" ] || [ "$FORCE_CONFIG_REGEN" = "true" ]; then
    echo "Generating CoreConfig.xml"
    python3 /home/takserver/utils/generateCoreConfig.py
    mv CoreConfig.xml /opt/tak/CoreConfig.xml
fi