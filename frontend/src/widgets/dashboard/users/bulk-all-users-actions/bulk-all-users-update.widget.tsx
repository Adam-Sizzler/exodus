import { Button, Group, Stack } from '@mantine/core'
import { useForm } from '@mantine/form'
import { modals } from '@mantine/modals'
import { BulkAllUpdateUsersCommand } from '@exodus/backend-contract'
import { zodResolver } from 'mantine-form-zod-resolver'
import { motion } from 'motion/react'
import { useTranslation } from 'react-i18next'
import { PiFloppyDiskDuotone } from 'react-icons/pi'

import { useBulkAllUpdateUsers, useGetExternalSquads, useGetUserTags } from '@shared/api/hooks'
import { BulkFormsUsersShared } from '@shared/ui/forms/users/bulk-forms-components'
import { ModalFooter } from '@shared/ui/modal-footer'
import { handleFormErrors } from '@shared/utils/misc'

const MotionWrapper = motion.div
const MotionStack = motion.create(Stack)

const containerVariants = {
    hidden: {},
    visible: {
        transition: {
            staggerChildren: 0.1
        }
    }
}

const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
        opacity: 1,
        y: 0,
        transition: { duration: 0.3 }
    }
}

interface IProps {
    isMobile: boolean
}

export const BulkAllUsersUpdateWidget = (props: IProps) => {
    const { isMobile } = props

    const { t } = useTranslation()

    const { data: externalSquads } = useGetExternalSquads()
    const { data: tags } = useGetUserTags()

    const form = useForm<BulkAllUpdateUsersCommand.Request>({
        mode: 'uncontrolled',
        name: 'bulk-all-user-actions-form',
        initialValues: {
            status: undefined,
            trafficLimitBytes: undefined,
            trafficLimitStrategy: undefined,
            expireAt: undefined,
            description: undefined,
            telegramId: undefined,
            tag: undefined,
            hwidDeviceLimit: undefined
        },
        validate: zodResolver(
            BulkAllUpdateUsersCommand.RequestSchema.omit({
                expireAt: true,
                telegramId: true,
                email: true
            })
        )
    })

    const { mutate: updateUsers, isPending: isUpdatePending } = useBulkAllUpdateUsers({
        mutationFns: {
            onSuccess: () => {
                form.reset()
                modals.closeAll()
            }
        }
    })

    const handleBulkUpdate = form.onSubmit(async (values) => {
        updateUsers(
            {
                variables: {
                    ...values,
                    trafficLimitBytes: values.trafficLimitBytes,
                    // @ts-expect-error - TODO: fix ZOD schema
                    telegramId: values.telegramId === '' ? null : values.telegramId,
                    email: values.email === '' ? null : values.email,
                    // @ts-expect-error - TODO: fix ZOD schema
                    expireAt: values.expireAt ? dayjs(values.expireAt).toISOString() : undefined,
                    tag: values.tag === '' ? null : values.tag,
                    // @ts-expect-error - TODO: fix ZOD schema
                    hwidDeviceLimit: values.hwidDeviceLimit === '' ? null : values.hwidDeviceLimit
                }
            },
            {
                onError: (error) => {
                    handleFormErrors(form, error)
                }
            }
        )
    })

    return (
        <motion.div
            animate={{ opacity: 1 }}
            initial={{ opacity: 0 }}
            transition={{
                duration: 0.4,
                ease: 'easeInOut'
            }}
        >
            {isMobile && (
                <MotionStack
                    animate="visible"
                    gap="md"
                    initial="hidden"
                    variants={containerVariants}
                >
                    <BulkFormsUsersShared.TrafficLimitsCard
                        cardVariants={cardVariants}
                        form={form}
                        motionWrapper={MotionWrapper}
                    />
                    <BulkFormsUsersShared.AccessSettingsCard
                        cardVariants={cardVariants}
                        externalSquads={externalSquads}
                        form={form}
                        motionWrapper={MotionWrapper}
                    />
                    <BulkFormsUsersShared.ContactInformationCard
                        cardVariants={cardVariants}
                        form={form}
                        motionWrapper={MotionWrapper}
                    />
                    <BulkFormsUsersShared.DeviceTagSettingsCard
                        cardVariants={cardVariants}
                        form={form}
                        motionWrapper={MotionWrapper}
                        tags={tags}
                    />
                </MotionStack>
            )}

            {!isMobile && (
                <Group align="flex-start" gap="md" grow={false} wrap="wrap">
                    <MotionStack
                        animate="visible"
                        gap="md"
                        initial="hidden"
                        style={{ flex: '1 1 450px' }}
                        variants={containerVariants}
                    >
                        <BulkFormsUsersShared.ContactInformationCard
                            cardVariants={cardVariants}
                            form={form}
                            motionWrapper={MotionWrapper}
                        />
                        <BulkFormsUsersShared.DeviceTagSettingsCard
                            cardVariants={cardVariants}
                            form={form}
                            motionWrapper={MotionWrapper}
                            tags={tags}
                        />
                    </MotionStack>

                    <MotionStack
                        animate="visible"
                        gap="md"
                        initial="hidden"
                        style={{ flex: '1 1 450px' }}
                        variants={containerVariants}
                    >
                        <BulkFormsUsersShared.TrafficLimitsCard
                            cardVariants={cardVariants}
                            form={form}
                            motionWrapper={MotionWrapper}
                        />
                        <BulkFormsUsersShared.AccessSettingsCard
                            cardVariants={cardVariants}
                            externalSquads={externalSquads}
                            form={form}
                            motionWrapper={MotionWrapper}
                        />
                    </MotionStack>
                </Group>
            )}

            <ModalFooter isMobile={isMobile}>
                <Button
                    color="cyan"
                    leftSection={<PiFloppyDiskDuotone size="16px" />}
                    loading={isUpdatePending}
                    onClick={() => {
                        handleBulkUpdate()
                    }}
                    size="md"
                    variant="light"
                >
                    {t('common.update')}
                </Button>
            </ModalFooter>
        </motion.div>
    )
}
