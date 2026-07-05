import { Group, Text } from '@mantine/core'
import { modals } from '@mantine/modals'
import { useMemo } from 'react'
import semver from 'semver'
import clsx from 'clsx'

import { useExodusInfo } from '@entities/dashboard/updates-store'
import { useGetExodusMetadata } from '@shared/api/hooks'

import { BaseOverlayHeader } from '../overlays/base-overlay-header'
import { SkeletonHeaderControl } from './SkeletonHeaderControl'
import { BuildInfoModal } from '../sidebar/build-info-modal'
import classes from './VersionControl.module.css'
import { HeaderControl } from './HeaderControl'
import { Logo } from '../logo'

const gitDescribePattern = /^(.+)-\d+-g[0-9a-f]{7,}(?:-dirty)?$/i
const customReleasePattern = /^(\d+\.\d+\.\d+)\.([A-Za-z0-9][A-Za-z0-9.-]*)$/

function normalizeComparableVersion(value: string | undefined) {
    const raw = value?.trim()
    if (!raw || raw === 'unknown' || raw === 'latest') {
        return null
    }

    const describeMatch = raw.match(gitDescribePattern)
    const candidate = describeMatch?.[1] ?? raw
    const stripped = candidate.replace(/^[vV]/, '')

    if (semver.valid(stripped)) {
        return stripped
    }

    const customMatch = stripped.match(customReleasePattern)
    if (customMatch) {
        return `${customMatch[1]}-${customMatch[2]}`
    }

    return semver.coerce(stripped)?.version ?? null
}

export function VersionControl() {
    const exodusInfo = useExodusInfo()
    const { data: exodusMetadata, isLoading } = useGetExodusMetadata()

    const [isNewVersionAvailable, isDev] = useMemo(() => {
        if (!exodusMetadata) return [false, false]

        const currentVersion = normalizeComparableVersion(exodusMetadata.version)
        const latest = normalizeComparableVersion(exodusInfo.latestVersion)
        const isNewVersionAvailable =
            currentVersion !== null && latest !== null ? semver.gt(latest, currentVersion) : false

        return [isNewVersionAvailable, exodusMetadata.git.backend.branch === 'dev']
    }, [exodusInfo.latestVersion, exodusMetadata])

    if (isLoading || !exodusMetadata) {
        return <SkeletonHeaderControl width={85} />
    }

    const handleClick = () => {
        modals.open({
            title: (
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={Logo}
                    iconVariant="soft"
                    title="Build Info"
                />
            ),
            centered: true,
            size: 'md',
            withCloseButton: true,
            children: (
                <BuildInfoModal
                    isNewVersionAvailable={isNewVersionAvailable}
                    exodusMetadata={exodusMetadata}
                />
            )
        })
    }

    return (
        <HeaderControl
            className={clsx(classes.version, {
                [classes.newVersion]: isNewVersionAvailable && !isDev,
                [classes.dev]: isDev
            })}
            onClick={handleClick}
            w="auto"
        >
            <Group gap={8} ml={10} mr={10} wrap="nowrap">
                <Logo size={20} />
                <Text ff="text" fw={600} size="sm">
                    {exodusMetadata.version}
                </Text>
            </Group>
        </HeaderControl>
    )
}
