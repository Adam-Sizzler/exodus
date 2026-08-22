import { ActionIcon, Tooltip } from '@mantine/core'
import { GetNodeCommand } from '@exodus/backend-contract'
import { memo, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { TbMapSearch } from 'react-icons/tb'
import semver from 'semver'

import { showModal } from '@shared/_modals/show-modal'

const MIN_NODE_VERSION = '2.7.0'

interface IProps {
    node: GetNodeCommand.Response['response']
}

const GetNodeGeocheckFeatureComponent = (props: IProps) => {
    const { node } = props
    const { t } = useTranslation()

    const nodeVersion = node.versions?.node || (node as any).nodeVersion

    const isSupported = useMemo(() => {
        if (!nodeVersion || nodeVersion === 'dev' || nodeVersion === 'unknown') return true
        const version = semver.coerce(nodeVersion)

        return version !== null ? semver.gte(version, MIN_NODE_VERSION) : true
    }, [nodeVersion])

    return (
        <Tooltip
            label={
                isSupported
                    ? t('node-geocheck.title')
                    : t('node-geocheck.requires-node-version', { version: MIN_NODE_VERSION })
            }
        >
            <ActionIcon
                disabled={!isSupported}
                color="indigo"
                onClick={() => {
                    showModal('nodes_nodeGeocheckModal', { node })
                }}
                size="lg"
                variant="soft"
            >
                <TbMapSearch size="22px" />
            </ActionIcon>
        </Tooltip>
    )
}

export const GetNodeGeocheckFeature = memo(GetNodeGeocheckFeatureComponent)
