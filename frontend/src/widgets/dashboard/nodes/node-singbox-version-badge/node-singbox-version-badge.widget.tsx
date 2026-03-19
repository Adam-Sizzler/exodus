import { Badge, Group } from '@mantine/core'
import { memo } from 'react'

import { SingboxLogo } from '@shared/ui/logos'

import { IProps } from './interface'

export const NodeSingboxVersionBadgeWidget = memo(({ node, fetchedNode, ...rest }: IProps) => {
    const nodeData = fetchedNode || node
    const nodeSingboxVersion = nodeData.singboxVersion

    return (
        <Group>
            <Badge color="grape" leftSection={<SingboxLogo size={18} />} size="lg" {...rest}>
                {nodeSingboxVersion ?? 'unknown'}
            </Badge>
        </Group>
    )
})
