import { Stack, StackProps } from "aws-cdk-lib";
import { IVpc, Vpc } from "aws-cdk-lib/aws-ec2";
import { Cluster, ICluster } from "aws-cdk-lib/aws-ecs";
import { Construct } from "constructs";
import { TakDatabase } from "../constructs/database";
import { TakFileSystem } from "../constructs/filesystem";
import { LoggingBucket } from "../constructs/logging-bucket";
import { NetworkFirewall } from "../constructs/network-firewall";
import { TakServerVpc } from "../constructs/vpc";
import { SecureBucket } from "../constructs/secure-bucket";
import { TakCluster } from "../constructs/ecs-cluster";
import { FirewallConfig } from "../types/firewallConfig";

interface CoreInfraStackProps extends StackProps {
  /**
   * ID of an existing VPC
   */
  vpcId?: string;

  /**
   * Name of an existing ECS Cluster
   */
  clusterName?: string;

  /**
   * Whether or not this is a development environment
   */
  isDevelopmentEnv?: boolean;

  /**
   * Firewall configuration
   */
  firewallConfig?: FirewallConfig;
}

export class CoreInfraStack extends Stack {
  public readonly vpc: IVpc | Vpc;
  public readonly cluster: Cluster | ICluster;
  public readonly bucket: SecureBucket;
  public readonly loggingBucket: LoggingBucket;
  public readonly filesystem: TakFileSystem;
  public readonly database: TakDatabase;
  public readonly firewall?: NetworkFirewall;

  constructor(scope: Construct, id: string, props?: CoreInfraStackProps) {
    super(scope, id, props);

    const isDevelopmentEnv = props?.isDevelopmentEnv || false;
    const isFirewallEnabled = props?.firewallConfig?.enabled || true;

    /**
     * VPC
     */
    this.vpc = new TakServerVpc(this, "TakServerVpc", {
      isDevelopmentEnv: isDevelopmentEnv,
      isFirewallEnabled: isFirewallEnabled,
    }).vpc;

    /**
     * ECS Cluster
     */
    this.cluster = new TakCluster(this, "TakCluster", {
      vpc: this.vpc,
    }).cluster;

    /**
     * S3 Buckets
     */
    this.loggingBucket = new LoggingBucket(this, "TakLoggingBucket", {
      isDevelopmentEnv,
    });
    this.bucket = new SecureBucket(this, "TakDataBucket", {
      serverAccessLogsBucket: this.loggingBucket,
      serverAccessLogsPrefix: "tak-data-bucket-logs/",
    });

    /**
     * EFS
     */
    this.filesystem = new TakFileSystem(this, "TakFileSystem", {
      vpc: this.vpc,
      isDevelopmentEnv,
    });

    /**
     * RDS Aurora
     */
    this.database = new TakDatabase(this, "TakDatabase", {
      vpc: this.vpc,
      isDevelopmentEnv,
    });

    /**
     * AWS Network Firewall (optional)
     */
    if (isFirewallEnabled) {
      this.firewall = new NetworkFirewall(this, "NetworkFirewall", {
        vpc: this.vpc as Vpc,
        isDevelopmentEnv,
        firewallConfig: props?.firewallConfig,
      });
    }
  }
}
