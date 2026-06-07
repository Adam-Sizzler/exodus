import { PiCheck, PiCopy, PiCpu, PiPencil, PiTrashDuotone } from 'react-icons/pi'
import { TbCopyCheck, TbEdit, TbPackage } from 'react-icons/tb'
import { generatePath, useNavigate } from 'react-router-dom'
import { CopyButton, Menu } from '@mantine/core'
import { useTranslation } from 'react-i18next'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'
import { WithDndSortable } from '@shared/hocs/with-dnd-sortable'
import { EntityCardShared } from '@shared/ui/entity-card'
import { NodePluginResponse } from '@shared/api/hooks'
import { ROUTES } from '@shared/constants'

interface IProps {
    handleCloneNodePlugin: (nodePluginUuid: string) => void
    handleDeleteNodePlugin: (nodePluginUuid: string) => void
    handleShowActiveNodes: (nodePluginUuid: string) => void
    isDragOverlay?: boolean
    nodePlugin: NodePluginResponse
}

export function NodePluginCardWidget(props: IProps) {
    const {
        handleCloneNodePlugin,
        handleDeleteNodePlugin,
        handleShowActiveNodes,
        isDragOverlay = false,
        nodePlugin
    } = props

    const { t } = useTranslation()
    const navigate = useNavigate()
    const openModalWithData = useModalsStoreOpenWithData()

    const navigateToNodePlugin = () => {
        navigate(
            generatePath(ROUTES.DASHBOARD.MANAGEMENT.NODE_PLUGINS.NODE_PLUGIN_BY_UUID, {
                uuid: nodePlugin.uuid
            })
        )
    }

    return (
        <WithDndSortable
            dragHandlePosition="top-right"
            id={nodePlugin.uuid}
            isDragOverlay={isDragOverlay}
        >
            <EntityCardShared.Root>
                <EntityCardShared.Header>
                    <EntityCardShared.Icon highlight={false} onClick={navigateToNodePlugin}>
                        <TbPackage size={24} />
                    </EntityCardShared.Icon>

                    <EntityCardShared.Content
                        subtitle={t('node-plugin-card.widget.plugin')}
                        title={nodePlugin.name}
                    />
                </EntityCardShared.Header>

                <EntityCardShared.Actions>
                    <EntityCardShared.Button
                        leftSection={<TbEdit size={16} />}
                        onClick={navigateToNodePlugin}
                    >
                        {t('common.edit')}
                    </EntityCardShared.Button>

                    <EntityCardShared.Menu>
                        <CopyButton timeout={2000} value={nodePlugin.uuid}>
                            {({ copied, copy }) => (
                                <Menu.Item
                                    color={copied ? 'teal' : undefined}
                                    leftSection={
                                        copied ? <PiCheck size={18} /> : <PiCopy size={18} />
                                    }
                                    onClick={copy}
                                >
                                    {t('common.copy-uuid')}
                                </Menu.Item>
                            )}
                        </CopyButton>

                        <Menu.Item
                            leftSection={<PiCpu size={18} />}
                            onClick={() => handleShowActiveNodes(nodePlugin.uuid)}
                        >
                            {t('node-plugin-card.widget.active-on-nodes')}
                        </Menu.Item>

                        <Menu.Item
                            leftSection={<PiPencil size={18} />}
                            onClick={() => {
                                openModalWithData(MODALS.RENAME_SQUAD_OR_CONFIG_PROFILE_MODAL, {
                                    name: nodePlugin.name,
                                    uuid: nodePlugin.uuid
                                })
                            }}
                        >
                            {t('common.rename')}
                        </Menu.Item>

                        <Menu.Item
                            leftSection={<TbCopyCheck size={18} />}
                            onClick={() => handleCloneNodePlugin(nodePlugin.uuid)}
                        >
                            {t('common.clone')}
                        </Menu.Item>

                        <Menu.Item
                            color="red"
                            leftSection={<PiTrashDuotone size={18} />}
                            onClick={(event) => {
                                event.stopPropagation()
                                handleDeleteNodePlugin(nodePlugin.uuid)
                            }}
                        >
                            {t('common.delete')}
                        </Menu.Item>
                    </EntityCardShared.Menu>
                </EntityCardShared.Actions>
            </EntityCardShared.Root>
        </WithDndSortable>
    )
}
