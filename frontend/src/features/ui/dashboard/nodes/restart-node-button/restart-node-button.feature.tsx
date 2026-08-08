import { Menu } from '@mantine/core'
import { modals } from '@mantine/modals'
import { GetNodeCommand } from '@exodus/backend-contract'
import { useTranslation } from 'react-i18next'
import { TbReload, TbRocket } from 'react-icons/tb'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { RestartNodeModalContentFeature } from './restart-node.modal-content.feature'

interface IProps {
    handleClose: () => void
    node: GetNodeCommand.Response['response']
}

export function RestartNodeButtonFeature(props: IProps) {
    const { t } = useTranslation()

    const { handleClose, node } = props

    if (!node) return null

    return (
        <Menu.Item
            color="teal"
            disabled={node.isDisabled}
            leftSection={<TbReload size="1rem" />}
            onClick={() => {
                modals.open({
                    title: (
                        <BaseOverlayHeader
                            iconColor="teal"
                            IconComponent={TbRocket}
                            iconVariant="soft"
                            title={t('restart-node-button.feature.restart')}
                        />
                    ),
                    centered: true,
                    size: 'md',
                    children: (
                        <RestartNodeModalContentFeature
                            handleClose={handleClose}
                            nodeUuid={node.uuid}
                        />
                    )
                })
            }}
        >
            {t('restart-node-button.feature.restart')}
        </Menu.Item>
    )
}
