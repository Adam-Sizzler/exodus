import {
    ActionIcon,
    Badge,
    Group,
    Loader,
    Paper,
    SimpleGrid,
    Text,
    ThemeIconProps,
    Tooltip
} from '@mantine/core'
import { PiCloudArrowUpDuotone, PiWarningCircle } from 'react-icons/pi'
import { TbPower, TbWifi, TbWifiOff } from 'react-icons/tb'
import { memo, useMemo } from 'react'
import { useTranslation } from 'react-i18next'

import { getSingboxUptimeUtil } from '@shared/utils/time-utils'
import {
    QueryKeys,
    SubscriptionConnectionResponse,
    useDisableSubscriptionConnection,
    useEnableSubscriptionConnection
} from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'
import { queryClient } from '@shared/api'
import { Logo } from '@shared/ui'

interface IProps {
    node: SubscriptionConnectionResponse
}

export const NodeDetailsCardWidget = memo((props: IProps) => {
    const { node } = props
    const nodeSingboxUptime = node.singboxUptime

    const { t } = useTranslation()

    const mutationParams = {
        route: {
            uuid: node.uuid
        },
        mutationFns: {
            onSuccess: async (node: SubscriptionConnectionResponse) => {
                await queryClient.setQueryData(
                    QueryKeys.subscriptionConnections.getNode({ uuid: node.uuid }).queryKey,
                    node
                )
            }
        }
    }

    const { mutate: disableNode, isPending: isDisableNodePending } =
        useDisableSubscriptionConnection(mutationParams)
    const { mutate: enableNode, isPending: isEnableNodePending } =
        useEnableSubscriptionConnection(mutationParams)

    const isConfigMissing = useMemo(() => {
        return !node.subpageConfigUuid
    }, [node.subpageConfigUuid])

    const { IconComponent, themeIconColor } = useMemo(() => {
        let IconComponent: React.ComponentType<{ size: number }>
        let themeIconColor: ThemeIconProps['color'] = 'red'

        if (isConfigMissing) {
            IconComponent = PiWarningCircle
            themeIconColor = 'red'
            return { IconComponent, themeIconColor }
        }

        if (node.isDisabled) {
            IconComponent = TbWifiOff
            themeIconColor = 'gray'
            return { IconComponent, themeIconColor }
        }

        if (node.isConnected) {
            IconComponent = TbWifi
            themeIconColor = 'teal'
        } else if (node.isConnecting) {
            IconComponent = PiCloudArrowUpDuotone
            themeIconColor = 'yellow'
        } else {
            IconComponent = PiWarningCircle
            themeIconColor = 'red'
        }

        return { IconComponent, themeIconColor }
    }, [node.isConnected, node.isConnecting, node.isDisabled, isConfigMissing])

    const handleToggleNodeStatus = () => {
        if (node.isDisabled) {
            enableNode({})
        } else {
            disableNode({})
        }
    }

    return (
        <SectionCard.Root>
            <SectionCard.Section>
                <Group align="flex-center" justify="space-between">
                    <BaseOverlayHeader
                        IconComponent={IconComponent}
                        iconColor={themeIconColor}
                        iconSize={20}
                        iconVariant="soft"
                        title={t('node-details-card.widget.node-details')}
                        titleOrder={5}
                    />

                    <Group gap="xs">
                        {node.isConnected && (
                            <Tooltip
                                label={t(
                                    'node-stats.card.represents-the-uptime-of-the-subscription-node-container'
                                )}
                            >
                                <Badge
                                    color="teal"
                                    h={28}
                                    leftSection={<Logo size={14} />}
                                    size="lg"
                                    variant="light"
                                    visibleFrom="sm"
                                >
                                    {getSingboxUptimeUtil(nodeSingboxUptime)}
                                </Badge>
                            </Tooltip>
                        )}
                        {!isConfigMissing && (
                            <Tooltip
                                label={
                                    node.isDisabled
                                        ? t('node-details-card.widget.enable-node')
                                        : t('node-details-card.widget.disable-node')
                                }
                            >
                                <ActionIcon
                                    color={node.isDisabled ? 'teal' : 'red'}
                                    disabled={isDisableNodePending || isEnableNodePending}
                                    onClick={handleToggleNodeStatus}
                                    size="md"
                                    style={{
                                        backgroundColor: node.isDisabled
                                            ? 'rgba(45, 212, 191, 0.15)'
                                            : 'rgba(239, 68, 68, 0.15)',
                                        border: `1px solid ${
                                            node.isDisabled
                                                ? 'rgba(45, 212, 191, 0.3)'
                                                : 'rgba(239, 68, 68, 0.3)'
                                        }`,
                                        boxShadow: `0 0 10px ${
                                            node.isDisabled
                                                ? 'rgba(45, 212, 191, 0.2)'
                                                : 'rgba(239, 68, 68, 0.2)'
                                        }`
                                    }}
                                    variant="light"
                                >
                                    {isDisableNodePending || isEnableNodePending ? (
                                        <Loader
                                            color={node.isDisabled ? 'teal' : 'red'}
                                            size="xs"
                                        />
                                    ) : (
                                        <TbPower
                                            size={16}
                                            style={{
                                                color: node.isDisabled
                                                    ? 'var(--mantine-color-teal-4)'
                                                    : 'var(--mantine-color-red-4)'
                                            }}
                                        />
                                    )}
                                </ActionIcon>
                            </Tooltip>
                        )}

                        {isConfigMissing && (
                            <Tooltip label={t('base-node-form.select-subscription-profile')}>
                                <ActionIcon
                                    color="gray"
                                    disabled
                                    size="md"
                                    style={{
                                        backgroundColor: 'rgba(107, 114, 128, 0.15)',
                                        border: `1px solid rgba(107, 114, 128, 0.3)`,
                                        boxShadow: `0 0 10px rgba(107, 114, 128, 0.2)`,
                                        opacity: 0.7
                                    }}
                                    variant="light"
                                >
                                    <TbPower
                                        size={16}
                                        style={{
                                            color: 'var(--mantine-color-teal-4)'
                                        }}
                                    />
                                </ActionIcon>
                            </Tooltip>
                        )}
                    </Group>
                </Group>
            </SectionCard.Section>

            {node.isConnected && (
                <SectionCard.Section>
                    <SimpleGrid
                        cols={{
                            base: 1,
                            xs: 2
                        }}
                        spacing="xs"
                    >
                        {nodeSingboxUptime !== 0 && (
                            <Paper
                                hiddenFrom="sm"
                                p="xs"
                                radius="md"
                                style={{
                                    background: 'rgba(20, 184, 166, 0.08)', // teal-500 at 8%
                                    border: '1px solid rgba(20, 184, 166, 0.2)' // teal-500 at 20%
                                }}
                            >
                                <Tooltip
                                    label={t(
                                        'node-stats.card.represents-the-uptime-of-the-subscription-node-container'
                                    )}
                                >
                                    <Group gap="xs" justify="center">
                                        <Logo color="var(--mantine-color-teal-5)" size={16} />
                                        <Text
                                            c="teal.5"
                                            fw={600}
                                            size="sm"
                                            style={{ textTransform: 'uppercase' }}
                                        >
                                            {getSingboxUptimeUtil(nodeSingboxUptime)}
                                        </Text>
                                    </Group>
                                </Tooltip>
                            </Paper>
                        )}

                        {node.nodeVersion && (
                            <Paper
                                p="xs"
                                radius="md"
                                style={{
                                    background: 'rgba(99, 102, 241, 0.08)',
                                    border: '1px solid rgba(99, 102, 241, 0.2)'
                                }}
                            >
                                <Tooltip label={t('node-details-card.widget.remnawave-node-version')}>
                                    <Group gap="xs" justify="center">
                                        <Logo color="var(--mantine-color-indigo-5)" size={16} />
                                        <Text c="indigo.5" fw={600} size="sm">
                                            {node.nodeVersion}
                                        </Text>
                                    </Group>
                                </Tooltip>
                            </Paper>
                        )}
                    </SimpleGrid>
                </SectionCard.Section>
            )}
        </SectionCard.Root>
    )
})
