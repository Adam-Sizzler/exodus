import { GetConfigProfilesCommand } from '@exodus/backend-contract'

export type NodePluginHaproxyInboundTagOption = {
    nodeNames: string[]
    tag: string
    type: string | null
}

export const HAPROXY_AUTH_ALL_INBOUNDS_TAG = '*'
export const HAPROXY_AUTH_SUPPORTED_INBOUND_TYPES = new Set(['vless', 'trojan', 'naive', 'anytls'])
 
export const getNodePluginHaproxyInboundTagOptions = (
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'] | undefined
): NodePluginHaproxyInboundTagOption[] | undefined => {
    if (!configProfiles) return undefined

    const tags = new Map<string, { nodeNames: Set<string>; tag: string; type: string | null }>()

    configProfiles.forEach((profile) => {
        profile.inbounds.forEach((inbound) => {
            const type = inbound.type.trim().toLowerCase()
            const tag = inbound.tag.trim()

            if (!tag || !HAPROXY_AUTH_SUPPORTED_INBOUND_TYPES.has(type)) return

            const existing = tags.get(tag) ?? {
                nodeNames: new Set<string>(),
                tag,
                type
            }
            if (profile.nodes && profile.nodes.length > 0) {
                profile.nodes.forEach((node) => {
                    existing.nodeNames.add(node.name)
                })
            }
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
