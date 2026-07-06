import { Divider, Group, Text } from '@mantine/core'

import { formatCurrencyWithIntl } from '@shared/utils/misc'

interface MonthDividerProps {
    label: string
    total: number
}

export function MonthDivider(props: MonthDividerProps) {
    const { label, total } = props

    return (
        <Group gap="sm" wrap="nowrap">
            <Text c="dimmed" fw={600} size="xs" style={{ whiteSpace: 'nowrap' }} tt="capitalize">
                {label}
            </Text>

            <Divider style={{ flex: 1 }} />

            <Text c="teal" fw={700} size="xs" style={{ whiteSpace: 'nowrap' }}>
                {formatCurrencyWithIntl(total)}
            </Text>
        </Group>
    )
}
