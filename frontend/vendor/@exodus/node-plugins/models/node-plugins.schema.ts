import { z } from 'zod';

const DOCS_LINK = `\n\n[📖 Documentation](https://docs.ex/docs/learn/node-plugins)`;

// https://github.com/colinhacks/zod/issues/5944
const IPV6 = z.regexes.ipv6.source.slice(1, -1);

const ipv6 = () => z.string().regex(new RegExp(`^(${IPV6})$`), { error: 'Invalid IPv6 address' });
const cidrv6 = () =>
    z.string().regex(new RegExp(`^(${IPV6})\\/(12[0-8]|1[01][0-9]|[1-9]?[0-9])$`), {
        error: 'Invalid IPv6 CIDR range',
    });

const IpCidrOrExtSchema = z
    .union([
        z.union([z.cidrv4(), cidrv6()]),
        z.union([z.ipv4(), ipv6()]),
        z.string().startsWith('ext:'),
    ])
    .meta({
        title: 'IP address or CIDR range',
        markdownDescription: `IP address or CIDR range. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**.${DOCS_LINK}`,
    });

const IpListSchema = z
    .object({
        type: z.literal('ipList').meta({
            title: 'IP List',
            markdownDescription: `A list of IP addresses and CIDR ranges.${DOCS_LINK}`,
        }),
        items: z.array(z.union([z.cidrv4(), cidrv6(), z.union([z.ipv4(), ipv6()])])).meta({
            title: 'IP addresses and CIDR ranges',
            markdownDescription: [
                'IPv4 and IPv6 addresses, plain or as CIDR ranges. Mixing both families in one list is allowed.',
                '',
                '```json',
                '{',
                '  "type": "ipList",',
                '  "items": ["1.1.1.1", "10.0.0.0/8", "2001:db8::1", "2001:db8::/32"]',
                '}',
                '```',
                DOCS_LINK,
            ].join('\n'),
        }),
    })
    .meta({
        title: 'IP List',
        markdownDescription: `Shared list of IP addresses and CIDR ranges.${DOCS_LINK}`,
    });

const AsListSchema = z
    .object({
        type: z.literal('asList').meta({
            title: 'AS List',
            markdownDescription: `A list of autonomous system numbers.${DOCS_LINK}`,
        }),
        items: z.array(z.int().min(1).max(4294967295)).meta({
            title: 'Autonomous system numbers',
            markdownDescription: [
                'ASN numbers without the `AS` prefix, from 1 to 4294967295.',
                '',
                '```json',
                '{',
                '  "type": "asList",',
                '  "items": [13335, 15169, 32934]',
                '}',
                '```',
                DOCS_LINK,
            ].join('\n'),
        }),
    })
    .meta({
        title: 'AS List',
        markdownDescription: `Shared list of autonomous system numbers.${DOCS_LINK}`,
    });

export const SharedListConfigSchema = z
    .discriminatedUnion('type', [IpListSchema, AsListSchema])
    .meta({
        title: 'Shared List',
        markdownDescription: `Shared list body. Pick **ipList** for IP addresses and CIDR ranges, or **asList** for autonomous system numbers.${DOCS_LINK}`,
    });

export const SharedListSchema = z.discriminatedUnion('type', [
    IpListSchema.extend({ name: z.string().startsWith('ext:') }),
    AsListSchema.extend({ name: z.string().startsWith('ext:') }),
]);

export const IngressFilterPluginSchema = z.object({
    enabled: z.boolean().meta({
        title: 'Enabled',
        markdownDescription: `If this plugin is enabled, all IP addresses specified in the **blockedIps** object will be blocked via nftables. **Use with caution.**${DOCS_LINK}`,
    }),
    blockedIps: z.array(IpCidrOrExtSchema).meta({
        title: 'Blocked IPs',
        markdownDescription: `List of IP addresses and CIDR ranges to block via nftables. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**.${DOCS_LINK}`,
    }),
});

