import json
import os

# Native Python xml library is being used to build, not parse, XML and is not vulnerable to the
# XXE attack flagged by semgrep.
# nosemgrep: use-defused-xml
import xml.etree.ElementTree as ET

from boto3.session import Session


def get_keystore_password():
    """Retrieve keystore password from AWS Secrets Manager"""
    try:
        secret_name = os.getenv("KEYSTORE_SECRET_NAME")
        if not secret_name:
            return "atakatak"  # fallback for local development

        session = Session()
        client = session.client("secretsmanager")
        response = client.get_secret_value(SecretId=secret_name)
        secret = json.loads(response["SecretString"])
        return secret["password"]
    except Exception:
        return "atakatak"  # fallback if secrets manager unavailable


# Get keystore password from Secrets Manager
keystore_password = get_keystore_password()

# Configuration as a Python dictionary
config = {
    "Configuration": {
        "@xmlns": "http://bbn.com/marti/xml/config",
        "@xmlns:xsi": "http://www.w3.org/2001/XMLSchema-instance",
        "@xsi:schemaLocation": "CoreConfig.xsd",
        "network": {
            "@multicastTTL": "5",
            "input": {
                "@_name": "stdssl",
                "@protocol": "tls",
                "@port": "8089",
                "@coreVersion": "2",
            },
            "connector": [
                {
                    "@port": "8443",
                    "@_name": "https",
                    "@enableAdminUI": "true",
                },
                {
                    "@port": "8446",
                    "@_name": "cert_https",
                    "@keystoreFile": f"/opt/tak/certs/files/{os.getenv('CERT_DOMAIN')}/letsencrypt.jks",
                    "@keystorePass": keystore_password,
                    "@enableAdminUI": "false",
                    "@clientAuth": "false",
                    "@auth": "ldap",
                },
            ],
        },
        "submission": {"@ignoreStaleMessages": "false", "@validateXml": "false"},
        "subscription": {"@reloadPersistent": "false"},
        "repository": {
            "@enable": "true",
            "@numDbConnections": "64",
            "@primaryKeyBatchSize": "500",
            "@insertionBatchSize": "500",
            "connection": {
                "@url": f"jdbc:postgresql://{os.getenv('DB_HOST')}:{os.getenv('DB_PORT')}/cot?ssl=verify-full",
                "@username": "postgres",
                "@password": os.getenv("DB_PASSWORD"),
            },
        },
        "repeater": {
            "@enable": "true",
            "@periodMillis": "3000",
            "@staleDelayMillis": "15000",
            "repeatableType": [
                {
                    "@_name": "911",
                    "@cancel-test": "/event/detail/emergency[@cancel='true']",
                    "@initiate-test": "/event/detail/emergency[@type='911 Alert']",
                },
                {
                    "@_name": "RingTheBell",
                    "@cancel-test": "/event/detail/emergency[@cancel='true']",
                    "@initiate-test": "/event/detail/emergency[@type='Ring The Bell']",
                },
                {
                    "@_name": "GeoFenceBreach",
                    "@cancel-test": "/event/detail/emergency[@cancel='true']",
                    "@initiate-test": "/event/detail/emergency[@type='Geo-fence Breached']",
                },
                {
                    "@_name": "TroopsInContact",
                    "@cancel-test": "/event/detail/emergency[@cancel='true']",
                    "@initiate-test": "/event/detail/emergency[@type='Troops In Contact']",
                },
            ],
        },
        "dissemination": {"@smartRetry": "false"},
        "filter": {
            "flowtag": {"@enable": "false", "@text": ""},
            "streamingbroker": {
                "@enable": "true",
            },
            "scrubber": {"@enable": "false", "@action": "overwrite"},
        },
        "buffer": {"latestSA": {"@enable": "true"}, "queue": None},
        "security": {
            "tls": {
                "@context": "TLSv1.2",
                "@keymanager": "SunX509",
                "@keystore": "JKS",
                "@keystoreFile": f"/opt/tak/certs/files/{os.getenv('CERT_SUBDOMAIN')}.{os.getenv('CERT_DOMAIN')}.jks",
                "@keystorePass": keystore_password,
                "@truststore": "JKS",
                "@truststoreFile": "/opt/tak/certs/files/truststore-tak.jks",
                "@truststorePass": keystore_password,
            }
        },
        "federation": {
            "@enableFederation": "true",
            "@allowFederatedDelete": "false",
            "@allowMissionFederation": "true",
            "@allowDataFeedFederation": "true",
            "@enableMissionFederationDisruptionTolerance": "true",
            "@missionFederationDisruptionToleranceRecencySeconds": "43200",
            "@federatedGroupMapping": "true",
            "@automaticGroupMapping": "false",
            "@enableDataPackageAndMissionFileFilter": "false",
            "federation-server": {
                "@port": "9000",
                "@webBaseUrl": f"https://{os.getenv('CERT_SUBDOMAIN')}.{os.getenv('CERT_DOMAIN')}/Marti",
                "tls": {
                    "@context": "TLSv1.2",
                    "@keymanager": "SunX509",
                    "@keystore": "JKS",
                    "@keystoreFile": f"/opt/tak/certs/files/{os.getenv('CERT_SUBDOMAIN')}.{os.getenv('CERT_DOMAIN')}.jks",
                    "@keystorePass": keystore_password,
                    "@truststore": "JKS",
                    "@truststoreFile": "/opt/tak/certs/files/fed-truststore.jks",
                    "@truststorePass": keystore_password,
                },
            },
        },
        "authorization": {
            "@default": "ldap",
            "ldap": {
                "@groupMapping": "true",
            },
        },
        "certificateSigning": {
            "@CA": "TAKServer",
            "certificateConfig": {
                "nameEntries": {
                    "nameEntry": [
                        {"@name": "O", "@value": os.getenv("ORGANIZATION", "TAK")},
                        {
                            "@name": "OU",
                            "@value": os.getenv("ORGANIZATIONAL_UNIT", "TAK"),
                        },
                    ]
                }
            },
            "TAKServerCAConfig": {
                "@keystore": "JKS",
                "@keystoreFile": "/opt/tak/certs/files/tak-signing.jks",
                "@keystorePass": keystore_password,
                "@validityDays": "30",
                "@signatureAlg": "SHA256WithRSA",
            },
        },
    }
}

