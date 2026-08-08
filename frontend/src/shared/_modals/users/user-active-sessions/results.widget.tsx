import { Center, Stack, Text } from '@mantine/core'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiEmptyDuotone } from 'react-icons/pi'

import { SectionCard } from '@shared/ui/section-card'

import { IpStatsWidget } from './ip-stats.widget'
import { NodeSessionsCardWidget } from './node-sessions-card.widget'
import { SummaryCardWidget } from './summary-card.widget'
import { ActiveSessionNode } from './use-user-active-sessions'

interface IProps {
    nodes: ActiveSessionNode[]
    onDropAll: () => void
    onDropIp: (ip: string, nodeUuid: string) => void
    onDropNode: (nodeUuid: string) => void
    onRefresh: () => void
}

export const ResultsWidget = ({ nodes, onDropAll, onDropIp, onDropNode, onRefresh }: IProps) => {
    const { t } = useTranslation()

    const { distinctIps, totalIps } = useMemo(() => {
        const allIps = nodes.flatMap((node) => node.ips.map((item) => item.ip))
        return { distinctIps: new Set(allIps).size, totalIps: allIps.length }
    }, [nodes])

    const hasNodes = nodes.length > 0

    return (
        <Stack gap="md">
            <SummaryCardWidget
                distinctIps={distinctIps}
                ipStats={hasNodes && <IpStatsWidget nodes={nodes} />}
                onDropAll={onDropAll}
                onRefresh={onRefresh}
                totalIps={totalIps}
            />

            {!hasNodes && (
                <SectionCard.Root gap="sm">
                    <SectionCard.Section>
                        <Center h="230">
                            <Stack align="center" gap="xs">
                                <PiEmptyDuotone color="var(--mantine-color-gray-5)" size="3rem" />
                                <Text c="dimmed" size="sm">
                                    {t('active-sessions-drawer.widget.no-active-sessions')}
                                </Text>
                            </Stack>
                        </Center>
                    </SectionCard.Section>
                </SectionCard.Root>
            )}

            {nodes.map((node) => (
                <NodeSessionsCardWidget
                    key={node.nodeUuid}
                    node={node}
                    onDropIp={(ip) => onDropIp(ip, node.nodeUuid)}
                    onDropNode={() => onDropNode(node.nodeUuid)}
                />
            ))}
        </Stack>
    )
}
