import { Button, Group, Stack, Text } from '@mantine/core'
import { UseFormReturnType } from '@mantine/form'
import { useTranslation } from 'node_modules/react-i18next'
import { PiArrowLeft, PiArrowRight } from 'react-icons/pi'

import { CreateSubscriptionConnectionRequest } from '@shared/api/hooks'
import { CopyDockerComposeWidget } from './copy-docker-compose.widget'

interface IProps {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    form: UseFormReturnType<CreateSubscriptionConnectionRequest, any>
    isCreating: boolean
    onPrev: () => void
    onCreate: () => void
    port?: number
}

export const CreateNodeStep2ApiToken = ({ form, isCreating, onPrev, onCreate, port }: IProps) => {
    const { t } = useTranslation()
    const apiSchema = form.getValues().apiSchema === 'tls' ? 'tls' : 'mtls'

    return (
        <Stack gap="md" mih={400}>
            <Text c="dimmed" size="sm">
                {t('create-node-step-2-api-token.description', {
                    defaultValue:
                        'Сгенерируйте compose и поднимите subscription-ноду с выбранным режимом защиты.'
                })}
            </Text>

            <CopyDockerComposeWidget
                apiPath={form.getValues().apiPath}
                apiSchema={apiSchema}
                port={port}
            />

            <Group justify="space-between" mt="auto">
                <Button
                    color="gray"
                    leftSection={<PiArrowLeft size={18} />}
                    onClick={onPrev}
                    size="md"
                >
                    {t('create-node-modal.widget.back')}
                </Button>

                <Button
                    color="teal"
                    loading={isCreating}
                    onClick={onCreate}
                    rightSection={<PiArrowRight size={18} />}
                    size="md"
                >
                    {t('create-node-modal.widget.create-subscription', {
                        defaultValue: 'Создать подписку'
                    })}
                </Button>
            </Group>
        </Stack>
    )
}
