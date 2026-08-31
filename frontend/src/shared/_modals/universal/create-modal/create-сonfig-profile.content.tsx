import { Stack, TextInput, Group, Button, Text } from '@mantine/core'
import { useField } from '@mantine/form'
import { CreateConfigProfileCommand } from '@exodus/backend-contract'
import { t } from 'i18next'
import { generatePath, NavigateFunction } from 'react-router'

import { queryClient } from '@shared/api'
import { useCreateConfigProfile } from '@shared/api/hooks/config-profiles/config-profiles.mutation.hooks'
import { QueryKeys } from '@shared/api/hooks/keys-factory'
import { ROUTES } from '@shared/constants/routes'

interface IProps {
    onClose: () => void
    navigate: NavigateFunction
}

const generateDefaultConfig = () => {
    const randomNumber = Math.floor(Math.random() * 999999) + 1

    return {
        log: {
            level: 'info'
        },
        dns: {
            servers: [
                {
                    tag: 'dns-remote',
                    type: 'udp',
                    server: '1.1.1.1',
                    detour: 'direct'
                }
            ]
        },
        inbounds: [
            {
                type: 'shadowsocks',
                tag: `ss-in-${randomNumber}`,
                listen: '127.0.0.1',
                listen_port: 2080,
                method: 'chacha20-ietf-poly1305'
            },
            {
                type: 'trojan',
                tag: `trojan-in-${randomNumber}`,
                listen: '127.0.0.1',
                listen_port: 2443,
                users: []
            }
        ],
        outbounds: [
            {
                type: 'direct',
                tag: 'direct'
            },
            {
                type: 'block',
                tag: 'block'
            }
        ],
        route: {
            final: 'direct'
        }
    }
}

export const CreateConfigProfileContent = (props: IProps) => {
    const { onClose, navigate } = props

    const handleUpdate = async () => {
        await queryClient.refetchQueries({
            queryKey: QueryKeys.configProfiles.getConfigProfiles.queryKey
        })
    }

    const nameField = useField<CreateConfigProfileCommand.RequestBody['name']>({
        initialValue: '',
        validateOnChange: true,
        validate: (value) => {
            const result = CreateConfigProfileCommand.RequestBodySchema.omit({
                config: true
            }).safeParse({ name: value })
            return result.success ? null : result.error.issues[0]?.message
        }
    })
    const { mutate: createConfigProfile, isPending } = useCreateConfigProfile({
        mutationFns: {
            onSuccess: (data) => {
                onClose()

                handleUpdate()
                navigate(
                    generatePath(ROUTES.DASHBOARD.MANAGEMENT.CONFIG_PROFILE_BY_UUID, {
                        uuid: data.uuid
                    })
                )
            }
        }
    })

    return (
        <form
            onSubmit={(e) => {
                e.preventDefault()
                createConfigProfile({
                    variables: {
                        name: nameField.getValue(),
                        config: generateDefaultConfig()
                    }
                })
            }}
        >
            <Stack gap="md">
                <Text size="sm">
                    {t(
                        'config-profiles-header-action-buttons.feature.create-a-new-config-profile-by-entering-a-name-below'
                    )}
                    <br />

                    {t(
                        'config-profiles-header-action-buttons.feature.you-can-customize-singbox-config-after-creation'
                    )}
                </Text>
                <TextInput
                    data-autofocus
                    label={t('config-profiles-header-action-buttons.feature.profile-name')}
                    placeholder={t(
                        'config-profiles-header-action-buttons.feature.enter-profile-name'
                    )}
                    required
                    {...nameField.getInputProps()}
                />
                <Group justify="flex-end">
                    <Button color="gray" onClick={onClose} variant="light">
                        {t('common.action.cancel')}
                    </Button>

                    <Button color="teal" loading={isPending} type="submit">
                        {t('common.action.create')}
                    </Button>
                </Group>
            </Stack>
        </form>
    )
}
