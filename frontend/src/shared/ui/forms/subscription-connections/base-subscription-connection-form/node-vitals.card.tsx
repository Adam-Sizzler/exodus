import { TbCertificate, TbPlugConnected, TbRoute2, TbUserCheck, TbWorld } from 'react-icons/tb'
import { ForwardRefComponent, HTMLMotionProps, Variants } from 'motion/react'
import { NumberInput, Select, SimpleGrid, Stack, TextInput } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { HiOutlineServer } from 'react-icons/hi'
import { useTranslation } from 'react-i18next'

import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'
import { SubscriptionConnectionKeygenResponse } from '@shared/api/hooks'

interface SubscriptionNodeVitalsForm {
    name?: string
    address?: string
    publicDomain?: string | null
    port?: number
    apiSchema?: 'mtls' | 'tls'
    apiPath?: string
}

interface IProps<T extends SubscriptionNodeVitalsForm> {
    cardVariants: Variants
    form: UseFormReturnType<T>
    motionWrapper: ForwardRefComponent<HTMLDivElement, HTMLMotionProps<'div'>>
    secretKey: SubscriptionConnectionKeygenResponse | undefined
}

export const NodeVitalsCard = <T extends SubscriptionNodeVitalsForm>(props: IProps<T>) => {
    const { t } = useTranslation()
    const { cardVariants, form, motionWrapper, secretKey } = props
    const apiSchema: 'mtls' | 'tls' = form.values.apiSchema === 'tls' ? 'tls' : 'mtls'
    const apiSchemaInputProps = form.getInputProps('apiSchema')
    const credentialLabel =
        apiSchema === 'tls'
            ? t('base-node-form.grpc-token-sub-grpc-token', {
                  defaultValue: 'gRPC Token (SUB_GRPC_TOKEN)'
              })
            : t('base-node-form.secret-key-sub-secret-key', {
                  defaultValue: 'Secret Key (SUB_SECRET_KEY)'
              })
    const credentialValue =
        apiSchema === 'tls'
            ? (secretKey?.grpcToken?.trim() ?? 'Error loading...')
            : (secretKey?.secretKey.trimEnd() ?? 'Error loading...')

    const MotionWrapper = motionWrapper

    return (
        <MotionWrapper variants={cardVariants}>
            <SectionCard.Root>
                <SectionCard.Section>
                    <BaseOverlayHeader
                        IconComponent={HiOutlineServer}
                        iconVariant="soft" iconColor="blue"
                        title={t('base-node-form.node-vitals')}
                        titleOrder={5}
                    />
                </SectionCard.Section>
                <SectionCard.Section>
                    <Stack gap="md">
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
                                label={t('base-node-form.node-port')}
                                {...form.getInputProps('port')}
                                allowDecimal={false}
                                allowNegative={false}
                                clampBehavior="strict"
                                decimalScale={0}
                                hideControls
                                max={65535}
                                placeholder={t('base-node-form.node-port-placeholder')}
                                required
                                styles={{
                                    label: { fontWeight: 500 }
                                }}
                            />
                        </SimpleGrid>

                        <TextInput
                            key={form.key('publicDomain')}
                            description={t('base-node-form.public-domain-description')}
                            label={t('base-node-form.public-domain')}
                            leftSection={<TbWorld size={16} />}
                            placeholder={t('base-node-form.public-domain-placeholder')}
                            size="sm"
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                            {...form.getInputProps('publicDomain')}
                        />

                        <Stack gap="xs">
                            <Select
                                key={form.key('apiSchema')}
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
                                key={form.key('apiPath')}
                                description={t('base-node-form.api-path-description')}
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
                    </Stack>
                </SectionCard.Section>
            </SectionCard.Root>
        </MotionWrapper>
    )
}
