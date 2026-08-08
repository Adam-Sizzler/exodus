import { Stack, Text } from '@mantine/core'
import { useTranslation } from 'react-i18next'

import { LottieGlobeShared } from '@shared/ui/lotties/globe'
import { SectionCard } from '@shared/ui/section-card'

export const ProgressWidget = () => {
    const { t } = useTranslation()

    return (
        <SectionCard.Root gap="md">
            <SectionCard.Section>
                <Stack align="center" gap="md" py="xl">
                    <div style={{ height: 120, display: 'flex', alignItems: 'center' }}>
                        <LottieGlobeShared />
                    </div>

                    <Text c="white" fw={600} size="md">
                        {t('active-sessions-drawer.widget.fetching')}
                    </Text>
                </Stack>
            </SectionCard.Section>
        </SectionCard.Root>
    )
}
