import { Button, Group, Stack, Text, Title } from '@mantine/core'
import {
    IconCrownFilled,
    IconHeadset,
    IconMessageCircle2,
    IconRoute,
    IconUsersGroup
} from '@tabler/icons-react'
import { useTranslation } from 'react-i18next'

import classes from './PrimeModal.module.css'

const PRIME_LINK = 'https://docs.exodus.dev/prime'

const FEATURES = [
    {
        icon: IconHeadset,
        titleKey: 'prime-modal.shared.priority-support-title',
        descriptionKey: 'prime-modal.shared.priority-support-description'
    },
    {
        icon: IconUsersGroup,
        titleKey: 'prime-modal.shared.private-community-title',
        descriptionKey: 'prime-modal.shared.private-community-description'
    },
    {
        icon: IconRoute,
        titleKey: 'prime-modal.shared.shape-the-roadmap-title',
        descriptionKey: 'prime-modal.shared.shape-the-roadmap-description'
    },
    {
        icon: IconMessageCircle2,
        titleKey: 'prime-modal.shared.direct-developer-access-title',
        descriptionKey: 'prime-modal.shared.direct-developer-access-description'
    }
] as const

export function PrimeModalContent() {
    const { t } = useTranslation()

    return (
        <Stack gap="xl" pb={4}>
            <div className={classes.hero}>
                <div className={classes.crownRing}>
                    <IconCrownFilled size={36} />
                </div>

                <Title className={classes.title} mt={4} order={2}>
                    RW Prime
                </Title>
                <Text c="dimmed" maw={380} mt={8} size="sm">
                    {t('prime-modal.shared.description')}
                </Text>
            </div>

            <Stack gap="md">
                {FEATURES.map((feature) => (
                    <Group align="flex-start" gap="md" key={feature.titleKey} wrap="nowrap">
                        <div className={classes.featureIcon}>
                            <feature.icon size={20} />
                        </div>
                        <div>
                            <Text fw={600} size="sm">
                                {t(feature.titleKey)}
                            </Text>
                            <Text c="dimmed" size="xs">
                                {t(feature.descriptionKey)}
                            </Text>
                        </div>
                    </Group>
                ))}
            </Stack>

            <Stack gap={10}>
                <Button
                    className={classes.joinButton}
                    component="a"
                    fullWidth
                    href={PRIME_LINK}
                    leftSection={<IconCrownFilled size={18} />}
                    radius="md"
                    rel="noopener noreferrer"
                    size="md"
                    target="_blank"
                >
                    {t('prime-modal.shared.join-button')}
                </Button>
            </Stack>
        </Stack>
    )
}
