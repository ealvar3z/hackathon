import { RemovalPolicy } from "aws-cdk-lib";
import {
  FlowLogDestination,
  GatewayVpcEndpointAwsService,
  IpAddresses,
  SubnetType,
  Vpc,
  VpcProps,
  CfnEIP,
  CfnNatGateway,
  CfnRoute,
  InterfaceVpcEndpoint,
  InterfaceVpcEndpointAwsService,
} from "aws-cdk-lib/aws-ec2";
import { Role, ServicePrincipal } from "aws-cdk-lib/aws-iam";
import { LogGroup } from "aws-cdk-lib/aws-logs";
import { Construct, IConstruct } from "constructs";

export interface TakVpcProps {
  /**
   * If set to true, all logging data will be deleted when the resource is destroyed
   *
   * @default false
   */
  isDevelopmentEnv?: boolean;

  /**
   * By default, the VPC is created with a flow log.
   *
   * The creation of the flow log can be suppressed by setting `noFlowLog` to `true`.
   *
   * @see https://docs.aws.amazon.com/vpc/latest/userguide/flow-logs.html
   *
   * @default false
   */
  noFlowLog?: boolean;

  /**
   *Enables deployment of an AWS Network Firewall in the VPC.
   *
   * If set to `true`, a Network Firewall is provisioned in each Availability Zone,
   * along with corresponding firewall subnets and routing configuration.
   *
   * @default true
   */
  isFirewallEnabled?: boolean;

  /**
   * Additional options to set in the VPC
   */
  vpcOptions?: VpcProps;
}

// Creates a new VPC with default configuration.
export class TakServerVpc extends Construct {
  readonly vpc: Vpc;

  constructor(scope: IConstruct, id: string, props: TakVpcProps) {
    super(scope, id);

    const removalPolicy = props.isDevelopmentEnv
      ? RemovalPolicy.DESTROY
      : RemovalPolicy.RETAIN;

    if (props.isFirewallEnabled == true) {
      this.vpc = new Vpc(this, "TakServerVPC", {
        ...props.vpcOptions,
        vpcName: "TakServerVPC",
        ipAddresses: IpAddresses.cidr("10.13.0.0/16"),
        maxAzs: 2,
        createInternetGateway: true,
        subnetConfiguration: [
          {
            cidrMask: 21,
            name: "Private",
            subnetType: SubnetType.PRIVATE_ISOLATED,
          },
          {
            cidrMask: 24,
            name: "Trust",
            subnetType: SubnetType.PRIVATE_ISOLATED,
          },
          {
            cidrMask: 24,
            name: "DMZ",
            subnetType: SubnetType.PUBLIC,
          },
        ],
        gatewayEndpoints: {
          S3: {
            service: GatewayVpcEndpointAwsService.S3,
          },
        },
        flowLogs: {},
      });

      // Get the Trust subnets
      const trustSubnets = this.vpc.selectSubnets({
        subnetGroupName: "Trust",
      }).subnets;

      // Get the Private subnets
      const privateSubnets = this.vpc.selectSubnets({
        subnetGroupName: "Private",
      }).subnets;

      // Store the NAT Gateways to assign to private subnets
      const natGateways: CfnNatGateway[] = [];

      // Create NAT Gateways in each isolated subnet
      trustSubnets.forEach((subnet, index) => {
        // Create an EIP for each NAT Gateway
        const eip = new CfnEIP(this, `TrustNatGatewayEIP${index}`, {
          domain: "vpc",
        });

        // Create NAT Gateway per AZ
        const natGateway = new CfnNatGateway(this, `TrustNatGateway${index}`, {
          subnetId: subnet.subnetId,
          allocationId: eip.attrAllocationId,
          tags: [
            {
              key: "Name",
              value: `Trust-NAT-Gateway-${index}`,
            },
          ],
        });

        natGateways.push(natGateway);
      });

      // Create a route in each private subnet to the respective NAT Gateway
      privateSubnets.forEach((subnet, index) => {
        const natGateway = natGateways[index];

        new CfnRoute(this, `PrivateSubnetRoute${index}`, {
          routeTableId: subnet.routeTable.routeTableId,
          destinationCidrBlock: "0.0.0.0/0",
          natGatewayId: natGateway.ref,
        });
      });

      // Create Secrets Manager Endpoint
      new InterfaceVpcEndpoint(this, "VPC Endpoint", {
        vpc: this.vpc,
        service: InterfaceVpcEndpointAwsService.SECRETS_MANAGER,
        subnets: {
          subnets: privateSubnets,
        },
        privateDnsEnabled: true,
      });
    } else {
      this.vpc = new Vpc(this, "TakServerVPC", {
        ...props.vpcOptions,
        vpcName: "TakServerVPC",
        ipAddresses: IpAddresses.cidr("10.13.0.0/16"),
        maxAzs: 2,
        natGateways: 2,
        subnetConfiguration: [
          {
            cidrMask: 21,
            name: "Private",
            subnetType: SubnetType.PRIVATE_WITH_EGRESS,
          },
          {
            cidrMask: 24,
            name: "Public",
            subnetType: SubnetType.PUBLIC,
          },
        ],
        gatewayEndpoints: {
          S3: {
            service: GatewayVpcEndpointAwsService.S3,
          },
        },
        flowLogs: {},
      });
    }

    if (!props.noFlowLog) {
      const logGroup = new LogGroup(this, "VpcFlowLogGroup", {
        removalPolicy,
      });

      const flowLogRole = new Role(this, "VpcFlowLogRole", {
        assumedBy: new ServicePrincipal("vpc-flow-logs.amazonaws.com"),
      });

      this.vpc.addFlowLog("VpcFlowLog", {
        destination: FlowLogDestination.toCloudWatchLogs(logGroup, flowLogRole),
      });
    }
  }
}
