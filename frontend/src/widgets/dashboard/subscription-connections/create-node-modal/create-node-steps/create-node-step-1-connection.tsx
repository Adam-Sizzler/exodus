import {
    Anchor,
    Button,
    Code,
    Divider,
    Group,
    NumberInput,
    Select,
    SimpleGrid,
    Stack,
    Text,
    TextInput
} from '@mantine/core'
import { TbCertificate, TbId, TbPlugConnected, TbRoute2, TbWorld } from 'react-icons/tb'
import { UseFormReturnType } from '@mantine/form'
import { useTranslation } from 'react-i18next'
import { PiArrowRight } from 'react-icons/pi'

import { CreateSubscriptionConnectionRequest } from '@shared/api/hooks'
import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'

interface IProps {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    form: UseFormReturnType<CreateSubscriptionConnectionRequest, any>
    onNext: () => void
    pubKey:
        | {
              pubKey: string
              grpcToken?: string
          }
        | undefined
}

export const CreateNodeStep1Connection = ({ form, onNext, pubKey }: IProps) => {
    const { t } = useTranslation()
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
            ? (pubKey?.grpcToken?.trim() ?? 'Error loading...')
            : (pubKey?.pubKey.trimEnd() ?? 'Error loading...')

    const handleNext = async () => {
        const nameErrors = form.validateField('name')
        const addressErrors = form.validateField('address')
        const publicDomainErrors = form.validateField('publicDomain')
        const portErrors = form.validateField('port')
        const apiSchemaErrors = form.validateField('apiSchema')
        const apiPathErrors = form.validateField('apiPath')

        if (
            nameErrors.hasError ||
            addressErrors.hasError ||
            publicDomainErrors.hasError ||
            portErrors.hasError ||
            apiSchemaErrors.hasError ||
            apiPathErrors.hasError
        ) {
            return
        }

        onNext()
    }

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault()
                handleNext()
            }}
        >
            <Stack gap="xs" mih={400}>
                <Text c="dimmed" size="sm">
                    {t('create-node-step-1-connection.copy-the')}{' '}
                    <Code c="white" color="gray.8">
                        docker-compose.yml
                    </Code>{' '}
                    {t('create-node-step-1-connection.content-for-the-exodus-node-below')}{' '}
                    <Anchor
                        fw="700"
                        href="https://docs.exodus.dev/docs/install/exodus-node"
                        inherit
                        rel="noopener noreferrer"
                        target="_blank"
                        underline="hover"
                    >
                        {t('create-node-step-1-connection.learn-more')}
                    </Anchor>
                </Text>

                <Divider />
                <Stack gap="xs">
                    <CopyableFieldShared
                        label={credentialLabel}
                        leftSection={<TbCertificate size={16} />}
                        size="sm"
                        value={credentialValue}
                    />

                    <TextInput
                        key={form.key('name')}
                        label={t('base-node-form.internal-name')}
                        leftSection={<TbId size={16} />}
                        placeholder={t('base-node-form.internal-name-placeholder')}
                        required
                        size="sm"
                        styles={{
                            label: { fontWeight: 500 }
                        }}
                        {...form.getInputProps('name')}
                    />

                    <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="xs">
                        <TextInput
                            key={form.key('address')}
                            label={t('base-node-form.address')}
                            {...form.getInputProps('address')}
                            leftSection={<TbWorld size={16} />}
                            placeholder={t('base-node-form.e-g-example-com')}
                            required
                            size="sm"
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
                            size="sm"
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
                </Stack>

                <Group justify="flex-end" mt="auto">
                    <Button
                        color="teal"
                        disabled={!pubKey}
                        rightSection={<PiArrowRight size={18} />}
                        size="md"
                        type="submit"
                    >
                        {t('common.next', { defaultValue: 'Далее' })}
                    </Button>
                </Group>
            </Stack>
        </form>
    )
}
