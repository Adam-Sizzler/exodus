import {
    ActionIcon,
    Badge,
    Box,
    Center,
    Checkbox,
    Divider,
    Group,
    MantineStyleProp,
    Stack,
    Text,
    ThemeIcon
} from '@mantine/core'
import { GetInfraBillingNodesCommand } from '@exodus/backend-contract'
import dayjs from 'dayjs'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { TbCalendar, TbCheck, TbCpu, TbCreditCard, TbExternalLink, TbServer } from 'react-icons/tb'

import { queryClient } from '@shared/api'
import { QueryKeys, useUpdateInfraBillingNode } from '@shared/api/hooks'
import { BaseOverlayHeader } from '@shared/ui/overlays/base-overlay-header'
import { SectionCard } from '@shared/ui/section-card'
import { formatTimeUtil } from '@shared/utils/time-utils'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'

type BillingNode = GetInfraBillingNodesCommand.Response['response']['billingNodes'][number]

function getNodeStatus(nextBillingAt: Date, language: string) {
    const now = dayjs().startOf('day')
    const target = dayjs(nextBillingAt).startOf('day')
    const isOverdue = target.isBefore(now)

    return {
        label: target.locale(language).fromNow(),
        color: isOverdue ? 'red' : 'teal',
        isOverdue
    }
}

interface IProps {
    nodes: GetInfraBillingNodesCommand.Response['response']['billingNodes']
    onToggleSelect?: (uuid: string) => void
    selectedUuids?: Set<string>
    selectMode?: boolean
    style: MantineStyleProp
}

