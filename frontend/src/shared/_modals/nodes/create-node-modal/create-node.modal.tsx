import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Group, Modal, Progress, Stack, Transition } from '@mantine/core'
import { useForm, schemaResolver } from '@mantine/form'
import { CreateNodeCommand } from '@exodus/backend-contract'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { TbCpu } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { queryClient } from '@shared/api'
import {
    configProfilesQueryKeys,
    QueryKeys,
    useCreateNode,
    useGetNodeSecretKey
} from '@shared/api/hooks'
import { generateGrpcAuthToken } from '@shared/utils/misc'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'

import { CreateNodeStep1Connection } from './create-node-steps/create-node-step-1-connection'
import { CreateNodeStep2ConfigProfiles } from './create-node-steps/create-node-step-2-config-profiles'
import { CreateNodeStep3Status } from './create-node-steps/create-node-step-3-status'

export const CreateNodeModal = NiceModal.create(() => {
    const { t } = useTranslation()

    const modal = useModal()
    const { modalProps, hide } = useNiceMantineModal({
        modal,
        onClose: () => {
            queryClient.refetchQueries({ queryKey: QueryKeys.nodes.getAllNodes.queryKey })
        }
    })

    const { data: secretKey } = useGetNodeSecretKey()

    const [activeStep, setActiveStep] = useState(0)
    const [createdNodeUuid, setCreatedNodeUuid] = useState<string>()
    const [selectedPort, setSelectedPort] = useState<number>(2222)

    const form = useForm<CreateNodeCommand.RequestBody>({
        name: 'create-node-form',
        mode: 'uncontrolled',
        initialValues: {
            apiSchema: 'mtls'
        } as any,
        validate: schemaResolver(CreateNodeCommand.RequestBodySchema)
    })

    const { mutate: createNode, isPending: isCreateNodePending } = useCreateNode({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.refetchQueries({
                    queryKey: configProfilesQueryKeys.getConfigProfiles.queryKey
                })

                setCreatedNodeUuid(data.uuid)
                setActiveStep(2) // Move to status step
            }
        }
    })

    const handleCreateNode = () => {
        const values = form.getValues()
        createNode({
            variables: {
                ...values,
                name: values.name.trim(),
                address: values.address.trim()
            }
        })
    }

    const nextStep = () => setActiveStep((current) => (current < 2 ? current + 1 : current))
    const prevStep = () => setActiveStep((current) => (current > 0 ? current - 1 : current))

    useEffect(() => {
        if (!form.getValues().port) {
            form.setValues({
                port: 2222
            })
        }
        if (!form.getValues().grpcAuthToken) {
            form.setFieldValue('grpcAuthToken', generateGrpcAuthToken())
        }
    }, [form])

    form.watch('port', ({ value }) => {
        if (value) {
            setSelectedPort(value)
        }
    })

    return (
        <Modal
            {...modalProps}
            size="md"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={TbCpu}
                    iconVariant="soft"
                    title={t('create-node-modal.widget.create-node')}
                />
            }
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
                                onNext={nextStep}
                                port={selectedPort}
                                secretKey={secretKey?.secretKey}
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
                            <CreateNodeStep2ConfigProfiles
                                form={form}
                                isCreating={isCreateNodePending}
                                onCreateNode={handleCreateNode}
                                onPrev={prevStep}
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
                            <CreateNodeStep3Status nodeUuid={createdNodeUuid} onClose={hide} />
                        </div>
                    )}
                </Transition>
            </Stack>
        </Modal>
    )
})