if os.getenv("AUTH_TYPE") == "LDAP":
    ldap_host = os.getenv("LDAP_HOST", "ldap-server.tak.local")
    ldap_port = os.getenv("LDAP_PORT", "389")
    ldap_base_dn = os.getenv("LDAP_BASE_DN", "dc=tak,dc=local")
    ldap_bind_dn = os.getenv("LDAP_BIND_DN", "uid=admin,ou=people,dc=tak,dc=local")
    ldap_bind_password = os.getenv("LDAP_BIND_PASSWORD", "admin")
    ldap_user_search_base = os.getenv(
        "LDAP_USER_SEARCH_BASE", "ou=people,dc=tak,dc=local"
    )

    config["Configuration"]["auth"] = {
        "@default": "ldap",
        "ldap": {
            "@url": f"ldap://{ldap_host}:{ldap_port}",
            "@userstring": "uid={username}," + ldap_user_search_base,
            "@updateinterval": "60",
            "@groupprefix": "",
            "@style": "DS",
            "@serviceAccountDN": ldap_bind_dn,
            "@serviceAccountCredential": ldap_bind_password,
            "@groupSearchBase": f"ou=groups,{ldap_base_dn}",
            "@groupObjectClass": "groupOfUniqueNames",
            "@groupSearchFilter": "(uniqueMember=uid={username}," + ldap_user_search_base,
        },
    }
else:
    config["Configuration"]["auth"] = {
        "File": {"@location": "UserAuthenticationFile.xml"}
    }


def dict_to_xml(tag, d):
    elem = ET.Element(tag)
    for key, val in d.items():
        if key.startswith("@"):
            elem.set(key[1:], val)
        elif val is None:
            elem.append(ET.Element(key))
        elif isinstance(val, dict):
            elem.append(dict_to_xml(key, val))
        elif isinstance(val, list):
            for item in val:
                elem.append(dict_to_xml(key, item))
        else:
            child = ET.Element(key)
            child.text = str(val)
            elem.append(child)
    return elem


root = dict_to_xml("Configuration", config["Configuration"])
tree = ET.ElementTree(root)
tree.write("CoreConfig.xml", encoding="UTF-8", xml_declaration=True)
