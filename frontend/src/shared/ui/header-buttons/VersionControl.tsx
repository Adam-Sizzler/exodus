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

export function VersionControl() {
    const exodusInfo = useExodusInfo()
    const { data: exodusMetadata, isLoading } = useGetExodusMetadata()

    const [isNewVersionAvailable, isDev] = useMemo(() => {
        if (!exodusMetadata) return [false, false]

        const currentVersion = semver.valid(exodusMetadata.version) || '0.0.0'
        const latest = semver.valid(exodusInfo.latestVersion || '') || '0.0.0'
        return [semver.gt(latest, currentVersion), exodusMetadata.git.backend.branch === 'dev']
    }, [exodusInfo.latestVersion, exodusMetadata])

    if (isLoading || !exodusMetadata) {
        return <SkeletonHeaderControl width={85} />
    }

    const handleClick = () => {
        modals.open({
            title: (
                <BaseOverlayHeader
                    IconComponent={Logo}
                    iconVariant="gradient-teal"
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
