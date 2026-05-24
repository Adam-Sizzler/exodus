import { Center, Stack, Text, ThemeIcon } from '@mantine/core'
import { TbServer } from 'react-icons/tb'
import { PiCpu } from 'react-icons/pi'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'
import { NodeResponse } from '@shared/api/hooks'

interface IProps {
    nodes: NodeResponse[]
}

export const ActivePluginsOnNodesModalShared = (props: IProps) => {
    const { nodes } = props

    if (nodes.length === 0) {
        return (
            <Center py="xl">
                <Stack align="center" gap="sm">
                    <ThemeIcon color="gray" size="xl" variant="light">
                        <PiCpu size={24} />
                    </ThemeIcon>
                    <Text c="dimmed" size="sm" ta="center">
                        This plugin is not active on any nodes.
                    </Text>
                </Stack>
            </Center>
        )
    }

    return (
        <Stack gap="md">
            <SectionCard.Root gap="xs">
                {nodes.map((node) => (
                    <SectionCard.Section key={node.uuid}>
                        <BaseOverlayHeader
                            countryCode={node.countryCode}
                            IconComponent={TbServer}
                            iconVariant="light"
                            subtitle={node.address}
                            title={node.name}
                            titleOrder={5}
                            withCopy={true}
                        />
                    </SectionCard.Section>
                ))}
            </SectionCard.Root>
        </Stack>
    )
}