export function MobileNodesListWidget(props: IProps) {
    const { nodes, style, selectMode = false, selectedUuids, onToggleSelect } = props
    const openModalWithData = useModalsStoreOpenWithData()
    const { t, i18n } = useTranslation()
    const [updatingUuids, setUpdatingUuids] = useState<Set<string>>(new Set())

    const { mutate: updateNode } = useUpdateInfraBillingNode({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(QueryKeys.infraBilling.getInfraBillingNodes.queryKey, data)
                setUpdatingUuids(new Set())
            },
            onError: () => setUpdatingUuids(new Set())
        }
    })

    const handleQuickUpdate = (uuid: string, currentDate: Date) => {
        setUpdatingUuids((prev) => new Set(prev).add(uuid))
        updateNode({
            variables: {
                uuids: [uuid],
                // @ts-expect-error - TODO: fix ZOD schema
                nextBillingAt: dayjs(currentDate).add(1, 'month').toISOString()
            }
        })
    }

    const handleClickBillingAt = (node: BillingNode) => {
        openModalWithData(MODALS.UPDATE_BILLING_DATE_MODAL, {
            uuids: [node.uuid],
            nextBillingAt: node.nextBillingAt
        })
    }

    const handleOpenNode = (node: BillingNode) => {
        if (!node.nodeUuid) return

        openModalWithData(MODALS.EDIT_NODE_BY_UUID_MODAL, {
            nodeUuid: node.nodeUuid
        })
    }

    const handleOpenProvider = (node: BillingNode) => {
        if (node.provider.loginUrl) {
            window.open(node.provider.loginUrl, '_blank', 'noopener,noreferrer')
        }
    }

    const groupedByMonth = useMemo((): { label: string; nodes: BillingNode[] }[] => {
        const sorted = [...nodes].sort(
            (a, b) => dayjs(a.nextBillingAt).valueOf() - dayjs(b.nextBillingAt).valueOf()
        )

        const groups = new Map<string, BillingNode[]>()

        for (const node of sorted) {
            const key = dayjs.utc(node.nextBillingAt).format('YYYY-MM')
            const existing = groups.get(key)
            if (existing) {
                existing.push(node)
            } else {
                groups.set(key, [node])
            }
        }

        return Array.from(groups.entries()).map(([key, groupNodes]) => ({
            label: dayjs.utc(key).locale(i18n.language).format('MMMM YYYY'),
            nodes: groupNodes
        }))
    }, [nodes, i18n.language])

    if (nodes.length === 0) {
        return (
            <SectionCard.Root p="xl">
                <SectionCard.Section>
                    <Center py="xl">
                        <Stack align="center" gap="lg">
                            <ThemeIcon color="gray" radius="xl" size={64} variant="soft">
                                <TbServer size={32} />
                            </ThemeIcon>

                            <Stack align="center" gap="xs">
                                <Text c="dimmed" fw={600} size="md" ta="center">
                                    {t('infra-billing-nodes.widget.no-nodes-found')}
                                </Text>
                            </Stack>
                        </Stack>
                    </Center>
                </SectionCard.Section>
            </SectionCard.Root>
        )
    }

    return (
        <Stack gap="md" style={style}>
            {groupedByMonth.map((group) => (
                <Stack gap="xs" key={group.label}>
                    <Divider
                        label={
                            <Text c="dimmed" fw={600} size="xs" tt="capitalize">
                                {group.label}
                            </Text>
                        }
                        labelPosition="center"
                    />

                    {group.nodes.map((node) => {
                        const status = getNodeStatus(node.nextBillingAt, i18n.language)
                        const isSelected = selectedUuids?.has(node.uuid) ?? false

                        return (
                            <Box
                                key={node.uuid}
                                onClick={selectMode ? () => onToggleSelect?.(node.uuid) : undefined}
                                style={{ cursor: selectMode ? 'pointer' : undefined }}
                            >
                                <SectionCard.Root
                                    style={
                                        selectMode && isSelected
                                            ? { borderColor: 'var(--mantine-color-cyan-filled)' }
                                            : undefined
                                    }
                                >
                                    <SectionCard.Section>
                                        <Group justify="space-between" wrap="nowrap">
                                            <Group gap="sm" wrap="nowrap">
                                                {selectMode && (
                                                    <Checkbox
                                                        checked={isSelected}
                                                        color="cyan"
                                                        readOnly
                                                    />
                                                )}

                                                {node.node && (
                                                    <BaseOverlayHeader
                                                        countryCode={node.node.countryCode}
                                                        hideIcon={true}
                                                        iconColor="blue"
                                                        IconComponent={TbServer}
                                                        iconVariant="soft"
                                                        subtitle={node.provider.name}
                                                        title={node.node.name}
                                                    />
                                                )}

                                                {!node.node && (
                                                    <BaseOverlayHeader
                                                        iconColor="blue"
                                                        IconComponent={TbServer}
                                                        iconVariant="soft"
                                                        subtitle={node.provider.name}
                                                        title={node.name!}
                                                    />
                                                )}
                                            </Group>

                                            {!selectMode && (
                                                <Group gap={4} wrap="nowrap">
                                                    {node.node && (
                                                        <ActionIcon
                                                            color="blue"
                                                            onClick={() => handleOpenNode(node)}
                                                            size="input-xs"
                                                            variant="soft"
                                                        >
                                                            <TbCpu size={18} />
                                                        </ActionIcon>
                                                    )}

                                                    {node.provider.loginUrl && (
                                                        <ActionIcon
                                                            color="gray"
                                                            onClick={() => handleOpenProvider(node)}
                                                            size="input-xs"
                                                            variant="soft"
                                                        >
                                                            <TbExternalLink size={18} />
                                                        </ActionIcon>
                                                    )}

                                                    <ActionIcon
                                                        color="teal"
                                                        loading={updatingUuids.has(node.uuid)}
                                                        onClick={() =>
                                                            handleQuickUpdate(
                                                                node.uuid,
                                                                node.nextBillingAt
                                                            )
                                                        }
                                                        size="input-xs"
                                                        variant="soft"
                                                    >
                                                        <TbCheck size={18} />
                                                    </ActionIcon>
                                                </Group>
                                            )}
                                        </Group>
                                    </SectionCard.Section>

                                    <SectionCard.Section>
                                        <Group justify="space-between" wrap="nowrap">
                                            <Group
                                                gap={6}
                                                onClick={
                                                    selectMode
                                                        ? undefined
                                                        : () => handleClickBillingAt(node)
                                                }
                                                style={{
                                                    cursor: selectMode ? 'inherit' : 'pointer'
                                                }}
                                                wrap="nowrap"
                                            >
                                                <TbCalendar
                                                    color={`var(--mantine-color-${status.color}-5)`}
                                                    size={16}
                                                />
                                                <Text c={status.color} fw={600} size="sm">
                                                    {formatTimeUtil({
                                                        time: node.nextBillingAt,
                                                        template: 'FULL_DATE',
                                                        language: i18n.language,
                                                        utc: true
                                                    })}
                                                </Text>
                                            </Group>

                                            <Badge
                                                color={status.color}
                                                leftSection={<TbCreditCard size={14} />}
                                                radius="sm"
                                                size="sm"
                                                variant="soft"
                                            >
                                                {status.label}
                                            </Badge>
                                        </Group>
                                    </SectionCard.Section>
                                </SectionCard.Root>
                            </Box>
                        )
                    })}
                </Stack>
            ))}
        </Stack>
    )
}
