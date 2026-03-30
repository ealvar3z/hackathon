import { CfnOutput, Stack, StackProps } from "aws-cdk-lib";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { ICluster } from "aws-cdk-lib/aws-ecs";
import { HostedZone, IHostedZone } from "aws-cdk-lib/aws-route53";
import { PrivateDnsNamespace } from "aws-cdk-lib/aws-servicediscovery";
import { Construct } from "constructs";
import { TakFileSystem } from "../constructs/filesystem";
import { LdapServer } from "../constructs/ldap-server";
import { EnvironmentConfig } from "../types/environmentConfig";
import { InfraConfig } from "../types/infrastructureConfig";
import { LdapConfig } from "../types/takConfig";

interface LdapStackProps extends StackProps {
  vpc: IVpc;
  cluster: ICluster;
  environmentConfig: EnvironmentConfig;
  infrastructureConfig?: InfraConfig;
  ldapConfig: LdapConfig;
  hostedZoneId?: string;
  domainName?: string;
}

export class LdapStack extends Stack {
  public readonly ldapServer: LdapServer;
  public readonly namespace: PrivateDnsNamespace;

  constructor(scope: Construct, id: string, props: LdapStackProps) {
    super(scope, id, props);

    const fileSystem = new TakFileSystem(this, "LdapFileSystem", {
      vpc: props.vpc,
      isDevelopmentEnv: props.environmentConfig.isDevelopmentEnvironment,
    });

    let hostedZone: IHostedZone | undefined;
    if (props.hostedZoneId && props.domainName) {
      hostedZone = HostedZone.fromHostedZoneAttributes(this, "HostedZone", {
        hostedZoneId: props.hostedZoneId,
        zoneName: props.domainName,
      });
    }

    this.ldapServer = new LdapServer(this, "LdapServer", {
      vpc: props.vpc,
      cluster: props.cluster,
      fileSystem: fileSystem.efsFileSystem,
      ldapConfig: props.ldapConfig,
      hostedZone,
      domainName: props.domainName,
      isDevelopmentEnv: props.environmentConfig.isDevelopmentEnvironment,
    });

    this.namespace = this.ldapServer.namespace;

    new CfnOutput(this, "LdapWebUrl", {
      value: this.ldapServer.webUrl,
      description: "URL for LLDAP web interface",
    });
  }
}
