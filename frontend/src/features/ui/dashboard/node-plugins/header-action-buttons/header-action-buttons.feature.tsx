import {
    ActionIcon,
    ActionIconGroup,
    Button,
    Group,
    Modal,
    Stack,
    TextInput,
    Tooltip
} from '@mantine/core'
import { TbBook, TbPackage, TbPlus, TbRefresh, TbTerminal } from 'react-icons/tb'
import { generatePath, useNavigate } from 'react-router-dom'
import { useDisclosure } from '@mantine/hooks'
import { useTranslation } from 'react-i18next'
import { useField } from '@mantine/form'

import { QueryKeys, useCreateNodePlugin, useGetNodePlugins } from '@shared/api/hooks'
import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { ROUTES } from '@shared/constants'
import { queryClient } from '@shared/api'

export function NodePluginsHeaderActionButtonsFeature() {
    const { t } = useTranslation()
    const [opened, { open, close }] = useDisclosure(false)
    const navigate = useNavigate()
    const openModalWithData = useModalsStoreOpenWithData()

    const { isFetching } = useGetNodePlugins()

    const nameField = useField<string>({
        initialValue: '',
        validateOnChange: true,
        validate: (value) =>
            value.trim().length > 0
                ? null
                : t('node-plugins-header-action-buttons.feature.name-is-required')
    })

    const { mutate: createNodePlugin, isPending } = useCreateNodePlugin({
        mutationFns: {
            onSuccess: async (plugin) => {
                await queryClient.invalidateQueries({
                    queryKey: QueryKeys.nodes.getNodePlugins.queryKey
                })
                nameField.reset()
                close()
                navigate(
                    generatePath(ROUTES.DASHBOARD.MANAGEMENT.NODE_PLUGINS.NODE_PLUGIN_BY_UUID, {
                        uuid: plugin.uuid
                    })
                )
            }
        }
    })

    const handleRefresh = async () => {
        await queryClient.invalidateQueries({
            queryKey: QueryKeys.nodes.getNodePlugins.queryKey
        })
    }

    return (
        <Group grow preventGrowOverflow={false} wrap="wrap">
            <ActionIcon
                color="lime"
                component="a"
                href="https://docs.rw/docs/learn/node-plugins"
                size="input-md"
                target="_blank"
                variant="soft"
            >
                <TbBook size={24} />
            </ActionIcon>

            <UniversalSpotlightActionIconShared />

            <ActionIconGroup>
                <Tooltip label={t('node-plugin-executor.drawer.executor')} withArrow>
                    <ActionIcon
                        color="grape"
                        onClick={() =>
                            openModalWithData(MODALS.NODE_PLUGIN_EXECUTOR_DRAWER, undefined)
                        }
                        size="input-md"
                        variant="soft"
                    >
                        <TbTerminal size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <Tooltip label={t('common.refresh')} withArrow>
                    <ActionIcon
                        loading={isFetching}
                        onClick={handleRefresh}
                        size="input-md"
                        variant="soft"
                    >
                        <TbRefresh size="24px" />
                    </ActionIcon>
                </Tooltip>
            </ActionIconGroup>

            <ActionIconGroup>
                <ActionIcon color="teal" onClick={open} size="input-md" variant="soft">
                    <TbPlus size="24px" />
                </ActionIcon>
            </ActionIconGroup>

            <Modal
                centered
                onClose={close}
                opened={opened}
                title={
                    <BaseOverlayHeader
                        IconComponent={TbPackage}
                        iconVariant="soft"
                        iconColor="teal"
                        title={t('common.create')}
                    />
                }
            >
                <form
                    onSubmit={(event) => {
                        event.preventDefault()
                        createNodePlugin({
                            variables: {
                                name: nameField.getValue()
                            }
                        })
                    }}
                >
                    <Stack gap="md">
                        <TextInput
                            data-autofocus
                            label={t('common.name')}
                            placeholder={t(
                                'node-plugins-header-action-buttons.feature.my-node-plugin'
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
