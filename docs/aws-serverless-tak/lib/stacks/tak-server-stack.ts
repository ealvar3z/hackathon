import { RemovalPolicy, Stack, StackProps } from "aws-cdk-lib";
import { IVpc } from "aws-cdk-lib/aws-ec2";
import { Cluster, ICluster } from "aws-cdk-lib/aws-ecs";
import { IFileSystem } from "aws-cdk-lib/aws-efs";
import { IDatabaseCluster } from "aws-cdk-lib/aws-rds";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { ISecret } from "aws-cdk-lib/aws-secretsmanager";
import { IPrivateDnsNamespace } from "aws-cdk-lib/aws-servicediscovery";
import { Construct } from "constructs";
import { SecureBucket } from "../constructs/secure-bucket";
import { LdapServer } from "../constructs/ldap-server";
import { TakServer } from "../constructs/tak-server";
import { FirewallConfig } from "../types/firewallConfig";
import { LdapConfig } from "../types/takConfig";

interface TakServerStackProps extends StackProps {
  vpc: IVpc;
  cluster: Cluster | ICluster;
  database: IDatabaseCluster;
  dbCredentials: ISecret;
  fileSystem: IFileSystem;
  dataBucket: Bucket;
  loggingBucket: Bucket;
  domainName: string;
  subdomainName: string;
  adminEmail: string;
  certState: string;
  certCity: string;
  certOrg: string;
  certOrgUnit: string;
  hostedZoneId?: string;
  authType: string;
  ldapConfig?: LdapConfig;
  ldapServer?: LdapServer;
  ldapNamespace?: IPrivateDnsNamespace;
  isDevelopmentEnv?: boolean;
  firewallConfig?: FirewallConfig;
  forceConfigRegen?: boolean;
}

export class TakServerStack extends Stack {
  constructor(scope: Construct, id: string, props: TakServerStackProps) {
    super(scope, id, props);

    // Create S3 bucket for TAK server certificates
    const certsBucket = new SecureBucket(this, "TakCertsBucket", {
      bucketName: `${this.account}-tak-server-certs`,
      serverAccessLogsBucket: props.loggingBucket,
      serverAccessLogsPrefix: "certs-access-logs/",
      removalPolicy: props.isDevelopmentEnv
        ? RemovalPolicy.DESTROY
        : RemovalPolicy.RETAIN,
      autoDeleteObjects: props.isDevelopmentEnv,
    });

    // Create the TakServer application resources (ECS Service, NLB, Route53 DNS record)
    const takServer = new TakServer(this, "TakServer", {
      vpc: props.vpc,
      cluster: props.cluster,
      database: props.database,
      dbCredentials: props.dbCredentials,
      fileSystem: props.fileSystem,
      domainName: props.domainName,
      subdomain: props.subdomainName,
      dataBucket: props.dataBucket,
      loggingBucket: props.loggingBucket,
      certsBucket: certsBucket,
      isDevelopmentEnv: props.isDevelopmentEnv,
      isFirewallEnabled: props.firewallConfig?.enabled || true,
      certInfo: {
        email: props.adminEmail,
        state: props.certState,
        city: props.certCity,
        org: props.certOrg,
        orgUnit: props.certOrgUnit,
      },
      hostedZoneId: props.hostedZoneId,
      authType: props.authType,
      ldapConfig: props.ldapConfig,
      ldapServer: props.ldapServer,
      ldapNamespace: props.ldapNamespace,
      forceConfigRegen: props.forceConfigRegen,
    });
  }
}
