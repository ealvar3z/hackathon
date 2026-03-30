import * as cdk from "aws-cdk-lib";
import { Vpc, CfnRouteTable, CfnRoute, ISubnet } from "aws-cdk-lib/aws-ec2";
import {
  CfnFirewall,
  CfnFirewallPolicy,
  CfnRuleGroup,
} from "aws-cdk-lib/aws-networkfirewall";
import { Construct } from "constructs";
import { FirewallConfig } from "../types/firewallConfig";

export interface NetworkFirewallProps {
  vpc: Vpc;
  isDevelopmentEnv?: boolean;
  firewallConfig?: FirewallConfig;
}

export class NetworkFirewall extends Construct {
  public readonly firewall: CfnFirewall;
  public readonly firewallPolicy: CfnFirewallPolicy;
  public readonly trustSubnets: ISubnet[];
  public readonly dmzSubnets: ISubnet[];

  constructor(
    scope: Construct,
    id: string,
    private props: NetworkFirewallProps
  ) {
    super(scope, id);

    // Set Trust subnets for Network Firewall
    this.trustSubnets = props.vpc.selectSubnets({
      subnetGroupName: "Trust",
    }).subnets;

    // Set DMZ subnets for Network Firewall
    this.dmzSubnets = props.vpc.selectSubnets({
      subnetGroupName: "DMZ",
    }).subnets;

    // Create NFW Policy and NFW
    this.firewallPolicy = this.createFirewallPolicy(props.firewallConfig);
    this.firewall = this.createFirewall();

    // Add IGW -> NFW and Trust -> NFW
    this.createIgwToNfwRoutes();
    this.createTrustToNfwRoutes();
  }

  // Create routes for IGW, CIDR destinations are trust subnets with targets of respective NFW endpoint
  private createIgwToNfwRoutes(): void {
    // New IGW RT for explicit edge association
    const igwRouteTable = new CfnRouteTable(this, "IgwRouteTable", {
      vpcId: this.props.vpc.vpcId,
      tags: [
        { key: "Name", value: "CoreInfraStack/TakServerVpc/TakServerVPC/IGW" },
      ],
    });

    // Create routes for each 'Trust' subnet
    this.trustSubnets.forEach((subnet, index) => {
      const publicSubnetCidr = subnet.ipv4CidrBlock;
      const firewallEndpointId = this.getFirewallEndpointForAz(index);

      new CfnRoute(this, `IgwToNfwRoute${index + 1}`, {
        routeTableId: igwRouteTable.ref,
        destinationCidrBlock: publicSubnetCidr,
        vpcEndpointId: firewallEndpointId,
      }).addDependency(this.firewall);
    });

    // Associate route table with IGW
    new cdk.CfnResource(this, "IgwEdgeAssociation", {
      type: "AWS::EC2::GatewayRouteTableAssociation",
      properties: {
        GatewayId: this.props.vpc.internetGatewayId,
        RouteTableId: igwRouteTable.ref,
      },
    });
  }

  // Create routes for Trust Subnets, internet bound traffic destined for respective firewall endpoint
  private createTrustToNfwRoutes(): void {
    // Route from Trust subnets to DMZ Firewall for outbound traffic
    this.trustSubnets.forEach((subnet, index) => {
      const firewallEndpointId = this.getFirewallEndpointForAz(index);

      new CfnRoute(this, `DmzToIgwRoute${index + 1}`, {
        routeTableId: subnet.routeTable.routeTableId,
        destinationCidrBlock: "0.0.0.0/0",
        vpcEndpointId: firewallEndpointId,
      }).addDependency(this.firewall);
    });
  }

  // Get function for Firewall Endpoints
  private getFirewallEndpointForAz(azIndex: number): string {
    // Select GWLB Endpoint by AZ Index
    return cdk.Fn.select(
      1,
      cdk.Fn.split(":", cdk.Fn.select(azIndex, this.firewall.attrEndpointIds))
    );
  }

  // Function to create NFW policy containing geo-blocking and rate limiting based off config
  private createFirewallPolicy(config?: FirewallConfig): CfnFirewallPolicy {
    // Create geo-blocking rule group if countries are specified
    const geoAllowRuleGroup = this.creategeoAllowRuleGroup(
      config?.geoAllow || ["US"]
    );

    // Collect all rule group references
    const statefulRuleGroupReferences: Array<{ resourceArn: string }> = [];
    const statelessRuleGroupReferences: Array<{
      resourceArn: string;
      priority: number;
    }> = [];

    if (geoAllowRuleGroup) {
      statefulRuleGroupReferences.push({
        resourceArn: geoAllowRuleGroup.attrRuleGroupArn,
      });
    }

    const firewallPolicy = new CfnFirewallPolicy(this, "CDKFirewallPolicy", {
      firewallPolicy: {
        statelessDefaultActions: ["aws:forward_to_sfe"],
        statelessFragmentDefaultActions: ["aws:forward_to_sfe"],
        statefulEngineOptions: {
          ruleOrder: "DEFAULT_ACTION_ORDER",
          flowTimeouts: {
            tcpIdleTimeoutSeconds: config?.connectionTimeout || 300,
          },
        },
        statefulRuleGroupReferences: statefulRuleGroupReferences,
        statelessRuleGroupReferences: statelessRuleGroupReferences,
      },
      firewallPolicyName: config?.name || "NetworkFirewallPolicy",
      description:
        "Network Firewall Policy with geo-blocking and connection limiting",
    });

    return firewallPolicy;
  }

  // Geo blocking rule group
  private creategeoAllowRuleGroup(
    geoAllowCountries?: string[]
  ): CfnRuleGroup | undefined {
    if (!geoAllowCountries || geoAllowCountries.length === 0) {
      return undefined;
    }

    // Validate country codes
    const validCountryCodes = geoAllowCountries
      .map((country) => country.toUpperCase())
      .filter((code) => code.length === 2); // Basic validation for 2-letter codes

    if (validCountryCodes.length === 0) {
      return undefined;
    }

    // Combine all country codes into one rule for efficiency
    const allowCountries = validCountryCodes.join(",");
    const geoAllowRule = {
      action: "DROP",
      header: {
        destinationPort: "ANY",
        direction: "ANY", // evaluate outbound traffic direction
        destination: "ANY",
        source: "ANY",
        sourcePort: "ANY",
        protocol: "IP",
      },
      ruleOptions: [
        { keyword: "sid", settings: ["1000001"] },
        { keyword: "geoip", settings: [`any,!${allowCountries}`] },
      ],
    };

    return new CfnRuleGroup(this, "geoAllowRuleGroup", {
      capacity: 100,
      ruleGroupName: "geoAllowRuleGroup",
      type: "STATEFUL",
      description: `Geo-allow rule group for countries: ${validCountryCodes.join(
        ", "
      )}`,
      ruleGroup: {
        rulesSource: {
          statefulRules: [geoAllowRule],
        },
      },
    });
  }

  // Create firewall with previously created firewall policy
  private createFirewall(): CfnFirewall {
    const subnetMappings = this.dmzSubnets.map((subnet) => ({
      subnetId: subnet.subnetId,
    }));

    const config = this.props.firewallConfig;

    const nfw = new CfnFirewall(this, "NFW", {
      firewallName: config?.name || "NetworkFirewall",
      firewallPolicyArn: this.firewallPolicy.attrFirewallPolicyArn,
      subnetMappings: subnetMappings,
      vpcId: this.props.vpc.vpcId,
      deleteProtection: !this.props.isDevelopmentEnv,
      tags: [{ key: "Name", value: config?.name || "Network-Firewall" }],
    });

    return nfw;
  }
}
