import { UpdateHostCommand } from '@exodus/backend-contract'
import { zodResolver } from 'mantine-form-zod-resolver'
import { notifications } from '@mantine/notifications'
import { memo, useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'
import { modals } from '@mantine/modals'
import { useForm } from '@mantine/form'
import { Drawer } from '@mantine/core'

import {
    QueryKeys,
    useCreateHost,
    useGetConfigProfiles,
    useGetInternalSquads,
    useGetNodes,
    useGetSubscriptionTemplates,
    useUpdateHost
} from '@shared/api/hooks'
import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { BaseHostForm } from '@shared/ui/forms/hosts/base-host-form'
import { cloneString } from '@shared/utils/misc/clone-string'
import { queryClient } from '@shared/api'
import {} from '@entities/dashboard'

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

export const EditHostModalWidget = memo(() => {
    const { t } = useTranslation()

    const { isOpen, internalState: host } = useModalState(MODALS.EDIT_HOST_MODAL)
    const close = useModalClose(MODALS.EDIT_HOST_MODAL)

    const [advancedOpened, setAdvancedOpened] = useState(false)

    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()
    const { data: templates } = useGetSubscriptionTemplates()
    const { data: internalSquads } = useGetInternalSquads()

    const form = useForm<UpdateHostCommand.Request>({
        name: 'edit-host-form',
        mode: 'uncontrolled',
        validateInputOnBlur: true,
        validate: zodResolver(UpdateHostCommand.RequestSchema.omit({ uuid: true }))
    })

    const handleClose = () => {
        close()

        setTimeout(() => {
            form.reset()
            form.resetDirty()
            form.resetTouched()
            setAdvancedOpened(false)
        }, 200)
    }

    const { mutate: updateHost, isPending: isUpdateHostPending } = useUpdateHost({
        mutationFns: {
            onSuccess: async () => {
                handleClose()
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.hosts.getAllTags.queryKey
                })
            }
        }
    })

    const { mutate: createHost } = useCreateHost({
        mutationFns: {
            onSuccess: async () => {
                handleClose()
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.hosts.getAllTags.queryKey
                })
            }
        }
    })

    useEffect(() => {
        if (host && configProfiles) {
            const hostAny = host as any
            let muxParamsParsed: null | object | string
            let singboxMuxParamsParsed: null | object | string
            let clashMuxParamsParsed: null | object | string
            let sockoptParamsParsed: null | object | string

            if (typeof hostAny.muxParams === 'object' && hostAny.muxParams !== null) {
                muxParamsParsed = JSON.stringify(hostAny.muxParams, null, 2)
            } else {
                muxParamsParsed = ''
            }

            if (typeof hostAny.singboxMuxParams === 'object' && hostAny.singboxMuxParams !== null) {
                singboxMuxParamsParsed = JSON.stringify(hostAny.singboxMuxParams, null, 2)
            } else {
                singboxMuxParamsParsed = ''
            }

            if (typeof hostAny.clashMuxParams === 'string') {
                clashMuxParamsParsed = hostAny.clashMuxParams
            } else if (
                typeof hostAny.clashMuxParams === 'object' &&
                hostAny.clashMuxParams !== null
            ) {
                clashMuxParamsParsed = JSON.stringify(hostAny.clashMuxParams, null, 2)
            } else {
                clashMuxParamsParsed = ''
            }

            if (typeof host.sockoptParams === 'object' && host.sockoptParams !== null) {
                sockoptParamsParsed = JSON.stringify(host.sockoptParams, null, 2)
            } else {
                sockoptParamsParsed = ''
            }

            form.setValues({
                remark: host.remark,
                address: host.address,
                port: host.port,
                securityLayer: host.securityLayer,
                isDisabled: !host.isDisabled,
                sni: host.sni ?? undefined,
                host: host.host ?? undefined,
                path: host.path ?? undefined,
                alpn: (host.alpn as UpdateHostCommand.Request['alpn']) ?? undefined,
                fingerprint:
                    (host.fingerprint as UpdateHostCommand.Request['fingerprint']) ?? undefined,
                inbound: {
                    configProfileUuid: host.inbound.configProfileUuid ?? '',
                    configProfileInboundUuid: host.inbound.configProfileInboundUuid ?? ''
                },
                serverDescription: host.serverDescription ?? undefined,
                muxParams: muxParamsParsed,
                singboxMuxParams: singboxMuxParamsParsed,
                clashMuxParams: clashMuxParamsParsed,
                sockoptParams: sockoptParamsParsed,
                tag: host.tag ?? undefined,
                isHidden: host.isHidden,
                overrideSniFromAddress: host.overrideSniFromAddress,
                keepSniBlank: host.keepSniBlank,
                overrideProtocolCredential: hostAny.overrideProtocolCredential ?? false,
                protocolCredential: hostAny.protocolCredential ?? undefined,
                allowInsecure: host.allowInsecure ?? undefined,
                shuffleHost: host.shuffleHost ?? undefined,
                selectorNodesFirst: hostAny.selectorNodesFirst ?? undefined,
                mihomoX25519: host.mihomoX25519 ?? undefined,
                nodes: host.nodes ?? undefined,
                xrayJsonTemplateUuid: host.xrayJsonTemplateUuid ?? undefined,
                excludedInternalSquads: host.excludedInternalSquads ?? undefined,
                excludeFromSubscriptionTypes: host.excludeFromSubscriptionTypes ?? undefined
            } as any)
        }
    }, [host, configProfiles])

    form.watch('inbound.configProfileInboundUuid', ({ value }) => {
        const { inbound } = form.getValues()
        if (!inbound?.configProfileUuid) {
            return
        }

        const configProfile = configProfiles?.configProfiles.find(
            (configProfile) => configProfile.uuid === inbound.configProfileUuid
        )
        if (configProfile) {
            form.setFieldValue(
                'port',
                configProfile.inbounds.find((inbound) => inbound.uuid === value)?.port ?? undefined
            )
        }
    })

    form.watch('allowInsecure', ({ value }) => {
        if (value === true) {
            modals.openConfirmModal({
                title: t('edit-host-modal.widget.are-you-sure'),
                children: t(
                    'edit-host-modal.widget.allowing-insecure-connections-can-lead-to-security-risks-we-do-not-recommend-enabling-this-option'
                ),
                centered: true,
                labels: {
                    confirm: t('edit-host-modal.widget.proceed'),
                    cancel: t('edit-host-modal.widget.cancel')
                },
                confirmProps: {
                    color: 'red'
                },
                onConfirm: () => {
                    form.setFieldValue('allowInsecure', true)
                },
                onCancel: () => {
                    form.setFieldValue('allowInsecure', false)
                }
            })
        }
    })

    const handleSubmit = form.onSubmit(async (values) => {
        if (!host) {
            return
        }

        const valuesAny = values as UpdateHostCommand.Request & {
            overrideProtocolCredential?: boolean
            protocolCredential?: string | null
        }

        let muxParams
        let singboxMuxParams
        let clashMuxParams
        let sockoptParams

        const muxParamsResult = parseOptionalJSONValue(values.muxParams)
        if (!muxParamsResult.ok) {
            notifications.show({
                title: t('edit-host-modal.widget.error'),
                message: t('base-host-form.invalid-json'),
                color: 'red'
            })
            return
        }
        muxParams = muxParamsResult.value

        const singboxMuxParamsResult = parseOptionalJSONValue(values.singboxMuxParams)
        if (!singboxMuxParamsResult.ok) {
            notifications.show({
                title: t('edit-host-modal.widget.error'),
                message: t('base-host-form.invalid-json'),
                color: 'red'
            })
            return
        }
        singboxMuxParams = singboxMuxParamsResult.value

        clashMuxParams =
            typeof values.clashMuxParams === 'string' && values.clashMuxParams.trim() !== ''
                ? values.clashMuxParams
                : null

        const sockoptParamsResult = parseOptionalJSONValue(values.sockoptParams)
        sockoptParams = sockoptParamsResult.ok ? sockoptParamsResult.value : null

        updateHost({
            variables: {
                ...values,
                isDisabled: !values.isDisabled,
                uuid: host.uuid,
                overrideProtocolCredential: Boolean(valuesAny.overrideProtocolCredential),
                protocolCredential: valuesAny.overrideProtocolCredential
                    ? valuesAny.protocolCredential || null
                    : null,
                muxParams,
                singboxMuxParams,
                clashMuxParams,
                sockoptParams,
                tag: values.tag === '' ? null : values.tag
            } as any
        })
    })

    const handleCloneHost = () => {
        if (!host) {
            return
        }

        if (!host.inbound.configProfileInboundUuid || !host.inbound.configProfileUuid) {
            notifications.show({
                title: t('edit-host-modal.widget.error'),
                message: t('edit-host-modal.widget.dangling-host-cannot-be-cloned'),
                color: 'red'
            })

            return
        }

        createHost({
            variables: {
                ...host,
                remark: cloneString(host.remark),
                port: host.port,

                isDisabled: true,
                path: host.path ?? undefined,
                sni: host.sni ?? undefined,
                host: host.host ?? undefined,
                alpn: (host.alpn as UpdateHostCommand.Request['alpn']) ?? undefined,
                muxParams: (host as any).muxParams ?? undefined,
                singboxMuxParams: (host as any).singboxMuxParams ?? undefined,
                clashMuxParams: (host as any).clashMuxParams ?? undefined,
                fingerprint:
                    (host.fingerprint as UpdateHostCommand.Request['fingerprint']) ?? undefined,
                inbound: {
                    configProfileUuid: host.inbound.configProfileUuid,
                    configProfileInboundUuid: host.inbound.configProfileInboundUuid
                },
                serverDescription: host.serverDescription ?? undefined,
                sockoptParams: host.sockoptParams ?? undefined,
                tag: host.tag ?? undefined,
                overrideSniFromAddress: host.overrideSniFromAddress,
                keepSniBlank: host.keepSniBlank,
                overrideProtocolCredential: (host as any).overrideProtocolCredential ?? undefined,
                protocolCredential: (host as any).overrideProtocolCredential
                    ? ((host as any).protocolCredential ?? undefined)
                    : null,
                allowInsecure: host.allowInsecure ?? undefined,
                selectorNodesFirst: (host as any).selectorNodesFirst ?? undefined,
                nodes: host.nodes ?? undefined,
                xrayJsonTemplateUuid: host.xrayJsonTemplateUuid ?? undefined,
                excludedInternalSquads: host.excludedInternalSquads ?? undefined
            } as any
        })
    }

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
                    subtitle={host?.uuid}
                    title={t('edit-host-modal.widget.edit-host')}
                    withCopy={true}
                />
            }
        >
            {host && (
                <BaseHostForm
                    advancedOpened={advancedOpened}
                    configProfiles={configProfiles?.configProfiles ?? []}
                    form={form}
                    handleCloneHost={handleCloneHost}
                    handleSubmit={handleSubmit}
                    internalSquads={internalSquads?.internalSquads ?? []}
                    isSubmitting={isUpdateHostPending}
                    nodes={nodes!}
                    setAdvancedOpened={setAdvancedOpened}
                    subscriptionTemplates={templates?.templates ?? []}
                />
            )}
        </Drawer>
    )
})
