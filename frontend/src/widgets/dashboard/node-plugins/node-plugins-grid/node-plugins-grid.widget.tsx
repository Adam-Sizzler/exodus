import { TbAlertTriangle, TbLogin, TbLogout, TbPackage, TbShieldLock } from 'react-icons/tb'
import { Badge, Center, Group, Stack, Text, ThemeIcon } from '@mantine/core'
import { useTranslation } from 'react-i18next'
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
    const { t } = useTranslation()
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
            title: t('common.confirm-action'),
            children: t(
                'node-plugins-grid.widget.delete-this-node-plugin-nodes-using-it-will-be-detached-from-the-plugin'
            ),
            labels: {
                confirm: t('common.delete'),
                cancel: t('common.cancel')
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
                    iconColor="teal"
                    iconVariant="soft"
                    title={t('node-plugin-card.widget.active-on-nodes')}
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
                        iconVariant="soft"
                        subtitle={t(
                            'node-plugins-grid.widget.node-plugins-are-an-advanced-feature-please-review-the-documentation-before-use'
                        )}
                        title={t('node-plugins-grid.widget.warning')}
                        titleOrder={4}
                    />
                </SectionCard.Section>

                <SectionCard.Section>
                    <Center py="xl">
                        <Stack align="center" gap="lg">
                            <ThemeIcon color="gray" radius="xl" size={64} variant="soft">
                                <TbPackage size={32} />
                            </ThemeIcon>

                            <Stack align="center" gap="xs">
                                <Text fw={600} size="lg" ta="center">
                                    {t('node-plugins-grid.widget.no-node-plugins-yet')}
                                </Text>
                                <Text c="dimmed" maw={460} size="sm" ta="center">
                                    {t(
                                        'node-plugins-grid.widget.create-a-plugin-to-extend-node-capabilities-with'
                                    )}
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
                                    {t('node-plugins-grid.widget.ingress-filter')}
                                </Badge>
                                <Badge
                                    color="orange"
                                    leftSection={<TbLogout size={16} />}
                                    radius="md"
                                    size="lg"
                                    variant="light"
                                >
                                    {t('node-plugins-grid.widget.egress-filter')}
                                </Badge>
                                <Badge
                                    color="teal"
                                    leftSection={<TbShieldLock size={16} />}
                                    radius="md"
                                    size="lg"
                                    variant="light"
                                >
                                    {t('node-plugins-grid.widget.haproxy-auth')}
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