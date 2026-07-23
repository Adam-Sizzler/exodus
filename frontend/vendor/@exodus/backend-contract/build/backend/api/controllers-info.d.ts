export declare const CONTROLLERS_INFO: {
    readonly AUTH: {
        readonly tag: "Auth Controller";
        readonly description: "Used to authenticate admin users.";
        readonly resource: "auth";
    };
    readonly PASSKEYS: {
        readonly tag: "Passkeys Controller";
        readonly description: "Management of Passkeys.";
        readonly resource: "passkeys";
    };
    readonly API_TOKENS: {
        readonly tag: "API Tokens Controller";
        readonly description: "Manage API tokens to use in your code. This controller can't be used with API token, only with Admin JWT token";
        readonly resource: "api-tokens";
    };
    readonly USERS: {
        readonly tag: "Users Controller";
        readonly description: "Manage users, change their status, reset traffic, etc.";
        readonly resource: "users";
    };
    readonly USERS_BULK_ACTIONS: {
        readonly tag: "Users Bulk Actions Controller";
        readonly description: "Bulk actions with users.";
        readonly resource: "users";
    };
    readonly HWID_USER_DEVICES: {
        readonly tag: "HWID User Devices Controller";
        readonly description: "";
        readonly resource: "hwid-user-devices";
    };
    readonly SUBSCRIPTION: {
        readonly tag: "[Public] Subscription Controller";
        readonly description: "Public Subscription Controller. Methods of this controller are not protected with auth. Use it only for public requests.";
        readonly resource: "subscription";
    };
    readonly SUBSCRIPTIONS: {
        readonly tag: "[Protected] Subscriptions Controller";
        readonly description: "Methods of this controller are protected with auth, most of them is returning the same informations as public Subscription Controller.";
        readonly resource: "subscriptions";
    };
    readonly NODES: {
        readonly tag: "Nodes Controller";
        readonly description: "";
        readonly resource: "nodes";
    };
    readonly NODE_PLUGINS: {
        readonly tag: "Node Plugins Controller";
        readonly description: "";
        readonly resource: "node-plugins";
    };
    readonly BANDWIDTH_STATS: {
        readonly tag: "Bandwidth Stats Controller";
        readonly description: "";
        readonly resource: "bandwidth-stats";
    };
    readonly IP_CONTROL: {
        readonly tag: "IP Management Controller";
        readonly description: "Management of IP addresses and connections.";
        readonly resource: "ip-control";
    };
    readonly CONFIG_PROFILES: {
        readonly tag: "Config Profiles Controller";
        readonly description: "Management of Config Profiles.";
        readonly resource: "config-profiles";
    };
    readonly INTERNAL_SQUADS: {
        readonly tag: "Internal Squads Controller";
        readonly description: "Management of Internal Squads.";
        readonly resource: "internal-squads";
    };
    readonly EXTERNAL_SQUADS: {
        readonly tag: "External Squads Controller";
        readonly description: "Management of External Squads.";
        readonly resource: "external-squads";
    };
    readonly HOSTS: {
        readonly tag: "Hosts Controller";
        readonly description: "";
        readonly resource: "hosts";
    };
    readonly HOSTS_BULK_ACTIONS: {
        readonly tag: "Hosts Bulk Actions Controller";
        readonly description: "";
        readonly resource: "hosts";
    };
    readonly SUBSCRIPTION_TEMPLATE: {
        readonly tag: "Subscription Template Controller";
        readonly description: "";
        readonly resource: "subscription-template";
    };
    readonly SUBSCRIPTION_SETTINGS: {
        readonly tag: "Subscription Settings Controller";
        readonly description: "";
        readonly resource: "subscription-settings";
    };
    readonly INFRA_BILLING: {
        readonly tag: "Infra Billing Controller";
        readonly description: "";
        readonly resource: "infra-billing";
    };
    readonly SYSTEM: {
        readonly tag: "System Controller";
        readonly description: "";
        readonly resource: "system";
    };
    readonly KEYGEN: {
        readonly tag: "Keygen Controller";
        readonly description: "Generation of SECRET_KEY for Remnawave Node.";
        readonly resource: "keygen";
    };
    readonly SUBSCRIPTION_REQUEST_HISTORY: {
        readonly tag: "Subscription Request History Controller";
        readonly description: "";
        readonly resource: "subscription-request-history";
    };
    readonly SNIPPETS: {
        readonly tag: "Snippets Controller";
        readonly description: "";
        readonly resource: "snippets";
    };
    readonly REMNAAWAVE_SETTINGS: {
        readonly tag: "Remnawave Settings Controller";
        readonly description: "";
        readonly resource: "remnawave-settings";
    };
    readonly SUBSCRIPTION_PAGE_CONFIGS: {
        readonly tag: "Subscription Page Configs Controller";
        readonly description: "";
        readonly resource: "subscription-page-configs";
    };
    readonly METADATA: {
        readonly tag: "Metadata Controller";
        readonly description: "Manage arbitrary metadata for Users and Nodes.";
        readonly resource: "metadata";
    };
};
//# sourceMappingURL=controllers-info.d.ts.map