import {
    ActionIcon,
    ActionIconGroup,
    Button,
    Group,
    Modal,
    Stack,
    Text,
    TextInput,
    Tooltip
} from '@mantine/core'
import { CreateConfigProfileCommand } from '@exodus/backend-contract'
import { generatePath, useNavigate } from 'react-router-dom'
import { TbPlus, TbRefresh } from 'react-icons/tb'
import { useDisclosure } from '@mantine/hooks'
import { useTranslation } from 'node_modules/react-i18next'
import { useField } from '@mantine/form'

import { QueryKeys, useCreateConfigProfile, useGetConfigProfiles } from '@shared/api/hooks'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { HelpActionIconShared } from '@shared/ui/help-drawer'
import { SingboxLogo } from '@shared/ui/logos'
import { ROUTES } from '@shared/constants'
import { queryClient } from '@shared/api'

interface IProps {
    configProfileCount: number
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
                    address: 'https://1.1.1.1/dns-query',
                    detour: 'direct'
                }
            ]
        },
        inbounds: [
            {
                type: 'mixed',
                tag: `mixed_${randomNumber}`,
                listen: '127.0.0.1',
                listen_port: 2080
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

export const ConfigProfilesHeaderActionButtonsFeature = (props: IProps) => {
    const { configProfileCount } = props
    const { isFetching } = useGetConfigProfiles()
    const { t } = useTranslation()

    const [opened, { open, close }] = useDisclosure(false)
    const navigate = useNavigate()

    const handleUpdate = async () => {
        await queryClient.refetchQueries({
            queryKey: QueryKeys.configProfiles.getConfigProfiles.queryKey
        })
    }

    const nameField = useField<CreateConfigProfileCommand.Request['name']>({
        initialValue: '',
        validateOnChange: true,
        validate: (value) => {
            const result = CreateConfigProfileCommand.RequestSchema.omit({
                config: true
            }).safeParse({ name: value })
            return result.success ? null : result.error.errors[0]?.message
        }
    })
    const { mutate: createConfigProfile, isPending } = useCreateConfigProfile({
        mutationFns: {
            onSuccess: (data) => {
                close()
                nameField.reset()
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
        <Group grow preventGrowOverflow={false} wrap="wrap">
            <HelpActionIconShared hidden={false} screen="PAGE_CONFIG_PROFILES" />

            {configProfileCount > 0 && <UniversalSpotlightActionIconShared />}

            <ActionIconGroup>
                <Tooltip label={t('common.update')} withArrow>
                    <ActionIcon
                        loading={isFetching}
                        onClick={handleUpdate}
                        size="input-md"
                        variant="light"
                    >
                        <TbRefresh size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <Tooltip
                    label={t('config-profiles-header-action-buttons.feature.create-config-profile')}
                    withArrow
                >
                    <ActionIcon color="teal" onClick={open} size="input-md" variant="light">
                        <TbPlus size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <Modal
                centered
                onClose={close}
                opened={opened}
                size="md"
                title={
                    <BaseOverlayHeader
                        IconComponent={SingboxLogo}
                        iconVariant="gradient-teal"
                        title={t(
                            'config-profiles-header-action-buttons.feature.create-config-profile'
                        )}
                    />
                }
            >
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
                                'config-profiles-header-action-buttons.feature.you-can-customize-xray-config-after-creation'
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
                            <Button color="gray" onClick={close} variant="light">
                                {t('common.cancel')}
                            </Button>

                            <Button color="teal" loading={isPending} type="submit">
                                {t('common.create')}
                            </Button>
                        </Group>
                    </Stack>
                </form>
            </Modal>
        </Group>
    )
}
