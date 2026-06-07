import { Badge, Button, CopyButton, Skeleton, Stack, Text } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { SiDocker } from 'react-icons/si'
import { PiCheck } from 'react-icons/pi'

import { useGetSubscriptionConnectionsPubKey } from '@shared/api/hooks'

interface IProps {
    port?: number
    apiPath?: string
    apiSchema?: 'mtls' | 'tls'
}

const normalizePath = (value?: string) => {
    const trimmed = (value ?? '').trim()
    if (!trimmed || trimmed === '/') {
        return '/'
    }
    return `/${trimmed.replace(/^\/+|\/+$/g, '')}/`
}

export const CopyDockerComposeWidget = ({ port, apiPath, apiSchema }: IProps) => {
    const { data: pubKey, isLoading: isPubKeyLoading } = useGetSubscriptionConnectionsPubKey()
    const { t } = useTranslation()

    if (isPubKeyLoading || !pubKey) {
        return <Skeleton height={78} />
    }

    const grpcPath = normalizePath(apiPath)
    const normalizedToken = (pubKey.grpcToken ?? '').trim()
    const subPort = 3010
    const grpcPort = port ?? 2222

    const composeBase = `services:\n  exodus-sub:\n    container_name: exodus-sub\n    hostname: exodus-sub\n    image: ghcr.io/Adam-Sizzler/exodus-sub:latest\n    restart: always\n    environment:\n      - APP_PORT_SUB=${subPort}\n      - SUB_GRPC_ADDRESS=0.0.0.0\n      - SUB_GRPC_PORT=${grpcPort}\n      - SUB_PATH=${grpcPath}`

    const composeAuth =
        apiSchema === 'tls'
            ? `\n      - SUB_GRPC_TOKEN=${normalizedToken}`
            : `\n      - SUB_SECRET_KEY=${pubKey.pubKey.trimEnd()}`

    const composePorts =
        apiSchema === 'tls' ? '' : `\n    ports:\n      - \"${grpcPort}:${grpcPort}\"`
    const compose = `${composeBase}${composeAuth}${composePorts}`

    return (
        <Stack gap="xs" mt="md">
            <Text c="dimmed" size="sm">
                {t('copy-docker-compose.widget.generated-from-selected-schema', {
                    defaultValue:
                        apiSchema === 'tls'
                            ? 'Сгенерировано для gRPC + TLS + gRPC token'
                            : 'Сгенерировано для gRPC + mTLS'
                })}
            </Text>

            <Badge color="teal" variant="soft">
                {apiSchema === 'tls' ? 'gRPC + TLS + token' : 'gRPC + mTLS'}
            </Badge>
            <Badge color="gray" variant="soft">
                SUB_PATH={grpcPath}
            </Badge>
            {apiSchema === 'tls' ? (
                <Badge color="gray" variant="soft">
                    SUB_GRPC_TOKEN={normalizedToken}
                </Badge>
            ) : (
                <Badge color="gray" variant="soft">
                    SUB_SECRET_KEY=...
                </Badge>
            )}

            <CopyButton timeout={2000} value={compose}>
                {({ copied, copy }) => (
                    <Button
                        color={copied ? 'teal' : 'gray'}
                        leftSection={copied ? <PiCheck size={18} /> : <SiDocker size={18} />}
                        onClick={copy}
                        size="md"
                        variant="soft"
                    >
                        {t('copy-docker-compose.widget.copy-direct-compose', {
                            defaultValue: 'Copy Compose'
                        })}
                    </Button>
                )}
            </CopyButton>
        </Stack>
    )
}
