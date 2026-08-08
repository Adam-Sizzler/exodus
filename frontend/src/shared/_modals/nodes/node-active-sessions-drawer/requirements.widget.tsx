import { CodeHighlight } from '@mantine/code-highlight'
import { Group, Stack, ThemeIcon, Text } from '@mantine/core'
import { Trans } from 'react-i18next'
import { TbBrandDocker, TbClock, TbRadar } from 'react-icons/tb'

const DOCKER_SNIPPET = `
    cap_add:
      - NET_ADMIN
`

const HIGHLIGHT_SPAN = <Text c="white" component="span" fw={600} size="sm" />

export const RequirementsWidget = () => {
    return (
        <Stack gap="md">
            <Stack gap="xs">
                <Group gap="sm" wrap="nowrap">
                    <ThemeIcon color="violet" size="md" variant="soft">
                        <TbBrandDocker size={16} />
                    </ThemeIcon>
                    <Text c="dimmed" size="sm">
                        <Trans
                            components={{
                                highlight: HIGHLIGHT_SPAN
                            }}
                            i18nKey="active-sessions-drawer.widget.warning-docker"
                        />
                    </Text>
                </Group>
                <CodeHighlight
                    background="rgba(22, 27, 35)"
                    code={DOCKER_SNIPPET}
                    language="yaml"
                    radius="md"
                    style={{
                        border: '1px solid rgba(255, 255, 255, 0.08)',
                        borderRadius: 'var(--mantine-radius-md)'
                    }}
                />
            </Stack>

            <Group gap="sm" wrap="nowrap">
                <ThemeIcon color="cyan" size="md" variant="soft">
                    <TbClock size={16} />
                </ThemeIcon>
                <Text c="dimmed" size="sm">
                    <Trans
                        components={{
                            highlight: HIGHLIGHT_SPAN
                        }}
                        i18nKey="active-sessions-drawer.widget.warning-activity"
                    />
                </Text>
            </Group>

            <Group gap="sm" wrap="nowrap">
                <ThemeIcon color="yellow" size="md" variant="soft">
                    <TbRadar size={16} />
                </ThemeIcon>
                <Text c="dimmed" size="sm">
                    <Trans
                        components={{
                            highlight: HIGHLIGHT_SPAN
                        }}
                        i18nKey="active-sessions-drawer.widget.warning-n-plus-one"
                    />
                </Text>
            </Group>
        </Stack>
    )
}
