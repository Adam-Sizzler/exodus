import { GetConfigProfilesCommand } from '@exodus/backend-contract'
import { ActionIcon, Avatar, Badge, Group, Text } from '@mantine/core'
import { DataTableColumn } from '@kastov/mantine-datatable'
import ReactCountryFlag from 'react-country-flag'
import { TbEdit } from 'react-icons/tb'
import { TFunction } from 'i18next'
import sortBy from 'lodash/sortBy'

import { faviconResolver } from '@shared/utils/misc'
import { SubscriptionConnectionResponse } from '@shared/api/hooks'

import { NodeStatusSimplfiedBadgeWidget } from '../node-status-simplfied-badge'

export function getNodesTableColumns(
    t: TFunction,
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'],
    handleViewNode: (nodeUuid: string) => void
): DataTableColumn<SubscriptionConnectionResponse>[] {
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
                    <TbEdit size={18} />
                </Group>
            ),
            width: '0%',
            textAlign: 'right',
            render: (record) => (
                <Group gap={4} justify="flex-end" wrap="nowrap">
                    <ActionIcon
                        color="teal"
                        onClick={() => handleViewNode(record.uuid)}
                        size="md"
                        variant="outline"
                    >
                        <TbEdit size={18} />
                    </ActionIcon>
                </Group>
            )
        },
        {
            accessor: 'isConnected',
            title: '',
            render: (record) => (
                <NodeStatusSimplfiedBadgeWidget
                    isConnected={record.isConnected}
                    isConnecting={record.isConnecting}
                    isDisabled={record.isDisabled}
                    nodeUuid={record.uuid}
                />
            )
        },

        {
            accessor: 'name',
            title: t('use-nodes-table-widget.name'),
            render: (record) => (
                <Group gap={6} wrap="nowrap">
                    {record.countryCode && record.countryCode !== 'XX' && (
                        <ReactCountryFlag
                            countryCode={record.countryCode}
                            style={{
                                fontSize: '1.1em',
                                borderRadius: '2px'
                            }}
                        />
                    )}
                    <Text
                        size="sm"
                        style={{
                            maxWidth: '100%',
                            overflow: 'hidden',
                            textOverflow: 'ellipsis',
                            whiteSpace: 'nowrap'
                        }}
                    >
                        {record.name}
                    </Text>
                </Group>
            )
        },
        {
            accessor: 'address',
            title: t('use-nodes-table-widget.address'),
            render: (record) => `${record.address}:${record.port}`
        },
        {
            accessor: 'configProfile.activeConfigProfileUuid',
            title: t('use-nodes-table-widget.config-profile'),
            render: (record) =>
                configProfiles?.find((profile) => profile.uuid === record.configProfile?.activeConfigProfileUuid)?.name
        },
        {
            accessor: 'configProfile.activeInbounds',
            title: t('use-nodes-table-widget.inbounds'),
            render: (record) =>
                sortBy(record.configProfile?.activeInbounds || [], 'tag')
                    .map((inbound) => inbound.tag)
                    .join(', ')
        },
        {
            accessor: 'nodeVersion',
            title: t('use-nodes-table-widget.node-v')
        },
        {
            accessor: 'provider.name',
            title: t('use-nodes-table-widget.provider'),
            render: (record) =>
                record.provider ? (
                    <Group gap="xs" wrap="nowrap">
                        <Avatar
                            alt={record.provider.name}
                            color="initials"
                            name={record.provider.name}
                            onLoad={(event) => {
                                const img = event.target as HTMLImageElement
                                if (img.naturalWidth <= 16 && img.naturalHeight <= 16) {
                                    img.src = ''
                                }
                            }}
                            radius="sm"
                            size={16}
                            src={faviconResolver(record.provider.faviconLink)}
                        />

                        <Text size="sm">{record.provider.name}</Text>
                    </Group>
                ) : null
        },
        {
            accessor: 'tags',
            title: t('use-nodes-table-widget.tags'),
            render: (record) => record.tags?.join(', ') ?? '-'
        },
        {
            accessor: 'totalRam',
            title: t('use-nodes-table-widget.total-ram')
        },
        {
            accessor: 'cpuModel',
            title: t('use-nodes-table-widget.cpu-model')
        }
    ]
}
