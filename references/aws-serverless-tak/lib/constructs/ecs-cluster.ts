import { Cluster, ClusterProps, ContainerInsights } from "aws-cdk-lib/aws-ecs";
import { Construct } from "constructs";

export class TakCluster extends Construct {
  public readonly cluster: Cluster;

  constructor(scope: Construct, id: string, props: ClusterProps) {
    super(scope, id);

    this.cluster = new Cluster(this, "TakCluster", {
      clusterName: "TAKCluster",
      containerInsightsV2: ContainerInsights.ENABLED,
      enableFargateCapacityProviders: true,
      vpc: props.vpc,
      defaultCloudMapNamespace: { name: "tak" },
      ...props,
    });
  }
}
