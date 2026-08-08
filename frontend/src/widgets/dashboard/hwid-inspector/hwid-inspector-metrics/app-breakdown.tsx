import { Group, Text } from '@mantine/core'

import { IPlatformApp } from './hwid-inspector-metrics.types'
import { formatCount } from './hwid-inspector-metrics.utils'

export function AppBreakdown({ apps }: { apps: IPlatformApp[] }) {
    if (apps.length === 0) return null

    return (
        <Group gap="md" wrap="wrap">
            {apps.map((item) => (
                <Group gap={6} key={item.app} wrap="nowrap">
                    <Text c="dimmed" lineClamp={1} size="xs" title={item.app}>
                        {item.app}
                    </Text>
                    <Text c="dimmed" fw={600} size="xs" style={{ opacity: 0.6 }}>
                        {formatCount(item.count)}
                    </Text>
                </Group>
            ))}
        </Group>
    )
}
