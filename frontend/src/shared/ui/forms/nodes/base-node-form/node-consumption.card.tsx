import {
    ActionIcon,
    Group,
    HoverCard,
    NumberInput,
    NumberInputHandlers,
    rem,
    Stack,
    Text
} from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { CreateNodeCommand, UpdateNodeCommand } from '@exodus/backend-contract'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { useRef } from 'react'
import { useTranslation } from 'react-i18next'
import { HiQuestionMarkCircle } from 'react-icons/hi'
import { TbChartLine, TbMinus, TbPlus } from 'react-icons/tb'

import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

interface IProps {
    cardVariants: Variants
    form: UseFormReturnType<CreateNodeCommand.RequestBody | UpdateNodeCommand.RequestBody>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
}

export const NodeConsumptionCard = (props: IProps) => {
    const { t } = useTranslation()
    const { cardVariants, form, motionWrapper } = props
    const consumptionMultiplierRef = useRef<NumberInputHandlers>(null)
    const nodeConsumptionRef = useRef<NumberInputHandlers>(null)

    const MotionWrapper = motionWrapper

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        iconColor="indigo"
                        IconComponent={TbChartLine}
                        iconVariant="soft"
                        title={t('base-node-form.consumption')}
                        titleOrder={5}
                    />
                </SectionCard.Section>
                <SectionCard.Section>
                    <Stack gap="md">
                        <NumberInput
                            allowDecimal
                            allowedDecimalSeparators={['.']}
                            allowNegative={false}
                            clampBehavior="strict"
                            decimalScale={1}
                            fixedDecimalScale
                            handlersRef={consumptionMultiplierRef}
                            hideControls
                            key={form.key('consumptionMultiplier')}
                            leftSection={
                                <ActionIcon
                                    color="red"
                                    onClick={() => consumptionMultiplierRef.current?.decrement()}
                                    radius="md"
                                    size={rem(44)}
                                    variant="light"
                                >
                                    <TbMinus size={16} />
                                </ActionIcon>
                            }
                            leftSectionPointerEvents="all"
                            leftSectionProps={{
                                style: {
                                    overflow: 'hidden'
                                }
                            }}
                            leftSectionWidth={40}
                            max={100.0}
                            min={0}
                            rightSection={
                                <ActionIcon
                                    color="teal"
                                    onClick={() => consumptionMultiplierRef.current?.increment()}
                                    radius="md"
                                    size={rem(44)}
                                    variant="light"
                                >
                                    <TbPlus size={16} />
                                </ActionIcon>
                            }
                            rightSectionPointerEvents="all"
                            rightSectionProps={{
                                style: {
                                    overflow: 'hidden'
                                }
                            }}
                            rightSectionWidth={40}
                            step={0.1}
                            styles={{
                                input: {
                                    textAlign: 'center',
                                    fontWeight: 600
                                }
                            }}
                            {...form.getInputProps('consumptionMultiplier')}
                            label={
                                <Group align="center" gap={3}>
                                    <HoverCard shadow="md" width={280} withArrow>
                                        <HoverCard.Target>
                                            <ActionIcon color="gray" size="xs" variant="subtle">
                                                <HiQuestionMarkCircle size={20} />
                                            </ActionIcon>
                                        </HoverCard.Target>
                                        <HoverCard.Dropdown>
                                            <Stack gap="sm">
                                                <Text c="dimmed" size="sm">
                                                    {t('base-node-form.consumption-m-line-1')}
                                                </Text>
                                                <Text c="dimmed" size="sm">
                                                    {t('base-node-form.consumption-m-line-2')}
                                                </Text>
                                            </Stack>
                                        </HoverCard.Dropdown>
                                    </HoverCard>
                                    <Text inherit>
                                        {t('node-consumption.card.user-consumption-multiplier')}
                                    </Text>
                                </Group>
                            }
                        />

                        <NumberInput
                            allowDecimal
                            allowedDecimalSeparators={['.']}
                            allowNegative={false}
                            clampBehavior="strict"
                            decimalScale={1}
                            fixedDecimalScale
                            handlersRef={nodeConsumptionRef}
                            hideControls
                            key={form.key('nodeConsumptionMultiplier')}
                            leftSection={
                                <ActionIcon
                                    color="red"
                                    onClick={() => nodeConsumptionRef.current?.decrement()}
                                    radius="md"
                                    size={rem(44)}
                                    variant="light"
                                >
                                    <TbMinus size={16} />
                                </ActionIcon>
                            }
                            leftSectionPointerEvents="all"
                            leftSectionProps={{
                                style: {
                                    overflow: 'hidden'
                                }
                            }}
                            leftSectionWidth={40}
                            max={100.0}
                            min={0}
                            rightSection={
                                <ActionIcon
                                    color="teal"
                                    onClick={() => nodeConsumptionRef.current?.increment()}
                                    radius="md"
                                    size={rem(44)}
                                    variant="light"
                                >
                                    <TbPlus size={16} />
                                </ActionIcon>
                            }
                            rightSectionPointerEvents="all"
                            rightSectionProps={{
                                style: {
                                    overflow: 'hidden'
                                }
                            }}
                            rightSectionWidth={40}
                            step={0.1}
                            styles={{
                                input: {
                                    textAlign: 'center',
                                    fontWeight: 600
                                }
                            }}
                            {...form.getInputProps('nodeConsumptionMultiplier')}
                            label={
                                <Group align="center" gap={3}>
                                    <HoverCard shadow="md" width={280} withArrow>
                                        <HoverCard.Target>
                                            <ActionIcon color="gray" size="xs" variant="subtle">
                                                <HiQuestionMarkCircle size={20} />
                                            </ActionIcon>
                                        </HoverCard.Target>
                                        <HoverCard.Dropdown>
                                            <Stack gap="sm">
                                                <Text c="dimmed" size="sm">
                                                    {t(
                                                        'node-consumption.card.node-consumption-multiplier-1'
                                                    )}
                                                </Text>
                                                <Text c="dimmed" size="sm">
                                                    {t(
                                                        'node-consumption.card.node-consumption-multiplier-2'
                                                    )}
                                                </Text>
                                            </Stack>
                                        </HoverCard.Dropdown>
                                    </HoverCard>
                                    <Text inherit>
                                        {t('node-consumption.card.node-consumption-multiplier')}
                                    </Text>
                                </Group>
                            }
                        />
                    </Stack>
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
