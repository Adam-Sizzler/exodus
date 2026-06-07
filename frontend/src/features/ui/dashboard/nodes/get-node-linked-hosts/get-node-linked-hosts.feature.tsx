import { useTranslation } from 'react-i18next'
import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import { TbServerCog } from 'react-icons/tb'
import { memo } from 'react'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

import { IProps } from './interfaces'

const GetNodeLinkedHostsFeatureComponent = (props: IProps) => {
    const { nodeUuid, renderAs = 'menu' } = props
    const { t } = useTranslation()

    const openModalWithData = useModalsStoreOpenWithData()
    const label = t('get-node-linked-hosts.feature.linked-hosts')
    const handleOpen = () => {
        openModalWithData(MODALS.SHOW_NODE_LINKED_HOSTS_DRAWER, {
            nodeUuid
        })
    }

    if (renderAs === 'action') {
        return (
            <Tooltip label={label} withArrow>
                <ActionIcon color="cyan" onClick={handleOpen} size="md" variant="soft">
                    <TbServerCog size="16px" />
                </ActionIcon>
            </Tooltip>
        )
    }

    return (
        <Menu.Item leftSection={<TbServerCog size="16px" />} onClick={handleOpen}>
            {label}
        </Menu.Item>
    )
}

export const GetNodeLinkedHostsFeature = memo(GetNodeLinkedHostsFeatureComponent)
