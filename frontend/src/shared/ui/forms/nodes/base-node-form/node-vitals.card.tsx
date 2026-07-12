import { ActionIcon, Button, Group, NumberInput, Select, Stack, TextInput, Tooltip } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import {
    CreateNodeCommand,
    GetNodePluginsCommand,
    GetPubKeyCommand,
    UpdateNodeCommand
} from '@exodus/backend-contract'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { useState } from 'react'
import { useTranslation } from 'react-i18next'
import { HiOutlineServer } from 'react-icons/hi'
import {
    TbCertificate,
    TbFingerprint,
    TbMapPin,
    TbNetwork,
    TbPackage,
    TbPlugConnected,
    TbRoute2,
    TbUserCheck,
    TbWorld
} from 'react-icons/tb'

import { generateGrpcAuthToken } from '@shared/utils/misc'
import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

import { COUNTRIES } from './constants'

interface IProps<T extends CreateNodeCommand.Request | UpdateNodeCommand.Request> {
    cardVariants: Variants
    form: UseFormReturnType<T>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
    nodePlugins: GetNodePluginsCommand.Response['response']['nodePlugins']
    nodeUuid: string
    pubKey: GetPubKeyCommand.Response['response'] | undefined
}

export const NodeVitalsCard = <T extends CreateNodeCommand.Request | UpdateNodeCommand.Request>(
    props: IProps<T>
) => {
    const { t } = useTranslation()
    const { cardVariants, form, motionWrapper, nodePlugins, pubKey, nodeUuid } = props

    const MotionWrapper = motionWrapper

    const apiSchema: 'mtls' | 'tls' = (form.values as any).apiSchema === 'tls' ? 'tls' : 'mtls'
    const grpcAuthToken: string = (form.values as any).grpcAuthToken || ''

    const apiSchemaInputProps = form.getInputProps('apiSchema')
    const credentialLabel =
        apiSchema === 'tls'
            ? t('base-node-form.grpc-token-grpc-auth-token', {
                defaultValue: 'gRPC Token (GRPC_AUTH_TOKEN)'
            })
            : t('base-node-form.secret-key-secret-key', { defaultValue: 'Secret Key (SECRET_KEY)' })
    const credentialValue =
        apiSchema === 'tls'
            ? (grpcAuthToken.trim() || 'Error loading...')
            : (pubKey?.pubKey.trimEnd() ?? 'Error loading...')

    const regenerateGrpcAuthToken = () => {
        const newToken = generateGrpcAuthToken()
        form.setFieldValue('grpcAuthToken' as never, newToken as never)
        form.setDirty({ grpcAuthToken: true } as never)
    }

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        iconColor="blue"
                        IconComponent={HiOutlineServer}
                        iconVariant="soft"
                        subtitle={nodeUuid}
                        title={t('base-node-form.node-vitals')}
                        titleOrder={5}
                        withCopy
                    />
                </SectionCard.Section>
                <SectionCard.Section>
                    <Stack gap="md">
                        <Select
                            key={form.key('countryCode')}
                            label={t('base-node-form.country')}
                            {...form.getInputProps('countryCode')}
                            data={COUNTRIES}
                            leftSection={<TbMapPin size={16} />}
                            placeholder={t('base-node-form.select-country')}
                            required
                            searchable
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />

                        <TextInput
                            key={form.key('name')}
                            label={t('base-node-form.internal-name')}
                            {...form.getInputProps('name')}
                            leftSection={<TbUserCheck size={16} />}
                            required
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />

                        <Group gap="xs" grow justify="space-between" w="100%">
                            <TextInput
                                key={form.key('address')}
                                label={t('base-node-form.address')}
                                {...form.getInputProps('address')}
                                leftSection={<TbWorld size={16} />}
                                placeholder={t('base-node-form.e-g-example-com')}
                                required
                                styles={{
                                    label: { fontWeight: 500 },
                                    root: { flex: '1 1 70%' }
                                }}
                            />

                            <NumberInput
                                key={form.key('port')}
                                label="Node Port"
                                {...form.getInputProps('port')}
                                allowDecimal={false}
                                allowNegative={false}
                                clampBehavior="strict"
                                decimalScale={0}
                                hideControls
                                max={65535}
                                placeholder="2222"
                                required
                                styles={{
                                    label: { fontWeight: 500 },
                                    root: { flex: '1 1 25%' }
                                }}
                            />
                        </Group>

                        <Stack gap="xs">
                            <Select
                                data={[
                                    {
                                        label: t('base-node-form.api-schema-mtls', {
                                            defaultValue: 'mTLS (SECRET_KEY)'
                                        }),
                                        value: 'mtls'
                                    },
                                    {
                                        label: t('base-node-form.api-schema-tls-token', {
                                            defaultValue: 'TLS + gRPC Token'
                                        }),
                                        value: 'tls'
                                    }
                                ]}
                                description={t('base-node-form.api-schema-description', {
                                    defaultValue: 'How the panel authenticates to this node over gRPC'
                                })}
                                key={form.key('apiSchema')}
                                label={t('base-node-form.api-schema', { defaultValue: 'API Schema' })}
                                leftSection={<TbPlugConnected size={16} />}
                                required
                                size="sm"
                                styles={{
                                    label: { fontWeight: 500 }
                                }}
                                {...apiSchemaInputProps}
                                onChange={(value) => {
                                    const val = value === 'tls' ? 'tls' : 'mtls'
                                    apiSchemaInputProps.onChange(val)
                                }}
                            />

                            {apiSchema === 'tls' && (
                                <TextInput
                                    description={t('base-node-form.api-path-description', {
                                        defaultValue: 'Path prefix the node listens on (PATH_PREFIX)'
                                    })}
                                    key={form.key('apiPath')}
                                    label={t('base-node-form.api-path', { defaultValue: 'API Path' })}
                                    leftSection={<TbRoute2 size={16} />}
                                    placeholder="/"
                                    required
                                    size="sm"
                                    styles={{
                                        label: { fontWeight: 500 }
                                    }}
                                    {...form.getInputProps('apiPath')}
                                />
                            )}
                        </Stack>

                        <Group align="flex-end" gap="xs">
                            <div style={{ flexGrow: 1 }}>
                                <CopyableFieldShared
                                    label={credentialLabel}
                                    leftSection={<TbCertificate size={16} />}
                                    size="sm"
                                    value={credentialValue}
                                />
                            </div>
                            {apiSchema === 'tls' && (
                                <Tooltip
                                    label={t('base-node-form.generate-credential', {
                                        defaultValue: 'Generate credential'
                                    })}
                                    withArrow
                                >
                                    <ActionIcon
                                        color="teal"
                                        h={36}
                                        onClick={regenerateGrpcAuthToken}
                                        style={{ flexShrink: 0 }}
                                        variant="soft"
                                        w={36}
                                    >
                                        <TbFingerprint size={18} />
                                    </ActionIcon>
                                </Tooltip>
                            )}
                        </Group>

                        <Select
                            key={form.key('activePluginUuid')}
                            label={t('node-vitals.card.plugin')}
                            {...form.getInputProps('activePluginUuid')}
                            allowDeselect
                            clearable
                            data={nodePlugins.map((nodePlugin) => ({
                                label: nodePlugin.name,
                                value: nodePlugin.uuid
                            }))}
                            description={t(
                                'node-vitals.card.review-documentation-for-more-information'
                            )}
                            leftSection={<TbPackage size={16} />}
                            nothingFoundMessage={t('node-vitals.card.nothing-found')}
                            placeholder={t('node-vitals.card.select-plugin')}
                            searchable
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />

                        <TextInput
                            key={form.key('proxyUrl')}
                            label={t('node-vitals.card.proxy-url')}
                            {...form.getInputProps('proxyUrl')}
                            description={t('node-vitals.card.proxy-url-description')}
                            leftSection={<TbNetwork size={16} />}
                            placeholder="socks5://user:pass@address:port"
                        />
                    </Stack>
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
