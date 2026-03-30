export interface LdapConfig {
  deployLdapServer?: boolean;
  host?: string;
  port?: number;
  baseDn?: string;
  bindDn?: string;
  bindPassword?: string;
  userSearchBase?: string;
  groupSearchBase?: string;
}

export interface TakServerConfig {
  adminEmail: string;
  authType: string;
  domainName: string;
  hostedZoneId: string;
  subdomainName: string;
  certCountry: string;
  certState: string;
  certCity: string;
  certOrg: string;
  certOrgUnit: string;
  ldap?: LdapConfig;
  forceConfigRegen?: boolean;
}
