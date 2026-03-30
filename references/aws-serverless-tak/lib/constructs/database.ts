import { Duration, RemovalPolicy } from "aws-cdk-lib";
import { IVpc, SecurityGroup } from "aws-cdk-lib/aws-ec2";
import { Key } from "aws-cdk-lib/aws-kms";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import {
  AuroraPostgresEngineVersion,
  ClusterInstance,
  Credentials,
  DatabaseCluster,
  DatabaseClusterEngine,
} from "aws-cdk-lib/aws-rds";
import { Secret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";

export interface TakDatabaseProps {
  vpc: IVpc;
  kmsKey?: Key;
  isDevelopmentEnv?: boolean;
}

export class TakDatabase extends Construct {
  public readonly cluster: DatabaseCluster;
  public readonly dbCredentials: Secret;

  constructor(scope: Construct, id: string, props: TakDatabaseProps) {
    super(scope, id);

    const removalPolicy = props.isDevelopmentEnv
      ? RemovalPolicy.DESTROY
      : RemovalPolicy.RETAIN;

    const kmsKey =
      props.kmsKey ??
      new Key(scope, `${id}-DatabaseKey`, {
        enableKeyRotation: true,
        removalPolicy: removalPolicy,
      });

    this.dbCredentials = new Secret(this, "DatabaseCredentials", {
      description: "TAK Database password",
      generateSecretString: {
        secretStringTemplate: JSON.stringify({ username: "postgres" }),
        generateStringKey: "password",
        passwordLength: 32,
        includeSpace: false,
        excludeCharacters: '<>"``&@/\\', // this password goes in the TAK CoreConfig.xml file so must be XML safe
      },
      secretName: "tak-DbCredentials",
    });

    const rdsSecurityGroup = new SecurityGroup(this, "RdsSecurityGroup", {
      vpc: props.vpc,
      description: "RDS Security Group",
      allowAllOutbound: true,
    });

    this.cluster = new DatabaseCluster(this, "AuroraDBCluster", {
      engine: DatabaseClusterEngine.auroraPostgres({
        version: AuroraPostgresEngineVersion.VER_15_6,
      }),
      storageEncryptionKey: kmsKey,
      credentials: Credentials.fromSecret(this.dbCredentials),
      securityGroups: [rdsSecurityGroup],
      writer: ClusterInstance.serverlessV2("TakDBWriter"),
      readers: [
        ClusterInstance.serverlessV2("TakDBReader", { scaleWithWriter: true }),
      ],
      serverlessV2MinCapacity: 0.5,
      serverlessV2MaxCapacity: 8,
      vpc: props.vpc,
      vpcSubnets: { subnetGroupName: "Private" },
      parameters: {
        "rds.force_ssl": "1",
        "rds.log_retention_period": "1440",
      },
      enableClusterLevelEnhancedMonitoring: true,
      monitoringInterval: Duration.seconds(60),
      cloudwatchLogsExports: ["postgresql"],
      cloudwatchLogsRetention: RetentionDays.ONE_WEEK,
      backup: { retention: Duration.days(7) },
      deletionProtection: !props.isDevelopmentEnv,
      removalPolicy,
    });
  }
}
