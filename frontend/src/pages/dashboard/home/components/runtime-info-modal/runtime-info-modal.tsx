import { Code, SimpleGrid, Stack, Text } from '@mantine/core'
import { Trans, useTranslation } from 'react-i18next'
import {
    TbBolt,
    TbBraces,
    TbCloud,
    TbCpu,
    TbGauge,
    TbPackage,
    TbPlug,
    TbServer,
    TbStack3,
    TbStopwatch
} from 'react-icons/tb'

import { MetricEntry } from './metric-entry'
import { ProcessCard } from './process-card'
import classes from './runtime-info-modal.module.css'
import { SectionShell } from './section-shell'

const TRANS_COMPONENTS = { highlight: <Code className={classes.highlight} /> }

export function RuntimeInfoModalContent() {
    const { t } = useTranslation()

    return (
        <Stack gap="lg">
            <Text c="dimmed" lh={1.6} size="sm">
                <Trans components={TRANS_COMPONENTS} i18nKey="home.runtime-info.intro" />
            </Text>

            <SectionShell
                accent="rgba(59, 130, 246, 0.25)"
                color="blue"
                description={t('home.runtime-info.instance-types-description')}
                Icon={TbCpu}
                title={t('home.runtime-info.instance-types-title')}
            >
                <SimpleGrid cols={1} spacing="sm">
                    <ProcessCard
                        color="blue"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.api-description"
                            />
                        }
                        Icon={TbCloud}
                        title="API"
                    />
                </SimpleGrid>
            </SectionShell>

            <SectionShell
                accent="rgba(45, 212, 191, 0.25)"
                color="teal"
                description={t('home.runtime-info.memory-description')}
                Icon={TbServer}
                title="Memory"
            >
                <Stack gap="md">
                    <MetricEntry
                        color="teal"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.heap-description"
                            />
                        }
                        Icon={TbStack3}
                        title="Heap"
                    />
                    <MetricEntry
                        color="cyan"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.rss-description"
                            />
                        }
                        Icon={TbServer}
                        title="RSS"
                    />
                    <MetricEntry
                        color="lime"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.external-description"
                            />
                        }
                        Icon={TbPackage}
                        title="Go Sys"
                    />
                    <MetricEntry
                        color="green"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.array-buffers-description"
                            />
                        }
                        Icon={TbBraces}
                        title="Stack"
                    />
                </Stack>
            </SectionShell>

            <SectionShell
                accent="rgba(139, 92, 246, 0.25)"
                color="violet"
                description={t('home.runtime-info.event-loop-description')}
                Icon={TbBolt}
                title="Scheduler"
            >
                <Stack gap="md">
                    <MetricEntry
                        color="violet"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.delay-description"
                            />
                        }
                        Icon={TbStopwatch}
                        title="Delay"
                    />
                    <MetricEntry
                        color="grape"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.p99-description"
                            />
                        }
                        Icon={TbGauge}
                        title="P99"
                    />
                    <MetricEntry
                        color="indigo"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.active-handles-description"
                            />
                        }
                        Icon={TbPlug}
                        title="Goroutines"
                    />
                    <MetricEntry
                        color="blue"
                        description={
                            <Trans
                                components={TRANS_COMPONENTS}
                                i18nKey="home.runtime-info.threads-description"
                            />
                        }
                        Icon={TbCpu}
                        title="Threads"
                    />
                </Stack>
            </SectionShell>
        </Stack>
    )
}
