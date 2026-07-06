import { useBandwidthTableColumns } from '@features/dashboard/nodes-bandwidth-table/bandwidth-table-columns/model/use-bandwidth-table-columns'
import { ListViewTable } from '@gfazioli/mantine-list-view-table'
import { Table } from '@mantine/core'
import { GetAllNodesCommand } from '@exodus/backend-contract'
import dayjs from 'dayjs'
import { useTranslation } from 'react-i18next'
import { TbChartArcs } from 'react-icons/tb'

import { useGetNodes } from '@shared/api/hooks'
import { DataTableShared } from '@shared/ui/table'

export function NodesBandwidthTableWidget() {
    const { data: nodes, isLoading } = useGetNodes({
        rQueryParams: {
            select: (data: unknown) => {
                const nodes = data as GetAllNodesCommand.Response['response']
                return nodes.filter((node) => node.isTrafficTrackingActive)
            }
        }
    })
    const { t } = useTranslation()

    const tableColumns = useBandwidthTableColumns()

    return (
        <DataTableShared.Container>
            <DataTableShared.Title
                icon={<TbChartArcs size={24} />}
                title={`${t('table.widget.today')}: ${dayjs().format('DD MMMM')}`}
            />
            <DataTableShared.Content>
                <Table.ScrollContainer minWidth={1200}>
                    <ListViewTable
                        columns={tableColumns}
                        data={nodes ?? []}
                        emptyText={t('table.widget.nodes-with-active-traffic-tracking-not-found')}
                        highlightOnHover
                        loading={isLoading}
                        rowKey="uuid"
                        stripedColor="cyan"
                        withColumnBorders={false}
                    />
                </Table.ScrollContainer>
            </DataTableShared.Content>
        </DataTableShared.Container>
    )
}
