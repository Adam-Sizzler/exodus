import {
    TbCertificate,
    TbMapPin,
    TbPackage,
    TbPlugConnected,
    TbRoute2,
    TbUserCheck,
    TbWorld
} from 'react-icons/tb'
import { NumberInput, Select, SimpleGrid, Stack, TextInput } from '@mantine/core'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { useTranslation } from 'node_modules/react-i18next'
import { UseFormReturnType } from '@mantine/form'
import { HiOutlineServer } from 'react-icons/hi'

import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { NodeKeygenResponse, NodePluginResponse } from '@shared/api/hooks'
import { SectionCard } from '@shared/ui/section-card'

import { COUNTRIES } from './constants'

interface NodeVitalsForm {
    activePluginUuid?: null | string
    address?: string
    apiPath?: string
    apiSchema?: 'mtls' | 'tls'
    countryCode?: string
    name?: string
    port?: number
}

interface IProps<T extends NodeVitalsForm> {
    cardVariants: Variants
    form: UseFormReturnType<T>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
    nodePlugins: NodePluginResponse[]
    pubKey: NodeKeygenResponse | undefined
}

export const NodeVitalsCard = <T extends NodeVitalsForm>(props: IProps<T>) => {
    const { t } = useTranslation()
    const { cardVariants, form, motionWrapper, nodePlugins, pubKey } = props
    const apiSchema: 'mtls' | 'tls' = form.values.apiSchema === 'tls' ? 'tls' : 'mtls'
    const apiSchemaInputProps = form.getInputProps('apiSchema')
    const credentialLabel =
        apiSchema === 'tls'
            ? t('base-node-form.grpc-token-node-grpc-token', {
                defaultValue: 'gRPC Token (NODE_GRPC_TOKEN)'
            })
            : t('base-node-form.secret-key-secret-key', {
                defaultValue: 'Secret Key (SECRET_KEY)'
            })
    const credentialValue =
        apiSchema === 'tls'
            ? (pubKey?.grpcToken?.trim() ?? 'Error loading...')
            : (pubKey?.pubKey.trimEnd() ?? 'Error loading...')

    const MotionWrapper = motionWrapper

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        IconComponent={HiOutlineServer}
                        iconVariant="gradient-blue"
                        title={t('base-node-form.node-vitals')}
                        titleOrder={5}
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

                        <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="xs">
                            <TextInput
                                key={form.key('address')}
                                label={t('base-node-form.address')}
                                {...form.getInputProps('address')}
                                leftSection={<TbWorld size={16} />}
                                placeholder={t('base-node-form.e-g-example-com')}
                                required
                                styles={{
                                    label: { fontWeight: 500 }
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
                                    label: { fontWeight: 500 }
                                }}
                            />
                        </SimpleGrid>

                        <Stack gap="xs">
                            <Select
                                data={[
                                    {
                                        label: t('base-node-form.api-schema-mtls'),
                                        value: 'mtls'
                                    },
                                    {
                                        label: t('base-node-form.api-schema-tls-token'),
                                        value: 'tls'
                                    }
                                ]}
                                description={t('base-node-form.api-schema-description')}
                                key={form.key('apiSchema')}
                                label={t('base-node-form.api-schema')}
                                leftSection={<TbPlugConnected size={16} />}
                                required
                                size="sm"
                                styles={{
                                    label: { fontWeight: 500 }
                                }}
                                {...apiSchemaInputProps}
                                onChange={(value) => {
                                    apiSchemaInputProps.onChange(value)
                                }}
                            />

                            <TextInput
                                description={t('base-node-form.api-path-description')}
                                key={form.key('apiPath')}
                                label={t('base-node-form.api-path')}
                                leftSection={<TbRoute2 size={16} />}
                                placeholder={t('base-node-form.api-path-placeholder')}
                                required
                                size="sm"
                                styles={{
                                    label: { fontWeight: 500 }
                                }}
                                {...form.getInputProps('apiPath')}
                            />
                        </Stack>

                        <CopyableFieldShared
                            label={credentialLabel}
                            leftSection={<TbCertificate size={16} />}
                            size="sm"
                            value={credentialValue}
                        />

                        <Select
                            key={form.key('activePluginUuid')}
                            label={t('node-vitals.card.plugin', { defaultValue: 'Plugin' })}
                            {...form.getInputProps('activePluginUuid')}
                            allowDeselect
                            clearable
                            data={nodePlugins.map((nodePlugin) => ({
                                label: nodePlugin.name,
                                value: nodePlugin.uuid
                            }))}
                            description={t(
                                'node-vitals.card.review-documentation-for-more-information',
                                {
                                    defaultValue:
                                        'Review documentation for more information about node plugins.'
                                }
                            )}
                            leftSection={<TbPackage size={16} />}
                            nothingFoundMessage={t('node-vitals.card.nothing-found', {
                                defaultValue: 'Nothing found'
                            })}
                            placeholder={t('node-vitals.card.select-plugin', {
                                defaultValue: 'Select plugin'
                            })}
                            searchable
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />
                    </Stack>
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