export const EgressFilterPluginSchema = z.object({
    enabled: z.boolean().meta({
        title: 'Enabled',
        markdownDescription: `If this plugin is enabled, outbound connections to specified IP addresses and ports will be blocked. **Use with caution.**${DOCS_LINK}`,
    }),
    blockedIps: z
        .array(IpCidrOrExtSchema)
        .optional()
        .meta({
            title: 'Blocked IPs',
            markdownDescription: `List of destination IP addresses and CIDR ranges to block. \n\n You can use lists from **sharedLists** in the format: **ext:list_name**. \n\n Example: \`["10.0.0.1", "ext:blocked_destinations"]\`${DOCS_LINK}`,
        }),
    blockedPorts: z
        .array(z.int().min(1).max(65535))
        .optional()
        .meta({
            title: 'Blocked Ports',
            markdownDescription: `List of destination ports to block. \n\n Example: \`[25, 465, 587]\` to block SMTP traffic.${DOCS_LINK}`,
        }),
});

const cleanupPathSchema = z
    .string()
    .trim()
    .min(1, { message: 'Path must not be empty' })
    .refine((v) => v.startsWith('/'), {
        message: 'Path must be absolute (start with "/")',
    })
    .refine((v) => !v.includes('\0'), {
        message: 'Path must not contain null bytes',
    });

export const preStartPluginSchema = z.object({
    enabled: z
        .boolean()
        .default(false)
        .meta({
            title: 'Enabled',
            markdownDescription: `Enables the pre-start stage. All enabled sections below run every time before the Xray-Core process starts — on node startup, on core restart, and after any configuration change that triggers a core reload. If a section fails, the failure is logged and the core still starts.${DOCS_LINK}`,
        }),
    cleanupSockets: z
        .object({
            enabled: z.boolean().meta({
                title: 'Enable socket cleanup',
                markdownDescription: `Removes stale unix socket files left behind by a previous core process that did not shut down cleanly. Such leftovers make Xray-Core fail to bind with \`address already in use\`. Only entries that are actually unix sockets are removed — regular files, directories and symlinks are always skipped.${DOCS_LINK}`,
            }),
            files: z
                .array(cleanupPathSchema)
                .max(64, { message: 'No more than 64 entries allowed' })
                .meta({
                    title: 'Files',
                    markdownDescription: `Absolute paths to socket files. Glob patterns are supported (\`*\`, \`?\`, \`[…]\`), for example \`/dev/shm/*.sock\`. Paths that do not exist are skipped silently.${DOCS_LINK}`,
                }),
        })
        .optional()
        .meta({
            title: 'Cleanup Sockets',
            markdownDescription: `Stale unix socket removal before the core starts.${DOCS_LINK}`,
        }),
});

export const HaproxyAuthPluginSchema = z.object({
    enabled: z
        .boolean()
        .default(false)
        .meta({
            title: 'Enabled',
            markdownDescription: `If this plugin is enabled, HAProxy user credentials table will be generated for participating inbounds.${DOCS_LINK}`,
        }),
    inboundTags: z
        .array(z.string())
        .optional()
        .default([])
        .meta({
            title: 'Inbound Tags',
            markdownDescription: `List of inbound tags that participate in HAProxy authentication. Use an empty array to disable HAProxy authentication for all inbounds. Use "*" to include every supported inbound assigned to the node.${DOCS_LINK}`,
        }),
});

export const NodePluginSchema = z.object({
    sharedLists: z
        .array(SharedListSchema)
        .optional()
        .default([])
        .meta({
            title: 'Shared Lists',
            markdownDescription: `Array of shared lists, which can be used in other plugins. Optional.${DOCS_LINK}`,
        }),
    ingressFilter: IngressFilterPluginSchema.optional().meta({
        title: 'Ingress Filter',
        markdownDescription: `Ingress Filter Plugin configuration. Optional.${DOCS_LINK}`,
    }),
    egressFilter: EgressFilterPluginSchema.optional().meta({
        title: 'Egress Filter',
        markdownDescription: `Egress Filter Plugin configuration. Optional.${DOCS_LINK}`,
    }),
    haproxyAuth: HaproxyAuthPluginSchema.optional().meta({
        title: 'HAProxy Auth',
        markdownDescription: `HAProxy Auth Plugin configuration. Optional.${DOCS_LINK}`,
    }),
    preStart: preStartPluginSchema.optional().meta({
        title: 'Pre-Start',
        markdownDescription: `Pre-Start Plugin configuration. Optional.${DOCS_LINK}`,
    }),
});

export const NodePluginEditorSchema = NodePluginSchema.omit({ sharedLists: true });

export type TSharedListConfig = z.infer<typeof SharedListConfigSchema>;
export type TNodePlugin = z.infer<typeof NodePluginSchema>;
export type TNodePluginEditor = z.infer<typeof NodePluginEditorSchema>;
