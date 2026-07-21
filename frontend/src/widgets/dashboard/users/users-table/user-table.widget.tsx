import { UserActionGroupFeature } from '@features/dashboard/users/users-action-group'
import { useUserTableColumns } from '@features/dashboard/users/users-table/model/use-table-columns'
import { UsersTableSelectionFeature } from '@features/ui/dashboard/users/users-table-selection/users-table-selection.feature'
import {
    MantineReactTable,
    MRT_ColumnFilterFnsState,
    MRT_ShowHideColumnsButton,
    MRT_SortingState,
    MRT_ToggleDensePaddingButton,
    MRT_ToggleFullScreenButton,
    useMantineReactTable
} from '@kastov/mantine-react-table-open'
import { Badge } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { useEffect, useLayoutEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiUsersDuotone } from 'react-icons/pi'
import { TbSearch, TbSearchOff } from 'react-icons/tb'
import { useSearchParams } from 'react-router'

import {
    useGetExternalSquads,
    useGetInternalSquads,
    useGetNodes,
    useGetUsersV2,
    useGetUserTags
} from '@shared/api/hooks'
import { SEARCH_PARAMS } from '@shared/constants/search-params'
import { DEFAULT_PAGINATION_STATE, useMrtTableBinding } from '@shared/lib/mrt-table-store'
import { DataTableShared } from '@shared/ui/table'
import { preventBackScrollTables } from '@shared/utils/misc'
import { sToMs } from '@shared/utils/time-utils'

import { useUserModalStoreActions } from '@entities/dashboard/user-modal-store'
import {
    useBulkUsersActionsStoreActions,
    useBulkUsersActionsStoreTableSelection
} from '@entities/dashboard/users/bulk-users-actions-store'
import { useUsersTableStore } from '@entities/dashboard/users/users-table-store'

