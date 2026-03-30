import { RemovalPolicy } from "aws-cdk-lib";
import { IVpc, SecurityGroup } from "aws-cdk-lib/aws-ec2";
import { FileSystem } from "aws-cdk-lib/aws-efs";
import { AnyPrincipal, Effect, PolicyStatement } from "aws-cdk-lib/aws-iam";
import { Key } from "aws-cdk-lib/aws-kms";
import { Construct } from "constructs";

export interface TakFilesystemProps {
  vpc: IVpc;
  kmsKey?: Key;
  isDevelopmentEnv?: boolean;
}

export class TakFileSystem extends Construct {
  public readonly efsFileSystem: FileSystem;

  constructor(scope: Construct, id: string, props: TakFilesystemProps) {
    super(scope, id);

    const removalPolicy = props.isDevelopmentEnv
      ? RemovalPolicy.DESTROY
      : RemovalPolicy.RETAIN;

    const efsSecurityGroup = new SecurityGroup(this, "EFSSecurityGroup", {
      vpc: props.vpc!,
      description: "EFS Security Group",
      allowAllOutbound: true,
    });

    // Create the EFS file system
    this.efsFileSystem = new FileSystem(this, "EfsFileSystem", {
      vpc: props.vpc,
      removalPolicy: removalPolicy,
      encrypted: true,
      kmsKey:
        props.kmsKey ??
        new Key(scope, `${id}-FilesystemKey`, {
          enableKeyRotation: true,
          removalPolicy: removalPolicy,
        }),
      vpcSubnets: { subnetGroupName: "Private" },
      securityGroup: efsSecurityGroup,
    });

    // Allow clients to mount and write to the file system
    this.efsFileSystem.addToResourcePolicy(
      new PolicyStatement({
        effect: Effect.ALLOW,
        actions: [
          "elasticfilesystem:ClientMount",
          "elasticfilesystem:ClientWrite",
        ],
        principals: [new AnyPrincipal()],
        conditions: {
          Bool: {
            "elasticfilesystem:AccessedViaMountTarget": "true",
          },
        },
      })
    );
  }
}
