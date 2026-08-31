import { ActionIcon, ActionIconGroup, Box, Flex, Group, Tooltip } from '@mantine/core'
import { GetNodePluginCommand } from '@exodus/backend-contract'
import { NodePluginEditorWidget } from '@widgets/dashboard/node-plugins/node-plugin-editor'
import { useTranslation } from 'react-i18next'
import { TbArrowBackUp, TbBook, TbList, TbPackage } from 'react-icons/tb'
import { useNavigate } from 'react-router'

import { showModal } from '@shared/_modals/show-modal'
import { OPEN_ENTITY, ROUTES } from '@shared/constants'
import { CopyEntityLinkButton, Page, PageHeaderShared } from '@shared/ui'

interface Props {
    plugin: GetNodePluginCommand.Response['response']
}

export const NodePluginEditorPageComponent = (props: Props) => {
    const { plugin } = props
    const { t } = useTranslation()

    const navigate = useNavigate()

    return (
        <Page title={plugin.name}>
            <PageHeaderShared
                actions={
                    <Group>
                        <CopyEntityLinkButton
                            entity={OPEN_ENTITY.NODE_PLUGIN}
                            iconSize={24}
                            id={plugin.uuid}
                            size="input-md"
                            variant="soft"
                        />

                        <Tooltip label={t('help-action-icon.shared.help-article')} withArrow>
                            <ActionIcon
                                color="lime"
                                onClick={() => showModal('helpDrawer', { screen: 'PAGE_NODE_PLUGINS' })}
                                size="input-md"
                                variant="soft"
                            >
                                <TbBook size={24} />
                            </ActionIcon>
                        </Tooltip>

                        <ActionIconGroup>
                            <Tooltip label={t('common.field.shared-lists')} withArrow>
                                <ActionIcon
                                    color="indigo"
                                    onClick={() => showModal('sharedLists_sharedListsModal')}
                                    size="input-md"
                                    variant="soft"
                                >
                                    <TbList size="24px" />
                                </ActionIcon>
                            </Tooltip>
                        </ActionIconGroup>

                        <ActionIcon
                            color="gray"
                            onClick={() => navigate(ROUTES.DASHBOARD.MANAGEMENT.NODE_PLUGINS.ROOT)}
                            size="input-md"
                            variant="soft"
                        >
                            <TbArrowBackUp size={24} />
                        </ActionIcon>
                    </Group>
                }
                description={plugin.uuid}
                icon={<TbPackage size={24} />}
                title={plugin.name}
            />
            <Flex gap="md">
                <Box style={{ flex: 1, minWidth: 0 }}>
                    <NodePluginEditorWidget
                        nodePlugin={plugin.pluginConfig}
                        pluginUuid={plugin.uuid}
                    />
                </Box>
            </Flex>
        </Page>
    )
}
