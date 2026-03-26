import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { Select, Skeleton, Stack } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { SiSecurityscorecard } from 'react-icons/si'
import { useTranslation } from 'node_modules/react-i18next'

import { useGetSubscriptionPageConfigs } from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

interface SubscriptionConfigForm {
    subpageConfigUuid?: string | null
}

interface IProps<T extends SubscriptionConfigForm> {
    cardVariants: Variants
    form: UseFormReturnType<T>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
}

export const SubscriptionConfigCard = <T extends SubscriptionConfigForm>(props: IProps<T>) => {
    const { t } = useTranslation()
    const { cardVariants, form, motionWrapper } = props
    const MotionWrapper = motionWrapper

    const { data: subpageConfigs, isLoading: isSubpageConfigsLoading } = useGetSubscriptionPageConfigs()

    const configOptions = (subpageConfigs?.configs ?? []).map((item) => ({
        value: item.uuid,
        label: item.name
    }))

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        IconComponent={SiSecurityscorecard}
                        iconVariant="gradient-teal"
                        title={t('base-node-form.subscription-configuration')}
                        titleOrder={5}
                    />
                </SectionCard.Section>
                <SectionCard.Section>
                    {isSubpageConfigsLoading ? (
                        <Stack gap="md">
                            <Skeleton height={24} width="40%" />
                            <Skeleton height={36} radius="sm" width="100%" />
                        </Stack>
                    ) : (
                        <Select
                            clearable
                            data={configOptions}
                            key={form.key('subpageConfigUuid')}
                            label={t('base-node-form.subscription-profile')}
                            placeholder={t('base-node-form.select-subscription-profile')}
                            searchable
                            size="sm"
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                            {...form.getInputProps('subpageConfigUuid')}
                        />
                    )}
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
