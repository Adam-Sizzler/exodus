import { Group, Stack, Text, ThemeIcon } from '@mantine/core'
import { ComponentType, ReactNode } from 'react'

interface IProps {
    color: string
    description: ReactNode
    Icon: ComponentType<{ size: number }>
    title: string
}

export const MetricEntry = (props: IProps) => {
    const { color, description, Icon, title } = props
    return (
        <Group align="flex-start" gap="md" wrap="nowrap">
            <ThemeIcon color={color} radius="md" size="md" variant="soft">
                <Icon size={16} />
            </ThemeIcon>
            <Stack gap={4} style={{ flex: 1, minWidth: 0 }}>
                <Text c="white" ff="monospace" fw={700} lh={1} size="sm">
                    {title}
                </Text>
                <Text c="dimmed" lh={1.55} size="sm">
                    {description}
                </Text>
            </Stack>
        </Group>
    )
}
