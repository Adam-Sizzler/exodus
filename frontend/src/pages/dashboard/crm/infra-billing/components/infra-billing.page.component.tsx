import { Stack } from '@mantine/core'
import { CreateInfraBillingNodeModalWidget } from '@widgets/dashboard/infra-billing/create-infra-billing-node-modal/create-infra-billing-node.modal.widget'
import { CreateInfraBillingRecordDrawerWidget } from '@widgets/dashboard/infra-billing/create-infra-billing-record-modal/create-infra-billing-record.modal.widget'
import { CreateInfraProviderDrawerWidget } from '@widgets/dashboard/infra-billing/create-infra-provider-drawer/create-infra-provider.drawer.widget'
import { DesktopColumnsInfraBillingWidget } from '@widgets/dashboard/infra-billing/desktop-columns'
import { MobileInfraBillingWidget } from '@widgets/dashboard/infra-billing/mobile'
import { StatsWidget } from '@widgets/dashboard/infra-billing/stats-widget/stats.widget'
import { UpdateBillingDateModalWidget } from '@widgets/dashboard/infra-billing/update-billing-date-modal'
import { ViewInfraProviderDrawerWidget } from '@widgets/dashboard/infra-billing/view-infra-provider-drawer/view-infra-provider.drawer.widget'
import { EditNodeByUuidModalWidget } from '@widgets/dashboard/nodes/edit-node-by-uuid-modal/edit-node-by-uuid-modal.widget'
import { LinkedHostsDrawer } from '@widgets/dashboard/nodes/linked-hosts-drawer/linked-hosts-drawer.widget'
import { NodeUsersUsageDrawer } from '@widgets/dashboard/nodes/node-users-usage-statistic/node-users-usage-drawer.widget'
import { NodesUsersUsageStatisticsDrawer } from '@widgets/dashboard/nodes/nodes-users-usage-statistics/nodes-users-usage-statistics.drawer'
import { motion } from 'motion/react'
import { useLayoutEffect } from 'react'
import { useTranslation } from 'react-i18next'

import { useIsMobile } from '@shared/hooks'
import { Page } from '@shared/ui/page'
import { preventBackScrollTables } from '@shared/utils/misc'

export const InfraBillingPageComponent = () => {
    const { t } = useTranslation()
    const isMobile = useIsMobile()

    useLayoutEffect(() => {
        if (isMobile) {
            return undefined
        }

        document.body.addEventListener('wheel', preventBackScrollTables, {
            passive: false
        })
        return () => {
            document.body.removeEventListener('wheel', preventBackScrollTables)
        }
    }, [isMobile])

    return (
        <Page title={t('constants.infra-billing')}>
            <motion.div
                animate={{ opacity: 1 }}
                initial={{ opacity: 0 }}
                transition={{ duration: 0.5 }}
            >
                {isMobile ? (
                    <MobileInfraBillingWidget />
                ) : (
                    <Stack>
                        <StatsWidget />

                        <DesktopColumnsInfraBillingWidget />
                    </Stack>
                )}
            </motion.div>

            <ViewInfraProviderDrawerWidget />
            <CreateInfraProviderDrawerWidget />
            <CreateInfraBillingRecordDrawerWidget />
            <UpdateBillingDateModalWidget />
            <CreateInfraBillingNodeModalWidget />

            <EditNodeByUuidModalWidget key="edit-node-by-uuid-modal" />
            <NodeUsersUsageDrawer key="node-users-usage-drawer" />
            <LinkedHostsDrawer key="linked-hosts-drawer" />
            <NodesUsersUsageStatisticsDrawer key="nodes-users-usage-statistics-drawer" />
        </Page>
    )
}
