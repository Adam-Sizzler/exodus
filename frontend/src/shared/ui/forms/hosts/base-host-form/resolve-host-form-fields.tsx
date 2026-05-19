import { ExternalSquadHostOverridesSchema } from '@exodus/backend-contract'
import { TFunction } from 'i18next'

import { HappLogo } from '@shared/ui/logos'

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
                leftSection: <HappLogo size={20} />,
                inputType: 'string'
            }

        default:
            return {
                label: 'Unknown setting'
            }
    }
}
