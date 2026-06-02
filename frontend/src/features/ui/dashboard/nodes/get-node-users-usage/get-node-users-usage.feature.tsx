import { useTranslation } from 'react-i18next'
import { ActionIcon, Menu, Tooltip } from '@mantine/core'
import { PiChartBarDuotone } from 'react-icons/pi'
import { memo } from 'react'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

import { IProps } from './interfaces'

const GetNodeUsersUsageFeatureComponent = (props: IProps) => {
    const { nodeUuid, renderAs = 'menu' } = props
    const { t } = useTranslation()

    const openModalWithData = useModalsStoreOpenWithData()
    const label = t('get-user-usage.feature.show-usage')
    const handleOpen = () => {
        openModalWithData(MODALS.SHOW_NODE_USERS_USAGE_DRAWER, {
            nodeUuid
        })
    }

    if (renderAs === 'action') {
        return (
            <Tooltip label={label} withArrow>
                <ActionIcon color="grape" onClick={handleOpen} size="md" variant="light">
                    <PiChartBarDuotone size="16px" />
                </ActionIcon>
            </Tooltip>
        )
    }

    return (
        <Menu.Item color="grape" leftSection={<PiChartBarDuotone size="16px" />} onClick={handleOpen}>
            {label}
        </Menu.Item>
    )
}

export const GetNodeUsersUsageFeature = memo(GetNodeUsersUsageFeatureComponent)
