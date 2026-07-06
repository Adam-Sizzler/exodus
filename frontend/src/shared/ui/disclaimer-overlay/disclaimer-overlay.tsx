import { Anchor, Button, Divider, List, Modal, Stack, Text, ThemeIcon, Title } from '@mantine/core'
import { Trans, useTranslation } from 'react-i18next'
import { PiCheckCircle, PiGavel, PiShieldWarning } from 'react-icons/pi'

import { useDisclaimerAccepted, useMiscStoreActions } from '@entities/dashboard/misc-store'

export function DisclaimerOverlay() {
    const { t } = useTranslation()

    const disclaimerAccepted = useDisclaimerAccepted()
    const actions = useMiscStoreActions()

    const handleAccept = () => {
        actions.setDisclaimerAccepted(true)
    }

    const highlightComponents = {
        highlight: <Text c="gray.1" component="span" fw={700} />,
        warning: <Text c="yellow.4" component="span" fw={700} />
    }

    return (
        <Modal
            centered
            closeOnClickOutside={false}
            closeOnEscape={false}
            onClose={handleAccept}
            opened={!disclaimerAccepted}
            padding="xl"
            size="lg"
            withCloseButton={false}
        >
            <Stack align="center" gap="lg">
                <ThemeIcon
                    color="yellow"
                    size={64}
                    style={{
                        background: 'rgba(250, 204, 21, 0.1)',
                        border: '2px solid rgba(250, 204, 21, 0.3)'
                    }}
                    variant="light"
                >
                    <PiShieldWarning size="36px" />
                </ThemeIcon>

                <Stack align="center" gap="xs">
                    <Title c="yellow.4" order={3} ta="center">
                        {t('disclaimer-overlay.title')}
                    </Title>

                    <Text c="gray.4" fw={700} size="sm" ta="center">
                        {t('disclaimer-overlay.subtitle')}
                    </Text>
                </Stack>

                <Divider color="yellow.4" opacity={0.3} variant="dashed" w="100%" />

                <Stack gap="md" w="100%">
                    <Text c="gray.3" size="sm">
                        <Trans
                            components={highlightComponents}
                            i18nKey="disclaimer-overlay.intro"
                        />
                    </Text>

                    <List
                        center
                        icon={
                            <ThemeIcon color="yellow" radius="xl" size={20} variant="light">
                                <PiGavel size="12px" />
                            </ThemeIcon>
                        }
                        size="sm"
                        spacing="sm"
                    >
                        <List.Item>
                            <Text c="gray.3" size="sm">
                                <Trans
                                    components={highlightComponents}
                                    i18nKey="disclaimer-overlay.responsibility"
                                />
                            </Text>
                        </List.Item>
                        <List.Item>
                            <Text c="gray.3" size="sm">
                                <Trans
                                    components={highlightComponents}
                                    i18nKey="disclaimer-overlay.compliance"
                                />
                            </Text>
                        </List.Item>
                        <List.Item>
                            <Text c="gray.3" size="sm">
                                <Trans
                                    components={highlightComponents}
                                    i18nKey="disclaimer-overlay.liability"
                                />
                            </Text>
                        </List.Item>
                    </List>

                    <Text c="gray.5" size="xs" ta="center">
                        <Trans
                            components={{
                                anchor: (
                                    <Anchor
                                        fw={600}
                                        href="https://github.com/exodus/backend/blob/main/LICENCE"
                                        rel="noopener noreferrer"
                                        size="xs"
                                        target="_blank"
                                    />
                                )
                            }}
                            i18nKey="disclaimer-overlay.license-acknowledgement"
                        />
                    </Text>
                </Stack>

                <Button
                    color="yellow"
                    fullWidth
                    leftSection={<PiCheckCircle size="18px" />}
                    onClick={handleAccept}
                    size="md"
                    variant="light"
                >
                    {t('disclaimer-overlay.accept-button')}
                </Button>
            </Stack>
        </Modal>
    )
}
