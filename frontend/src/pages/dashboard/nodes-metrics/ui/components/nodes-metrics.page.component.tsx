import { EditNodeByUuidModalWidget } from '@widgets/dashboard/nodes/edit-node-by-uuid-modal'
import { LinkedHostsDrawer } from '@widgets/dashboard/nodes/linked-hosts-drawer'
import { NodeUsersUsageDrawer } from '@widgets/dashboard/nodes/node-users-usage-statistic'
import { NodeMetricsWidget } from '@widgets/dashboard/nodes/nodes-metrics'
import { useTranslation } from 'react-i18next'
import { PiSpeedometer } from 'react-icons/pi'

import { LoadingScreen, Page, PageHeaderShared } from '@shared/ui'

import { IProps } from './interfaces'

export default function NodesMetricsPageComponent(props: IProps) {
    const { t } = useTranslation()
    const { isLoading } = props

    return (
        <Page title={t('constants.nodes-metrics')}>
            <PageHeaderShared icon={<PiSpeedometer size={24} />} title="Metrics Overview" />

            {isLoading ? <LoadingScreen height="80vh" /> : <NodeMetricsWidget />}

            <EditNodeByUuidModalWidget />
            <NodeUsersUsageDrawer key="node-users-usage-drawer" />
            <LinkedHostsDrawer key="linked-hosts-drawer" />
        </Page>
    )
}
