import { Select, Stack } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { CreateUserCommand, UpdateUserCommand } from '@exodus/backend-contract'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { useTranslation } from 'react-i18next'
import { PiClockDuotone } from 'react-icons/pi'
import { TbChartLine } from 'react-icons/tb'

import { resetDataStrategy } from '@shared/constants/forms'
import { TrafficLimitInput } from '@shared/ui/forms/traffic-limit-input'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

interface IProps<T extends CreateUserCommand.RequestBody | UpdateUserCommand.RequestBody> {
    cardVariants: Variants
    form: UseFormReturnType<T>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
}

export const TrafficLimitsCard = <
    T extends CreateUserCommand.RequestBody | UpdateUserCommand.RequestBody
>(
    props: IProps<T>
) => {
    const { t } = useTranslation()

    const { cardVariants, motionWrapper, form } = props

    const MotionWrapper = motionWrapper

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        iconColor="violet"
                        IconComponent={TbChartLine}
                        iconSize={20}
                        iconVariant="soft"
                        title={t('traffic-limits-card.traffic-and-limits')}
                        titleOrder={5}
                    />
                </SectionCard.Section>
                <SectionCard.Section>
                    <Stack gap="md">
                        <TrafficLimitInput
                            description={t('traffic-limits-card.traffic-limit-description')}
                            key={form.key('trafficLimitBytes')}
                            label={t('traffic-limits-card.traffic-limit')}
                            leftSection={<TbChartLine size={16} />}
                            {...form.getInputProps('trafficLimitBytes')}
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />

                        <Select
                            allowDeselect={false}
                            comboboxProps={{
                                transitionProps: { transition: 'fade', duration: 200 }
                            }}
                            data={resetDataStrategy(t)}
                            defaultValue={form.values.trafficLimitStrategy}
                            description={t(
                                'create-user-modal.widget.traffic-reset-strategy-description'
                            )}
                            key={form.key('trafficLimitStrategy')}
                            label={t('create-user-modal.widget.traffic-reset-strategy')}
                            leftSection={<PiClockDuotone size="16px" />}
                            placeholder={t('create-user-modal.widget.pick-value')}
                            {...form.getInputProps('trafficLimitStrategy')}
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />
                    </Stack>
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
