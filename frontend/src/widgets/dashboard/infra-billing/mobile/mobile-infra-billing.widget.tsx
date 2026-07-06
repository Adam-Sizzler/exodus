import { ActionIcon, Group, Stack, Tabs, Transition } from '@mantine/core'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { TbCloud, TbCreditCard, TbPlus, TbRefresh, TbServer } from 'react-icons/tb'

import {
    useGetInfraBillingHistoryRecordsInfinite,
    useGetInfraBillingNodes,
    useGetInfraProviders
} from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

import styles from './mobile-infra-billing.module.css'
import { MobileNodesListWidget } from './mobile-nodes-list.widget'
import { MobileProvidersListWidget } from './mobile-providers-list.widget'
import { MobileStatsWidget } from './mobile-stats.widget'
import { VirtualizedRecordsList } from './virtualized-records-list.widget'

const RECORDS_HEIGHT = 'calc(100vh - 280px)'

type TabValue = 'nodes' | 'providers' | 'records'

export function MobileInfraBillingWidget() {
    const [activeTab, setActiveTab] = useState<TabValue>('nodes')
    const {
        data: infraProviders,
        isLoading: isInfraProvidersLoading,
        refetch: refetchInfraProviders,
        isRefetching: isInfraProvidersRefetching
    } = useGetInfraProviders()
    const {
        data: infraBillingNodes,
        isLoading: isInfraBillingNodesLoading,
        refetch: refetchInfraBillingNodes,
        isRefetching: isInfraBillingNodesRefetching
    } = useGetInfraBillingNodes()
    const {
        data: infraBillingRecordsData,
        refetch: refetchRecords,
        isLoading: isInfraBillingRecordsLoading,
        isRefetching: isInfraBillingRecordsRefetching,
        fetchNextPage: fetchNextRecordsPage,
        hasNextPage: hasNextRecordsPage,
        isFetchingNextPage: isFetchingNextRecordsPage
    } = useGetInfraBillingHistoryRecordsInfinite()

    const billingRecords = infraBillingRecordsData?.pages.flatMap((page) => page.records) ?? []

    const handleLoadMoreRecords = () => {
        if (hasNextRecordsPage && !isFetchingNextRecordsPage) {
            fetchNextRecordsPage()
        }
    }

    const openModalWithData = useModalsStoreOpenWithData()
    const { t } = useTranslation()

    const handleAdd = () => {
        switch (activeTab) {
            case 'nodes':
                openModalWithData(MODALS.CREATE_INFRA_BILLING_NODE_MODAL, undefined)
                break
            case 'providers':
                openModalWithData(MODALS.CREATE_INFRA_PROVIDER_DRAWER, undefined)
                break
            case 'records':
                openModalWithData(MODALS.CREATE_INFRA_BILLING_RECORD_DRAWER, undefined)
                break
            default:
                break
        }
    }

    const handleRefetch = () => {
        refetchRecords()
        refetchInfraBillingNodes()
        refetchInfraProviders()
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

    return (
        <Stack gap="xs" pos="relative">
            <MobileStatsWidget />

            <Tabs
                classNames={{
                    tab: styles.tab,
                    tabLabel: styles.tabLabel
                }}
                color="cyan"
                keepMountedMode="display-none"
                onChange={(value) => {
                    if (value) {
                        setActiveTab(value as TabValue)
                    }
                }}
                value={activeTab}
                variant="unstyled"
            >
                <Tabs.List grow mb="md">
                    <Tabs.Tab leftSection={<TbServer size={16} />} value="nodes">
                        {t('constants.nodes')}
                    </Tabs.Tab>
                    <Tabs.Tab leftSection={<TbCreditCard size={16} />} value="records">
                        {t('mobile-infra-billing.widget.history')}
                    </Tabs.Tab>
                    <Tabs.Tab leftSection={<TbCloud size={16} />} value="providers">
                        {t('mobile-infra-billing.widget.providers')}
                    </Tabs.Tab>
                </Tabs.List>

                <Tabs.Panel value="nodes">
                    <Transition
                        duration={200}
                        keepMounted
                        mounted={activeTab === 'nodes'}
                        timingFunction="linear"
                        transition="fade"
                    >
                        {(styles) => (
                            <MobileNodesListWidget
                                nodes={infraBillingNodes.billingNodes}
                                style={styles}
                            />
                        )}
                    </Transition>
                </Tabs.Panel>

                <Tabs.Panel value="records">
                    <Transition
                        duration={200}
                        keepMounted
                        mounted={activeTab === 'records'}
                        timingFunction="linear"
                        transition="fade"
                    >
                        {(styles) => (
                            <VirtualizedRecordsList
                                height={RECORDS_HEIGHT}
                                isLoadingMore={isFetchingNextRecordsPage}
                                onReachBottom={handleLoadMoreRecords}
                                records={billingRecords}
                                refetchRecords={refetchRecords}
                                style={styles}
                            />
                        )}
                    </Transition>
                </Tabs.Panel>

                <Tabs.Panel value="providers">
                    <Transition
                        duration={200}
                        keepMounted
                        mounted={activeTab === 'providers'}
                        timingFunction="linear"
                        transition="fade"
                    >
                        {(styles) => (
                            <MobileProvidersListWidget
                                providers={infraProviders.providers}
                                style={styles}
                            />
                        )}
                    </Transition>
                </Tabs.Panel>
            </Tabs>

            <Group
                bottom={16}
                gap={8}
                pos="fixed"
                right={15}
                style={{
                    zIndex: 100,
                    backgroundColor: 'var(--mantine-color-dark-6)',
                    padding: 6,
                    borderRadius: 'var(--mantine-radius-md)',
                    boxShadow: '0 4px 12px rgba(0, 0, 0, 0.4)'
                }}
            >
                <ActionIcon
                    color="cyan"
                    loading={
                        isInfraProvidersRefetching ||
                        isInfraBillingNodesRefetching ||
                        isInfraBillingRecordsRefetching
                    }
                    onClick={handleRefetch}
                    size={40}
                    variant="soft"
                >
                    <TbRefresh size={22} />
                </ActionIcon>
                <ActionIcon color="teal" onClick={handleAdd} size={40} variant="soft">
                    <TbPlus size={22} />
                </ActionIcon>
            </Group>
        </Stack>
    )
}
