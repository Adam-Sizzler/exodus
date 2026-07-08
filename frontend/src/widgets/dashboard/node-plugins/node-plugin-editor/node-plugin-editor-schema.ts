import { NodePluginSchema } from '@exodus/node-plugins'
import zodToJsonSchema, { jsonDescription } from 'zod-to-json-schema'

type JsonSchemaNode = Record<string, unknown>

type NodePluginEditorNode = {
    activePluginUuid: string | null
    name: string
    configProfile?: {
        activeInbounds?: ReadonlyArray<{
            tag: string
            type: string
        }>
    }
}

export type NodePluginHaproxyInboundTagOption = {
    nodeNames: string[]
    tag: string
    type: string | null
}

const HAPROXY_AUTH_ALL_INBOUNDS_TAG = '*'
const HAPROXY_AUTH_SUPPORTED_INBOUND_TYPES = new Set(['vless', 'trojan', 'naive', 'anytls'])

export const getNodePluginHaproxyInboundTagOptions = (
    nodes: ReadonlyArray<NodePluginEditorNode> | undefined,
    pluginUuid: string
): NodePluginHaproxyInboundTagOption[] | undefined => {
    if (!nodes) return undefined

    const tags = new Map<string, { nodeNames: Set<string>; tag: string; type: string | null }>()

    nodes
        .filter((node) => node.activePluginUuid === pluginUuid)
        .forEach((node) => {
            node.configProfile?.activeInbounds?.forEach((inbound) => {
                const type = inbound.type.trim().toLowerCase()
                const tag = inbound.tag.trim()

                if (!tag || !HAPROXY_AUTH_SUPPORTED_INBOUND_TYPES.has(type)) return

                const existing = tags.get(tag) ?? {
                    nodeNames: new Set<string>(),
                    tag,
                    type
                }
                existing.nodeNames.add(node.name)
                tags.set(tag, existing)
            })
        })

    return Array.from(tags.values())
        .map((item) => ({
            nodeNames: Array.from(item.nodeNames).sort((a, b) => a.localeCompare(b)),
            tag: item.tag,
            type: item.type
        }))
        .sort((a, b) => a.tag.localeCompare(b.tag))
}

export const buildNodePluginEditorSchema = (
    haproxyInboundTagOptions?: NodePluginHaproxyInboundTagOption[]
): JsonSchemaNode => {
    const schema = zodToJsonSchema(NodePluginSchema, {
        name: 'Node Plugin Schema',
        applyRegexFlags: true,
        errorMessages: true,
        postProcess: jsonDescription
    }) as JsonSchemaNode

    if (haproxyInboundTagOptions) {
        patchNodePluginHaproxyAuthSchema(schema, haproxyInboundTagOptions)
    }

    return schema
}

const getNodePluginHaproxyAuthSchema = (
    inboundTagOptions: NodePluginHaproxyInboundTagOption[]
): JsonSchemaNode => {
    const tags = Array.from(
        new Set(
            inboundTagOptions
                .map((option) => option.tag.trim())
                .filter((tag) => tag.length > 0 && tag !== HAPROXY_AUTH_ALL_INBOUNDS_TAG)
        )
    ).sort((a, b) => a.localeCompare(b))

    const descriptionsByTag = inboundTagOptions.reduce<Record<string, string>>((acc, option) => {
        const tag = option.tag.trim()
        if (!tag) return acc

        const nodeNames = option.nodeNames.filter(Boolean).join(', ')
        const metadata = [option.type ? `type: ${option.type}` : null, nodeNames ? `nodes: ${nodeNames}` : null]
            .filter(Boolean)
            .join('; ')

        acc[tag] = metadata ? `Inbound tag \`${tag}\` (${metadata}).` : `Inbound tag \`${tag}\`.`
        return acc
    }, {})

    const enumValues = [HAPROXY_AUTH_ALL_INBOUNDS_TAG, ...tags]

    return {
        type: 'object',
        additionalProperties: false,
        markdownDescription:
            'HAProxy Auth Plugin configuration. Optional. Use inboundTags to explicitly select which inbound tags participate in HAProxy authentication.',
        properties: {
            inboundTags: {
                type: 'array',
                default: [],
                uniqueItems: true,
                markdownDescription:
                    'List of inbound tags that participate in HAProxy authentication. Empty array disables HAProxy authentication. Use "*" to include every supported inbound assigned to the node.',
                items: {
                    type: 'string',
                    enum: enumValues,
                    markdownEnumDescriptions: [
                        'All supported inbounds assigned to the node.',
                        ...tags.map((tag) => descriptionsByTag[tag] ?? `Inbound tag \`${tag}\`.`)
                    ]
                }
            }
        }
    }
}

const patchNodePluginHaproxyAuthSchema = (
    schema: JsonSchemaNode,
    inboundTagOptions: NodePluginHaproxyInboundTagOption[]
) => {
    const haproxyAuthSchema = getNodePluginHaproxyAuthSchema(inboundTagOptions)

    const patchNode = (node: unknown, visited = new WeakSet<object>()) => {
        if (!node || typeof node !== 'object') return
        if (visited.has(node)) return
        visited.add(node)

        const current = node as JsonSchemaNode
        const properties = current.properties

        if (properties && typeof properties === 'object') {
            const typedProperties = properties as Record<string, unknown>
            if (Object.prototype.hasOwnProperty.call(typedProperties, 'haproxyAuth')) {
                typedProperties.haproxyAuth = haproxyAuthSchema
            }
        }

        for (const value of Object.values(current)) {
            if (value && typeof value === 'object') {
                patchNode(value, visited)
            }
        }
    }

    patchNode(schema)
}
