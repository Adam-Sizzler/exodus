"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.HostMapperSchema = exports.HostMapperOperationsSchema = exports.Base64HostMapperOperationsSchema = exports.SingBoxHostMapperOperationsSchema = exports.MihomoHostMapperOperationsSchema = exports.XrayJsonHostMapperOperationsSchema = void 0;
const zod_1 = require("zod");
const buildFromSchema = () => zod_1.z
    .string()
    .min(1)
    .max(512)
    .meta({
    title: 'Source path',
    markdownDescription: [
        'Dot-separated path inside the **raw inbound** of the config profile this host belongs to. Array items are addressed by index.',
        '',
        'With a `$host.` prefix the value is read from the **host** itself, after overrides have been resolved – `$host.securityOptions.serverName` holds the SNI that actually went into the config.',
        '',
        'Only a part of the host is readable. Each entry below also opens everything nested under it:',
        '',
        '`address`, `port`, `finalRemark`, `protocol`, `transport`, `security`, `protocolOptions`, `transportOptions`, `securityOptions`, `mux`, `streamOverrides.finalMask`, `streamOverrides.sockopt`, `clientOverrides.serverDescription`, `metadata.remark`, `metadata.tags`, `metadata.inboundTag`',
        '',
        '> If the path resolves to nothing – or points outside the list above – the operation is skipped and the target key is **not** created.',
    ].join('\n'),
    examples: [
        'streamSettings.tlsSettings.cipherSuites',
        'streamSettings.tlsSettings.alpn',
        'streamSettings.realitySettings.serverNames.0',
        '$host.address',
        '$host.securityOptions.serverName',
        '$host.transportOptions.path',
        '$host.mux.smux',
    ],
});
const buildToSchema = (client) => zod_1.z
    .string()
    .min(1)
    .max(512)
    .meta({
    title: 'Target path',
    markdownDescription: [
        `Dot-separated path inside ${client.target}, counted from its root.`,
        '',
        'Missing objects along the path are created. Array items are addressed by index, and keys with dashes are written as is (`ip-version`).',
        '',
        '> A path running through a value that is not an object – for example `path.foo`, where `path` holds a string – is refused, so an existing value is never destroyed.',
    ].join('\n'),
    examples: client.examples,
});
const buildOperationsSchema = (client) => zod_1.z
    .discriminatedUnion('op', [
    zod_1.z
        .object({
        op: zod_1.z.literal('copy').meta({
            title: 'Copy',
            markdownDescription: 'Take a value from the raw inbound.',
        }),
        from: buildFromSchema(),
        to: buildToSchema(client),
    })
        .meta({
        title: 'Copy from the inbound',
        markdownDescription: [
            'Reads a value from the raw inbound and writes it into the generated config.',
            '',
            '```json',
            JSON.stringify({ op: 'copy', ...client.copySnippet }, null, 2),
            '```',
        ].join('\n'),
        defaultSnippets: [
            {
                label: 'copy',
                description: 'Copy a value from the raw inbound',
                body: { op: 'copy', from: '$1', to: '$2' },
            },
        ],
    }),
    zod_1.z
        .object({
        op: zod_1.z.literal('set').meta({
            title: 'Set',
            markdownDescription: 'Write a fixed value.',
        }),
        value: zod_1.z
            .union([
            zod_1.z.string(),
            zod_1.z.number(),
            zod_1.z.boolean(),
            zod_1.z.array(zod_1.z.json()),
            zod_1.z.record(zod_1.z.string(), zod_1.z.json()),
        ])
            .meta({
            title: 'Value',
            markdownDescription: 'The value to write. Strings, numbers, booleans, arrays and objects are allowed.',
        }),
        to: buildToSchema(client),
    })
        .meta({
        title: 'Set a fixed value',
        markdownDescription: [
            'Writes a literal value into the generated config, creating the key if it is missing.',
            '',
            '```json',
            JSON.stringify({ op: 'set', ...client.setSnippet }, null, 2),
            '```',
        ].join('\n'),
        defaultSnippets: [
            {
                label: 'set',
                description: 'Write a fixed value',
                body: { op: 'set', to: '$1', value: '$2' },
            },
        ],
    }),
    zod_1.z
        .object({
        op: zod_1.z.literal('unset').meta({
            title: 'Unset',
            markdownDescription: 'Remove a field.',
        }),
        to: buildToSchema(client),
    })
        .meta({
        title: 'Remove a field',
        markdownDescription: [
            'Removes a field the generator produced by itself.',
            '',
            '```json',
            JSON.stringify({ op: 'unset', ...client.unsetSnippet }, null, 2),
            '```',
            '',
            '> Only the field itself is removed. A parent object left empty stays in the config as `{}`.',
        ].join('\n'),
        defaultSnippets: [
            {
                label: 'unset',
                description: 'Remove a field',
                body: { op: 'unset', to: '$1' },
            },
        ],
    }),
])
    .meta({
    title: 'Operation',
    markdownDescription: [
        'Operations run one after another, in the order they are listed.',
        '',
        'A later operation sees what the previous ones wrote, so an `unset` can remove a key an earlier `set` added.',
    ].join('\n'),
});
exports.XrayJsonHostMapperOperationsSchema = buildOperationsSchema({
    target: 'the generated **outbound**',
    examples: [
        'streamSettings.tlsSettings.enableSessionResumption',
        'streamSettings.tlsSettings.cipherSuites',
    ],
    copySnippet: {
        from: 'streamSettings.tlsSettings.cipherSuites',
        to: 'streamSettings.tlsSettings.cipherSuites',
    },
    setSnippet: { to: 'streamSettings.tlsSettings.enableSessionResumption', value: true },
    unsetSnippet: { to: 'mux' },
});
exports.MihomoHostMapperOperationsSchema = buildOperationsSchema({
    target: 'the generated **proxy node**',
    examples: ['ip-version', 'client-fingerprint', 'tfo', 'reality-opts.support-x25519mlkem768'],
    copySnippet: { from: 'streamSettings.realitySettings.serverNames.0', to: 'servername' },
    setSnippet: { to: 'reality-opts.support-x25519mlkem768', value: true },
    unsetSnippet: { to: 'smux' },
});
exports.SingBoxHostMapperOperationsSchema = buildOperationsSchema({
    target: 'the generated **outbound**',
    examples: [
        'domain_resolver',
        'multiplex.protocol',
        'packet_encoding',
        'tcp_fast_open',
        'tls.insecure',
        'tls.utls.fingerprint',
    ],
    copySnippet: { from: 'streamSettings.realitySettings.serverNames.0', to: 'tls.server_name' },
    setSnippet: { to: 'tls.utls.fingerprint', value: 'chrome' },
    unsetSnippet: { to: 'multiplex' },
});
exports.Base64HostMapperOperationsSchema = buildOperationsSchema({
    target: 'the query string of the generated **share link**',
    examples: [
        'alpn',
        'authority',
        'cs',
        'encryption',
        'extra',
        'flow',
        'fm',
        'fp',
        'headerType',
        'heartbeatPeriod',
        'host',
        'mode',
        'mtu',
        'obfs',
        'obfs-password',
        'path',
        'pbk',
        'pcs',
        'pinSHA256',
        'pqv',
        'security',
        'serviceName',
        'sid',
        'sni',
        'spx',
        'tti',
        'type',
        'vcn',
    ],
    copySnippet: { from: 'streamSettings.realitySettings.serverNames.0', to: 'sni' },
    setSnippet: { to: 'fp', value: 'chrome' },
    unsetSnippet: { to: 'fm' },
});
exports.HostMapperOperationsSchema = exports.XrayJsonHostMapperOperationsSchema;
exports.HostMapperSchema = zod_1.z
    .object({
    xrayJson: zod_1.z
        .array(exports.XrayJsonHostMapperOperationsSchema)
        .optional()
        .meta({
        title: 'Xray JSON',
        markdownDescription: [
            'Operations applied to the **outbound** generated for this host in the Xray JSON subscription. Paths are counted from the root of the outbound.',
            '',
            '```json',
            JSON.stringify([
                {
                    op: 'copy',
                    from: 'streamSettings.tlsSettings.cipherSuites',
                    to: 'streamSettings.tlsSettings.cipherSuites',
                },
                { op: 'unset', to: 'mux' },
            ], null, 2),
            '```',
        ].join('\n'),
    }),
    mihomo: zod_1.z
        .array(exports.MihomoHostMapperOperationsSchema)
        .optional()
        .meta({
        title: 'Mihomo',
        markdownDescription: [
            'Operations applied to the **proxy node** generated for this host in the Mihomo subscription.',
            '',
            'Paths are counted from the root of the node, where Mihomo keys are written in kebab-case.',
            '',
            '```json',
            JSON.stringify([{ op: 'set', to: 'ip-version', value: 'ipv4' }], null, 2),
            '```',
        ].join('\n'),
    }),
    base64: zod_1.z
        .array(exports.Base64HostMapperOperationsSchema)
        .optional()
        .meta({
        title: 'Base64',
        markdownDescription: [
            'Operations applied to the **query string** of the share link generated for this host in the Base64 subscription.',
            '',
            '`to` is a plain query parameter name, not a path – dots carry no meaning here.',
            '',
            '```json',
            JSON.stringify([
                { op: 'set', to: 'fp', value: 'chrome' },
                { op: 'unset', to: 'fm' },
            ], null, 2),
            '```',
        ].join('\n'),
    }),
    singbox: zod_1.z
        .array(exports.SingBoxHostMapperOperationsSchema)
        .optional()
        .meta({
        title: 'sing-box',
        markdownDescription: [
            'Operations applied to the **outbound** generated for this host in the sing-box subscription. Paths are counted from the root of the outbound.',
            '',
            'sing-box keys are written in snake_case, and the outbound of a Hysteria2 host has a shape of its own.',
            '',
            '```json',
            JSON.stringify([
                { op: 'set', to: 'tls.utls.fingerprint', value: 'chrome' },
                { op: 'unset', to: 'multiplex' },
            ], null, 2),
            '```',
        ].join('\n'),
    }),
})
    .meta({
    title: 'Host Mapper',
    markdownDescription: [
        'Rewrites the config generated for this host, per client type.',
        '',
        'Operations run **after** the generator has finished, so they can change or remove anything it produced.',
        '',
        'The source for `copy` is the raw inbound of the config profile this host belongs to, or the host itself when the path starts with `$host.`.',
        '',
        '> `to` is never checked against the target client. A misspelled key is written exactly like a real one.',
    ].join('\n'),
});
