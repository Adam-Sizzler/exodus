import { Accordion, Button, Center, Container, Group, Switch, Text, ThemeIcon } from '@mantine/core'
import { notifications } from '@mantine/notifications'
import { useMutation, useQuery } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { TbPackages } from 'react-icons/tb'
import { useState } from 'react'

import { HelpActionIconShared } from '@shared/ui/help-drawer'
import { SettingsCardShared } from '@shared/ui/settings-card'
import { PageHeaderShared, Page } from '@shared/ui'
import { instance } from '@shared/api/axios'
import { createUrl } from '@shared/api/helpers/create-url'

type ModulesSettingsResponse = {
    response?: {
        haproxy?: {
            enabled?: boolean
        }
    }
}

const HaproxySvgIcon = () => (
    <svg
        height="24"
        viewBox="0 0 64 64"
        width="24"
        xmlns="http://www.w3.org/2000/svg"
        xmlnsXlink="http://www.w3.org/1999/xlink"
    >
        <g fill="none" strokeMiterlimit={10}>
            <g stroke="#284a6a">
                <path d="M31.973 23.69l-6.65-7.518" strokeWidth=".26" />
                <path d="M11 17.908l5.205 7.6" strokeWidth=".13" />
                <g strokeWidth=".26">
                    <path d="M31.973 23.7l6.578-7.6m1.374 15.987l7.807-6.578m-7.807 6.568l7.807 6.867m-9.108 8.676l-6.65-7.373m-6.65 7.373l6.65-7.373m-15.76-1.592l7.518-6.578m-7.517-6.58l7.518 6.578m8.24-8.386l-15.76 1.807m15.76-1.806l15.76 1.807" />
                    <path d="M39.925 32.077L38.55 16.1m1.375 15.977l-1.3 15.542" />
                    <path d="M47.732 38.944l-15.76 1.3m-15.758-1.59l15.76 1.6" />
                    <path d="M25.323 47.62l-1.6-15.542m1.6-15.905l-1.6 15.903" />
                </g>
            </g>
            <path
                d="M19.862 35.835l.053-7.6 7.605.053-.053 7.6zm8.228-8.372l.053-7.6 7.605.053-.053 7.6zm-.002 16.67l.053-7.6 7.605.053-.053 7.6zm8.3-8.37l.053-7.6 7.605.053-.053 7.6z"
                fill="#256ea5"
                stroke="none"
            />
            <path
                d="M35.95 18.744l.037-5.277 5.3.037-.037 5.277zm-13.25-.001l.037-5.277 5.3.037L28 18.78zm-9.172 9.4l.037-5.277 5.3.037-.037 5.277zM13.522 41.4l.055-5.277 5.298.055-.055 5.277zm31.598-.1l.037-5.277 5.3.037-.037 5.277zm.107-13.12l.055-5.277 5.298.055-.055 5.277z"
                fill="#3378bc"
                stroke="none"
            />
            <path
                d="M9.233 19.63l.024-3.47 3.463.024-.024 3.47zm6.77-7.062l.024-3.47 3.463.024-.024 3.47zm9.32-3.493l.024-3.47 3.463.024-.024 3.47zM7.194 29.24l.024-3.47 3.44.024-.024 3.47zm43.993-13.032l3.47-.024.024 3.463-3.47.024zm-6.773-7.18l3.47-.024.024 3.463-3.47.024zm-9.35-3.406l3.47-.024.024 3.463-3.47.024zm18.14 20.186l3.47-.024.024 3.463-3.47.024z"
                fill="#169bd6"
                stroke="none"
            />
            <path
                d="M22.627 50.4l.037-5.277 5.3.037-.037 5.277zm13.25 0l.037-5.277 5.3.037-.037 5.277z"
                fill="#3378bc"
                stroke="none"
            />
            <path
                d="M51.164 47.73l.024-3.47 3.463.024-.024 3.47zm-6.698 7.134l.024-3.47 3.463.024-.024 3.47zm-9.32 3.42l.024-3.47 3.463.024-.024 3.47zM53.202 38.12l.024-3.47 3.463.024-.024 3.47zM9.24 44.26l3.47-.024.024 3.463-3.47.024zm6.772 7.108l3.47-.024.024 3.463-3.47.024zm9.276 3.405l3.47-.024.024 3.463-3.47.024zM7.22 34.66l3.47-.024.024 3.463-3.47.024z"
                fill="#169bd6"
                stroke="none"
            />
        </g>
    </svg>
)

