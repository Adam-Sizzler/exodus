import type { BooleanFilterValue } from './use-hosts-table-widget'

import {
    ActionIcon,
    Group,
    MultiSelect,
    SegmentedControl,
    Stack,
    Text,
    TextInput,
    ThemeIcon
} from '@mantine/core'
import { ReactNode } from 'react'
import { PiCheckCircleDuotone, PiXCircleDuotone } from 'react-icons/pi'
import { TbSearch, TbX } from 'react-icons/tb'

export function EllipsisCell({ children }: { children: ReactNode }) {
    return (
        <Group gap={6} wrap="nowrap">
            <Text
                size="sm"
                style={{
                    maxWidth: '100%',
                    overflow: 'hidden',
                    textOverflow: 'ellipsis',
                    whiteSpace: 'nowrap'
                }}
            >
                {children}
            </Text>
        </Group>
    )
}

export function BooleanCell({ value }: { value: boolean }) {
    return value ? (
        <ThemeIcon color="teal" size="md" variant="light">
            <PiCheckCircleDuotone size={18} />
        </ThemeIcon>
    ) : (
        <ThemeIcon color="gray" size="md" variant="light">
            <PiXCircleDuotone size={18} />
        </ThemeIcon>
    )
}

export function TextSearchFilter({
    label,
    onChange,
    value
}: {
    label: string
    onChange: (value: string) => void
    value: string
}) {
    return (
        <TextInput
            label={label}
            leftSection={<TbSearch size={16} />}
            onChange={(e) => onChange(e.currentTarget.value)}
            rightSection={
                value ? (
                    <ActionIcon
                        c="dimmed"
                        onClick={() => onChange('')}
                        size="sm"
                        variant="transparent"
                    >
                        <TbX size={14} />
                    </ActionIcon>
                ) : null
            }
            value={value}
        />
    )
}

export function BooleanFilterControl({
    label,
    onChange,
    value
}: {
    label: string
    onChange: (value: BooleanFilterValue) => void
    value: BooleanFilterValue
}) {
    return (
        <Stack gap={6}>
            <Text fw={500} size="sm">
                {label}
            </Text>
            <SegmentedControl
                data={[
                    { label: 'All', value: 'all' },
                    { label: 'Yes', value: 'yes' },
                    { label: 'No', value: 'no' }
                ]}
                fullWidth
                onChange={(next) => onChange(next as BooleanFilterValue)}
                size="xs"
                value={value}
            />
        </Stack>
    )
}

export function SelectFilter({
    data,
    label,
    onChange,
    value
}: {
    data: { label: string; value: string }[]
    label: string
    onChange: (value: string[]) => void
    value: string[]
}) {
    return (
        <MultiSelect
            clearable
            comboboxProps={{ withinPortal: false }}
            data={data}
            label={label}
            leftSection={<TbSearch size={16} />}
            onChange={onChange}
            searchable
            value={value}
        />
    )
}
