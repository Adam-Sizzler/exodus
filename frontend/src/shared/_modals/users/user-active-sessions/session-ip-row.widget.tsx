import { ActionIcon, Box, Group, Stack, Text, Tooltip } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import {
    TbClockCheck,
    TbClockExclamation,
    TbClockPause,
    TbExternalLink,
    TbUnlink
} from 'react-icons/tb'

import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'
import { formatRelativeDateUtil, formatTimeUtil } from '@shared/utils/time-utils'

import { ActiveSessionNode } from './use-user-active-sessions'

const getLastSeenIndicator = (lastSeen: Date | string) => {
    const diffMs = Date.now() - new Date(lastSeen).getTime()
    const diffMinutes = diffMs / 60_000
    if (diffMinutes <= 5) return { color: 'var(--mantine-color-teal-6)', Icon: TbClockCheck }
    if (diffMinutes <= 60) return { color: 'var(--mantine-color-yellow-6)', Icon: TbClockPause }
    return { color: 'var(--mantine-color-red-6)', Icon: TbClockExclamation }
}

interface IProps {
    ip: ActiveSessionNode['ips'][number]
    onDrop: () => void
}

export const SessionIpRowWidget = ({ ip, onDrop }: IProps) => {
    const { t, i18n } = useTranslation()

    const { color, Icon } = getLastSeenIndicator(ip.lastSeen)

    return (
        <Group align="center" gap="xs" wrap="nowrap">
            <ActionIcon
                color="cyan"
                component="a"
                href={`https://ipinfo.io/${ip.ip}`}
                rel="noopener noreferrer"
                size="input-sm"
                target="_blank"
                variant="soft"
            >
                <TbExternalLink size={18} />
            </ActionIcon>

            <Box style={{ flex: 1 }}>
                <CopyableFieldShared
                    leftSection={
                        <Tooltip
                            color="dark.6"
                            label={
                                <Stack gap={2} p={4}>
                                    <Text c="white" fw={600} size="xs">
                                        {formatRelativeDateUtil(ip.lastSeen, t, i18n.language)}
                                    </Text>
                                    <Text c="dimmed" ff="monospace" size="xs">
                                        {formatTimeUtil({
                                            time: ip.lastSeen,
                                            template: 'TIME_FIRST_DATETIME',
                                            language: i18n.language
                                        })}
                                    </Text>
                                </Stack>
                            }
                            radius="md"
                            styles={{
                                tooltip: {
                                    border: '1px solid var(--mantine-color-dark-4)',
                                    backdropFilter: 'blur(8px)'
                                }
                            }}
                        >
                            <Box style={{ display: 'flex', cursor: 'help', color }}>
                                <Icon size={16} />
                            </Box>
                        </Tooltip>
                    }
                    size="sm"
                    value={ip.ip}
                />
            </Box>

            <Tooltip
                label={t('user-active-session-drawer.widget.drop-this-connection-on-this-node')}
            >
                <ActionIcon color="orange" onClick={onDrop} size="lg" variant="soft">
                    <TbUnlink size={20} />
                </ActionIcon>
            </Tooltip>
        </Group>
    )
}
