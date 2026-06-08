import { Badge, Button, CopyButton, Skeleton, Stack, Text } from '@mantine/core'
import { useTranslation } from 'react-i18next'
import { SiDocker } from 'react-icons/si'
import { PiCheck } from 'react-icons/pi'

import { NodeKeygenResponse } from '@shared/api/hooks'

interface IProps {
    apiPath?: string
    apiSchema?: 'mtls' | 'tls'
    port?: number
    pubKey: NodeKeygenResponse | undefined
}

const normalizePath = (value?: string) => {
    const trimmed = (value ?? '').trim()
    if (!trimmed || trimmed === '/') {
        return '/'
    }
    return `/${trimmed.replace(/^\/+|\/+$/g, '')}/`
}

export const CopyDockerComposeWidget = ({ port, apiPath, apiSchema, pubKey }: IProps) => {
    const { t } = useTranslation()

    if (!pubKey) {
        return <Skeleton height={40} />
    }

    const grpcPath = normalizePath(apiPath)
    const normalizedToken = (pubKey.grpcToken ?? '').trim()
    const grpcPort = port ?? 2222

    const composeBase = `services:
  exodus-node:
    container_name: exodus-node
    hostname: exodus-node
    image: ghcr.io/Adam-Sizzler/exodus-node:latest
    restart: always
    network_mode: host
    cap_add:
      - NET_ADMIN
    ulimits:
      nofile:
        soft: 1048576
        hard: 1048576
    environment:
      - NODE_GRPC_ADDRESS=0.0.0.0
      - NODE_GRPC_PORT=${grpcPort}
      - NODE_GRPC_PATH=${grpcPath}`

    const composeAuth =
        apiSchema === 'tls'
            ? `\n      - NODE_GRPC_TOKEN=${normalizedToken}`
            : `\n      - SECRET_KEY=${pubKey.pubKey.trimEnd()}`

    const compose = `${composeBase}${composeAuth}`

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
                NODE_GRPC_PATH={grpcPath}
            </Badge>
            {apiSchema === 'tls' ? (
                <Badge color="gray" variant="soft">
                    NODE_GRPC_TOKEN={normalizedToken}
                </Badge>
            ) : (
                <Badge color="gray" variant="soft">
                    SECRET_KEY=...
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
                        {t('copy-docker-compose.widget.copy-docker-compose-yml')}
                    </Button>
                )}
            </CopyButton>
        </Stack>
    )
}
