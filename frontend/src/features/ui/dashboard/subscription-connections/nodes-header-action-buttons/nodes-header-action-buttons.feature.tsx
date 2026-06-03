import {
    TbAlertCircle,
    TbCards,
    TbPlus,
    TbRefresh,
    TbRocket,
    TbSearch,
    TbTable
} from 'react-icons/tb'
import { ActionIcon, ActionIconGroup, Group, Stack, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { spotlight } from '@mantine/spotlight'
import { PiSpiral } from 'react-icons/pi'
import { modals } from '@mantine/modals'

import { useSubscriptionConnectionsStoreActions } from '@entities/dashboard/subscription-connections/nodes-store/nodes-store'
import { NodesViewMode } from '@pages/dashboard/subscription-connections/ui/components/interfaces'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import {
    useGetSubscriptionConnections,
    useRestartAllSubscriptionConnections
} from '@shared/api/hooks'
import { ActionCardShared } from '@shared/ui'

interface IProps {
    setViewMode: (viewMode: NodesViewMode) => void
    viewMode: NodesViewMode
}

export const NodesHeaderActionButtonsFeature = (props: IProps) => {
    const { setViewMode, viewMode } = props

    const { t } = useTranslation()

    const actions = useSubscriptionConnectionsStoreActions()

    const handleCreate = () => {
        actions.toggleCreateModal(true)
    }

    const {
        isLoading: isGetNodesPending,
        refetch: refetchNodes,
        isPending,
        isRefetching
    } = useGetSubscriptionConnections()
    const { mutate: restartAllNodes, isPending: isRestartAllNodesPending } =
        useRestartAllSubscriptionConnections()

    const openRestartAllNodesModal = () => {
        modals.open({
            title: (
                <BaseOverlayHeader
                    IconComponent={TbRocket}
                    iconVariant="soft"
                    iconColor="teal"
                    title={t('nodes-header-action-buttons.feature.restart-all-nodes')}
                />
            ),
            centered: true,
            size: 'md',
            children: (
                <Stack gap="sm">
                    <ActionCardShared
                        description={t(
                            'nodes-header-action-buttons.feature.force-restart-description'
                        )}
                        icon={<TbAlertCircle size={22} />}
                        isLoading={isPending}
                        onClick={() => {
                            restartAllNodes({
                                variables: {
                                    forceRestart: true
                                }
                            })
                            modals.closeAll()
                        }}
                        title={t('nodes-header-action-buttons.feature.force')}
                        iconColor="red"
                        variant="soft"
                    />

                    <ActionCardShared
                        description={t(
                            'nodes-header-action-buttons.feature.graceful-restart-description-1'
                        )}
                        icon={<TbRocket size={22} />}
                        isLoading={isPending}
                        onClick={() => {
                            restartAllNodes({
                                variables: {
                                    forceRestart: false
                                }
                            })
                            modals.closeAll()
                        }}
                        title={t('nodes-header-action-buttons.feature.graceful')}
                        iconColor="teal"
                        variant="soft"
                    />
                </Stack>
            )
        })
    }

    return (
        <Group grow preventGrowOverflow={false} wrap="wrap">
            {viewMode === NodesViewMode.CARDS && (
                <ActionIconGroup>
                    <Tooltip label={t('nodes-header-action-buttons.feature.search-nodes')}>
                        <ActionIcon
                            color="gray"
                            onClick={spotlight.open}
                            size="input-md"
                            variant="soft"
                        >
                            <TbSearch size="24px" />
                        </ActionIcon>
                    </Tooltip>
                </ActionIconGroup>
            )}

            <ActionIconGroup>
                <Tooltip label="Toggle view mode">
                    <ActionIcon
                        color="gray"
                        onClick={() =>
                            setViewMode(
                                viewMode === NodesViewMode.TABLE
                                    ? NodesViewMode.CARDS
                                    : NodesViewMode.TABLE
                            )
                        }
                        size="input-md"
                        variant="soft"
                    >
                        {viewMode === NodesViewMode.CARDS ? (
                            <TbTable size="24px" />
                        ) : (
                            <TbCards size="24px" />
                        )}
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
            <ActionIconGroup>
                <Tooltip
                    label={t('nodes-header-action-buttons.feature.restart-all-nodes')}
                    withArrow
                >
                    <ActionIcon
                        color="grape"
                        loading={isRestartAllNodesPending}
                        onClick={() => {
                            openRestartAllNodesModal()
                        }}
                        size="input-md"
                        variant="soft"
                    >
                        <PiSpiral size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
            <ActionIconGroup>
                <Tooltip label={t('common.update')} withArrow>
                    <ActionIcon
                        loading={isGetNodesPending || isPending || isRefetching}
                        onClick={() => refetchNodes()}
                        size="input-md"
                        variant="soft"
                    >
                        <TbRefresh size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
            <ActionIconGroup>
                <Tooltip
                    label={t('nodes-header-action-buttons.feature.create-new-subscription', {
                        defaultValue: 'Создать подписку'
                    })}
                    withArrow
                >
                    <ActionIcon color="teal" onClick={handleCreate} size="input-md" variant="soft">
                        <TbPlus size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>
        </Group>
    )
}
