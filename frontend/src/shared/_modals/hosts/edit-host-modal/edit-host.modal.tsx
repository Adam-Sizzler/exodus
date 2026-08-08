import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Drawer } from '@mantine/core'
import { useForm, schemaResolver } from '@mantine/form'
import { UpdateHostCommand } from '@exodus/backend-contract'
import { useEffect, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { queryClient } from '@shared/api'
import {
    QueryKeys,
    useGetConfigProfiles,
    useGetHostTags,
    useGetInternalSquads,
    useGetNodes,
    useGetSubscriptionTemplates,
    useUpdateHost
} from '@shared/api/hooks'
import { BaseHostForm } from '@shared/ui/forms/hosts/base-host-form'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { parseJsonField, parseYamlField, stringifyJsonField, stringifyYamlField } from '@shared/utils/misc'

interface IProps {
    host: UpdateHostCommand.Response['response']
}

export const EditHostDrawer = NiceModal.create((props: IProps) => {
    const { host } = props
    const { t } = useTranslation()

    const modal = useModal()
    const { modalProps, hide } = useNiceMantineModal({
        modal,
        drawer: true,
        onClose() {
            queryClient.refetchQueries({
                queryKey: QueryKeys.hosts.getAllTags.queryKey
            })
            queryClient.refetchQueries({
                queryKey: QueryKeys.hosts.getAllHosts.queryKey
            })
        }
    })

    const [advancedOpened, setAdvancedOpened] = useState(false)

    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()
    const { data: templates } = useGetSubscriptionTemplates()
    const { data: internalSquads } = useGetInternalSquads()
    const { data: hostTags } = useGetHostTags()

    const form = useForm<UpdateHostCommand.RequestBody>({
        name: 'edit-host-form',
        mode: 'uncontrolled',
        validateInputOnBlur: true,
        onValuesChange: (values) => {
            if (typeof values.vlessRouteId === 'string' && values.vlessRouteId === '') {
                form.setFieldValue('vlessRouteId', null)
            }
        },
        validate: schemaResolver(UpdateHostCommand.RequestBodySchema.omit({ uuid: true }))
    })

    const { mutate: updateHost, isPending: isUpdateHostPending } = useUpdateHost({
        mutationFns: {
            onSuccess: async () => {
                hide()
            }
        }
    })

    useEffect(() => {
        if (configProfiles) {
            const hostAny = host as HostType & {
                clashMuxParams?: unknown
                overrideProtocolCredential?: boolean
                protocolCredential?: string | null
                singboxMuxParams?: unknown
                singboxCustomParams?: unknown
                mihomoCustomParams?: string | null
            }

            form.initialize({
                uuid: host.uuid,
                remark: host.remark,
                address: host.address,
                port: host.port,
                securityLayer: host.securityLayer,
                isDisabled: !host.isDisabled,
                sni: host.sni ?? undefined,
                host: host.host ?? undefined,
                path: host.path ?? undefined,
                alpn: host.alpn ?? undefined,
                fingerprint: host.fingerprint ?? undefined,
                inbound: {
                    configProfileUuid: host.inbound.configProfileUuid ?? '',
                    configProfileInboundUuid: host.inbound.configProfileInboundUuid ?? ''
                },
                serverDescription: host.serverDescription ?? undefined,
                xhttpExtraParams: stringifyJsonField(host.xhttpExtraParams),
                muxParams: stringifyJsonField(host.muxParams),
                singboxMuxParams: stringifyJsonField(hostAny.singboxMuxParams),
                clashMuxParams: stringifyYamlField(hostAny.clashMuxParams),
                singboxCustomParams: stringifyJsonField(hostAny.singboxCustomParams),
                mihomoCustomParams: hostAny.mihomoCustomParams ?? '',
                sockoptParams: stringifyJsonField(host.sockoptParams),
                finalMask: stringifyJsonField(host.finalMask),
                tags: host.tags ?? undefined,
                isHidden: host.isHidden,
                overrideSniFromAddress: host.overrideSniFromAddress,
                keepSniBlank: host.keepSniBlank,
                overrideProtocolCredential: hostAny.overrideProtocolCredential ?? undefined,
                protocolCredential: hostAny.overrideProtocolCredential
                    ? (hostAny.protocolCredential ?? undefined)
                    : undefined,
                vlessRouteId: host.vlessRouteId ?? undefined,
                pinnedPeerCertSha256: host.pinnedPeerCertSha256 ?? undefined,
                verifyPeerCertByName: host.verifyPeerCertByName ?? undefined,
                shuffleHost: host.shuffleHost ?? undefined,
                mihomoX25519: host.mihomoX25519 ?? undefined,
                mihomoIpVersion: host.mihomoIpVersion ?? undefined,
                nodes: host.nodes ?? undefined,
                xrayJsonTemplateUuid: host.xrayJsonTemplateUuid ?? undefined,
                excludedInternalSquads: host.excludedInternalSquads ?? undefined,
                excludeFromSubscriptionTypes: host.excludeFromSubscriptionTypes ?? undefined
            })
        }
    }, [configProfiles])

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

    const handleSubmit = form.onSubmit(async (values) => {
        const valuesAny = values as typeof values & {
            clashMuxParams?: unknown
            overrideProtocolCredential?: boolean
            protocolCredential?: string | null
            singboxMuxParams?: unknown
            singboxCustomParams?: unknown
            mihomoCustomParams?: string | null
        }

        const variables = {
            ...valuesAny,
            isDisabled: !values.isDisabled,
            uuid: host.uuid,
            xhttpExtraParams: parseJsonField(values.xhttpExtraParams),
            muxParams: parseJsonField(values.muxParams),
            singboxMuxParams: parseJsonField(valuesAny.singboxMuxParams),
            clashMuxParams: parseYamlField(valuesAny.clashMuxParams),
            singboxCustomParams: parseJsonField(valuesAny.singboxCustomParams),
            mihomoCustomParams: valuesAny.mihomoCustomParams || null,
            sockoptParams: parseJsonField(values.sockoptParams),
            finalMask: parseJsonField(values.finalMask),
            overrideProtocolCredential: Boolean(valuesAny.overrideProtocolCredential),
            protocolCredential: valuesAny.overrideProtocolCredential
                ? valuesAny.protocolCredential || null
                : null
        } as unknown as UpdateHostCommand.Request

        updateHost({ variables })
    })

    return (
        <Drawer
            {...modalProps}
            padding="lg"
            position="right"
            size="lg"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={PiListChecks}
                    iconVariant="soft"
                    subtitle={host.uuid}
                    title={t('edit-host-modal.widget.edit-host')}
                    withCopy={true}
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
                isSubmitting={isUpdateHostPending}
                nodes={nodes!}
                hostUuid={host.uuid}
                setAdvancedOpened={setAdvancedOpened}
                subscriptionTemplates={templates?.templates ?? []}
            />
        </Drawer>
    )
})
