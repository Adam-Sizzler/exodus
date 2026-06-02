import { CreateHostCommand, SECURITY_LAYERS } from '@exodus/backend-contract'
import { zodResolver } from 'mantine-form-zod-resolver'
import { notifications } from '@mantine/notifications'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'
import { useForm } from '@mantine/form'
import { Drawer } from '@mantine/core'
import { useState } from 'react'

import {
    QueryKeys,
    useCreateHost,
    useGetConfigProfiles,
    useGetInternalSquads,
    useGetNodes,
    useGetSubscriptionTemplates
} from '@shared/api/hooks'
import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { BaseHostForm } from '@shared/ui/forms/hosts/base-host-form'
import { queryClient } from '@shared/api'

type OptionalJSONParseResult = { ok: false } | { ok: true; value: null | unknown }

const parseOptionalJSONValue = (value: unknown): OptionalJSONParseResult => {
    if (value === null || value === undefined) {
        return {
            ok: true,
            value: null
        }
    }

    if (typeof value === 'string') {
        const trimmed = value.trim()
        if (trimmed === '' || trimmed === 'null') {
            return {
                ok: true,
                value: null
            }
        }

        try {
            return {
                ok: true,
                value: JSON.parse(trimmed)
            }
        } catch {
            return {
                ok: false
            }
        }
    }

    if (typeof value === 'object') {
        if (!Array.isArray(value) && Object.keys(value).length === 0) {
            return {
                ok: true,
                value: null
            }
        }

        return {
            ok: true,
            value
        }
    }

    return {
        ok: false
    }
}

export const CreateHostModalWidget = () => {
    const { t } = useTranslation()

    const { isOpen } = useModalState(MODALS.CREATE_HOST_MODAL)
    const close = useModalClose(MODALS.CREATE_HOST_MODAL)

    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()
    const { data: internalSquads } = useGetInternalSquads()
    const { data: templates } = useGetSubscriptionTemplates()

    const [advancedOpened, setAdvancedOpened] = useState(false)

    const form = useForm<CreateHostCommand.Request>({
        mode: 'uncontrolled',
        name: 'create-host-form',
        validateInputOnBlur: true,
        validate: zodResolver(CreateHostCommand.RequestSchema),

        initialValues: {
            securityLayer: SECURITY_LAYERS.DEFAULT,
            port: 0,
            remark: '',
            address: '',
            selectorNodesFirst: false,
            overrideProtocolCredential: false,
            protocolCredential: null,
            inbound: {
                configProfileUuid: '',
                configProfileInboundUuid: ''
            }
        } as CreateHostCommand.Request
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
            }
        }
    })

    const handleSubmit = form.onSubmit(async (values) => {
        const valuesAny = values as CreateHostCommand.Request & {
            overrideProtocolCredential?: boolean
            protocolCredential?: string | null
        }

        if (!values.inbound.configProfileInboundUuid || !values.inbound.configProfileUuid) {
            notifications.show({
                title: t('create-host-modal.widget.error'),
                message: t('create-host-modal.widget.please-select-the-config-profile-and-inbound'),
                color: 'red'
            })

            return null
        }

        let muxParams
        let singboxMuxParams
        let clashMuxParams
        let sockoptParams

        const muxParamsResult = parseOptionalJSONValue(values.muxParams)
        if (!muxParamsResult.ok) {
            notifications.show({
                title: t('create-host-modal.widget.error'),
                message: t('base-host-form.invalid-json'),
                color: 'red'
            })
            return null
        }
        muxParams = muxParamsResult.value

        const singboxMuxParamsResult = parseOptionalJSONValue(values.singboxMuxParams)
        if (!singboxMuxParamsResult.ok) {
            notifications.show({
                title: t('create-host-modal.widget.error'),
                message: t('base-host-form.invalid-json'),
                color: 'red'
            })
            return null
        }
        singboxMuxParams = singboxMuxParamsResult.value

        clashMuxParams =
            typeof values.clashMuxParams === 'string' && values.clashMuxParams.trim() !== ''
                ? values.clashMuxParams
                : null

        const sockoptParamsResult = parseOptionalJSONValue(values.sockoptParams)
        sockoptParams = sockoptParamsResult.ok ? sockoptParamsResult.value : null

        createHost({
            variables: {
                ...values,
                isDisabled: !values.isDisabled,
                overrideProtocolCredential: Boolean(valuesAny.overrideProtocolCredential),
                protocolCredential: valuesAny.overrideProtocolCredential
                    ? valuesAny.protocolCredential || null
                    : null,
                sockoptParams,
                muxParams,
                singboxMuxParams,
                clashMuxParams,
                inbound: {
                    configProfileInboundUuid: values.inbound.configProfileInboundUuid,
                    configProfileUuid: values.inbound.configProfileUuid
                }
            } as any
        })

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
                    IconComponent={PiListChecks}
                    iconVariant="gradient-teal"
                    title={t('create-host-modal.widget.new-host')}
                />
            }
        >
            <BaseHostForm
                advancedOpened={advancedOpened}
                configProfiles={configProfiles?.configProfiles ?? []}
                form={form}
                handleSubmit={handleSubmit}
                internalSquads={internalSquads?.internalSquads ?? []}
                isSubmitting={isCreateHostPending}
                nodes={nodes!}
                setAdvancedOpened={setAdvancedOpened}
                subscriptionTemplates={templates?.templates ?? []}
            />
        </Drawer>
    )
}
