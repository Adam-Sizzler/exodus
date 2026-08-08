import { Group, Paper, Text } from '@mantine/core'
import { useMemo } from 'react'
import { TbClockCheck, TbClockExclamation, TbClockPause } from 'react-icons/tb'

import { formatInt } from '@shared/utils/misc'

import { ActiveSessionUser } from './use-node-active-sessions'

interface IProps {
    users: ActiveSessionUser[]
}

export const IpStatsWidget = ({ users }: IProps) => {
    const stats = useMemo(() => {
        const now = new Date()
        let active = 0
        let idle = 0
        let stale = 0

        for (const user of users) {
            for (const ip of user.ips) {
                const diffMinutes = (now.getTime() - new Date(ip.lastSeen).getTime()) / 60_000
                if (diffMinutes <= 5) active++
                else if (diffMinutes <= 60) idle++
                else stale++
            }
        }

        return { active, idle, stale }
    }, [users])

    return (
        <Group gap="xs" grow>
            <Paper
                bd="1px solid rgba(45, 212, 191, 0.2)"
                bg="rgba(45, 212, 191, 0.08)"
                p="xs"
                radius="md"
            >
                <Group gap={4} justify="center">
                    <TbClockCheck color="var(--mantine-color-teal-6)" size={18} />
                    <Text c="teal.6" fw={700} size="sm">
                        {formatInt(stats.active)}
                    </Text>
                </Group>
            </Paper>

            <Paper
                bd="1px solid rgba(251, 191, 36, 0.2)"
                bg="rgba(251, 191, 36, 0.08)"
                p="xs"
                radius="md"
            >
                <Group gap={4} justify="center">
                    <TbClockPause color="var(--mantine-color-yellow-6)" size={18} />
                    <Text c="yellow.6" fw={700} size="sm">
                        {formatInt(stats.idle)}
                    </Text>
                </Group>
            </Paper>

            <Paper
                bd="1px solid rgba(239, 68, 68, 0.2)"
                bg="rgba(239, 68, 68, 0.08)"
                p="xs"
                radius="md"
            >
                <Group gap={4} justify="center">
                    <TbClockExclamation color="var(--mantine-color-red-6)" size={18} />
                    <Text c="red.6" fw={700} size="sm">
                        {formatInt(stats.stale)}
                    </Text>
                </Group>
            </Paper>
        </Group>
    )
}
