import { Group, Text } from '@mantine/core'
import { modals } from '@mantine/modals'
import clsx from 'clsx'
import { useMemo } from 'react'

import { useGetExodusMetadata } from '@shared/api/hooks'
import { isStableVersionGreater } from '@shared/utils/version-utils'

import { useExodusInfo } from '@entities/dashboard/updates-store'

import { Logo } from '../logo'
import { BaseOverlayHeader } from '../overlays/base-overlay-header'
import { BuildInfoModal } from '../sidebar/build-info-modal'
import { HeaderControl } from './HeaderControl'
import { SkeletonHeaderControl } from './SkeletonHeaderControl'
import classes from './VersionControl.module.css'

export function VersionControl() {
    const exodusInfo = useExodusInfo()
    const { data: exodusMetadata, isLoading } = useGetExodusMetadata()

    const [isNewVersionAvailable, isDev] = useMemo(() => {
        if (!exodusMetadata) return [false, false]

        const isNewVersionAvailable = isStableVersionGreater(
            exodusInfo.latestVersion,
            exodusMetadata.version
        )

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
