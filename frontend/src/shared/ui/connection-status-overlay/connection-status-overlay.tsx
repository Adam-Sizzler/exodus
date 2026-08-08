import { Box, Group, Loader, Paper, Stack, Text, ThemeIcon, Transition } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { PiWifiSlash } from 'react-icons/pi'

import { useConnectionProbe } from './use-connection-probe'

export function ConnectionStatusOverlay() {
    const { t } = useTranslation()

    const { isOnline } = useConnectionProbe()

    return (
        <Transition
            duration={150}
            mounted={!isOnline}
            timingFunction="ease"
            transition="slide-down"
        >
            {(transitionStyles) => (
                <Box
                    style={{
                        display: 'flex',
                        inset: '16px 0 auto 0',
                        justifyContent: 'center',
                        pointerEvents: 'none',
                        position: 'fixed',
                        zIndex: 'var(--mantine-z-index-max)'
                    }}
                >
                    <Paper
                        p="sm"
                        radius="md"
                        shadow="xl"
                        style={{
                            ...transitionStyles,
                            background:
                                'linear-gradient(135deg, rgba(30, 30, 35, 0.98) 0%, rgba(45, 25, 25, 0.98) 100%)',
                            border: '1px solid rgba(250, 82, 82, 0.35)',
                            pointerEvents: 'auto'
                        }}
                    >
                        <Group gap="sm" wrap="nowrap">
                            <ThemeIcon
                                color="red"
                                size="lg"
                                style={{
                                    background: 'rgba(250, 82, 82, 0.1)',
                                    border: '1px solid rgba(250, 82, 82, 0.3)'
                                }}
                                variant="light"
                            >
                                <PiWifiSlash size="20px" />
                            </ThemeIcon>
                            <Stack gap={2}>
                                <Text fw={600} size="sm">
                                    {t('connection-status-overlay.connection-lost')}
                                </Text>
                                <Group gap={6}>
                                    <Loader color="red" size={12} />
                                    <Text c="dimmed" size="xs">
                                        {t('connection-status-overlay.reconnecting')}
                                    </Text>
                                </Group>
                            </Stack>
                        </Group>
                    </Paper>
                </Box>
            )}
        </Transition>
    )
}
