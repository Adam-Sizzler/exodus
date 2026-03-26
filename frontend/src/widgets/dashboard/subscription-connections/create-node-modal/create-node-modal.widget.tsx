import { em, Group, Modal, Progress, Stack, Transition } from '@mantine/core'
import { zodResolver } from 'mantine-form-zod-resolver'
import { useTranslation } from 'node_modules/react-i18next'
import { useMediaQuery } from '@mantine/hooks'
import { useEffect, useState } from 'react'
import { useForm } from '@mantine/form'
import { TbCpu } from 'react-icons/tb'

import { useSubscriptionConnectionsStoreActions, useSubscriptionConnectionsStoreCreateModalIsOpen } from '@entities/dashboard/subscription-connections'
import {
    CreateSubscriptionConnectionRequest,
    createSubscriptionConnectionSchema,
    useCreateSubscriptionConnection,
    useGetSubscriptionConnectionsPubKey
} from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { CreateNodeStep1Connection } from './create-node-steps/create-node-step-1-connection'
import { CreateNodeStep2ApiToken } from './create-node-steps/create-node-step-2-api-token'
import { CreateNodeStep3Status } from './create-node-steps/create-node-step-3-status'

export const CreateNodeModalWidget = () => {
    const { t } = useTranslation()

    const isModalOpen = useSubscriptionConnectionsStoreCreateModalIsOpen()
    const actions = useSubscriptionConnectionsStoreActions()

    const { data: pubKey } = useGetSubscriptionConnectionsPubKey()

    const isMobile = useMediaQuery(`(max-width: ${em(768)})`)

    const [activeStep, setActiveStep] = useState(0)
    const [createdNodeUuid, setCreatedNodeUuid] = useState<string>()
    const [selectedPort, setSelectedPort] = useState<number>(2222)

    const form = useForm<CreateSubscriptionConnectionRequest>({
        name: 'create-node-form',
        mode: 'uncontrolled',
        validate: zodResolver(createSubscriptionConnectionSchema)
    })

    const handleClose = () => {
        actions.toggleCreateModal(false)

        setTimeout(() => {
            form.reset()
            form.resetDirty()
            form.resetTouched()
            setActiveStep(0)
            setCreatedNodeUuid(undefined)
        }, 300)
    }

    const { mutate: createNode, isPending: isCreateNodePending } = useCreateSubscriptionConnection({
        mutationFns: {
            onSuccess: (data) => {
                setCreatedNodeUuid(data.uuid)
                setActiveStep(2)
            }
        }
    })

    const handleCreateNode = () => {
        const values = form.getValues()
        const schema = values.apiSchema === 'tls' ? 'tls' : 'mtls'

        createNode({
            variables: {
                ...values,
                name: values.name.trim(),
                address: values.address.trim(),
                apiSchema: schema,
                apiPath: values.apiPath.trim() || '/',
                subpageConfigUuid: (values.subpageConfigUuid ?? '').trim() || null
            }
        })
    }

    useEffect(() => {
        if (form.getValues().port) {
            return
        }

        form.setValues({
            port: 2222,
            apiSchema: 'mtls',
            apiPath: '/',
            subpageConfigUuid: null
        })
    }, [form])

    form.watch('port', ({ value }) => {
        if (typeof value === 'number' && value > 0) {
            setSelectedPort(value)
        }
    })

    return (
        <Modal
            centered
            fullScreen={isMobile}
            onClose={handleClose}
            opened={isModalOpen}
            size="md"
            title={
                <BaseOverlayHeader
                    IconComponent={TbCpu}
                    iconVariant="gradient-teal"
                    title={t('create-node-modal.widget.create-subscription', {
                        defaultValue: 'Создать подписку'
                    })}
                />
            }
            transitionProps={isMobile ? { transition: 'fade', duration: 200 } : undefined}
        >
            <Stack gap="xl">
                <Group gap="xs" grow>
                    <Progress
                        animated
                        color="teal"
                        radius="sm"
                        size="md"
                        striped
                        transitionDuration={300}
                        value={activeStep >= 0 ? 100 : 0}
                    />
                    <Progress
                        animated
                        color="teal"
                        radius="sm"
                        size="md"
                        striped
                        transitionDuration={300}
                        value={activeStep >= 1 ? 100 : 0}
                    />
                    <Progress
                        animated
                        color="teal"
                        radius="sm"
                        size="md"
                        striped
                        transitionDuration={300}
                        value={activeStep >= 2 ? 100 : 0}
                    />
                </Group>

                <Transition
                    duration={300}
                    exitDuration={0}
                    mounted={activeStep === 0}
                    timingFunction="ease"
                    transition="fade"
                >
                    {(styles) => (
                        <div style={styles}>
                            <CreateNodeStep1Connection
                                form={form}
                                onNext={() => setActiveStep(1)}
                                pubKey={pubKey}
                            />
                        </div>
                    )}
                </Transition>

                <Transition
                    duration={300}
                    exitDuration={0}
                    mounted={activeStep === 1}
                    timingFunction="ease"
                    transition="fade"
                >
                    {(styles) => (
                        <div style={styles}>
                            <CreateNodeStep2ApiToken
                                form={form}
                                isCreating={isCreateNodePending}
                                onCreate={handleCreateNode}
                                onPrev={() => setActiveStep(0)}
                                port={selectedPort}
                            />
                        </div>
                    )}
                </Transition>

                <Transition
                    duration={300}
                    exitDuration={0}
                    mounted={activeStep === 2}
                    timingFunction="ease"
                    transition="fade"
                >
                    {(styles) => (
                        <div style={styles}>
                            <CreateNodeStep3Status
                                nodeUuid={createdNodeUuid}
                                onClose={handleClose}
                            />
                        </div>
                    )}
                </Transition>
            </Stack>
        </Modal>
    )
}
