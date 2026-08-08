import { CodeHighlight } from '@mantine/code-highlight'
import { Group, Stack, Text, ThemeIcon } from '@mantine/core'
import { Trans, useTranslation } from 'react-i18next'
import { TbBrandDocker, TbClock, TbHourglass, TbTag } from 'react-icons/tb'

const DOCKER_SNIPPET = `
    cap_add:
      - NET_ADMIN
`

const HIGHLIGHT_SPAN = <Text c="white" component="span" fw={600} size="sm" />

export const RequirementsWidget = () => {
    const { t } = useTranslation()

    return (
        <Stack gap="md">
            <Group gap="sm" wrap="nowrap">
                <ThemeIcon color="teal" size="md" variant="soft">
                    <TbTag size={16} />
                </ThemeIcon>
                <Text c="dimmed" size="sm">
                    <Trans
                        components={{ highlight: HIGHLIGHT_SPAN }}
                        i18nKey="active-sessions-drawer.widget.warning-version"
                        values={{ version: '2.7.0' }}
                    />
                </Text>
            </Group>

            <Stack gap="xs">
                <Group gap="sm" wrap="nowrap">
                    <ThemeIcon color="violet" size="md" variant="soft">
                        <TbBrandDocker size={16} />
                    </ThemeIcon>
                    <Text c="dimmed" size="sm">
                        <Trans
                            components={{ highlight: HIGHLIGHT_SPAN }}
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
                        components={{ highlight: HIGHLIGHT_SPAN }}
                        i18nKey="active-sessions-drawer.widget.warning-activity"
                    />
                </Text>
            </Group>

            <Group gap="sm" wrap="nowrap">
                <ThemeIcon color="orange" size="md" variant="soft">
                    <TbHourglass size={16} />
                </ThemeIcon>
                <Text c="dimmed" size="sm">
                    {t('active-sessions-drawer.widget.warning-patience')}
                </Text>
            </Group>
        </Stack>
    )
}
