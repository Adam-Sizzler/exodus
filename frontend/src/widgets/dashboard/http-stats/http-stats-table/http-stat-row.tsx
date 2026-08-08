import { Badge, Group, RollingNumber, Text } from '@mantine/core'
import { GetHttpStatsCommand } from '@exodus/backend-contract'
import { memo } from 'react'

import classes from './http-stats-table.module.css'
import { getMethodColor } from './http-stats.utils'

interface HttpStatRowProps {
    percentOfTotal: number
    route: GetHttpStatsCommand.Response['response']['routes'][number]
}

function HttpStatRowBase({ route, percentOfTotal }: HttpStatRowProps) {
    const color = getMethodColor(route.method)

    return (
        <div className={classes.row}>
            <div
                className={classes.fill}
                style={{
                    background: `var(--mantine-color-${color}-6)`,
                    width: `${Math.max(percentOfTotal, 1.5)}%`
                }}
            />

            <Group className={classes.content} gap="sm" justify="space-between" wrap="nowrap">
                <Group gap="sm" miw={0} wrap="nowrap">
                    <Badge
                        className={classes.methodBadge}
                        color={color}
                        radius="sm"
                        size="lg"
                        variant="soft"
                    >
                        {route.method}
                    </Badge>
                    <Text className={classes.path} title={route.route} truncate="end">
                        {route.route}
                    </Text>
                </Group>

                <Group align="baseline" gap="xs" wrap="nowrap">
                    <Text className={classes.percent} c="dimmed" size="xs">
                        {percentOfTotal.toFixed(1)}%
                    </Text>
                    <RollingNumber value={route.count} thousandSeparator fw={700} />
                </Group>
            </Group>
        </div>
    )
}

export const HttpStatRow = memo(HttpStatRowBase)
