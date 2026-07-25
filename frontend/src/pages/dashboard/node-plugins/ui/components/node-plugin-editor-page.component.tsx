import { ActionIcon, Box, Flex, Group } from '@mantine/core'
import { GetNodePluginCommand } from '@exodus/backend-contract'
import { NodePluginEditorWidget } from '@widgets/dashboard/node-plugins/node-plugin-editor'
import { TbArrowBackUp, TbBook, TbPackage } from 'react-icons/tb'
import { useNavigate } from 'react-router'

import { ROUTES } from '@shared/constants'
import { Page, PageHeaderShared } from '@shared/ui'
import { HelpActionIconShared } from '@shared/ui/help-drawer'

interface Props {
    plugin: GetNodePluginCommand.Response['response']
}

export const NodePluginEditorPageComponent = (props: Props) => {
    const { plugin } = props

    const navigate = useNavigate()

    return (
        <Page title={plugin.name}>
            <PageHeaderShared
                actions={
                    <Group>
                        <HelpActionIconShared hidden={false} screen="PAGE_NODE_PLUGINS" />

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
