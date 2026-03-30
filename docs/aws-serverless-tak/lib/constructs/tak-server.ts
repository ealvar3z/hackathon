import { CfnOutput, Duration } from "aws-cdk-lib";
import {
  IVpc,
  Peer,
  Port,
  SecurityGroup,
  SubnetType,
  SubnetSelection,
} from "aws-cdk-lib/aws-ec2";
import { Platform } from "aws-cdk-lib/aws-ecr-assets";
import {
  ContainerImage,
  Protocol as ecsProtocol,
  FargateService,
  FargateTaskDefinition,
  ICluster,
  LogDrivers,
  Secret as EcsSecret,
} from "aws-cdk-lib/aws-ecs";
import { AccessPoint, IFileSystem } from "aws-cdk-lib/aws-efs";
import {
  NetworkLoadBalancer,
  NetworkTargetGroup,
  Protocol,
  TargetType,
} from "aws-cdk-lib/aws-elasticloadbalancingv2";
import {
  Effect,
  PolicyDocument,
  PolicyStatement,
  Role,
  ServicePrincipal,
} from "aws-cdk-lib/aws-iam";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import { IDatabaseCluster } from "aws-cdk-lib/aws-rds";
import { ARecord, HostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { LoadBalancerTarget } from "aws-cdk-lib/aws-route53-targets";
import { Bucket } from "aws-cdk-lib/aws-s3";
import { ISecret, Secret } from "aws-cdk-lib/aws-secretsmanager";
import { IPrivateDnsNamespace } from "aws-cdk-lib/aws-servicediscovery";
import { Construct } from "constructs";
import { LdapServer } from "./ldap-server";
import { LdapConfig as LdapConfigType } from "../types/takConfig";

interface CertificateInfo {
  email: string;
  state: string;
  city: string;
  org: string;
  orgUnit: string;
}

interface LdapConfig {
  deployLdapServer?: boolean;
  host?: string;
  port?: number;
  baseDn?: string;
  bindDn?: string;
  bindPassword?: string;
  userSearchBase?: string;
  groupSearchBase?: string;
}

interface TakServerProps {
  vpc: IVpc;
  cluster: ICluster;
  database: IDatabaseCluster;
  dbCredentials: ISecret;
  fileSystem: IFileSystem;
  dataBucket: Bucket;
  loggingBucket: Bucket;
  certsBucket: Bucket;
  domainName: string;
  subdomain: string;
  hostedZoneId?: string;
  certInfo: CertificateInfo;
  authType: string;
  ldapConfig?: LdapConfigType;
  ldapServer?: LdapServer;
  ldapNamespace?: IPrivateDnsNamespace;
  isDevelopmentEnv?: boolean;
  isFirewallEnabled: boolean;
  forceConfigRegen?: boolean;
}

export class TakServer extends Construct {
  public readonly service: FargateService;
  public readonly loadBalancer: NetworkLoadBalancer;
  public readonly keystoreSecret: ISecret;

  constructor(scope: Construct, id: string, props: TakServerProps) {
    super(scope, id);

    const isDevelopmentEnv = props.isDevelopmentEnv ?? false;

    /**
     * Keystore Password Secret
     */
    this.keystoreSecret = new Secret(this, "KeystoreSecret", {
      description: "TAK Server keystore passwords",
      generateSecretString: {
        secretStringTemplate: JSON.stringify({}),
        generateStringKey: "password",
        excludeCharacters: '"@/\\',
        passwordLength: 32,
      },
    });

    /**
     * IAM Role for ECS Task
     */
    const taskRole = new Role(this, "TakServerTaskRole", {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
      roleName: "TakServerTaskRole",
      inlinePolicies: {
        AllowSecrets: new PolicyDocument({
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: [
                "secretsmanager:GetResourcePolicy",
                "secretsmanager:GetSecretValue",
                "secretsmanager:DescribeSecret",
                "secretsmanager:ListSecretVersionIds",
              ],
              resources: [
                props.dbCredentials.secretArn,
                this.keystoreSecret.secretArn,
              ],
            }),
          ],
        }),
        AllowS3: new PolicyDocument({
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: [
                "s3:GetObject",
                "s3:ListBucket",
                "s3:PutObject",
                "s3:DeleteObject",
                "s3:GetBucketLocation",
              ],
              resources: [
                props.dataBucket.bucketArn,
                `${props.dataBucket.bucketArn}/*`,
              ],
            }),
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["kms:Decrypt", "kms:GenerateDataKey"],
              resources: [props.dataBucket.encryptionKey!.keyArn],
            }),
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"],
              resources: [
                props.certsBucket.bucketArn,
                `${props.certsBucket.bucketArn}/*`,
              ],
            }),
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: ["kms:Decrypt", "kms:GenerateDataKey"],
              resources: [props.certsBucket.encryptionKey!.keyArn],
            }),
          ],
        }),
        AllowDNSCertValidation: new PolicyDocument({
          statements: [
            new PolicyStatement({
              effect: Effect.ALLOW,
              actions: [
                "route53:ChangeResourceRecordSets",
                "route53:GetChange",
                "route53:ListHostedZones",
              ],
              resources: ["*"],
            }),
          ],
        }),
        ...(isDevelopmentEnv
          ? {
              AllowECSExec: new PolicyDocument({
                statements: [
                  new PolicyStatement({
                    effect: Effect.ALLOW,
                    actions: [
                      "ssmmessages:CreateControlChannel",
                      "ssmmessages:CreateDataChannel",
                      "ssmmessages:OpenControlChannel",
                      "ssmmessages:OpenDataChannel",
                    ],
                    resources: ["*"],
                  }),
                ],
              }),
            }
          : {}),
      },
    });

    /**
     * Ports used by TAK Server application
     */
    const takPorts = [
      { name: "mobileClients", port: 8089 },
      { name: "adminPortal", port: 8443 },
      { name: "webAuth", port: 8446 },
      { name: "federation1", port: 9000 },
      { name: "federation2", port: 9001 },
    ];

    /**
     * EFS Access point
     */
    const efsAccessPoint = new AccessPoint(this, "TakServerEfsAccessPoint", {
      fileSystem: props.fileSystem,
      createAcl: {
        ownerGid: "1001",
        ownerUid: "1001",
        permissions: "755",
      },
      path: "/tak-data",
    });

    /**
     * ECS Task Definition
     */
    const takTaskDefinition = new FargateTaskDefinition(
      this,
      "TakTaskDefinition",
      {
        taskRole: taskRole,
        memoryLimitMiB: 8192,
        cpu: 2048,
        volumes: [
          {
            name: "tak-data",
            efsVolumeConfiguration: {
              fileSystemId: props.fileSystem.fileSystemId,
              transitEncryption: "ENABLED",
              rootDirectory: "/",
              authorizationConfig: {
                accessPointId: efsAccessPoint.accessPointId,
              },
            },
          },
        ],
      }
    );

    // Build environment variables
    const environment: { [key: string]: string } = {
      BUCKET_NAME: props.dataBucket.bucketName,
      CERTS_BUCKET_NAME: props.certsBucket.bucketName,
      DB_HOST: props.database.clusterEndpoint.hostname,
      DB_PORT: `${props.database.clusterEndpoint.port}`,
      CERT_DOMAIN: props.domainName,
      CERT_SUBDOMAIN: props.subdomain,
      CERT_EMAIL: props.certInfo.email,
      STATE: props.certInfo.state,
      CITY: props.certInfo.city,
      ORGANIZATION: props.certInfo.org,
      ORGANIZATIONAL_UNIT: props.certInfo.orgUnit,
      AUTH_TYPE: props.authType,
      KEYSTORE_SECRET_NAME: this.keystoreSecret.secretName,
      FORCE_CONFIG_REGEN: `${props.forceConfigRegen}`,
      // Override Spring Boot SSL client auth configuration
      SERVER_SSL_CLIENT_AUTH: "NONE",
    };

    // Add LDAP configuration if using LDAP auth
    if (props.authType === "LDAP" && props.ldapConfig) {
      if (props.ldapConfig.deployLdapServer && props.ldapServer) {
        environment.LDAP_HOST = "ldap-server.tak.local";
        environment.LDAP_PORT = "3890";
        environment.LDAP_BASE_DN = props.ldapConfig.baseDn || "dc=tak,dc=local";
        environment.LDAP_BIND_DN =
          "uid=admin,ou=people," +
          (props.ldapConfig.baseDn || "dc=tak,dc=local");
        // LDAP_BIND_PASSWORD will be injected as a secret below
        environment.LDAP_USER_SEARCH_BASE =
          "ou=people," + (props.ldapConfig.baseDn || "dc=tak,dc=local");
        environment.LDAP_GROUP_SEARCH_BASE =
          "ou=groups," + (props.ldapConfig.baseDn || "dc=tak,dc=local");
      } else {
        environment.LDAP_HOST = props.ldapConfig.host || "";
        environment.LDAP_PORT = (props.ldapConfig.port || 389).toString();
        environment.LDAP_BASE_DN = props.ldapConfig.baseDn || "";
        environment.LDAP_BIND_DN = props.ldapConfig.bindDn || "";
        environment.LDAP_USER_SEARCH_BASE =
          props.ldapConfig.userSearchBase || "";
        environment.LDAP_GROUP_SEARCH_BASE =
          props.ldapConfig.groupSearchBase || "";
        if (props.ldapConfig.bindPassword) {
          environment.LDAP_BIND_PASSWORD = props.ldapConfig.bindPassword;
        }
      }
    }

    const container = takTaskDefinition.addContainer("TakServerContainer", {
      image: ContainerImage.fromAsset("tak-server", {
        platform: Platform.LINUX_AMD64,
      }),
      portMappings: takPorts.map((config) => ({
        containerPort: config.port,
        protocol: ecsProtocol.TCP,
      })),
      cpu: 2048,
      memoryLimitMiB: 8192,
      logging: LogDrivers.awsLogs({
        streamPrefix: "tak-server",
        logRetention: RetentionDays.ONE_WEEK,
      }),
      environment,
      healthCheck: {
        command: ["CMD-SHELL", "pgrep -f takserver || exit 1"],
        retries: 5,
        interval: Duration.seconds(30),
        startPeriod: Duration.seconds(60),
      },
    });

    // Inject RDS database password
    container.addSecret(
      "DB_PASSWORD",
      EcsSecret.fromSecretsManager(props.dbCredentials, "password")
    );

    // Inject LDAP admin password if using deployed LDAP server
    if (
      props.authType === "LDAP" &&
      props.ldapConfig?.deployLdapServer &&
      props.ldapServer
    ) {
      container.addSecret(
        "LDAP_BIND_PASSWORD",
        EcsSecret.fromSecretsManager(props.ldapServer.adminSecret, "password")
      );
    }

    // Pass keystore secret name for direct access
    environment.KEYSTORE_SECRET_NAME = this.keystoreSecret.secretName;

    // Mount the EFS volume
    container.addMountPoints({
      sourceVolume: "tak-data",
      containerPath: "/opt/tak",
      readOnly: false,
    });

    /**
     * ECS Service
     */
    const serviceSecurityGroup = new SecurityGroup(
      this,
      "TakServerServiceSecurityGroup",
      {
        vpc: props.vpc,
        allowAllOutbound: false,
        description: "Security group for the TAK Server Service",
      }
    );

    // Add specific egress rules
    serviceSecurityGroup.addEgressRule(
      Peer.anyIpv4(),
      Port.tcp(443),
      "HTTPS outbound"
    );
    serviceSecurityGroup.addEgressRule(
      Peer.anyIpv4(),
      Port.tcp(53),
      "DNS TCP outbound"
    );
    serviceSecurityGroup.addEgressRule(
      Peer.anyIpv4(),
      Port.udp(53),
      "DNS UDP outbound"
    );

    // Add LDAP egress rule if using LDAP
    if (props.authType === "LDAP") {
      serviceSecurityGroup.addEgressRule(
        Peer.anyIpv4(),
        Port.tcp(3890),
        "LDAP outbound"
      );
    }

    // Configure service discovery for TAK server if LDAP namespace is provided
    const serviceOptions: any = {
      cluster: props.cluster,
      vpcSubnets: { subnetGroupName: "Private" },
      taskDefinition: takTaskDefinition,
      healthCheckGracePeriod: Duration.seconds(600),
      desiredCount: 1,
      enableExecuteCommand: isDevelopmentEnv,
      securityGroups: [serviceSecurityGroup],
      minHealthyPercent: 100,
    };

    // If LDAP namespace is provided, register TAK server in the same namespace
    if (props.ldapNamespace) {
      serviceOptions.cloudMapOptions = {
        name: "tak-server",
        cloudMapNamespace: props.ldapNamespace,
      };
    }

    this.service = new FargateService(this, "TakServerService", serviceOptions);

    // Create autoscaling target for the ECS service
    const scalableTarget = this.service.autoScaleTaskCount({
      minCapacity: 1,
      maxCapacity: 3,
    });

    // Add CPU-based scaling policy
    scalableTarget.scaleOnCpuUtilization("CpuScaling", {
      targetUtilizationPercent: 80,
      scaleInCooldown: Duration.seconds(300),
      scaleOutCooldown: Duration.seconds(300),
    });

    // Allow NFS access to EFS filesystem
    this.service.connections.allowTo(
      props.fileSystem,
      Port.tcp(2049),
      "Allow TakServer to EFS"
    );
    props.fileSystem.grantRootAccess(
      this.service.taskDefinition.taskRole.grantPrincipal
    );

    // Allow access to RDS database
    this.service.connections.allowTo(
      props.database,
      Port.tcp(props.database.clusterEndpoint.port),
      "Allow TakServer to RDS"
    );

    // Allow access to LDAP server if deployed
    if (props.ldapServer) {
      this.service.connections.allowTo(
        props.ldapServer.service,
        Port.tcp(props.ldapServer.port),
        "Allow TakServer to LDAP"
      );
    }

    /**
     * Network Load Balancer
     */
    const nlbSecurityGroup = new SecurityGroup(
      this,
      "LoadBalancerSecurityGroup",
      {
        vpc: props.vpc,
        allowAllOutbound: true,
        description: "Security group for the Network Load Balancer",
      }
    );

    // Trust Subnet if NFW is enabled, Public Subnet if not
    const nlbSubnets: SubnetSelection = props.isFirewallEnabled
      ? { subnetGroupName: "Trust" }
      : { subnetType: SubnetType.PUBLIC };

    this.loadBalancer = new NetworkLoadBalancer(this, "TakLoadBalancer", {
      vpc: props.vpc,
      internetFacing: true,
      vpcSubnets: nlbSubnets,
      securityGroups: [nlbSecurityGroup],
    });

    // Create Listeners and Target Groups for NLB
    takPorts.forEach((config) => {
      // Adjust health check settings for slow-starting ports
      const isWebPort = config.port === 8443 || config.port === 8446;

      // Create Target Group with appropriate health check
      const targetGroupConfig: any = {
        vpc: props.vpc,
        port: config.port,
        protocol: Protocol.TCP,
        targetType: TargetType.IP,
      };

      if (isWebPort) {
        // Web ports: 3 failures × 60s = 3 minutes to detect failure
        targetGroupConfig.healthCheck = {
          healthyThresholdCount: 2,
          unhealthyThresholdCount: 3,
          interval: Duration.seconds(60),
          timeout: Duration.seconds(10),
        };
      } else {
        // Other ports: 3 failures × 30s = 90 seconds to detect failure
        targetGroupConfig.healthCheck = {
          healthyThresholdCount: 2,
          unhealthyThresholdCount: 3,
          interval: Duration.seconds(30),
          timeout: Duration.seconds(10),
        };
      }

      const targetGroup = new NetworkTargetGroup(
        this,
        `${config.name}TargetGroup`,
        targetGroupConfig
      );

      // Add target
      targetGroup.addTarget(
        this.service.loadBalancerTarget({
          containerName: "TakServerContainer",
          containerPort: config.port,
        })
      );

      // Create Listener
      this.loadBalancer.addListener(`${config.name}Listener`, {
        port: config.port,
        protocol: Protocol.TCP,
        defaultTargetGroups: [targetGroup],
      });

      // Allow inbound to loadbalancer from internet
      nlbSecurityGroup.connections.allowFromAnyIpv4(
        Port.tcp(config.port),
        config.name
      );

      // Allow ingress to container from loadbalancer
      this.service.connections.allowFrom(
        nlbSecurityGroup,
        Port.tcp(config.port)
      );
    });

    // Enable access logging for the load balancer
    this.loadBalancer.logAccessLogs(props.loggingBucket, "nlb-access-logs");

    /**
     * Route53
     *
     * Automatically add an A record if domain is hosted in Route53
     */
    if (props.hostedZoneId) {
      const zone = HostedZone.fromHostedZoneAttributes(this, "HostedZone", {
        hostedZoneId: props.hostedZoneId,
        zoneName: props.domainName,
      });

      new ARecord(this, "TakDNSRecord", {
        zone,
        target: RecordTarget.fromAlias(
          new LoadBalancerTarget(this.loadBalancer)
        ),
        recordName: `${props.subdomain}.${props.domainName}`,
      });

      new CfnOutput(this, "TakServerUrl", {
        description: "TAK Server URL",
        value: `${props.subdomain}.${props.domainName}`,
      });

      new CfnOutput(this, "TakAdminUrl", {
        description: "TAK Server admin console URL",
        value: `https://${props.subdomain}.${props.domainName}:8443`,
      });
    }

    /**
     * Outputs
     */
    new CfnOutput(this, "BucketName", {
      description: "Name of the S3 bucket for TAK certs and data",
      value: props.dataBucket.bucketName,
    });
  }
}
