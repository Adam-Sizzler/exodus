import { Button, Center, Modal, SegmentedControl, Stack, TextInput } from '@mantine/core'
import { DatePickerInput } from '@mantine/dates'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { CreateInfraBillingNodeCommand } from '@exodus/backend-contract'
import dayjs from 'dayjs'
import { zodResolver } from 'mantine-form-zod-resolver'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { HiCalendar } from 'react-icons/hi'
import { TbCursorText, TbServer } from 'react-icons/tb'

import { queryClient } from '@shared/api'
import { QueryKeys, useCreateInfraBillingNode } from '@shared/api/hooks'
import { SelectBillingNodeShared } from '@shared/ui/infra-billing/select-billing-node/select-billing-node.shared'
import { SelectInfraProviderShared } from '@shared/ui/infra-billing/select-infra-provider/select-infra-provider.shared'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { handleFormErrors } from '@shared/utils/misc'
import { toUtcDayISO } from '@shared/utils/time-utils'

import { MODALS, useModalClose, useModalIsOpen } from '@entities/dashboard/modal-store'

enum Mode {
    NAME = 'name',
    NODE = 'node'
}

export function CreateInfraBillingNodeModalWidget() {
    const isOpen = useModalIsOpen(MODALS.CREATE_INFRA_BILLING_NODE_MODAL)
    const close = useModalClose(MODALS.CREATE_INFRA_BILLING_NODE_MODAL)

    const { t, i18n } = useTranslation()

    const [mode, setMode] = useState<Mode>(Mode.NODE)

    const form = useForm<CreateInfraBillingNodeCommand.Request>({
        name: 'create-infra-billing-node-form',
        mode: 'uncontrolled',
        validate: zodResolver(
            CreateInfraBillingNodeCommand.RequestSchema.omit({
                providerUuid: true,
                nextBillingAt: true,
                nodeUuid: true
            })
        ),
        initialValues: {
            name: null,
            nodeUuid: null,
            // @ts-expect-error - ignore
            providerUuid: undefined,
            nextBillingAt: new Date()
        }
    })

    const resetForm = () => {
        form.reset()
        setMode(Mode.NODE)
    }

    const handleModeChange = (value: string) => {
        const nextMode = value as Mode
        setMode(nextMode)

        if (nextMode === Mode.NODE) {
            form.setFieldValue('name', null)
        } else {
            form.setFieldValue('nodeUuid', null)
        }
    }

    const { mutate: createInfraBillingNode, isPending: isCreateInfraBillingNodePending } =
        useCreateInfraBillingNode({
            mutationFns: {
                onSuccess: (data) => {
                    queryClient.setQueryData(
                        QueryKeys.infraBilling.getInfraBillingNodes.queryKey,
                        data
                    )

                    resetForm()

                    close()
                },
                onError: (error) => {
                    handleFormErrors(form, error)
                }
            }
        })

    const handleSubmit = form.onSubmit(async (values) => {
        const trimmedName = values.name?.trim() || null
        const hasNode = mode === Mode.NODE && Boolean(values.nodeUuid)
        const hasName = mode === Mode.NAME && Boolean(trimmedName)

        if (!values.providerUuid || (!hasNode && !hasName)) {
            notifications.show({
                title: t('create-infra-billing-node.modal.widget.error'),
                message: t(
                    'create-infra-billing-node.modal.widget.please-select-a-provider-and-billing-node'
                ),
                color: 'red'
            })

            return
        }
        createInfraBillingNode({
            variables: {
                providerUuid: values.providerUuid,
                nodeUuid: hasNode ? values.nodeUuid : null,
                name: hasName ? trimmedName : null,
                // @ts-expect-error - TODO: fix ZOD schema
                nextBillingAt: toUtcDayISO(values.nextBillingAt)
            }
        })
    })

    return (
        <Modal
            centered
            keepMounted={false}
            onClose={() => {
                resetForm()
                close()
            }}
            opened={isOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            padding="lg"
            size="md"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbServer}
                    iconVariant="soft"
                    title={t('create-infra-billing-node.modal.widget.billing-node')}
                />
            }
        >
            <form onSubmit={handleSubmit}>
                <Stack>
                    <SegmentedControl
                        data={[
                            {
                                label: (
                                    <Center style={{ gap: 10 }}>
                                        <TbServer size={16} />
                                        <span>
                                            {t(
                                                'create-infra-billing-node.modal.widget.existing-node'
                                            )}
                                        </span>
                                    </Center>
                                ),
                                value: 'node'
                            },
                            {
                                label: (
                                    <Center style={{ gap: 10 }}>
                                        <TbCursorText size={16} />
                                        <span>
                                            {t(
                                                'create-infra-billing-node.modal.widget.custom-name'
                                            )}
                                        </span>
                                    </Center>
                                ),
                                value: 'name'
                            }
                        ]}
                        fullWidth
                        onChange={handleModeChange}
                        transitionDuration={150}
                        value={mode}
                    />

                    {mode === Mode.NODE ? (
                        <SelectBillingNodeShared
                            selectedBillingNodeUuid={form.getValues().nodeUuid}
                            setSelectedBillingNodeUuid={(nodeUuid) => {
                                form.setValues({
                                    nodeUuid: nodeUuid ?? undefined
                                })
                                form.setTouched({
                                    nodeUuid: true
                                })
                                form.setDirty({
                                    nodeUuid: true
                                })
                            }}
                        />
                    ) : (
                        <TextInput
                            data-autofocus
                            description={t(
                                'create-infra-billing-node.modal.widget.custom-name-description'
                            )}
                            key={form.key('name')}
                            label={t('create-infra-billing-node.modal.widget.custom-name')}
                            leftSection={<TbCursorText size="16px" />}
                            placeholder="Management Server"
                            {...form.getInputProps('name')}
                        />
                    )}

                    <Stack gap="md">
                        <SelectInfraProviderShared
                            selectedInfraProviderUuid={form.getValues().providerUuid}
                            setSelectedInfraProviderUuid={(providerUuid) => {
                                form.setValues({
                                    providerUuid: providerUuid ?? undefined
                                })
                                form.setTouched({
                                    providerUuid: true
                                })
                                form.setDirty({
                                    providerUuid: true
                                })
                            }}
                        />
                    </Stack>

                    <DatePickerInput
                        data-autofocus
                        key={form.key('nextBillingAt')}
                        label={t('create-infra-billing-node.modal.widget.next-billing-at')}
                        locale={i18n.language}
                        required
                        valueFormat="D MMMM, YYYY"
                        {...form.getInputProps('nextBillingAt')}
                        description={t(
                            'create-infra-billing-node.modal.widget.next-billing-at-description'
                        )}
                        highlightToday
                        leftSection={<HiCalendar size="16px" />}
                        minDate={dayjs().subtract(1, 'day').toDate()}
                    />

                    <Button loading={isCreateInfraBillingNodePending} type="submit" variant="soft">
                        {t('common.create')}
                    </Button>
                </Stack>
            </form>
        </Modal>
    )
}
