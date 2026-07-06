import {
    MantineReactTable,
    type MRT_ColumnDef,
    useMantineReactTable
} from '@kastov/mantine-react-table-open'
import { Anchor } from '@mantine/core'
import { GetUserSubscriptionRequestHistoryCommand } from '@exodus/backend-contract'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { TbSearch } from 'react-icons/tb'

import { formatTimeUtil } from '@shared/utils/time-utils'

type TRecord = GetUserSubscriptionRequestHistoryCommand.Response['response']['records'][number]

interface IProps {
    isLoading: boolean
    records: TRecord[] | undefined
}

export const UserSubscriptionRequestsTable = (props: IProps) => {
    const { records, isLoading } = props
    const { t, i18n } = useTranslation()

    const columns = useMemo<MRT_ColumnDef<TRecord>[]>(
        () => [
            {
                accessorKey: 'id',
                header: 'ID',
                size: 80,
                accessorFn: (row) => row.id
            },
            {
                accessorKey: 'requestIp',
                header: t('get-user-subscription-request-history.feature.ip-address'),
                accessorFn: (row) => row.requestIp || '–',
                Cell: ({ row }) => {
                    const ip = row.original.requestIp
                    if (!ip) return '–'
                    return (
                        <Anchor
                            c="cyan"
                            ff="monospace"
                            href={`https://ipinfo.io/${ip}`}
                            rel="noopener noreferrer"
                            size="sm"
                            style={{ display: 'inline-flex', alignItems: 'center', gap: 4 }}
                            target="_blank"
                            underline="never"
                        >
                            {ip}
                        </Anchor>
                    )
                }
            },
            {
                accessorKey: 'userAgent',
                header: t('get-user-subscription-request-history.feature.user-agent'),
                enableClickToCopy: true,
                size: 400,
                accessorFn: (row) => row.userAgent || '–'
            },
            {
                accessorKey: 'requestAt',
                header: t('get-user-subscription-request-history.feature.request-at'),
                sortingFn: 'datetime',
                size: 250,
                Cell: ({ row }) =>
                    formatTimeUtil({
                        time: row.original.requestAt,
                        template: 'TIME_FIRST_DATETIME',
                        language: i18n.language
                    }),
                mantineTableBodyCellProps: {
                    align: 'left',
                    ff: 'monospace'
                }
            }
        ],
        [t, i18n.language]
    )

    const table = useMantineReactTable({
        columns,
        data: records ?? [],
        state: {
            showSkeletons: isLoading,
            showProgressBars: isLoading
        },
        enableGrouping: true,
        rowCount: records?.length ?? 0,
        enablePagination: false,
        enableRowVirtualization: false,
        enableBottomToolbar: false,
        enableTopToolbar: true,
        enableColumnFilters: false,
        enableGlobalFilter: true,
        enableColumnActions: true,
        enableDensityToggle: false,
        enableFullScreenToggle: false,
        enableHiding: false,
        enableColumnResizing: true,
        enableStickyHeader: true,
        enableSortingRemoval: true,
        initialState: {
            showGlobalFilter: true,
            density: 'xs',
            sorting: [{ id: 'requestAt', desc: true }],
            pagination: {
                pageIndex: 0,
                pageSize: 24
            }
        },
        mantineTableContainerProps: { style: { maxHeight: '500px', minHeight: '500px' } },
        getRowId: (row) => (row.id ? row.id.toString() : ''),
        displayColumnDefOptions: {
            'mrt-row-actions': { size: 60, header: '' }
        },
        mantineSearchTextInputProps: {
            placeholder: t('user-hwid-devices.drawer.widget.enter-search-query'),
            style: { minWidth: '350px' },
            variant: 'default',
            leftSection: <TbSearch size={16} />
        }
    })

    return <MantineReactTable table={table} />
}
