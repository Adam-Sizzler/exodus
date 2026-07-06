import {
    MantineReactTable,
    type MRT_ColumnDef,
    useMantineReactTable
} from '@kastov/mantine-react-table-open'
import { ActionIcon, Anchor, Group, Text, ThemeIcon } from '@mantine/core'
import { GetUserHwidDevicesCommand } from '@exodus/backend-contract'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiDeviceMobile, PiLinuxLogo } from 'react-icons/pi'
import {
    TbBrandAndroid,
    TbBrandApple,
    TbBrandFinder,
    TbBrandWindows,
    TbSearch,
    TbTrash
} from 'react-icons/tb'

import { formatTimeUtil } from '@shared/utils/time-utils'

type TDevice = GetUserHwidDevicesCommand.Response['response']['devices'][number]

interface IProps {
    devices: TDevice[] | undefined
    isLoading: boolean
    onDelete: (hwid: string) => void
}

const resolvePlatformIcon = (platform: null | string) => {
    switch (platform?.toLowerCase()) {
        case 'android':
            return <TbBrandAndroid size={24} />
        case 'ios':
            return <TbBrandApple size={24} />
        case 'linux':
            return <PiLinuxLogo size={24} />
        case 'macos':
            return <TbBrandFinder size={24} />
        case 'windows':
            return <TbBrandWindows size={24} />
        default:
            return <PiDeviceMobile size={24} />
    }
}

export const UserHwidDevicesTable = (props: IProps) => {
    const { devices, onDelete, isLoading } = props
    const { t, i18n } = useTranslation()

    const columns = useMemo<MRT_ColumnDef<TDevice>[]>(
        () => [
            {
                accessorKey: 'platform',
                header: 'Platform',
                accessorFn: (row) => row.platform || 'Unknown',
                size: 150,
                Cell: ({ row }) => (
                    <Group justify="center" wrap="nowrap">
                        <ThemeIcon color="indigo" size="lg" variant="soft">
                            {resolvePlatformIcon(row.original.platform)}
                        </ThemeIcon>
                        <Text size="sm">{row.original.platform || 'Unknown'}</Text>
                    </Group>
                )
            },
            {
                accessorKey: 'hwid',
                header: 'HWID',
                enableClickToCopy: true
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
                accessorKey: 'osVersion',
                header: t('get-hwid-user-devices.feature.os-version'),
                size: 100,
                accessorFn: (row) => row.osVersion || '–'
            },
            {
                accessorKey: 'deviceModel',
                header: t('get-hwid-user-devices.feature.model'),

                accessorFn: (row) => row.deviceModel || '–'
            },
            {
                accessorKey: 'userAgent',
                header: t('get-hwid-user-devices.feature.user-agent'),
                size: 170,
                enableClickToCopy: true,
                accessorFn: (row) => row.userAgent || '–'
            },
            {
                accessorKey: 'createdAt',
                header: t('get-hwid-user-devices.feature.added'),
                sortingFn: 'datetime',
                Cell: ({ row }) =>
                    formatTimeUtil({
                        time: row.original.createdAt,
                        template: 'TIME_FIRST_DATETIME',
                        language: i18n.language
                    }),
                mantineTableBodyCellProps: {
                    align: 'left',
                    ff: 'monospace'
                },
                minSize: 250
            },
            {
                accessorKey: 'updatedAt',
                header: 'Updated',
                sortingFn: 'datetime',
                Cell: ({ row }) =>
                    formatTimeUtil({
                        time: row.original.updatedAt,
                        template: 'TIME_FIRST_DATETIME',
                        language: i18n.language
                    }),
                mantineTableBodyCellProps: {
                    align: 'left',
                    ff: 'monospace'
                },
                minSize: 250
            }
        ],
        [t, i18n.language]
    )

    const table = useMantineReactTable({
        columns,
        data: devices ?? [],
        state: {
            showSkeletons: isLoading,
            showProgressBars: isLoading
        },
        enableGrouping: true,
        rowCount: devices?.length ?? 0,
        enablePagination: false,
        enableRowVirtualization: true,
        rowVirtualizerOptions: { overscan: 8 },
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
            sorting: [{ id: 'createdAt', desc: true }],
            pagination: {
                pageIndex: 0,
                pageSize: 20
            }
        },
        mantineTableContainerProps: { style: { maxHeight: '500px', minHeight: '500px' } },
        getRowId: (row) => row.hwid,
        enableRowActions: true,
        positionActionsColumn: 'last',
        displayColumnDefOptions: {
            'mrt-row-actions': { size: 60, header: '' }
        },
        mantineSearchTextInputProps: {
            placeholder: t('user-hwid-devices.drawer.widget.enter-search-query'),
            style: { minWidth: '350px' },
            variant: 'default',
            leftSection: <TbSearch size={16} />
        },
        renderRowActions: ({ row }) => (
            <ActionIcon
                aria-label={t('get-hwid-user-devices.feature.delete-device')}
                color="red"
                onClick={() => onDelete(row.original.hwid)}
                size="md"
                variant="soft"
            >
                <TbTrash size={16} />
            </ActionIcon>
        )
    })

    return <MantineReactTable table={table} />
}
