export type SubscriptionConnectionLike = {
    name?: string | null
    address?: string | null
    publicDomain?: string | null
    apiPath?: string | null
    isDisabled?: boolean
}

export type SubscriptionLink = {
    nodeName: string
    url: string
}

const normalizeSubscriptionDomain = (value: string): string | null => {
    const firstPart = value.split(',')[0] ?? ''
    const trimmed = firstPart.trim()
    if (trimmed === '') {
        return null
    }

    const normalized = trimmed.includes('://') ? trimmed : `https://${trimmed}`

    try {
        const parsed = new URL(normalized)
        if (parsed.host.trim() === '') {
            return null
        }

        parsed.username = ''
        parsed.password = ''
        parsed.pathname = ''
        parsed.search = ''
        parsed.hash = ''

        return parsed.toString().replace(/\/+$/, '')
    } catch {
        return null
    }
}

const normalizeSubscriptionAPIPath = (value: string): string => {
    const trimmed = value.trim()
    if (trimmed === '' || trimmed === '/') {
        return '/'
    }

    return `/${trimmed.replace(/^\/+|\/+$/g, '')}/`
}

export const buildSubscriptionLinksFromNodes = (
    nodes: SubscriptionConnectionLike[],
    shortUUID: string,
    fallbackURL?: string
): SubscriptionLink[] => {
    const cleanShortUUID = shortUUID.trim()
    const links: SubscriptionLink[] = []
    const uniqueURLs = new Set<string>()

    if (cleanShortUUID !== '') {
        nodes.forEach((node, index) => {
            if (!node || node.isDisabled) {
                return
            }

            const domain = normalizeSubscriptionDomain(node.publicDomain ?? node.address ?? '')
            if (!domain) {
                return
            }

            const path = normalizeSubscriptionAPIPath(node.apiPath ?? '/')
            const url = `${domain}${path}${cleanShortUUID}`
            if (uniqueURLs.has(url)) {
                return
            }

            uniqueURLs.add(url)
            links.push({
                nodeName:
                    node.name?.trim() ||
                    node.publicDomain?.trim() ||
                    node.address?.trim() ||
                    `Node ${index + 1}`,
                url
            })
        })
    }

    if (links.length > 0) {
        return links
    }

    const fallback = fallbackURL?.trim()
    if (!fallback) {
        return []
    }

    return [{ nodeName: 'Subscription', url: fallback }]
}
