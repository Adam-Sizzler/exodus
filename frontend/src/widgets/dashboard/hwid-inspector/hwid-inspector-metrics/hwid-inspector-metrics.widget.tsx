import { Card, Center, Group, Loader, SimpleGrid, Stack, Text, ThemeIcon } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import {
    PiCalculatorDuotone,
    PiChartPieDuotone,
    PiDeviceMobileDuotone,
    PiDevicesDuotone
} from 'react-icons/pi'
import { TbChartArcs3 } from 'react-icons/tb'

import { useGetHwidDevicesStats } from '@shared/api/hooks'
import { useIsMobile } from '@shared/hooks'
import { IMetricCardProps, MetricCardShared } from '@shared/ui/metrics/metric-card'

import { buildPlatformData } from './hwid-inspector-metrics.utils'
import { PlatformRow } from './platform-row'

export function HwidInspectorMetrics() {
    const { t } = useTranslation()

    const { data: stats, isLoading } = useGetHwidDevicesStats()
    const isMobile = useIsMobile()

    const platformData = stats?.byPlatform ? buildPlatformData(stats.byPlatform) : []
    const total = platformData.reduce((sum, item) => sum + item.value, 0)

    const metricCards: IMetricCardProps[] = [
        {
            IconComponent: PiDevicesDuotone,
            title: t('hwid-inspector-metrics.widget.total-unique-devices'),
            value: stats?.stats.totalUniqueDevices ?? 0,
            iconVariant: 'soft',
            iconColor: 'blue'
        },
        {
            IconComponent: PiDeviceMobileDuotone,
            title: t('hwid-inspector-metrics.widget.total-hwid-devices'),
            value: stats?.stats.totalHwidDevices ?? 0,
            iconVariant: 'soft',
            iconColor: 'teal'
        },
        {
            IconComponent: PiCalculatorDuotone,
            title: t('hwid-inspector-metrics.widget.avg-devices-per-user'),
            value: stats?.stats.averageHwidDevicesPerUser ?? 0,
            iconVariant: 'soft',
            iconColor: 'orange'
        }
    ]

    return (
        <Stack gap="md" mb={0}>
            <SimpleGrid cols={{ base: 1, sm: 2, xl: 3 }}>
                {metricCards.map((card) => (
                    <MetricCardShared key={card.title} {...card} />
                ))}
            </SimpleGrid>

            <Card
                p="lg"
                radius="md"
                style={{
                    border: '1px solid rgb(255, 255, 255, 0.08)',
                    background: 'rgb(255, 255, 255, 0.02)'
                }}
                withBorder
            >
                <Group align="center" gap="sm" justify="space-between" mb="lg" wrap="nowrap">
                    <Group align="center" gap="sm" style={{ minWidth: 0 }} wrap="nowrap">
                        <ThemeIcon color="cyan" radius="md" size="lg" variant="soft">
                            <TbChartArcs3 size="20px" />
                        </ThemeIcon>
                        <Text fw={600} lineClamp={1} size="lg">
                            {t('hwid-inspector-metrics.widget.platform-distribution')}
                        </Text>
                    </Group>
                </Group>

                {isLoading ? (
                    <Center mih={280}>
                        <Loader size="lg" />
                    </Center>
                ) : platformData.length === 0 ? (
                    <Center mih={280}>
                        <Stack align="center" gap="sm">
                            <PiChartPieDuotone color="var(--mantine-color-gray-6)" size="48px" />
                            <Text c="dimmed" size="sm">
                                {t('hwid-inspector-metrics.widget.no-platform-data')}
                            </Text>
                        </Stack>
                    </Center>
                ) : (
                    <SimpleGrid cols={{ base: 1, sm: 2, xl: 3 }} spacing="sm">
                        {platformData.map((platform) => (
                            <PlatformRow
                                key={platform.name}
                                platform={platform}
                                total={total}
                                isMobile={isMobile}
                            />
                        ))}
                    </SimpleGrid>
                )}
            </Card>
        </Stack>
    )
}
