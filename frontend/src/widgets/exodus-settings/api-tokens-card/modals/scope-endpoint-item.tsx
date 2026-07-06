import { Badge, Box, Checkbox, Group, Stack, Text } from '@mantine/core'
import { memo } from 'react'

import classes from '../api-token-card.module.css'
import { getMethodColor, ScopeEndpoint } from './scopes.utils'

interface IProps {
    checked: boolean
    endpoint: ScopeEndpoint
    onToggle?: (key: string) => void
    readOnly?: boolean
}

export const ScopeEndpointItem = memo(({ checked, endpoint, onToggle, readOnly }: IProps) => {
    const isWrite = endpoint.kind === 'write'

    return (
        <Box
            className={classes.endpointRow}
            data-checked={checked || undefined}
            data-readonly={readOnly || undefined}
            onClick={readOnly ? undefined : () => onToggle?.(endpoint.key)}
        >
            <Group gap="sm" wrap="nowrap">
                {!readOnly && <Checkbox checked={checked} radius="sm" size="xs" tabIndex={-1} />}

                <Badge
                    className={classes.methodBadge}
                    color={getMethodColor(endpoint.method)}
                    ff="monospace"
                    radius="sm"
                    variant="soft"
                >
                    {endpoint.method}
                </Badge>

                <Stack gap={0} style={{ flex: 1, minWidth: 0 }}>
                    <Text className={classes.endpointPath} truncate="end">
                        {endpoint.path}
                    </Text>
                    <Text c="dimmed" size="xs" truncate="end">
                        {endpoint.description}
                    </Text>
                </Stack>

                <Badge
                    color={isWrite ? 'orange' : 'blue'}
                    ff="monospace"
                    radius="sm"
                    variant="soft"
                >
                    {isWrite ? 'Write' : 'Read'}
                </Badge>
            </Group>
        </Box>
    )
})
