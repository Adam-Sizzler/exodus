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
import { TbCertificate, TbId, TbMapPin, TbPlugConnected, TbRoute2, TbWorld } from 'react-icons/tb'
import { UseFormReturnType } from '@mantine/form'
import { useTranslation } from 'react-i18next'
import { PiArrowRight } from 'react-icons/pi'

import { CreateNodeRequest, NodeKeygenResponse } from '@shared/api/hooks'
import { CopyableFieldShared } from '@shared/ui/copyable-field/copyable-field'
import { COUNTRIES } from '@shared/ui/forms/nodes/base-node-form/constants'

interface IProps {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    form: UseFormReturnType<CreateNodeRequest, any>
    onNext: () => void
    pubKey: NodeKeygenResponse | undefined
}

export const CreateNodeStep1Connection = ({ form, onNext, pubKey }: IProps) => {
    const { t } = useTranslation()
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

    const handleNext = async () => {
        const nameErrors = form.validateField('name')
        const countryCodeErrors = form.validateField('countryCode')
        const addressErrors = form.validateField('address')
        const portErrors = form.validateField('port')
        const apiSchemaErrors = form.validateField('apiSchema')
        const apiPathErrors = form.validateField('apiPath')

        if (
            nameErrors.hasError ||
            countryCodeErrors.hasError ||
            addressErrors.hasError ||
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

                    <Select
                        key={form.key('countryCode')}
                        label={t('base-node-form.country')}
                        {...form.getInputProps('countryCode')}
                        data={COUNTRIES}
                        leftSection={<TbMapPin size={16} />}
                        placeholder={t('base-node-form.select-country')}
                        searchable
                        size="sm"
                        styles={{
                            label: { fontWeight: 500 }
                        }}
                    />

                    <SimpleGrid cols={{ base: 1, sm: 2 }} spacing="xs">
                        <TextInput
                            key={form.key('address')}
                            label={t('create-node-step-1-connection.domain-or-ip')}
                            {...form.getInputProps('address')}
                            leftSection={<TbWorld size={16} />}
                            placeholder="192.168.1.1"
                            required
                            size="sm"
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
                            size="sm"
                            styles={{
                                label: { fontWeight: 500 }
                            }}
                        />
                    </SimpleGrid>

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
                        rightSection={<PiArrowRight size={18} />}
                        size="md"
                        type="submit"
                    >
                        {t('create-node-modal.widget.next')}
                    </Button>
                </Group>
            </Stack>
        </form>
    )
}
