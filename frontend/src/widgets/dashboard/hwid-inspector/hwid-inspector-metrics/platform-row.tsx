import { Box, Group, Paper, Progress, Text } from '@mantine/core'

import { AppBreakdown } from './app-breakdown'
import { IPlatformDatum } from './hwid-inspector-metrics.types'
import { formatCount, formatShare } from './hwid-inspector-metrics.utils'
import { openPlatformAppsModal } from './platform-apps-modal'
import classes from './platform-row.module.css'

export function PlatformRow({
    platform,
    total,
    isMobile
}: {
    platform: IPlatformDatum
    total: number
    isMobile: boolean
}) {
    return (
        <Paper
            className={classes.item}
            onClick={() => openPlatformAppsModal(platform, isMobile)}
            p="sm"
            radius="md"
        >
            <Group gap="xs" justify="space-between" mb={8} wrap="nowrap">
                <Group gap="xs" style={{ minWidth: 0 }} wrap="nowrap">
                    <Box
                        bg={platform.color}
                        h={12}
                        style={{ borderRadius: 3, flexShrink: 0 }}
                        w={12}
                    />
                    <Text fw={600} lineClamp={1} size="sm" title={platform.name}>
                        {platform.name}
                    </Text>
                </Group>
                <Group align="baseline" gap={6} wrap="nowrap">
                    <Text fw={700} size="sm">
                        {formatCount(platform.value)}
                    </Text>
                    <Text c="dimmed" size="xs">
                        {formatShare(platform.value, total)}
                    </Text>
                </Group>
            </Group>

            <Progress
                color={platform.color}
                mb={10}
                radius="xl"
                size={6}
                value={total > 0 ? (platform.value / total) * 100 : 0}
            />

            <AppBreakdown apps={platform.apps} />
        </Paper>
    )
}
