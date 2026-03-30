import { EnvironmentConfig } from "./environmentConfig";
import { FirewallConfig } from "./firewallConfig";
import { InfraConfig } from "./infrastructureConfig";
import { TakServerConfig } from "./takConfig";

export interface Config {
  env?: EnvironmentConfig;
  infra?: InfraConfig;
  tak: TakServerConfig;
  firewall?: FirewallConfig;
}
