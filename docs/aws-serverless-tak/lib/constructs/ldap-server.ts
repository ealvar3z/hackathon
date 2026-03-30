import { IVpc, Peer, Port, SecurityGroup } from "aws-cdk-lib/aws-ec2";
import {
  ContainerImage,
  Protocol as ecsProtocol,
  FargateService,
  FargateTaskDefinition,
  ICluster,
  LogDrivers,
  Secret,
} from "aws-cdk-lib/aws-ecs";
import {
  ApplicationLoadBalancer,
  ApplicationProtocol,
  ApplicationTargetGroup,
  TargetType,
  Protocol,
  ListenerAction,
} from "aws-cdk-lib/aws-elasticloadbalancingv2";
import {
  BlockPublicAccess,
  Bucket,
  BucketEncryption,
} from "aws-cdk-lib/aws-s3";
import { RemovalPolicy, Duration } from "aws-cdk-lib";
import { IFileSystem } from "aws-cdk-lib/aws-efs";
import {
  Role,
  ServicePrincipal,
  PolicyDocument,
  PolicyStatement,
  Effect,
} from "aws-cdk-lib/aws-iam";
import { RetentionDays } from "aws-cdk-lib/aws-logs";
import { ARecord, IHostedZone, RecordTarget } from "aws-cdk-lib/aws-route53";
import { LoadBalancerTarget } from "aws-cdk-lib/aws-route53-targets";
import {
  Certificate,
  CertificateValidation,
} from "aws-cdk-lib/aws-certificatemanager";
import {
  DnsRecordType,
  PrivateDnsNamespace,
} from "aws-cdk-lib/aws-servicediscovery";
import { Secret as SecretsManagerSecret } from "aws-cdk-lib/aws-secretsmanager";
import { Construct } from "constructs";
import { LdapConfig } from "../types/takConfig";

interface LdapServerProps {
  vpc: IVpc;
  cluster: ICluster;
  fileSystem: IFileSystem;
  ldapConfig: LdapConfig;
  hostedZone?: IHostedZone;
  domainName?: string;
  isDevelopmentEnv?: boolean;
}

export class LdapServer extends Construct {
  public readonly service: FargateService;
  public readonly host: string;
  public readonly port: number = 3890; // LLDAP LDAP port
  public readonly webPort: number = 17170; // LLDAP web interface port
  public readonly namespace: PrivateDnsNamespace;
  public readonly webUrl: string;
  public readonly loadBalancer: ApplicationLoadBalancer;
  public readonly adminSecret: SecretsManagerSecret;
  public readonly jwtSecret: SecretsManagerSecret;
  public readonly keySeed: SecretsManagerSecret;

