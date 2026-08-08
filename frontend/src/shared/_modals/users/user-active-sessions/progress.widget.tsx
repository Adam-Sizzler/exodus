import { Box, Group, Progress, Stack, Text } from '@mantine/core'
import { useTranslation } from 'react-i18next'

import { LottieGlobeShared } from '@shared/ui/lotties/globe'
import { SectionCard } from '@shared/ui/section-card'

interface IProps {
    completed: number
    percent: number
    total: number
}

export const ProgressWidget = ({ completed, percent, total }: IProps) => {
    const { t } = useTranslation()

    return (
        <SectionCard.Root gap="md">
            <SectionCard.Section>
                <Stack align="center" gap="md" py="xl">
                    <div style={{ height: 120, display: 'flex', alignItems: 'center' }}>
                        <LottieGlobeShared />
                    </div>

                    <Stack align="center" gap={4}>
                        <Text c="white" fw={600} size="md">
                            {t('active-sessions-drawer.widget.fetching')}
                        </Text>
                        <Text c="dimmed" size="sm">
                            {t('active-sessions-drawer.widget.progress', { completed, total })}
                        </Text>
                    </Stack>

                    <Box mt={6} w="100%">
                        <Progress.Root
                            radius="md"
                            size={18}
                            style={{
                                backgroundColor: 'rgba(255, 255, 255, 0.05)'
                            }}
                        >
                            <Progress.Section
                                animated
                                color="teal"
                                style={{
                                    transition: 'width 400ms ease'
                                }}
                                value={percent}
                            />
                        </Progress.Root>
                        <Group justify="center" mt={6}>
                            <Text c="white" ff="monospace" fw={700} size="sm">
                                {Math.round(percent)}%
                            </Text>
                        </Group>
                    </Box>
                </Stack>
            </SectionCard.Section>
        </SectionCard.Root>
    )
}