export default function ModulesPageComponent() {
    const { i18n } = useTranslation()
    const isRu = i18n.language.startsWith('ru')

    const text = {
        pageTitle: isRu ? 'Модули' : 'Modules',
        cardTitle: isRu ? 'Настройка модулей' : 'Module settings',
        cardDescription: isRu
            ? 'Управление дополнительными модулями панели.'
            : 'Manage additional panel modules.',
        haproxyTitle: 'HAPROXY',
        save: isRu ? 'Сохранить' : 'Save',
        updated: isRu ? 'Настройки модулей обновлены' : 'Modules settings updated',
        failed: isRu ? 'Ошибка обновления настроек модулей' : 'Failed to update modules settings'
    }

    const [enabled, setEnabled] = useState(false)

    const { isLoading } = useQuery({
        queryKey: ['modules-settings'],
        queryFn: async () => {
            const response = await instance.get<ModulesSettingsResponse>(
                createUrl('/api/modules-settings')
            )
            const nextEnabled = Boolean(response.data?.response?.haproxy?.enabled)
            setEnabled(nextEnabled)
            return response.data
        }
    })

    const { mutate: saveSettings, isPending } = useMutation({
        mutationFn: async () => {
            await instance.patch(createUrl('/api/modules-settings'), {
                haproxy: {
                    enabled
                }
            })
        },
        onSuccess: () => {
            notifications.show({
                title: 'Success',
                message: text.updated,
                color: 'teal'
            })
        },
        onError: (error) => {
            notifications.show({
                title: 'Error',
                message: error instanceof Error ? error.message : text.failed,
                color: 'red'
            })
        }
    })

    return (
        <Page title={text.pageTitle}>
            <PageHeaderShared
                icon={
                    <ThemeIcon color="dark" size="lg" variant="soft">
                        <TbPackages size={20} />
                    </ThemeIcon>
                }
                title={text.pageTitle}
            />

            <Container fluid p={0} size="xl">
                <SettingsCardShared.Container>
                    <SettingsCardShared.Header
                        description={text.cardDescription}
                        icon={
                            <ThemeIcon color="dark" size="lg" variant="soft">
                                <TbPackages size={20} />
                            </ThemeIcon>
                        }
                        title={text.cardTitle}
                    />

                    <SettingsCardShared.Content>
                        <Accordion multiple variant="separated">
                            <Accordion.Item key="haproxy" value="haproxy">
                                <Center>
                                    <Accordion.Control
                                        disabled
                                        icon={
                                            <ThemeIcon
                                                size="lg"
                                                style={{
                                                    '--ti-bg':
                                                        'linear-gradient(135deg, rgba(255, 255, 255, 0.98) 0%, rgba(240, 248, 255, 0.98) 100%)',
                                                    '--ti-color': 'var(--mantine-color-blue-8)',
                                                    '--ti-bd': '1px solid transparent'
                                                }}
                                                variant="default"
                                            >
                                                <HaproxySvgIcon />
                                            </ThemeIcon>
                                        }
                                        style={{ opacity: 1 }}
                                        styles={{ chevron: { display: 'none' } }}
                                    >
                                        <Group justify="space-between">
                                            <Text fw={500}>{text.haproxyTitle}</Text>
                                        </Group>
                                    </Accordion.Control>
                                    <Group gap="xs" justify="flex-end" pr="xs" wrap="nowrap">
                                        <HelpActionIconShared
                                            actionIconProps={{ size: 'input-xs' }}
                                            iconProps={{ size: 20 }}
                                            screen="MODULES_HAPROXY"
                                        />
                                        <Switch
                                            checked={enabled}
                                            color="teal.8"
                                            onClick={(e) => e.stopPropagation()}
                                            onChange={(e) => setEnabled(e.currentTarget.checked)}
                                            size="md"
                                        />
                                    </Group>
                                </Center>
                            </Accordion.Item>
                        </Accordion>
                    </SettingsCardShared.Content>

                    <SettingsCardShared.Bottom>
                        <Group justify="flex-end">
                            <Button
                                color="teal"
                                disabled={isLoading}
                                loading={isPending}
                                onClick={() => saveSettings()}
                                size="md"
                            >
                                {text.save}
                            </Button>
                        </Group>
                    </SettingsCardShared.Bottom>
                </SettingsCardShared.Container>
            </Container>
        </Page>
    )
}