  constructor(scope: Construct, id: string, props: LdapServerProps) {
    super(scope, id);

    // Create secrets for LDAP configuration
    const ldapAdminSecret = new SecretsManagerSecret(this, "LdapAdminSecret", {
      description: "LLDAP admin password",
      generateSecretString: {
        secretStringTemplate: JSON.stringify({ username: "admin" }),
        generateStringKey: "password",
        excludeCharacters: '"@/\\',
        passwordLength: 32,
      },
    });

    const jwtSecret = new SecretsManagerSecret(this, "LdapJwtSecret", {
      generateSecretString: {
        passwordLength: 32,
        includeSpace: false,
        excludePunctuation: true,
      },
      description: "JWT secret for LDAP server",
    });

    const keySeed = new SecretsManagerSecret(this, "LdapKeySeed", {
      generateSecretString: {
        passwordLength: 16,
        includeSpace: false,
        excludePunctuation: true,
      },
      description: "Key seed for LDAP server",
    });

    const taskRole = new Role(this, "LdapTaskRole", {
      assumedBy: new ServicePrincipal("ecs-tasks.amazonaws.com"),
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
                ldapAdminSecret.secretArn,
                jwtSecret.secretArn,
                keySeed.secretArn,
              ],
            }),
          ],
        }),
      },
    });

    const taskDefinition = new FargateTaskDefinition(
      this,
      "LdapTaskDefinition",
      {
        taskRole,
        memoryLimitMiB: 512,
        cpu: 256,
        volumes: [
          {
            name: "ldap-data",
            efsVolumeConfiguration: {
              fileSystemId: props.fileSystem.fileSystemId,
              transitEncryption: "ENABLED",
            },
          },
        ],
      }
    );

    const container = taskDefinition.addContainer("LldapContainer", {
      image: ContainerImage.fromRegistry("nitnelave/lldap:stable"),
      portMappings: [
        { containerPort: 3890, protocol: ecsProtocol.TCP }, // LDAP port
        { containerPort: 17170, protocol: ecsProtocol.TCP }, // Web interface port
      ],
      logging: LogDrivers.awsLogs({
        streamPrefix: "lldap-server",
        logRetention: RetentionDays.ONE_WEEK,
      }),
      environment: {
        LLDAP_LDAP_BASE_DN: props.ldapConfig.baseDn || "dc=tak,dc=local",
        LLDAP_HTTP_PORT: "17170",
        LLDAP_HTTP_HOST: "0.0.0.0",
        LLDAP_LDAP_PORT: "3890",
      },
      healthCheck: {
        command: ["CMD-SHELL", "curl -f http://localhost:17170/ || exit 1"],
        retries: 3,
        interval: Duration.seconds(30),
        startPeriod: Duration.seconds(60),
      },
    });

    // Inject secrets from Secrets Manager
    container.addSecret(
      "LLDAP_LDAP_USER_PASS",
      Secret.fromSecretsManager(ldapAdminSecret, "password")
    );
    container.addSecret(
      "LLDAP_JWT_SECRET",
      Secret.fromSecretsManager(jwtSecret)
    );
    container.addSecret("LLDAP_KEY_SEED", Secret.fromSecretsManager(keySeed));

    container.addMountPoints({
      sourceVolume: "ldap-data",
      containerPath: "/data", // LLDAP data directory
      readOnly: false,
    });

    const serviceSecurityGroup = new SecurityGroup(this, "LldapSecurityGroup", {
      vpc: props.vpc,
      allowAllOutbound: false,
      description: "Security group for LLDAP Server",
    });

    // Add egress rules for container image pulls
    serviceSecurityGroup.addEgressRule(
      Peer.anyIpv4(),
      Port.tcp(443),
      "HTTPS outbound for image pulls"
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

    // Create service discovery namespace
    this.namespace = new PrivateDnsNamespace(this, "LdapNamespace", {
      name: "tak.local",
      vpc: props.vpc,
    });

    this.service = new FargateService(this, "LldapService", {
      cluster: props.cluster,
      vpcSubnets: { subnetGroupName: "Private" },
      taskDefinition,
      desiredCount: 1,
      securityGroups: [serviceSecurityGroup],
      minHealthyPercent: 100,
      cloudMapOptions: {
        name: "ldap-server", // Keep same name for TAK Server compatibility
        cloudMapNamespace: this.namespace,
        dnsRecordType: DnsRecordType.A,
      },
      circuitBreaker: { rollback: true },
    });

    this.service.connections.allowTo(
      props.fileSystem,
      Port.tcp(2049),
      "Allow LLDAP to EFS"
    );

    // Create S3 bucket for server access logs
    const serverAccessLogsBucket = new Bucket(
      this,
      "LldapWebServerAccessLogsBucket",
      {
        encryption: BucketEncryption.S3_MANAGED,
        blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
        enforceSSL: true,
        removalPolicy: props.isDevelopmentEnv
          ? RemovalPolicy.DESTROY
          : RemovalPolicy.RETAIN,
        autoDeleteObjects: props.isDevelopmentEnv,
      }
    );

    // Create S3 bucket for ALB access logs (ALB doesn't support KMS encryption)
    const accessLogsBucket = new Bucket(this, "LldapWebAccessLogsBucket", {
      encryption: BucketEncryption.S3_MANAGED,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
      enforceSSL: true,
      removalPolicy: props.isDevelopmentEnv
        ? RemovalPolicy.DESTROY
        : RemovalPolicy.RETAIN,
      autoDeleteObjects: props.isDevelopmentEnv,
      serverAccessLogsBucket: serverAccessLogsBucket,
      serverAccessLogsPrefix: "access-logs/",
    });

    // Create security group for ALB
    const albSecurityGroup = new SecurityGroup(
      this,
      "LldapWebAlbSecurityGroup",
      {
        vpc: props.vpc,
        allowAllOutbound: false,
        description: "Security group for LLDAP web ALB",
      }
    );

    albSecurityGroup.addIngressRule(
      Peer.anyIpv4(),
      Port.tcp(443),
      "Allow HTTPS from internet"
    );

    albSecurityGroup.addIngressRule(
      Peer.anyIpv4(),
      Port.tcp(80),
      "Allow HTTP from internet for redirect"
    );

    albSecurityGroup.addEgressRule(
      serviceSecurityGroup,
      Port.tcp(17170),
      "Allow connection to LLDAP web interface"
    );

    // Allow ALB to connect to LLDAP web interface
    serviceSecurityGroup.addIngressRule(
      albSecurityGroup,
      Port.tcp(17170),
      "Allow ALB to web interface"
    );

    // Allow TAK server to connect to LDAP port
    serviceSecurityGroup.addIngressRule(
      Peer.ipv4(props.vpc.vpcCidrBlock),
      Port.tcp(3890),
      "Allow TAK server LDAP connections"
    );

    // Create ALB for LLDAP web interface
    const webLoadBalancer = new ApplicationLoadBalancer(
      this,
      "LldapWebLoadBalancer",
      {
        vpc: props.vpc,
        internetFacing: true,
        securityGroup: albSecurityGroup,
      }
    );

    webLoadBalancer.logAccessLogs(accessLogsBucket, "lldap-web-alb");

    const webTargetGroup = new ApplicationTargetGroup(
      this,
      "LldapWebTargetGroup",
      {
        vpc: props.vpc,
        port: 17170,
        protocol: ApplicationProtocol.HTTP,
        targetType: TargetType.IP,
        healthCheck: {
          path: "/",
          port: "17170",
          protocol: Protocol.HTTP,
          healthyThresholdCount: 2,
          unhealthyThresholdCount: 3,
          timeout: Duration.seconds(10),
          interval: Duration.seconds(30),
        },
      }
    );

    // Create certificate and Route53 record if hosted zone provided
    let certificate;
    if (props.hostedZone && props.domainName) {
      const ldapSubdomain = `ldap.${props.domainName}`;

      certificate = new Certificate(this, "LldapCertificate", {
        domainName: ldapSubdomain,
        validation: CertificateValidation.fromDns(props.hostedZone),
      });

      new ARecord(this, "LldapDnsRecord", {
        zone: props.hostedZone,
        recordName: ldapSubdomain,
        target: RecordTarget.fromAlias(new LoadBalancerTarget(webLoadBalancer)),
      });
      this.webUrl = `https://${ldapSubdomain}`;
    } else {
      this.webUrl = `https://${webLoadBalancer.loadBalancerDnsName}`;
    }

    // Add HTTP listener with redirect to HTTPS
    webLoadBalancer.addListener("HttpListener", {
      port: 80,
      protocol: ApplicationProtocol.HTTP,
      defaultAction: ListenerAction.redirect({
        protocol: "HTTPS",
        port: "443",
        permanent: true,
      }),
    });

    webLoadBalancer.addListener("WebListener", {
      port: 443,
      protocol: ApplicationProtocol.HTTPS,
      defaultTargetGroups: [webTargetGroup],
      certificates: certificate ? [certificate] : undefined,
    });

    webTargetGroup.addTarget(
      this.service.loadBalancerTarget({
        containerName: "LldapContainer",
        containerPort: 17170,
      })
    );

    this.host = this.service.serviceName;
    this.adminSecret = ldapAdminSecret;
    this.jwtSecret = jwtSecret;
    this.keySeed = keySeed;
  }
}
