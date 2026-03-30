import { RemovalPolicy } from "aws-cdk-lib";
import {
  AnyPrincipal,
  Effect,
  PolicyStatement,
  ServicePrincipal,
} from "aws-cdk-lib/aws-iam";
import { Key } from "aws-cdk-lib/aws-kms";
import {
  BlockPublicAccess,
  Bucket,
  BucketEncryption,
  BucketProps,
  ObjectOwnership,
} from "aws-cdk-lib/aws-s3";
import { Construct } from "constructs";

export type SecureBucketProps = Omit<
  BucketProps,
  "blockPublicAccess" | "encryption"
>;

/**
 * Secure S3 Bucket
 *
 * This Construct creates a bucket that:
 * * Uses SSE-KMS with an S3 Bucket Key
 * * Has all public access blocked
 * * Enforces TLS access only
 * * Rejects clients that do not use SigV4
 * * Disables object ACLs by default. Bucket owner owns all objects.
 */
export class SecureBucket extends Bucket {
  constructor(scope: Construct, id: string, props?: SecureBucketProps) {
    const secureProps: BucketProps = {
      ...props,
      encryptionKey:
        props?.encryptionKey ??
        new Key(scope, `${id}-SecureBucketKey`, {
          enableKeyRotation: true,
          removalPolicy: props?.removalPolicy ?? RemovalPolicy.RETAIN,
        }),
      encryption: BucketEncryption.KMS,
      bucketKeyEnabled: props?.bucketKeyEnabled ?? true,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
      enforceSSL: true,
      objectOwnership:
        props?.objectOwnership ?? ObjectOwnership.BUCKET_OWNER_ENFORCED,
    };

    super(scope, id, secureProps);

    addTlsAndSigV4ResourcePolicies(this);
  }
}

export function addTlsAndSigV4ResourcePolicies(bucket: Bucket) {
  bucket.addToResourcePolicy(
    new PolicyStatement({
      sid: "Deny requests that do not use TLS",
      effect: Effect.DENY,
      principals: [new AnyPrincipal()],
      actions: ["s3:*"],
      resources: [bucket.bucketArn, bucket.arnForObjects("*")],
      conditions: {
        Bool: {
          "aws:SecureTransport": "false",
        },
      },
    })
  );

  bucket.addToResourcePolicy(
    new PolicyStatement({
      sid: "Deny requests that do not use SigV4",
      effect: Effect.DENY,
      notPrincipals: [new ServicePrincipal("s3")],
      actions: ["s3:*"],
      resources: [bucket.bucketArn, bucket.arnForObjects("*")],
      conditions: {
        StringNotEquals: {
          "s3:signatureversion": ["AWS4-HMAC-SHA256", "AWS4-ECDSA-P256-SHA256"],
        },
      },
    })
  );
}
