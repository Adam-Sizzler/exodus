import { PiClockDuotone, PiCloudDuotone, PiGearSixDuotone, PiQueueDuotone } from 'react-icons/pi'
import { Badge, Card, Group, Progress, SimpleGrid, Stack, Text, ThemeIcon } from '@mantine/core'
import { TFunction } from 'i18next'

import {
    ExodusMetric,
    formatBytes,
    formatMilliseconds,
    getHeapUsedPercent,
    getHeapUsedValue
} from '../metrics/exodus-metrics'

import classes from './exodus-runtime-detail-card.module.css'

const PROCESS_CONFIG: Record<
    string,
    {
        accentColor: string
        accentGlow: string
        color: string
        Icon: typeof PiGearSixDuotone
        name: string
    }
> = {
    api: {
        accentColor: 'var(--mantine-color-blue-6)',
        accentGlow: 'rgba(59, 130, 246, 0.2)',
        color: 'blue',
        Icon: PiCloudDuotone,
        name: 'API'
    },
    scheduler: {
        accentColor: 'var(--mantine-color-violet-6)',
        accentGlow: 'rgba(139, 92, 246, 0.2)',
        color: 'violet',
        Icon: PiClockDuotone,
        name: 'Scheduler'
    },
    worker: {
        accentColor: 'var(--mantine-color-teal-6)',
        accentGlow: 'rgba(20, 184, 166, 0.2)',
        color: 'teal',
        Icon: PiQueueDuotone,
        name: 'Workers'
    },
    workers: {
        accentColor: 'var(--mantine-color-teal-6)',
        accentGlow: 'rgba(20, 184, 166, 0.2)',
        color: 'teal',
        Icon: PiQueueDuotone,
        name: 'Workers'
    },
    processor: {
        accentColor: 'var(--mantine-color-teal-6)',
        accentGlow: 'rgba(20, 184, 166, 0.2)',
        color: 'teal',
        Icon: PiQueueDuotone,
        name: 'Workers'
    }
}

const DEFAULT_PROCESS = {
    accentColor: 'var(--mantine-color-gray-6)',
    accentGlow: 'rgba(108, 117, 125, 0.2)',
    color: 'gray',
    Icon: PiGearSixDuotone,
    name: 'Runtime'
}

interface IProps {
    metric: ExodusMetric
    t: TFunction
}

const RuntimeValue = ({ label, value }: { label: string; value: number | string }) => (
    <Stack gap={0}>
        <Text className={classes.statLabel}>{label}</Text>
        <Text className={classes.statValue}>{value}</Text>
    </Stack>
)

export function ExodusRuntimeDetailCard({ metric, t }: IProps) {
    const instanceType = metric.instanceType?.toLowerCase() ?? 'api'
    const config = PROCESS_CONFIG[instanceType] ?? DEFAULT_PROCESS
    const heapPercent = getHeapUsedPercent(metric)
    const scheduler = metric.scheduler as
        | (ExodusMetric['scheduler'] & { schedulerDelayMs?: number })
        | undefined

    return (
        <Card
            className={classes.card}
            padding="md"
            radius="md"
            style={
                {
                    '--accent-color': config.accentColor,
                    '--accent-glow': config.accentGlow
                } as React.CSSProperties
            }
        >
            <div className={classes.topAccent} />

            <Stack gap="sm">
                <Group justify="space-between" wrap="nowrap">
                    <Group gap="xs" wrap="nowrap">
                        <ThemeIcon color={config.color} radius="md" size="lg" variant="soft">
                            <config.Icon size={18} />
                        </ThemeIcon>
                        <Stack gap={2}>
                            <Text c="white" ff="monospace" fw={700} lh={1} size="sm">
                                {config.name}-{metric.instanceId ?? 0}
                            </Text>
                            {typeof metric.pid === 'number' && (
                                <Badge color={config.color} ff="monospace" size="xs" variant="soft">
                                    PID: {metric.pid}
                                </Badge>
                            )}
                        </Stack>
                    </Group>
                </Group>

                <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="sm">
                    <div className={classes.memorySection}>
                        <Stack gap={6}>
                            <Text c="dimmed" fw={600} lh={1} lts={1} size="10px" tt="uppercase">
                                {t('runtime-metrics.memory')}
                            </Text>

                            <div>
                                <Group gap="xs" justify="space-between" mb={4} wrap="nowrap">
                                    <Text className={classes.statLabel}>
                                        {t('runtime-metrics.heap')}
                                    </Text>
                                    <Text className={classes.heapValues}>
                                        {getHeapUsedValue(metric)}{' '}
                                        <span className={classes.heapPercent}>
                                            ({heapPercent.toFixed(0)}%)
                                        </span>
                                    </Text>
                                </Group>
                                <Progress
                                    color={config.color}
                                    radius="xl"
                                    size="sm"
                                    value={heapPercent}
                                />
                            </div>

                            <Group grow>
                                <RuntimeValue
                                    label={t('runtime-metrics.rss')}
                                    value={formatBytes(metric.memory?.rssBytes)}
                                />
                                <RuntimeValue
                                    label={t('runtime-metrics.go-sys')}
                                    value={formatBytes(metric.memory?.sysBytes)}
                                />
                                <RuntimeValue
                                    label={t('runtime-metrics.stack')}
                                    value={formatBytes(metric.memory?.stackInuseBytes)}
                                />
                            </Group>
                        </Stack>
                    </div>

                    <div className={classes.perfSection}>
                        <Stack gap={6}>
                            <Text c="dimmed" fw={600} lh={1} lts={1} size="10px" tt="uppercase">
                                {t('runtime-metrics.scheduler')}
                            </Text>

                            <Stack gap={6} hiddenFrom="sm">
                                <Group grow>
                                    <RuntimeValue
                                        label={t('runtime-metrics.delay')}
                                        value={formatMilliseconds(scheduler?.schedulerDelayMs)}
                                    />
                                    <RuntimeValue
                                        label={t('runtime-metrics.p99')}
                                        value={formatMilliseconds(scheduler?.schedulerP99Ms)}
                                    />
                                </Group>
                                <Group grow>
                                    <RuntimeValue
                                        label={t('runtime-metrics.goroutines')}
                                        value={scheduler?.goroutines ?? 0}
                                    />
                                    <RuntimeValue
                                        label={t('runtime-metrics.threads')}
                                        value={metric.process?.threads ?? 0}
                                    />
                                </Group>
                            </Stack>

                            <Stack gap={6} visibleFrom="sm">
                                <Group grow>
                                    <RuntimeValue
                                        label={t('runtime-metrics.delay')}
                                        value={formatMilliseconds(scheduler?.schedulerDelayMs)}
                                    />
                                    <RuntimeValue
                                        label={t('runtime-metrics.p99')}
                                        value={formatMilliseconds(scheduler?.schedulerP99Ms)}
                                    />
                                </Group>
                                <Group grow>
                                    <RuntimeValue
                                        label={t('runtime-metrics.goroutines')}
                                        value={scheduler?.goroutines ?? 0}
                                    />
                                    <RuntimeValue
                                        label={t('runtime-metrics.threads')}
                                        value={metric.process?.threads ?? 0}
                                    />
                                </Group>
                            </Stack>
                        </Stack>
                    </div>
                </SimpleGrid>
            </Stack>
        </Card>
    )
}
