import { Duration, RemovalPolicy } from "aws-cdk-lib";
import { PolicyStatement, ServicePrincipal } from "aws-cdk-lib/aws-iam";
import {
  BlockPublicAccess,
  Bucket,
  BucketAccessControl,
  BucketEncryption,
  BucketProps,
  ObjectOwnership,
} from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";
import { addTlsAndSigV4ResourcePolicies } from "./secure-bucket";

export interface LoggingBucketProps extends BucketProps {
  /**
   * S3 prefix to use for cloudtrail logging
   *
   * @default 'cloudtrail-trails/'
   */
  cloudTrailLogsBucketPrefix?: string;

  /**
   * Expiration of the logs
   *
   * @default Duration.days(2)
   */
  logExpiration?: Duration;

  /**
   * Applies more relaxed security settings for a development environment
   *
   * @default false
   */
  isDevelopmentEnv?: boolean;
}

const defaultCloudTrailPrefix = "cloudtrail-trails/";

/**
 * A bucket configured to take CloudTrail Logs
 */
export class LoggingBucket extends Bucket {
  constructor(
    scope: Construct,
    id: string,
    private props?: LoggingBucketProps
  ) {
    super(scope, id, {
      ...props,
      accessControl: BucketAccessControl.LOG_DELIVERY_WRITE,
      // BucketAccessControl.LOG_DELIVERY_WRITE Does not currently support KMS Key encryption.
      // If this changes in future, we may want to have LoggingBucket extend SecureBucket for KMS encryption.
      encryption: BucketEncryption.S3_MANAGED,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
      lifecycleRules: [
        {
          id: "CloudTrail Expiration Rule",
          prefix: props?.cloudTrailLogsBucketPrefix ?? defaultCloudTrailPrefix,
          enabled: true,
          expiration: props?.logExpiration ?? Duration.days(2),
        },
      ],
      // At this time, ACLs cannot be completely disabled because they are used by CloudTrail to write the logs.
      objectOwnership: ObjectOwnership.BUCKET_OWNER_PREFERRED,
      ...envDependentProps(!!props?.isDevelopmentEnv),
    });
    this.constructResourcePolicy();
  }

  private constructResourcePolicy() {
    addTlsAndSigV4ResourcePolicies(this);

    this.addToResourcePolicy(
      new PolicyStatement({
        sid: "AWS CloudTrail Acl Check",
        principals: [new ServicePrincipal("cloudtrail.amazonaws.com")],
        actions: ["s3:GetBucketAcl"],
        resources: [this.bucketArn],
      })
    );

    this.addToResourcePolicy(
      new PolicyStatement({
        sid: "AWS CloudTrail Write",
        principals: [new ServicePrincipal("cloudtrail.amazonaws.com")],
        actions: ["s3:PutObject"],
        resources: [
          `${this.bucketArn}/${
            this.props?.cloudTrailLogsBucketPrefix ?? defaultCloudTrailPrefix
          }*`,
        ],
        conditions: {
          StringEquals: {
            "s3:x-amz-acl": "bucket-owner-full-control",
          },
        },
      })
    );
  }
}

/**
 * In DEV:
 *  - The bucket will be auto-destroyed
 *  - The bucket does not log access to itself to ensure that it can be destroyed cleanly
 *
 * In DEMO/PROD:
 *  - The bucket logs accesses to itself
 *  - The bucket has to be deleted manually when the solution is destroyed
 *  - The access logging to itself prevents the deletion, so in order to
 *    destroy the bucket you must first manually disable access logging
 *
 * @private
 */
const envDependentProps = (isDevelopmentEnv: boolean): Partial<BucketProps> =>
  isDevelopmentEnv
    ? {
        removalPolicy: RemovalPolicy.DESTROY,
        autoDeleteObjects: true,
      }
    : {
        serverAccessLogsPrefix: "logging-logs/",
        removalPolicy: RemovalPolicy.RETAIN,
      };
