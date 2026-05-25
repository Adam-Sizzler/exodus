import { GetConfigProfilesCommand } from '@exodus/backend-contract'
import { ActionIcon, Avatar, Badge, Group, Text } from '@mantine/core'
import { DataTableColumn } from 'mantine-datatable'
import ReactCountryFlag from 'react-country-flag'
import { PiUsersDuotone } from 'react-icons/pi'
import { TbEdit } from 'react-icons/tb'
import { TFunction } from 'i18next'
import sortBy from 'lodash/sortBy'

import { prettyBytesUtil, prettySiBytesUtil, prettySiRealtimeBytesUtil } from '@shared/utils/bytes'
import { getSingboxUptimeUtil, formatDurationUtil } from '@shared/utils/time-utils'
import { NodePluginResponse, NodeResponse } from '@shared/api/hooks'
import { faviconResolver } from '@shared/utils/misc'

import { NodeStatusSimplfiedBadgeWidget } from '../node-status-simplfied-badge'

export function getNodesTableColumns(
    t: TFunction,
    configProfiles: GetConfigProfilesCommand.Response['response']['configProfiles'],
    nodePlugins: NodePluginResponse[],
    handleViewNode: (nodeUuid: string) => void
): DataTableColumn<NodeResponse>[] {
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
            render: ({ uuid }) => (
                <Group gap={4} justify="flex-end" wrap="nowrap">
                    <ActionIcon
                        color="teal"
                        onClick={() => handleViewNode(uuid)}
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
            render: ({ isConnected, isConnecting, isDisabled, uuid }) => (
                <NodeStatusSimplfiedBadgeWidget
                    isConnected={isConnected}
                    isConnecting={isConnecting}
                    isDisabled={isDisabled}
                    nodeUuid={uuid}
                />
            )
        },

        {
            accessor: 'usersOnline',
            title: t('use-nodes-table-widget.online'),
            render: ({ usersOnline }) => (
                <Badge
                    color={(usersOnline ?? 0) > 0 ? 'teal' : 'gray'}
                    leftSection={<PiUsersDuotone size={14} />}
                    miw="10ch"
                    size="lg"
                    variant="outline"
                >
                    {usersOnline}
                </Badge>
            )
        },

        {
            accessor: 'name',
            title: t('use-nodes-table-widget.name'),
            render: ({ name, countryCode }) => (
                <Group gap={6} wrap="nowrap">
                    {countryCode && countryCode !== 'XX' && (
                        <ReactCountryFlag
                            countryCode={countryCode}
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
                        {name}
                    </Text>
                </Group>
            )
        },
        {
            accessor: 'address',
            title: t('use-nodes-table-widget.address'),
            render: ({ address, port }) => `${address}:${port}`
        },
        {
            accessor: 'trafficUsedBytes',
            title: t('use-nodes-table-widget.traffic-used'),
            render: ({ trafficUsedBytes }) => prettyBytesUtil(trafficUsedBytes, false)
        },
        {
            accessor: 'configProfile.activeConfigProfileUuid',
            title: t('use-nodes-table-widget.config-profile'),
            render: ({ configProfile: { activeConfigProfileUuid } }) =>
                configProfiles.find((profile) => profile.uuid === activeConfigProfileUuid)?.name
        },
        {
            accessor: 'configProfile.activeInbounds',
            title: t('use-nodes-table-widget.inbounds'),
            render: ({ configProfile: { activeInbounds } }) =>
                sortBy(activeInbounds, 'tag')
                    .map((inbound) => inbound.tag)
                    .join(', ')
        },
        {
            accessor: 'singboxVersion',
            title: t('use-nodes-table-widget.singbox-v'),
            render: ({ versions, singboxVersion }) => versions?.singbox ?? singboxVersion ?? '-'
        },
        {
            accessor: 'nodeVersion',
            title: t('use-nodes-table-widget.node-v'),
            render: ({ versions, nodeVersion }) => versions?.node ?? nodeVersion ?? '-'
        },
        {
            accessor: 'singboxUptime',
            title: 'Sing-box Uptime',
            render: ({ singboxUptime }) =>
                singboxUptime && singboxUptime !== '0' ? getSingboxUptimeUtil(singboxUptime) : '-'
        },
        {
            accessor: 'provider.name',
            title: t('use-nodes-table-widget.provider'),
            render: ({ provider }) =>
                provider ? (
                    <Group gap="xs" wrap="nowrap">
                        <Avatar
                            alt={provider.name}
                            color="initials"
                            name={provider.name}
                            onLoad={(event) => {
                                const img = event.target as HTMLImageElement
                                if (img.naturalWidth <= 16 && img.naturalHeight <= 16) {
                                    img.src = ''
                                }
                            }}
                            radius="sm"
                            size={16}
                            src={faviconResolver(provider.faviconLink)}
                        />

                        <Text size="sm">{provider.name}</Text>
                    </Group>
                ) : null
        },
        {
            accessor: 'tags',
            title: t('use-nodes-table-widget.tags'),
            render: ({ tags }) => tags?.join(', ') ?? '-'
        },
        {
            accessor: 'activePluginUuid',
            title: 'Plugin',
            render: ({ activePluginUuid }) =>
                nodePlugins.find((plugin) => plugin.uuid === activePluginUuid)?.name ?? '-'
        },
        {
            accessor: 'system.info.cpus',
            title: 'CPU Cores',
            render: ({ system, cpuCount }) => system?.info.cpus ?? cpuCount ?? '-'
        },
        {
            accessor: 'system.stats.memoryFree',
            title: 'Free RAM',
            render: ({ system }) => (system ? prettyBytesUtil(system.stats.memoryFree, false) : '-')
        },
        {
            accessor: 'system.stats.memoryUsed',
            title: 'Used RAM',
            render: ({ system }) => (system ? prettyBytesUtil(system.stats.memoryUsed, false) : '-')
        },
        {
            accessor: 'totalRam',
            title: t('use-nodes-table-widget.total-ram'),
            render: ({ system, totalRam }) =>
                system ? prettyBytesUtil(system.info.memoryTotal, false) : (totalRam ?? '-')
        },
        {
            accessor: 'cpuModel',
            title: t('use-nodes-table-widget.cpu-model'),
            render: ({ system, cpuModel }) => system?.info.cpuModel ?? cpuModel ?? '-'
        },
        {
            accessor: 'system.stats.uptime',
            title: 'Server Uptime',
            render: ({ system }) => (system ? formatDurationUtil(system.stats.uptime) : '-')
        },
        {
            accessor: 'system.info.networkInterfaces',
            title: 'Network Interfaces',
            render: ({ system }) => (system ? system.info.networkInterfaces.join(', ') : '-')
        },
        {
            accessor: 'system.stats.interface.rxBytesPerSec',
            title: 'RX Speed',
            render: ({ system }) =>
                system?.stats.interface
                    ? prettySiRealtimeBytesUtil(system.stats.interface.rxBytesPerSec, true, true)
                    : '-'
        },
        {
            accessor: 'system.stats.interface.txBytesPerSec',
            title: 'TX Speed',
            render: ({ system }) =>
                system?.stats.interface
                    ? prettySiRealtimeBytesUtil(system.stats.interface.txBytesPerSec, true, true)
                    : '-'
        },
        {
            accessor: 'system.stats.interface.rxTotal',
            title: 'RX Total',
            render: ({ system }) =>
                system?.stats.interface
                    ? prettySiBytesUtil(system.stats.interface.rxTotal, true)
                    : '-'
        },
        {
            accessor: 'system.stats.interface.txTotal',
            title: 'TX Total',
            render: ({ system }) =>
                system?.stats.interface
                    ? prettySiBytesUtil(system.stats.interface.txTotal, true)
                    : '-'
        },
        {
            accessor: 'system.info.release',
            title: 'OS Release',
            render: ({ system }) => system?.info.release ?? '-'
        }
    ]
}
