import { ActionIcon, Badge, SimpleGrid, Tooltip } from '@mantine/core'
import { modals } from '@mantine/modals'
import { MobileNodesListWidget } from '@widgets/dashboard/infra-billing/mobile/mobile-nodes-list.widget'
import { MobileProvidersListWidget } from '@widgets/dashboard/infra-billing/mobile/mobile-providers-list.widget'
import { VirtualizedRecordsList } from '@widgets/dashboard/infra-billing/mobile/virtualized-records-list.widget'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import {
    TbCalendar,
    TbCloud,
    TbCreditCard,
    TbListCheck,
    TbServer,
    TbTrash,
    TbX
} from 'react-icons/tb'

import { queryClient } from '@shared/api'
import {
    QueryKeys,
    useDeleteInfraBillingNode,
    useGetInfraBillingHistoryRecordsInfinite,
    useGetInfraBillingNodes,
    useGetInfraProviders
} from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

import { AddButton } from './add-button'
import { Column } from './column'
import { RefreshButton } from './refresh-button'

const COLUMN_HEIGHT = 'calc(100vh - 320px)'

export function DesktopColumnsInfraBillingWidget() {
    const { t } = useTranslation()
    const openModalWithData = useModalsStoreOpenWithData()

    const [selectMode, setSelectMode] = useState(false)
    const [selectedUuids, setSelectedUuids] = useState<Set<string>>(new Set())

    const {
        data: infraProviders,
        isLoading: isInfraProvidersLoading,
        isFetching: isInfraProvidersFetching,
        refetch: refetchProviders
    } = useGetInfraProviders()
    const {
        data: infraBillingNodes,
        isLoading: isInfraBillingNodesLoading,
        isFetching: isInfraBillingNodesFetching,
        refetch: refetchNodes
    } = useGetInfraBillingNodes()
    const {
        data: infraBillingRecordsData,
        refetch: refetchRecords,
        isLoading: isInfraBillingRecordsLoading,
        isFetching: isInfraBillingRecordsFetching,
        fetchNextPage: fetchNextRecordsPage,
        hasNextPage: hasNextRecordsPage,
        isFetchingNextPage: isFetchingNextRecordsPage
    } = useGetInfraBillingHistoryRecordsInfinite()

    const billingRecords = infraBillingRecordsData?.pages.flatMap((page) => page.records) ?? []
    const billingRecordsTotal = infraBillingRecordsData?.pages[0]?.total ?? 0

    const handleLoadMoreRecords = () => {
        if (hasNextRecordsPage && !isFetchingNextRecordsPage) {
            fetchNextRecordsPage()
        }
    }

    const { mutate: deleteNode } = useDeleteInfraBillingNode({
        mutationFns: {
            onSuccess: () => {
                queryClient.refetchQueries({
                    queryKey: QueryKeys.infraBilling.getInfraBillingNodes.queryKey
                })
            }
        }
    })

    const exitSelectMode = () => {
        setSelectMode(false)
        setSelectedUuids(new Set())
    }

    const toggleSelect = (uuid: string) => {
        setSelectedUuids((prev) => {
            const next = new Set(prev)
            if (next.has(uuid)) {
                next.delete(uuid)
            } else {
                next.add(uuid)
            }
            return next
        })
    }

    const handleDeleteSelected = () =>
        modals.openConfirmModal({
            title: t('common.delete'),
            children: t('common.confirm-action-description'),
            labels: { confirm: t('common.delete'), cancel: t('common.cancel') },
            centered: true,
            confirmProps: { color: 'red' },
            onConfirm: () => {
                for (const uuid of selectedUuids) {
                    deleteNode({ route: { uuid } })
                }
                exitSelectMode()
            }
        })

    const handleUpdateSelected = () => {
        openModalWithData(MODALS.UPDATE_BILLING_DATE_MODAL, {
            uuids: Array.from(selectedUuids),
            callback: exitSelectMode
        })
    }

    if (
        isInfraBillingNodesLoading ||
        isInfraProvidersLoading ||
        isInfraBillingRecordsLoading ||
        !infraBillingNodes ||
        !infraProviders ||
        !infraBillingRecordsData
    ) {
        return <LoadingScreen />
    }

    const nodesActions = selectMode ? (
        <>
            <Badge color="cyan" radius="sm" size="lg" variant="light">
                {selectedUuids.size}
            </Badge>

            <Tooltip label={t('common.delete')} withArrow>
                <ActionIcon
                    color="red"
                    disabled={selectedUuids.size === 0}
                    onClick={handleDeleteSelected}
                    size="input-xs"
                    variant="soft"
                >
                    <TbTrash size={18} />
                </ActionIcon>
            </Tooltip>

            <Tooltip label={t('infra-billing-nodes.widget.update-multiple-nodes')} withArrow>
                <ActionIcon
                    color="teal"
                    disabled={selectedUuids.size === 0}
                    onClick={handleUpdateSelected}
                    size="input-xs"
                    variant="soft"
                >
                    <TbCalendar size={18} />
                </ActionIcon>
            </Tooltip>

            <Tooltip label={t('common.close')} withArrow>
                <ActionIcon color="gray" onClick={exitSelectMode} size="input-xs" variant="subtle">
                    <TbX size={18} />
                </ActionIcon>
            </Tooltip>
        </>
    ) : (
        <>
            <Tooltip label={t('common.select')} withArrow>
                <ActionIcon
                    color="gray"
                    onClick={() => setSelectMode(true)}
                    size="input-xs"
                    variant="subtle"
                >
                    <TbListCheck size={18} />
                </ActionIcon>
            </Tooltip>

            <RefreshButton loading={isInfraBillingNodesFetching} onClick={refetchNodes} />

            <AddButton
                onClick={() => openModalWithData(MODALS.CREATE_INFRA_BILLING_NODE_MODAL, undefined)}
            />
        </>
    )

    return (
        <SimpleGrid cols={{ base: 1, md: 2, lg: 3 }} spacing="md">
            <Column
                actions={nodesActions}
                count={infraBillingNodes.billingNodes.length}
                height={COLUMN_HEIGHT}
                icon={<TbServer size={18} />}
                title={t('constants.nodes')}
            >
                <MobileNodesListWidget
                    nodes={infraBillingNodes.billingNodes}
                    onToggleSelect={toggleSelect}
                    selectedUuids={selectedUuids}
                    selectMode={selectMode}
                    style={{}}
                />
            </Column>

            <Column
                actions={
                    <>
                        <RefreshButton
                            loading={isInfraBillingRecordsFetching && !isFetchingNextRecordsPage}
                            onClick={refetchRecords}
                        />

                        <AddButton
                            onClick={() =>
                                openModalWithData(
                                    MODALS.CREATE_INFRA_BILLING_RECORD_DRAWER,
                                    undefined
                                )
                            }
                        />
                    </>
                }
                count={billingRecordsTotal}
                height={COLUMN_HEIGHT}
                icon={<TbCreditCard size={18} />}
                scrollable={false}
                title={t('mobile-infra-billing.widget.history')}
            >
                <VirtualizedRecordsList
                    height={COLUMN_HEIGHT}
                    isLoadingMore={isFetchingNextRecordsPage}
                    onReachBottom={handleLoadMoreRecords}
                    records={billingRecords}
                    refetchRecords={refetchRecords}
                />
            </Column>

            <Column
                actions={
                    <>
                        <RefreshButton
                            loading={isInfraProvidersFetching}
                            onClick={refetchProviders}
                        />

                        <AddButton
                            onClick={() =>
                                openModalWithData(MODALS.CREATE_INFRA_PROVIDER_DRAWER, undefined)
                            }
                        />
                    </>
                }
                count={infraProviders.providers.length}
                height={COLUMN_HEIGHT}
                icon={<TbCloud size={18} />}
                title={t('mobile-infra-billing.widget.providers')}
            >
                <MobileProvidersListWidget providers={infraProviders.providers} style={{}} />
            </Column>
        </SimpleGrid>
    )
}
