const gitDescribePattern = /^(.+)-\d+-g[0-9a-f]{7,}(?:-dirty)?$/i
const comparableVersionPattern = /^[vV]?(\d+\.\d+\.\d+(?:\.\d+)?)(?:[-+._]?[0-9A-Za-z][0-9A-Za-z.-]*)?$/

export function normalizeStableVersion(value: string | undefined) {
    const raw = value?.trim()
    if (!raw || raw === 'unknown' || raw === 'latest') {
        return null
    }

    const describeMatch = raw.match(gitDescribePattern)
    const candidate = describeMatch?.[1] ?? raw
    const versionMatch = candidate.match(comparableVersionPattern)

    if (!versionMatch) {
        return null
    }

    return versionMatch[1]
}

export function compareStableVersions(first: string, second: string) {
    const firstParts = first.split('.').map(Number)
    const secondParts = second.split('.').map(Number)
    const length = Math.max(firstParts.length, secondParts.length)

    for (let index = 0; index < length; index += 1) {
        const firstPart = firstParts[index] ?? 0
        const secondPart = secondParts[index] ?? 0

        if (firstPart > secondPart) return 1
        if (firstPart < secondPart) return -1
    }

    return 0
}

export function isStableVersionGreater(first: string | undefined, second: string | undefined) {
    const normalizedFirst = normalizeStableVersion(first)
    const normalizedSecond = normalizeStableVersion(second)

    if (!normalizedFirst || !normalizedSecond) {
        return false
    }

    return compareStableVersions(normalizedFirst, normalizedSecond) > 0
}

export function findLatestStableVersionTag(tags: string[]) {
    return tags
        .map(normalizeStableVersion)
        .filter((version): version is string => Boolean(version))
        .sort(compareStableVersions)
        .at(-1)
}
