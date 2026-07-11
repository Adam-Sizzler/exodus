import { Drawer } from '@mantine/core'
import { useForm } from '@mantine/form'
import { notifications } from '@mantine/notifications'
import { CreateHostCommand, SECURITY_LAYERS } from '@exodus/backend-contract'
import { zodResolver } from 'mantine-form-zod-resolver'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'

import { queryClient } from '@shared/api'
import {
    QueryKeys,
    useCreateHost,
    useGetConfigProfiles,
    useGetHostTags,
    useGetInternalSquads,
    useGetNodes,
    useGetSubscriptionTemplates
} from '@shared/api/hooks'
import { BaseHostForm } from '@shared/ui/forms/hosts/base-host-form'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { parseJsonField, parseYamlField } from '@shared/utils/misc'

import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'

export const CreateHostModalWidget = () => {
    const { t } = useTranslation()

    const { isOpen } = useModalState(MODALS.CREATE_HOST_MODAL)
    const close = useModalClose(MODALS.CREATE_HOST_MODAL)

    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()
    const { data: internalSquads } = useGetInternalSquads()
    const { data: templates } = useGetSubscriptionTemplates()
    const { data: hostTags } = useGetHostTags()

    const [advancedOpened, setAdvancedOpened] = useState(false)

    const form = useForm<CreateHostCommand.Request>({
        mode: 'uncontrolled',
        name: 'create-host-form',
        validateInputOnBlur: true,
        onValuesChange: (values) => {
            if (typeof values.vlessRouteId === 'string' && values.vlessRouteId === '') {
                form.setFieldValue('vlessRouteId', null)
            }
        },
        validate: zodResolver(CreateHostCommand.RequestSchema),

        initialValues: {
            securityLayer: SECURITY_LAYERS.DEFAULT,
            port: 0,
            remark: '',
            address: '',
            inbound: {
                configProfileUuid: '',
                configProfileInboundUuid: ''
            }
        }
    })

    const handleClose = () => {
        close()
        setAdvancedOpened(false)

        form.reset()
        form.resetDirty()
        form.resetTouched()
    }

    const { mutate: createHost, isPending: isCreateHostPending } = useCreateHost({
        mutationFns: {
            onSuccess: async () => {
                handleClose()
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.hosts.getAllTags.queryKey
                })

                notifications.show({
                    title: 'Success',
                    message: 'Host created successfully',
                    color: 'teal'
                })
            }
        }
    })

    const handleSubmit = form.onSubmit(async (values) => {
        if (!values.inbound.configProfileInboundUuid || !values.inbound.configProfileUuid) {
            notifications.show({
                title: t('create-host-modal.widget.error'),
                message: t('create-host-modal.widget.please-select-the-config-profile-and-inbound'),
                color: 'red'
            })

            return null
        }

        let xhttpExtraParams
        let muxParams
        let singboxMuxParams
        let clashMuxParams
        let sockoptParams
        let finalMask
        const valuesAny = values as typeof values & {
            clashMuxParams?: unknown
            overrideProtocolCredential?: boolean
            protocolCredential?: string | null
            singboxMuxParams?: unknown
        }

        try {
            if (values.xhttpExtraParams === '') {
                xhttpExtraParams = null
            } else {
                xhttpExtraParams = JSON.parse(values.xhttpExtraParams as unknown as string)
            }
        } catch {
            xhttpExtraParams = null
            // silence
        }

        try {
            if (values.muxParams === '') {
                muxParams = null
            } else {
                muxParams = JSON.parse(values.muxParams as unknown as string)
            }
        } catch {
            muxParams = null
            // silence
        }

        singboxMuxParams = parseJsonField(valuesAny.singboxMuxParams)
        clashMuxParams = parseYamlField(valuesAny.clashMuxParams)

        try {
            if (values.sockoptParams === '') {
                sockoptParams = null
            } else {
                sockoptParams = JSON.parse(values.sockoptParams as unknown as string)
            }
        } catch {
            sockoptParams = null
            // silence
        }

        try {
            if (values.finalMask === '') {
                finalMask = null
            } else {
                finalMask = JSON.parse(values.finalMask as unknown as string)
            }
        } catch {
            finalMask = null
            // silence
        }

        const variables = {
            ...valuesAny,
            isDisabled: !values.isDisabled,
            sockoptParams,
            muxParams,
            singboxMuxParams,
            clashMuxParams,
            xhttpExtraParams,
            finalMask,
            overrideProtocolCredential: Boolean(valuesAny.overrideProtocolCredential),
            protocolCredential: valuesAny.overrideProtocolCredential
                ? valuesAny.protocolCredential || null
                : null,
            inbound: {
                configProfileInboundUuid: values.inbound.configProfileInboundUuid,
                configProfileUuid: values.inbound.configProfileUuid
            }
        } as unknown as CreateHostCommand.Request

        createHost({ variables })

        return null
    })

    form.watch('inbound.configProfileInboundUuid', ({ value }) => {
        const { configProfileUuid } = form.getValues().inbound
        if (!configProfileUuid) {
            return
        }

        const configProfile = configProfiles?.configProfiles.find(
            (configProfile) => configProfile.uuid === configProfileUuid
        )
        if (configProfile) {
            form.setFieldValue(
                'port',
                configProfile.inbounds.find((inbound) => inbound.uuid === value)?.port ?? 0
            )
        }
    })

    return (
        <Drawer
            keepMounted={false}
            onClose={handleClose}
            opened={isOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            padding="lg"
            position="right"
            size="lg"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={PiListChecks}
                    iconVariant="soft"
                    title={t('create-host-modal.widget.new-host')}
                />
            }
        >
            <BaseHostForm
                advancedOpened={advancedOpened}
                configProfiles={configProfiles?.configProfiles ?? []}
                form={form}
                handleSubmit={handleSubmit}
                hostTags={hostTags?.tags ?? []}
                internalSquads={internalSquads?.internalSquads ?? []}
                isSubmitting={isCreateHostPending}
                nodes={nodes!}
                setAdvancedOpened={setAdvancedOpened}
                subscriptionTemplates={templates?.templates ?? []}
            />
        </Drawer>
    )
}
