import { ExternalSquadHostOverridesSchema } from '@exodus/backend-contract'
import { ActionIcon, HoverCard, px, Stack, Text } from '@mantine/core'
import { PiIdentificationBadge } from 'react-icons/pi'
import { HiQuestionMarkCircle } from 'react-icons/hi'
import { TbFileDescription } from 'react-icons/tb'
import { TFunction } from 'i18next'

const hoverCard = (text: string) => {
    return (
        <HoverCard shadow="md" width={280} withArrow>
            <HoverCard.Target>
                <ActionIcon color="gray" size="xs" variant="subtle">
                    <HiQuestionMarkCircle size={20} />
                </ActionIcon>
            </HoverCard.Target>
            <HoverCard.Dropdown>
                <Stack gap="md">
                    <Stack gap="sm">
                        <Text c="dimmed" size="sm">
                            {text}
                        </Text>
                    </Stack>
                </Stack>
            </HoverCard.Dropdown>
        </HoverCard>
    )
}

export function resolveHostFormFields(
    field: keyof typeof ExternalSquadHostOverridesSchema.shape,
    t: TFunction
): {
    description?: string
    hoverCard?: React.ReactNode
    inputType?: 'boolean' | 'number' | 'string' | 'textarea'
    label: string
    leftSection?: React.ReactNode
    rightSection?: React.ReactNode
} {
    switch (field) {
        case 'serverDescription':
            return {
                label: t('base-host-form.server-description-header'),
                leftSection: <TbFileDescription size={20} />,
                inputType: 'string'
            }

        default:
            return {
                label: 'Unknown setting'
            }
    }
}
