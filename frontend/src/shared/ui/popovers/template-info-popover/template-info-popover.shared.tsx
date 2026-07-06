import { ActionIcon, SimpleGrid, Stack, Text } from '@mantine/core'
import { modals } from '@mantine/modals'
import { TEMPLATE_KEYS, TemplateKeys } from '@exodus/backend-contract'
import { TSubscriptionPageTemplateKey } from '@exodus/subscription-page-types'
import { useTranslation } from 'react-i18next'
import { TbInfoSquare } from 'react-icons/tb'

import { useIsMobile } from '@shared/hooks'
import { CopyableCodeBlock } from '@shared/ui/copyable-code-block'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

interface IProps {
    templateKeys?: readonly TemplateKeys[] | readonly TSubscriptionPageTemplateKey[]
}

export const TemplateInfoPopoverShared = (props: IProps) => {
    const { templateKeys = TEMPLATE_KEYS } = props

    const isMobile = useIsMobile()

    const { t } = useTranslation()

    const handleClick = () => {
        modals.open({
            children: (
                <Stack>
                    <Text size="sm">
                        {t(
                            'template-info-popover.shared.you-can-use-template-variables-in-this-field'
                        )}
                        <br />
                        {t('template-info-popover.shared.available-variables-are-listed-below')}
                    </Text>

                    <SimpleGrid cols={{ base: 1, xs: 2 }} key="template-keys" spacing="xs">
                        {templateKeys.map((key) => (
                            <CopyableCodeBlock key={key} size="small" value={`{{${key}}}`} />
                        ))}
                    </SimpleGrid>
                </Stack>
            ),
            size: 'auto',
            fullScreen: isMobile,
            title: (
                <BaseOverlayHeader
                    iconColor="lime"
                    IconComponent={TbInfoSquare}
                    iconVariant="soft"
                    title={t('template-info-popover.shared.template-variables')}
                />
            )
        })
    }

    return (
        <ActionIcon color="lime" onClick={handleClick} variant="transparent">
            <TbInfoSquare size="20px" />
        </ActionIcon>
    )
}
