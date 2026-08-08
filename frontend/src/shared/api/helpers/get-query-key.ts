import { QueryKey } from '../types'

export function getQueryKey(
    queryKey: QueryKey,
    route: Record<string, unknown> = {},
    query: Record<string, unknown> = {}
) {
    const [mainKey, otherKeys = {}] = queryKey
    const objectKeys = typeof otherKeys === 'object' ? otherKeys : {}

    return [mainKey, { ...objectKeys, ...route, ...query }]
}
