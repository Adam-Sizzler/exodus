/* eslint-disable no-nested-ternary */
import { useTranslation } from 'react-i18next'
import { Grid, Stack } from '@mantine/core'
import { HiServer } from 'react-icons/hi'
import { motion } from 'motion/react'
import { useState } from 'react'

import { LinkedHostsDrawer } from '@widgets/dashboard/subscription-connections/linked-hosts-drawer/linked-hosts-drawer.widget'
import { NodesHeaderActionButtonsFeature } from '@features/ui/dashboard/subscription-connections/nodes-header-action-buttons'
import { NodesDataTableWidget } from '@widgets/dashboard/subscription-connections/nodes-datatable/nodes-datatable.widget'
import { EditNodeByUuidModalWidget } from '@widgets/dashboard/subscription-connections/edit-node-by-uuid-modal'
import { CreateNodeModalWidget } from '@widgets/dashboard/subscription-connections/create-node-modal'
import { NodesTableWidget } from '@widgets/dashboard/subscription-connections/nodes-table'
import { LoadingScreen, Page, PageHeaderShared } from '@shared/ui'
import { SubscriptionConnectionResponse } from '@shared/api/hooks'

import { IProps, NodesViewMode } from './interfaces'

export default function NodesPageComponent(props: IProps) {
    const { nodes, isLoading } = props

    const { t } = useTranslation()

    const [viewMode, setViewMode] = useState<NodesViewMode>(NodesViewMode.CARDS)
    const [selectedRecords, setSelectedRecords] = useState<SubscriptionConnectionResponse[]>([])

    return (
        <Page title={t('constants.subscription-connections')}>
            <Grid>
                <Grid.Col span={12}>
                    <Stack>
                        <PageHeaderShared
                            actions={
                                <NodesHeaderActionButtonsFeature
                                    setViewMode={setViewMode}
                                    viewMode={viewMode}
                                />
                            }
                            icon={<HiServer size={24} />}
                            title={t('constants.subscription-connections')}
                        />
                    </Stack>

                    {isLoading ? (
                        <LoadingScreen height="60vh" />
                    ) : viewMode === NodesViewMode.TABLE ? (
                        <motion.div
                            animate={{ opacity: 1 }}
                            initial={{ opacity: 0 }}
                            transition={{
                                duration: 0.5
                            }}
                        >
                            <NodesDataTableWidget
                                nodes={nodes}
                                selectedRecords={selectedRecords}
                                setSelectedRecords={setSelectedRecords}
                            />
                        </motion.div>
                    ) : (
                        <NodesTableWidget nodes={nodes} />
                    )}
                </Grid.Col>
            </Grid>

            <EditNodeByUuidModalWidget key="edit-node-by-uuid-modal" />
            <CreateNodeModalWidget key="create-node-widget" />
            <LinkedHostsDrawer key="linked-hosts-drawer" />
        </Page>
    )
}
