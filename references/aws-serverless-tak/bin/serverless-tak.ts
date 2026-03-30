#!/usr/bin/env node
import { App } from "aws-cdk-lib";
import { readFileSync } from "fs";
import { load as loadYaml } from "js-yaml";
import { CoreInfraStack } from "../lib/stacks/core-infra-stack";
import { LdapStack } from "../lib/stacks/ldap-stack";
import { TakServerStack } from "../lib/stacks/tak-server-stack";
import { Config } from "../lib/types/config";

const config: Config = loadYaml(readFileSync("config.yaml", "utf8")) as Config;

const env = {
  account: config.env?.account || process.env.CDK_DEFAULT_ACCOUNT,
  region: config.env?.region || process.env.CDK_DEFAULT_REGION,
};

const app = new App();

// Create core infrastructure stack
const infraStack = new CoreInfraStack(app, "CoreInfraStack", {
  env,
  isDevelopmentEnv: config.env?.isDevelopmentEnvironment ?? false,
  firewallConfig: config.firewall,
  ...config.infra,
});

// Create LDAP stack if needed
let ldapStack: LdapStack | undefined;
if (config.tak.authType === "LDAP" && config.tak.ldap?.deployLdapServer) {
  ldapStack = new LdapStack(app, "LdapStack", {
    env,
    vpc: infraStack.vpc,
    cluster: infraStack.cluster,
    environmentConfig: config.env!,
    infrastructureConfig: config.infra,
    ldapConfig: config.tak.ldap!,
    hostedZoneId: config.tak.hostedZoneId,
    domainName: config.tak.domainName,
  });
}

// Create TAK Server stack
const takServerStack = new TakServerStack(app, "TakServerStack", {
  env,
  isDevelopmentEnv: config.env?.isDevelopmentEnvironment ?? false,
  vpc: infraStack.vpc,
  cluster: infraStack.cluster,
  database: infraStack.database.cluster,
  dbCredentials: infraStack.database.dbCredentials,
  fileSystem: infraStack.filesystem.efsFileSystem,
  dataBucket: infraStack.bucket,
  loggingBucket: infraStack.loggingBucket,
  firewallConfig: config.firewall,
  ldapServer: ldapStack?.ldapServer,
  ldapNamespace: ldapStack?.namespace,
  ...config.tak,
  ldapConfig: config.tak.ldap,
  forceConfigRegen: config.tak.forceConfigRegen ?? false,
});

if (ldapStack) {
  takServerStack.addDependency(ldapStack);
}
