import { Badge, Group, Stack, Text } from '@mantine/core'

import { FadeScrollArea } from './fade-scroll-area'

interface ColumnProps {
    actions: React.ReactNode
    children: React.ReactNode
    count: number
    height: string
    icon: React.ReactNode

    scrollable?: boolean
    title: string
}

export function Column(props: ColumnProps) {
    const { actions, children, count, height, icon, scrollable = true, title } = props

    return (
        <Stack gap="xs">
            <Group justify="space-between" mih={34} wrap="nowrap">
                <Group gap="xs" wrap="nowrap">
                    {icon}
                    <Text fw={600} size="sm">
                        {title}
                    </Text>
                    <Badge color="gray" radius="sm" size="md" variant="soft">
                        {count}
                    </Badge>
                </Group>

                <Group gap="xs" wrap="nowrap">
                    {actions}
                </Group>
            </Group>

            {scrollable ? <FadeScrollArea height={height}>{children}</FadeScrollArea> : children}
        </Stack>
    )
}
