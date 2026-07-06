import { NumberInput, Stack, TextInput } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { BulkUpdateUsersCommand } from '@exodus/backend-contract'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { useTranslation } from 'react-i18next'
import { PiEnvelopeDuotone, PiTelegramLogoDuotone } from 'react-icons/pi'
import { TbMail } from 'react-icons/tb'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

interface IProps {
    cardVariants: Variants
    form: UseFormReturnType<BulkUpdateUsersCommand.Request>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
}

export function BulkContactInformationCard(props: IProps) {
    const { t } = useTranslation()

    const { cardVariants, motionWrapper, form } = props

    const MotionWrapper = motionWrapper

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        iconColor="teal"
                        IconComponent={TbMail}
                        iconSize={20}
                        iconVariant="soft"
                        title={t('contact-information-card.contact-information')}
                        titleOrder={5}
                    />
                </SectionCard.Section>

                <SectionCard.Section>
                    <Stack gap="md">
                        <NumberInput
                            allowDecimal={false}
                            allowNegative={false}
                            hideControls
                            key={form.key('fields.telegramId')}
                            label="Telegram ID"
                            leftSection={<PiTelegramLogoDuotone size="16px" />}
                            placeholder="Enter user's Telegram ID (optional)"
                            {...form.getInputProps('fields.telegramId')}
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />

                        <TextInput
                            key={form.key('fields.email')}
                            label="Email"
                            leftSection={<PiEnvelopeDuotone size="16px" />}
                            placeholder="Enter user's email (optional)"
                            {...form.getInputProps('fields.email')}
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
