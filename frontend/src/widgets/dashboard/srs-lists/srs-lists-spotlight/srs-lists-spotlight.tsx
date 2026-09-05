import { Badge, Group } from '@mantine/core'
import { PiListChecks } from 'react-icons/pi'

import { UniversalSpotlightContentShared } from '@shared/ui/universal-spotlight'

import { SRSListItem } from '../srs-list-card'

interface IProps {
    onEditItem?: (item: SRSListItem) => void
    srsLists: SRSListItem[]
}

export function SRSListsSpotlightWidget(props: IProps) {
    const { srsLists, onEditItem } = props

    return (
        <UniversalSpotlightContentShared
            actions={srsLists.map((item) => ({
                id: item.uuid,
                label: item.fileName || item.shortName || item.url,
                description: item.url,
                leftSection: <PiListChecks color="var(--mantine-color-gray-5)" size={16} />,
                rightSection: (
                    <Group gap="xs" wrap="nowrap">
                        {item.tags && item.tags.length > 0 && (
                            <Badge color="blue" size="sm" variant="light">
                                {item.tags[0]}
                            </Badge>
                        )}
                        <Badge
                            color={
                                !item.isEnabled
                                    ? 'gray'
                                    : item.isAvailable
                                      ? 'teal'
                                      : 'red'
                            }
                            size="sm"
                            variant="light"
                        >
                            {!item.isEnabled ? 'Disabled' : item.isAvailable ? 'Available' : 'Error'}
                        </Badge>
                    </Group>
                ),
                onClick: () => onEditItem?.(item)
            }))}
        />
    )
}
