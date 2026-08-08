import NiceModal, { useModal } from '@ebay/nice-modal-react'
import { Center, Drawer, Stack, Text, ThemeIcon } from '@mantine/core'
import { HostCardWidget } from '@widgets/dashboard/hosts/host-card/host-card.widget'
import { useMemo } from 'react'
import { useTranslation } from 'react-i18next'
import { PiListChecks } from 'react-icons/pi'
import { TbListCheck } from 'react-icons/tb'

import { useNiceMantineModal } from '@shared/_modals/use-nice-modal'
import { useGetConfigProfiles, useGetHosts, useGetNodes } from '@shared/api/hooks'
import { LoadingScreen } from '@shared/ui'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'

interface IProps {
    nodeUuid: string
}

export const LinkedHostsDrawer = NiceModal.create((props: IProps) => {
    const { nodeUuid } = props

    const modal = useModal()
    const { modalProps } = useNiceMantineModal({
        modal,
        drawer: true
    })

    const { t } = useTranslation()

    const { data: hosts } = useGetHosts()
    const { data: configProfiles } = useGetConfigProfiles()
    const { data: nodes } = useGetNodes()

    const nodesByUuid = useMemo(
        () => new Map((nodes ?? []).map((node) => [node.uuid, node] as const)),
        [nodes]
    )

    const renderLinkedHosts = () => {
        if (!nodeUuid || !hosts || !configProfiles) return null

        const linkedHosts = hosts.filter((host) => host.nodes.includes(nodeUuid))

        return (
            <>
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
                            viewOnly
                            key={host.uuid}
                        />
                    )
                })}
            </>
        )
    }

    return (
        <Drawer
            {...modalProps}
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
            {!nodeUuid || !hosts || !configProfiles ? (
                <LoadingScreen />
            ) : (
                <Stack gap={0}>{renderLinkedHosts()}</Stack>
            )}
        </Drawer>
    )
})