export function UserTableWidget() {
    const { t } = useTranslation()

    const { data: internalSquads } = useGetInternalSquads()
    const { data: externalSquads } = useGetExternalSquads()
    const { data: nodes } = useGetNodes()
    const { data: tags } = useGetUserTags()

    const tableColumns = useUserTableColumns(internalSquads, externalSquads, nodes)
    const bulkUsersActionsStoreActions = useBulkUsersActionsStoreActions()
    const tableSelection = useBulkUsersActionsStoreTableSelection()
    const userModalActions = useUserModalStoreActions()
    const [searchParams, setSearchParams] = useSearchParams()

    const { state: persistedTableState, handlers: persistedTableHandlers } =
        useMrtTableBinding(useUsersTableStore)

    const [sorting, setSorting] = useState<MRT_SortingState>([])

    const defaultFilterFns: Record<string, string> = {
        hwidDeviceLimit: 'equals',
        tag: 'equals',
        trafficLimitBytes: 'between'
    }

    const [columnFilterFns, setColumnFilterFns] = useState<MRT_ColumnFilterFnsState>(() =>
        Object.fromEntries(
            tableColumns.map(({ accessorKey }) => [
                accessorKey,
                defaultFilterFns[accessorKey!] ?? 'contains'
            ])
        )
    )

    useLayoutEffect(() => {
        document.body.addEventListener('wheel', preventBackScrollTables, {
            passive: false
        })
        return () => {
            document.body.removeEventListener('wheel', preventBackScrollTables)
        }
    }, [])

    const params = {
        start: persistedTableState.pagination.pageIndex * persistedTableState.pagination.pageSize,
        size: persistedTableState.pagination.pageSize,
        filters: persistedTableState.columnFilters.filter(({ value }) =>
            Array.isArray(value)
                ? value.some((bound) => bound !== null && bound !== undefined && bound !== '')
                : value !== null && value !== undefined && value !== ''
        ),
        filterModes: columnFilterFns,
        sorting
    }

    const {
        data: usersResponse,
        isError,
        isFetching,
        isLoading,
        refetch
    } = useGetUsersV2({
        query: params,
        rQueryParams: {
            // enabled: bulkUsersActionsStoreActions.getUuidLength() === 0,
            refetchInterval: bulkUsersActionsStoreActions.getUuidLength() === 0 ? sToMs(25) : false
        }
    })

    useEffect(() => {
        if (!isLoading && searchParams.get(SEARCH_PARAMS.USER)) {
            userModalActions.setUserUuid(searchParams.get(SEARCH_PARAMS.USER)!)
            userModalActions.changeModalState(true)

            searchParams.delete(SEARCH_PARAMS.USER)
            setSearchParams(searchParams)
        }
    }, [searchParams, isLoading])

    const table = useMantineReactTable({
        columns: tableColumns,
        data: usersResponse?.users ?? [],
        enableFacetedValues: true,
        getFacetedUniqueValues: (_table, columnId) => () => {
            if (columnId === 'tag') {
                return new Map<string, number>(tags?.tags.map((tag) => [tag, 0]) ?? [])
            }
            if (columnId === 'status') {
                return new Map<string, number>(
                    ['ACTIVE', 'DISABLED', 'LIMITED', 'EXPIRED'].map((status) => [status, 0]) ?? []
                )
            }
            return new Map<string, number>()
        },
        columnFilterDisplayMode: 'subheader',
        mantineFilterSelectProps: ({ column }) => {
            const value = column.getFilterValue()
            return {
                clearable: value !== undefined && value !== null && value !== ''
            }
        },
        mantineFilterMultiSelectProps: ({ column }) => {
            const value = column.getFilterValue()
            const count = Array.isArray(value) ? value.length : 0
            return {
                clearable: count > 0,
                renderPill: () => null,
                ...(count > 0 && {
                    leftSection: <Badge variant="soft">{count}</Badge>,
                    placeholder: '',
                    clearSectionMode: 'clear'
                })
            }
        },
        icons: {
            // oxlint-disable-next-line
            IconFilter: (props: any) => <TbSearch size={24} {...props} />,
            // oxlint-disable-next-line
            IconFilterOff: (props: any) => <TbSearchOff size={24} {...props} />
        },
        // mantineTableBodyCellProps: { style: { padding: '2px 6px' } },
        enableFullScreenToggle: true,
        enableSortingRemoval: true,
        enableGlobalFilter: false,
        enableClickToCopy: false,
        enableColumnFilterModes: true,
        enableColumnOrdering: true,
        columnFilterModeOptions: ['contains'],
        initialState: {
            density: 'xxs',
            pagination: DEFAULT_PAGINATION_STATE
        },
        mantineTopToolbarProps: {
            style: {
                '--mrt-base-background-color': '#1b2027'
            }
        },
        mantineTableHeadProps: {
            style: {
                '--mrt-base-background-color': '#1b2027'
            }
        },
        mantineBottomToolbarProps: {
            style: {
                '--mrt-base-background-color': '#1b2027'
            }
        },
        enableDensityToggle: true,
        manualFiltering: true,
        manualPagination: true,
        manualSorting: true,
        // mantinePaginationProps: {
        //     rowsPerPageOptions: ['25', '50', '100']
        // },

        // icons: customIcons,
        enableColumnResizing: true,

        /* prettier-ignore */
        mantineToolbarAlertBannerProps: isError ? {
            color: 'red',
            children: t('user-table.widget.error-loading-data')
        } : undefined,

        ...persistedTableHandlers,
        onColumnFilterFnsChange: setColumnFilterFns,
        onSortingChange: setSorting,
        mantinePaperProps: {
            style: {
                '--paper-radius': 'var(--mantine-radius-xs)'
            },
            withBorder: false
        },
        rowCount: usersResponse?.total ?? 0,
        enableRowSelection: true,
        mantineSelectCheckboxProps: {
            size: 'md',
            color: 'cyan',
            variant: 'outline'
        },
        mantineSelectAllCheckboxProps: {
            size: 'md',
            color: 'cyan',
            variant: 'outline'
        },
        enableColumnPinning: true,
        positionToolbarAlertBanner: 'top',
        renderToolbarAlertBannerContent: () => {
            return (
                <UsersTableSelectionFeature
                    resetRowSelection={table.resetRowSelection}
                    toggleAllPageRowsSelected={table.toggleAllPageRowsSelected}
                />
            )
        },
        selectAllMode: 'page',
        renderToolbarInternalActions: ({ table: t }) => (
            <>
                <MRT_ShowHideColumnsButton table={t} />
                <MRT_ToggleDensePaddingButton table={t} />
                <MRT_ToggleFullScreenButton table={t} />
            </>
        ),
        state: {
            ...persistedTableState,
            columnFilterFns,
            isLoading,
            showAlertBanner: isError,
            showColumnFilters: true,
            showProgressBars: isFetching,
            sorting,
            rowSelection: tableSelection
        },
        mantineTableBodyRowProps: ({ row }) => ({
            onClick: async () => {
                if (row.id === 'mrt-row-empty' || row.original.uuid === undefined) {
                    notifications.show({
                        title: 'Nice try!',
                        message: 'Nothing to show...',
                        color: 'indigo'
                    })
                    return
                }
                await userModalActions.setUserUuid(row.original.uuid)
                userModalActions.changeModalState(true)
            },
            style: {
                cursor: 'pointer'
            }
        }),
        onRowSelectionChange: bulkUsersActionsStoreActions.setTableSelection,
        getRowId: (originalRow) => originalRow.uuid
    })

    return (
        <DataTableShared.Container>
            <DataTableShared.Title
                actions={
                    <UserActionGroupFeature isLoading={isLoading} refetch={refetch} table={table} />
                }
                icon={<PiUsersDuotone size={24} />}
                title={t('user-table.widget.table-title')}
            />

            <DataTableShared.Content>
                <MantineReactTable table={table} />
            </DataTableShared.Content>
        </DataTableShared.Container>
    )
}
