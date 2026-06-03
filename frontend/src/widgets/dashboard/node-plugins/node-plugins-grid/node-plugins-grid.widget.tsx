import { TbAlertTriangle, TbLogin, TbLogout, TbPackage, TbShieldLock } from 'react-icons/tb'
import { Badge, Center, Group, Stack, Text, ThemeIcon } from '@mantine/core'
import { modals } from '@mantine/modals'

import {
    NodePluginResponse,
    NodeResponse,
    QueryKeys,
    useCloneNodePlugin,
    useDeleteNodePlugin,
    useReorderNodePlugins
} from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { VirtualizedDndGrid } from '@shared/ui/virtualized-dnd-grid'
import { queryClient } from '@shared/api/query-client'
import { SectionCard } from '@shared/ui/section-card'

import { ActivePluginsOnNodesModalShared } from '../active-on-nodes-modal/adtive-on-nodes.modal.shared'
import { NodePluginCardWidget } from '../node-plugin-card/node-plugin-card.widget'

interface IProps {
    nodes: NodeResponse[]
    plugins: NodePluginResponse[]
}

export function NodePluginsGridWidget(props: IProps) {
    const { nodes, plugins } = props

    const invalidate = async () => {
        await queryClient.invalidateQueries({
            queryKey: QueryKeys.nodes.getNodePlugins.queryKey
        })
        await queryClient.invalidateQueries({
            queryKey: QueryKeys.nodes.getAllNodes.queryKey
        })
    }

    const { mutate: deleteNodePlugin } = useDeleteNodePlugin({
        mutationFns: {
            onSuccess: invalidate
        }
    })
    const { mutate: cloneNodePlugin } = useCloneNodePlugin({
        mutationFns: {
            onSuccess: invalidate
        }
    })
    const { mutate: reorderNodePlugins } = useReorderNodePlugins({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(QueryKeys.nodes.getNodePlugins.queryKey, data)
            }
        }
    })

    const handleDeleteNodePlugin = (nodePluginUuid: string) => {
        modals.openConfirmModal({
            centered: true,
            title: 'Confirm action',
            children: 'Delete this node plugin? Nodes using it will be detached from the plugin.',
            labels: {
                confirm: 'Delete',
                cancel: 'Cancel'
            },
            confirmProps: { color: 'red' },
            onConfirm: () => {
                deleteNodePlugin({
                    route: {
                        uuid: nodePluginUuid
                    }
                })
            }
        })
    }

    const handleCloneNodePlugin = (nodePluginUuid: string) => {
        cloneNodePlugin({
            variables: {
                cloneFromUuid: nodePluginUuid
            }
        })
    }

    const handleShowActiveNodes = (nodePluginUuid: string) => {
        const activeOnNodes = nodes.filter((node) => node.activePluginUuid === nodePluginUuid)

        modals.open({
            centered: true,
            size: 'lg',
            title: (
                <BaseOverlayHeader
                    IconComponent={TbPackage}
                    iconVariant="soft"
                    iconColor="teal"
                    title="Active on nodes"
                />
            ),
            children: <ActivePluginsOnNodesModalShared nodes={activeOnNodes} />
        })
    }

    const handleReorder = (reorderedItems: typeof plugins) => {
        reorderNodePlugins({
            variables: {
                items: reorderedItems.map((plugin, position) => ({
                    uuid: plugin.uuid,
                    viewPosition: position
                }))
            }
        })
    }

    if (!plugins.length) {
        return (
            <SectionCard.Root p="xl">
                <SectionCard.Section>
                    <BaseOverlayHeader
                        themeIconProps={{ color: 'orange' }}
                        IconComponent={TbAlertTriangle}
                        iconVariant="light"
                        subtitle="Node Plugins are an advanced feature. Please review the documentation before use."
                        title="Warning"
                        titleOrder={4}
                    />
                </SectionCard.Section>

                <SectionCard.Section>
                    <Center py="xl">
                        <Stack align="center" gap="lg">
                            <ThemeIcon color="gray" radius="xl" size={64} variant="light">
                                <TbPackage size={32} />
                            </ThemeIcon>

                            <Stack align="center" gap="xs">
                                <Text fw={600} size="lg" ta="center">
                                    No node plugins yet
                                </Text>
                                <Text c="dimmed" maw={460} size="sm" ta="center">
                                    Create a plugin to enable per-node capabilities. This build
                                    supports ingress filter, egress filter, shared lists and HAProxy
                                    Auth.
                                </Text>
                            </Stack>

                            <Group gap="sm" justify="center">
                                <Badge
                                    color="teal"
                                    leftSection={<TbLogin size={16} />}
                                    radius="md"
                                    size="lg"
                                    variant="light"
                                >
                                    Ingress Filter
                                </Badge>
                                <Badge
                                    color="orange"
                                    leftSection={<TbLogout size={16} />}
                                    radius="md"
                                    size="lg"
                                    variant="light"
                                >
                                    Egress Filter
                                </Badge>
                                <Badge
                                    color="teal"
                                    leftSection={<TbShieldLock size={16} />}
                                    radius="md"
                                    size="lg"
                                    variant="light"
                                >
                                    HAProxy Auth
                                </Badge>
                            </Group>
                        </Stack>
                    </Center>
                </SectionCard.Section>
            </SectionCard.Root>
        )
    }

    return (
        <VirtualizedDndGrid
            enableDnd={true}
            items={plugins}
            key="node-plugins-grid-widget"
            onReorder={handleReorder}
            renderDragOverlay={(nodePlugin) => (
                <NodePluginCardWidget
                    handleCloneNodePlugin={handleCloneNodePlugin}
                    handleDeleteNodePlugin={handleDeleteNodePlugin}
                    handleShowActiveNodes={handleShowActiveNodes}
                    isDragOverlay
                    nodePlugin={nodePlugin}
                />
            )}
            renderItem={(nodePlugin) => (
                <NodePluginCardWidget
                    handleCloneNodePlugin={handleCloneNodePlugin}
                    handleDeleteNodePlugin={handleDeleteNodePlugin}
                    handleShowActiveNodes={handleShowActiveNodes}
                    nodePlugin={nodePlugin}
                />
            )}
            useWindowScroll={true}
        />
    )
}
