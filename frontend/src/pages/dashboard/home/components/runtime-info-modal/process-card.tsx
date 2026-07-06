import { Box, Group, Text, ThemeIcon } from '@mantine/core'
import { ComponentType, ReactNode } from 'react'

import classes from './runtime-info-modal.module.css'

interface IProps {
    color: string
    description: ReactNode
    Icon: ComponentType<{ size: number }>
    title: string
}

export const ProcessCard = (props: IProps) => {
    const { color, description, Icon, title } = props
    return (
        <Box
            className={classes.processCard}
            style={{ '--process-color': `var(--mantine-color-${color}-6)` } as React.CSSProperties}
        >
            <Group gap="xs" mb={6} wrap="nowrap">
                <ThemeIcon color={color} radius="md" size="md" variant="soft">
                    <Icon size={16} />
                </ThemeIcon>
                <Text c="white" ff="monospace" fw={700} size="sm">
                    {title}
                </Text>
            </Group>
            <Text c="dimmed" lh={1.55} size="sm">
                {description}
            </Text>
        </Box>
    )
}
