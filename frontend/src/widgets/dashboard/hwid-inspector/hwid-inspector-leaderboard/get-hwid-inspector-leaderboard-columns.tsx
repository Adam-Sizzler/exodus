import { ActionIcon, Group, Text } from '@mantine/core'
import { GetTopUsersByHwidDevicesCommand } from '@exodus/backend-contract'
import { TFunction } from 'i18next'
import { DataTableColumn } from 'mantine-datatable'
import { PiUserCircle } from 'react-icons/pi'
import { TbSum, TbUser } from 'react-icons/tb'

export function getHwidInspectorLeaderboardColumns(
    t: TFunction,
    handleViewUser: (userUuid: string) => void
): DataTableColumn<GetTopUsersByHwidDevicesCommand.Response['response']['users'][number]>[] {
    return [
        {
            accessor: 'actions',
            draggable: false,

            titleStyle: {
                backgroundColor: 'var(--mantine-color-dark-7)'
            },
            cellsStyle: () => {
                return {
                    backgroundColor: 'var(--mantine-color-dark-7)'
                }
            },
            title: (
                <Group c="dimmed" gap={4} justify="flex-end" pr={4} wrap="nowrap">
                    <TbUser size={18} />
                </Group>
            ),
            width: '0%',
            textAlign: 'right',
            render: ({ userUuid }) => (
                <Group gap={4} justify="flex-end" wrap="nowrap">
                    <ActionIcon
                        onClick={() => handleViewUser(userUuid)}
                        size="input-sm"
                        variant="soft"
                    >
                        <PiUserCircle size="1.5rem" />
                    </ActionIcon>
                </Group>
            )
        },
        {
            accessor: 'username',
            title: t('detailed-user-info-drawer.widget.username'),
            render: ({ username }) => username
        },

        {
            accessor: 'userUuid',
            title: t('use-hwid-inspector-table-columns.user-uuid'),
            render: ({ userUuid }) => userUuid
        },
        {
            accessor: 'id',
            title: 'ID',
            render: ({ id }) => id
        },
        {
            accessor: 'devicesCount',
            title: <TbSum size={18} />,
            render: ({ devicesCount }) => (
                <Text c="white" fw={600} size="sm">
                    {devicesCount}
                </Text>
            )
        }
    ]
}
