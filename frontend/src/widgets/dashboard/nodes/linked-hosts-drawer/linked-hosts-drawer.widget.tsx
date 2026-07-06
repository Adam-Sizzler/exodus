import { Center, Drawer, Stack, Text, ThemeIcon } from '@mantine/core'
import { HostCardWidget } from '@widgets/dashboard/hosts/host-card/host-card.widget'
import { memo, useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'
import { TbListCheck } from 'react-icons/tb'

import { useGetConfigProfiles, useGetHosts, useGetNodes } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

import { MODALS, useModalClose, useModalState } from '@entities/dashboard/modal-store'

export const LinkedHostsDrawer = memo(() => {
    const { isOpen, internalState: nodeUuid } = useModalState(MODALS.SHOW_NODE_LINKED_HOSTS_DRAWER)
    const close = useModalClose(MODALS.SHOW_NODE_LINKED_HOSTS_DRAWER)

    const { t } = useTranslation()

    const { data: hosts } = useGetHosts()
    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()

    const nodesByUuid = useMemo(
        () => new Map((nodes ?? []).map((node) => [node.uuid, node] as const)),
        [nodes]
    )

    if (!nodeUuid || !hosts || !configProfiles) {
        return (
            <Drawer
                keepMounted={false}
                onClose={close}
                opened={isOpen}
                overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
                padding="lg"
                position="right"
                size="500px"
                title={
                    <BaseOverlayHeader
                        iconColor="teal"
                        IconComponent={PiListChecks}
                        iconVariant="soft"
                        title={t('linked-hosts-drawer.widget.assigned-hosts')}
                    />
                }
            >
                <LoadingScreen />
            </Drawer>
        )
    }

    const linkedHosts = hosts.filter((host) => host.nodes.includes(nodeUuid.nodeUuid))

    return (
        <Drawer
            keepMounted={false}
            onClose={close}
            opened={isOpen}
            overlayProps={{ backgroundOpacity: 0.6, blur: 0 }}
            padding="lg"
            position="right"
            size="800px"
            title={
                <BaseOverlayHeader
                    iconColor="teal"
                    IconComponent={PiListChecks}
                    iconVariant="soft"
                    title={t('linked-hosts-drawer.widget.assigned-hosts')}
                />
            }
        >
            <Stack gap={0}>
                {linkedHosts.length === 0 && (
                    <SectionCard.Root p="xl">
                        <SectionCard.Section>
                            <Center py="xl">
                                <Stack align="center" gap="lg">
                                    <ThemeIcon color="gray" radius="xl" size={64} variant="soft">
                                        <TbListCheck size={32} />
                                    </ThemeIcon>

                                    <Stack align="center" gap="xs">
                                        <Text c="dimmed" fw={600} size="md" ta="center">
                                            {t(
                                                'linked-hosts-drawer.widget.no-hosts-assigned-to-this-node'
                                            )}
                                        </Text>
                                    </Stack>
                                </Stack>
                            </Center>
                        </SectionCard.Section>
                    </SectionCard.Root>
                )}

                {linkedHosts.map((host) => {
                    return (
                        <HostCardWidget
                            configProfiles={configProfiles.configProfiles}
                            item={host}
                            nodesByUuid={nodesByUuid}
                            onSelect={() => {}}
                            openExternal
                            viewOnly
                            key={host.uuid}
                        />
                    )
                })}
            </Stack>
        </Drawer>
    )
})
