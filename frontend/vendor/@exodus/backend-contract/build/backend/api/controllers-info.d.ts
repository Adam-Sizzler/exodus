export declare const CONTROLLERS_INFO: {
    readonly AUTH: {
        readonly tag: "Auth Controller";
        readonly description: "Used to authenticate admin users.";
    };
    readonly PASSKEYS: {
        readonly tag: "Passkeys Controller";
        readonly description: "Management of Passkeys.";
    };
    readonly API_TOKENS: {
        readonly tag: "API Tokens Controller";
        readonly description: "Manage API tokens to use in your code. This controller can't be used with API token, only with Admin JWT token";
    };
    readonly USERS: {
        readonly tag: "Users Controller";
        readonly description: "Manage users, change their status, reset traffic, etc.";
    };
    readonly USERS_BULK_ACTIONS: {
        readonly tag: "Users Bulk Actions Controller";
        readonly description: "Bulk actions with users.";
    };
    readonly HWID_USER_DEVICES: {
        readonly tag: "HWID User Devices Controller";
        readonly description: "";
    };
    readonly SUBSCRIPTION: {
        readonly tag: "[Public] Subscription Controller";
        readonly description: "Public Subscription Controller. Methods of this controller are not protected with auth. Use it only for public requests.";
    };
    readonly SUBSCRIPTIONS: {
        readonly tag: "[Protected] Subscriptions Controller";
        readonly description: "Methods of this controller are protected with auth, most of them is returning the same informations as public Subscription Controller.";
    };
    readonly NODES: {
        readonly tag: "Nodes Controller";
        readonly description: "";
    };
    readonly NODE_PLUGINS: {
        readonly tag: "Node Plugins Controller";
        readonly description: "";
    };
    readonly BANDWIDTH_STATS: {
        readonly tag: "Bandwidth Stats Controller";
        readonly description: "";
    };
    readonly IP_CONTROL: {
        readonly tag: "IP Management Controller";
        readonly description: "Management of IP addresses and connections.";
    };
    readonly CONFIG_PROFILES: {
        readonly tag: "Config Profiles Controller";
        readonly description: "Management of Config Profiles.";
    };
    readonly INTERNAL_SQUADS: {
        readonly tag: "Internal Squads Controller";
        readonly description: "Management of Internal Squads.";
    };
    readonly EXTERNAL_SQUADS: {
        readonly tag: "External Squads Controller";
        readonly description: "Management of External Squads.";
    };
    readonly HOSTS: {
        readonly tag: "Hosts Controller";
        readonly description: "";
    };
    readonly HOSTS_BULK_ACTIONS: {
        readonly tag: "Hosts Bulk Actions Controller";
        readonly description: "";
    };
    readonly SUBSCRIPTION_TEMPLATE: {
        readonly tag: "Subscription Template Controller";
        readonly description: "";
    };
    readonly SUBSCRIPTION_SETTINGS: {
        readonly tag: "Subscription Settings Controller";
        readonly description: "";
    };
    readonly INFRA_BILLING: {
        readonly tag: "Infra Billing Controller";
        readonly description: "";
    };
    readonly SYSTEM: {
        readonly tag: "System Controller";
        readonly description: "";
    };
    readonly KEYGEN: {
        readonly tag: "Keygen Controller";
        readonly description: "Generation of SSL_CERT for Exodus Node.";
    };
    readonly SUBSCRIPTION_REQUEST_HISTORY: {
        readonly tag: "Subscription Request History Controller";
        readonly description: "";
    };
    readonly SNIPPETS: {
        readonly tag: "Snippets Controller";
        readonly description: "";
    };
    readonly EXODUS_SETTINGS: {
        readonly tag: "Exodus Settings Controller";
        readonly description: "";
    };
    readonly SUBSCRIPTION_PAGE_CONFIGS: {
        readonly tag: "Subscription Page Configs Controller";
        readonly description: "";
    };
};
//# sourceMappingURL=controllers-info.d.ts.map