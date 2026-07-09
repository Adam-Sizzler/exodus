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

export const HAPROXY_AUTH_ALL_INBOUNDS_TAG = '*'
export const HAPROXY_AUTH_SUPPORTED_INBOUND_TYPES = new Set(['vless', 'trojan', 'naive', 'anytls'])
 
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
