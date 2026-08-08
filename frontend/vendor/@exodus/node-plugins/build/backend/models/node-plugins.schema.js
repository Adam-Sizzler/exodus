"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.NodePluginSchema = exports.HaproxyAuthPluginSchema = exports.EgressFilterPluginSchema = exports.IngressFilterPluginSchema = exports.SharedListSchema = void 0;
const zod_1 = require("zod");
const DOCS_LINK = `\n\n[📖 Documentation](https://docs.exodus.dev/docs/learn/node-plugins)`;
const IpCidrOrExtSchema = zod_1.z
    .union([
    zod_1.z.union([zod_1.z.cidrv4(), zod_1.z.cidrv6()]),
    zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()]),
    zod_1.z.string().startsWith('ext:'),
])
    .meta({
    title: 'IP address or CIDR range',
    markdownDescription: `IP address or CIDR range. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**.${DOCS_LINK}`,
});
exports.SharedListSchema = zod_1.z.discriminatedUnion('type', [
    zod_1.z.object({
        name: zod_1.z.string().startsWith('ext:'),
        type: zod_1.z.literal('ipList'),
        items: zod_1.z.array(zod_1.z.union([zod_1.z.cidrv4(), zod_1.z.cidrv6(), zod_1.z.union([zod_1.z.ipv4(), zod_1.z.ipv6()])])),
    }),
    zod_1.z.object({
        name: zod_1.z.string().startsWith('ext:'),
        type: zod_1.z.literal('asList'),
        items: zod_1.z.array(zod_1.z.int().min(1).max(4294967295)),
    }),
]);
exports.IngressFilterPluginSchema = zod_1.z.object({
    enabled: zod_1.z.boolean().meta({
        title: 'Enabled',
        markdownDescription: `If this plugin is enabled, all IP addresses specified in the **blockedIps** object will be blocked via nftables. **Use with caution.**${DOCS_LINK}`,
    }),
    blockedIps: zod_1.z.array(IpCidrOrExtSchema).meta({
        title: 'Blocked IPs',
        markdownDescription: `List of IP addresses and CIDR ranges to block via nftables. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**.${DOCS_LINK}`,
    }),
});
exports.EgressFilterPluginSchema = zod_1.z.object({
    enabled: zod_1.z.boolean().meta({
        title: 'Enabled',
        markdownDescription: `If this plugin is enabled, outbound connections to specified IP addresses and ports will be blocked. **Use with caution.**${DOCS_LINK}`,
    }),
    blockedIps: zod_1.z
        .array(IpCidrOrExtSchema)
        .optional()
        .meta({
        title: 'Blocked IPs',
        markdownDescription: `List of destination IP addresses and CIDR ranges to block. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**. \n\n Example: \`["10.0.0.1", "ext:blocked_destinations"]\`${DOCS_LINK}`,
    }),
    blockedPorts: zod_1.z
        .array(zod_1.z.int().min(1).max(65535))
        .optional()
        .meta({
        title: 'Blocked Ports',
        markdownDescription: `List of destination ports to block. \n\n Example: \`[25, 465, 587]\` to block SMTP traffic.${DOCS_LINK}`,
    }),
});
exports.HaproxyAuthPluginSchema = zod_1.z.object({
    inboundTags: zod_1.z
        .array(zod_1.z.string())
        .optional()
        .default([])
        .meta({
        title: 'Inbound Tags',
        markdownDescription: `List of inbound tags that participate in HAProxy authentication. Use an empty array to disable HAProxy authentication for all inbounds. Use "*" to include every supported inbound assigned to the node.${DOCS_LINK}`,
    }),
});
exports.NodePluginSchema = zod_1.z.object({
    sharedLists: zod_1.z
        .array(exports.SharedListSchema)
        .optional()
        .default([])
        .meta({
        title: 'Shared Lists',
        markdownDescription: `Array of shared lists, which can be used in other plugins. Optional.${DOCS_LINK}`,
    }),
    ingressFilter: exports.IngressFilterPluginSchema.optional().meta({
        title: 'Ingress Filter',
        markdownDescription: `Ingress Filter Plugin configuration. Optional.${DOCS_LINK}`,
    }),
    egressFilter: exports.EgressFilterPluginSchema.optional().meta({
        title: 'Egress Filter',
        markdownDescription: `Egress Filter Plugin configuration. Optional.${DOCS_LINK}`,
    }),
    haproxyAuth: exports.HaproxyAuthPluginSchema.optional().meta({
        title: 'HAProxy Auth',
        markdownDescription: `HAProxy Auth Plugin configuration. Optional.${DOCS_LINK}`,
    }),
});
