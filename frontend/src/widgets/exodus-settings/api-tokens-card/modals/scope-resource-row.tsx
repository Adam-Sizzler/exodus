import { Box, Collapse, Group, Stack, Text, ThemeIcon, UnstyledButton } from '@mantine/core'
import clsx from 'clsx'
import { TbChevronRight } from 'react-icons/tb'

import classes from '../api-token-card.module.css'
import { KindPill } from './kind-pill'
import { ScopeEndpointItem } from './scope-endpoint-item'
import {
    countSelected,
    getKindState,
    getReadKeys,
    getWriteKeys,
    humanizeResource,
    renderResourceIcon,
    ScopeEndpoint,
    ScopeResource
} from './scopes.utils'

interface IProps {
    endpoints: ScopeEndpoint[]
    expanded: boolean
    onToggleEndpoint?: (key: string) => void
    onToggleExpand: () => void
    onToggleKind?: (kind: 'read' | 'write') => void
    readOnly?: boolean
    resource: ScopeResource
    selectedEndpoints: Set<string>
}

export const ScopeResourceRow = (props: IProps) => {
    const {
        endpoints,
        expanded,
        onToggleEndpoint,
        onToggleExpand,
        onToggleKind,
        readOnly,
        resource,
        selectedEndpoints
    } = props

    const total = resource.endpoints.length
    const selectedCount = countSelected(resource.endpoints, selectedEndpoints)
    const readState = getKindState(getReadKeys(resource), selectedEndpoints)
    const writeState = getKindState(getWriteKeys(resource), selectedEndpoints)

    return (
        <Box className={classes.resourceRow}>
            <Group className={classes.resourceHeader} gap="xs" wrap="nowrap">
                <UnstyledButton className={classes.resourceTrigger} onClick={onToggleExpand}>
                    <Group gap="sm" wrap="nowrap">
                        <TbChevronRight
                            className={clsx(classes.chevron, expanded && classes.chevronOpen)}
                            size={12}
                        />
                        <ThemeIcon
                            color={selectedCount > 0 ? 'teal' : 'gray'}
                            size="md"
                            variant={selectedCount > 0 ? 'soft' : 'default'}
                        >
                            {renderResourceIcon(resource.resource, 15)}
                        </ThemeIcon>
                        <Stack gap={0} style={{ minWidth: 0 }}>
                            <Text className={classes.resourceName} truncate="end">
                                {humanizeResource(resource.resource)}
                            </Text>
                            <Text c="dimmed" ff="monospace" size="xs">
                                {selectedCount}/{total}
                            </Text>
                        </Stack>
                    </Group>
                </UnstyledButton>

                {getReadKeys(resource).length > 0 && (
                    <KindPill
                        label="Read"
                        onClick={() => onToggleKind?.('read')}
                        readOnly={readOnly}
                        state={readState}
                    />
                )}
                {getWriteKeys(resource).length > 0 && (
                    <KindPill
                        label="Write"
                        onClick={() => onToggleKind?.('write')}
                        readOnly={readOnly}
                        state={writeState}
                    />
                )}
            </Group>

            <Collapse expanded={expanded}>
                <Stack className={classes.endpointList} gap={4}>
                    {endpoints.map((endpoint) => (
                        <ScopeEndpointItem
                            checked={selectedEndpoints.has(endpoint.key)}
                            endpoint={endpoint}
                            key={endpoint.key}
                            onToggle={onToggleEndpoint}
                            readOnly={readOnly}
                        />
                    ))}
                </Stack>
            </Collapse>
        </Box>
    )
}
