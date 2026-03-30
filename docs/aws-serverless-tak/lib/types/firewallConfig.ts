export interface FirewallConfig {
  enabled?: boolean;
  name?: string;
  geoAllow?: string[];
  connectionTimeout?: number;
}
