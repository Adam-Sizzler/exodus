import {
    Affix,
    Badge,
    Button,
    Card,
    Group,
    Modal,
    Paper,
    Stack,
    Text,
    TextInput,
    Title,
    Transition
} from '@mantine/core'
import { useDisclosure } from '@mantine/hooks'
import { modals } from '@mantine/modals'
import { useMemo, useState } from 'react'
import { useTranslation } from 'react-i18next'
import { PiEmpty } from 'react-icons/pi'
import { TbPlus } from 'react-icons/tb'

import {
    useSectionActiveTag,
    useViewPreferencesStoreActions
} from '@entities/dashboard/view-preferences-store'
import {
    QueryKeys,
    useBulkDeleteSRSLists,
    useBulkDisableSRSLists,
    useBulkEnableSRSLists,
    useBulkSetIntervalSRSLists,
    useCheckSRSLists,
    useReorderSRSLists
} from '@shared/api/hooks'
import { queryClient } from '@shared/api/query-client'
import { filterByTag, TagFilterBar } from '@shared/ui'
import { VirtualizedDndGrid } from '@shared/ui/virtualized-dnd-grid'

import { SRSListCardWidget, SRSListItem } from '../srs-list-card'
import { IProps } from './interfaces'

export function SRSListsGridWidget(props: IProps) {
    const { srsLists, onEditItem, onCreateItem } = props
    const { t } = useTranslation()
    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const activeTag = useSectionActiveTag('srsLists')
    const { setSectionActiveTag } = useViewPreferencesStoreActions()

    const [selected, setSelected] = useState<string[]>([])
    const [intervalModalOpened, intervalModalHandlers] = useDisclosure(false)
    const [bulkUpdateInterval, setBulkUpdateInterval] = useState('1d')

    const itemsWithTags = useMemo(
        () =>
            (srsLists ?? []).map((item) => ({
                ...item,
                tags: item.tags ?? []
            })),
        [srsLists]
    )

    const visibleItems = useMemo(
        () => filterByTag(itemsWithTags, activeTag),
        [itemsWithTags, activeTag]
    )

    const selectedSet = useMemo(() => new Set(selected), [selected])
    const selectedCount = selected.length

    const mutationRefetch = async () => {
        await queryClient.refetchQueries({ queryKey: QueryKeys.srsLists.getSRSLists.queryKey })
        setSelected([])
    }

    const { mutate: reorderSRSLists } = useReorderSRSLists({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(QueryKeys.srsLists.getSRSLists.queryKey, data)
            }
        }
    })

    const { mutate: bulkDelete } = useBulkDeleteSRSLists({
        mutationFns: { onSuccess: mutationRefetch }
    })
    const { mutate: bulkEnable } = useBulkEnableSRSLists({
        mutationFns: { onSuccess: mutationRefetch }
    })
    const { mutate: bulkDisable } = useBulkDisableSRSLists({
        mutationFns: { onSuccess: mutationRefetch }
    })
    const { mutate: bulkSetInterval } = useBulkSetIntervalSRSLists({
        mutationFns: {
            onSuccess: async () => {
                intervalModalHandlers.close()
                await mutationRefetch()
            }
        }
    })

    const { mutate: checkLists } = useCheckSRSLists({
        mutationFns: {
            onSuccess: async () => {
                await queryClient.refetchQueries({
                    queryKey: QueryKeys.srsLists.getSRSLists.queryKey
                })
            }
        }
    })

    const handleReorder = (reorderedItems: typeof visibleItems) => {
        reorderSRSLists({
            variables: {
                items: reorderedItems.map((item, index) => ({
                    uuid: item.uuid,
                    viewPosition: index
                }))
            }
        })
    }

    const toggleSelect = (uuid: string) => {
        setSelected((prev) =>
            prev.includes(uuid) ? prev.filter((id) => id !== uuid) : [...prev, uuid]
        )
    }

    const selectAll = () => setSelected(visibleItems.map((item) => item.uuid))
    const clearSelection = () => setSelected([])

    const handleSingleDelete = (uuid: string, name: string) => {
        modals.openConfirmModal({
            centered: true,
            title: tr('common.action.confirm-action', 'Confirm action'),
            children: `${tr('common.action.delete', 'Delete')} "${name}"?`,
            labels: {
                confirm: tr('common.action.delete', 'Delete'),
                cancel: tr('common.action.cancel', 'Cancel')
            },
            cancelProps: { variant: 'subtle' },
            confirmProps: { color: 'red', variant: 'soft' },
            onConfirm: () => bulkDelete({ variables: { uuids: [uuid] } })
        })
    }

    const handleBulkDelete = () => {
        if (selectedCount === 0) return
        modals.openConfirmModal({
            centered: true,
            title: tr('srs-lists.feature.delete-selected-title', 'Delete selected SRS lists'),
            children: `${tr('srs-lists.feature.selected-count', 'Selected')}: ${selectedCount}`,
            labels: {
                confirm: tr('common.action.delete', 'Delete'),
                cancel: tr('common.action.cancel', 'Cancel')
            },
            cancelProps: { variant: 'subtle' },
            confirmProps: { color: 'red', variant: 'soft' },
            onConfirm: () => bulkDelete({ variables: { uuids: selected } })
        })
    }

    const applyBulkInterval = () => {
        if (selectedCount === 0) return
        bulkSetInterval({
            variables: {
                uuids: selected,
                updateInterval: bulkUpdateInterval.trim() || '1d'
            }
        })
    }

    if (!srsLists || srsLists.length === 0) {
        return (
            <Card p="xl" withBorder>
                <Stack align="center" gap="md">
                    <PiEmpty size={48} style={{ opacity: 0.5 }} />
                    <div>
                        <Title c="dimmed" order={4} ta="center">
                            {tr('srs-lists.feature.no-srs-lists-found', 'SRS lists not found')}
                        </Title>
                        <Text c="dimmed" mt="xs" size="sm" ta="center">
                            {tr(
                                'srs-lists.feature.create-your-first-srs-list-to-get-started',
                                'Add your first SRS list to get started'
                            )}
                        </Text>
                    </div>
                    {onCreateItem && (
                        <Button
                            color="teal"
                            leftSection={<TbPlus size={16} />}
                            onClick={onCreateItem}
                            variant="soft"
                        >
                            {tr('srs-lists.feature.add-links', 'Add SRS Lists')}
                        </Button>
                    )}
                </Stack>
            </Card>
        )
    }

    return (
        <>
            <VirtualizedDndGrid
                enableDnd={activeTag === null}
                header={
                    <TagFilterBar
                        activeTag={activeTag}
                        items={itemsWithTags}
                        onChange={(tag) => setSectionActiveTag('srsLists', tag)}
                    />
                }
                items={visibleItems}
                onReorder={handleReorder}
                renderDragOverlay={(item) => (
                    <SRSListCardWidget
                        disableReordering
                        isDragOverlay
                        isSelected={selectedSet.has(item.uuid)}
                        onCheck={(uuid) => checkLists({ variables: { uuids: [uuid] } })}
                        onDelete={handleSingleDelete}
                        onEdit={onEditItem}
                        onToggleEnable={(uuid, enabled) =>
                            enabled
                                ? bulkEnable({ variables: { uuids: [uuid] } })
                                : bulkDisable({ variables: { uuids: [uuid] } })
                        }
                        srsList={item}
                    />
                )}
                renderItem={(item) => (
                    <SRSListCardWidget
                        disableReordering={activeTag !== null}
                        isSelected={selectedSet.has(item.uuid)}
                        onCheck={(uuid) => checkLists({ variables: { uuids: [uuid] } })}
                        onDelete={handleSingleDelete}
                        onEdit={onEditItem}
                        onToggleEnable={(uuid, enabled) =>
                            enabled
                                ? bulkEnable({ variables: { uuids: [uuid] } })
                                : bulkDisable({ variables: { uuids: [uuid] } })
                        }
                        onToggleSelect={toggleSelect}
                        srsList={item}
                    />
                )}
                useWindowScroll={true}
            />

            <Affix position={{ bottom: 20, right: 20 }} zIndex={100}>
                <Transition mounted={selectedCount > 0} transition="slide-up">
                    {(styles) => (
                        <Paper
                            p="md"
                            shadow="md"
                            style={{
                                ...styles,
                                width: '320px',
                                maxWidth: '100vw'
                            }}
                            withBorder
                        >
                            <Stack gap="sm">
                                <Group justify="space-between">
                                    <Badge color="blue" size="lg" variant="light">
                                        {tr('srs-lists.feature.selected-count', 'Selected')}:{' '}
                                        {selectedCount}
                                    </Badge>
                                    <Group gap="xs">
                                        <Button
                                            onClick={clearSelection}
                                            size="xs"
                                            variant="subtle"
                                        >
                                            {tr('common.action.clear', 'Clear')}
                                        </Button>
                                        <Button
                                            onClick={selectAll}
                                            size="xs"
                                            variant="subtle"
                                        >
                                            {tr('common.action.select-all', 'Select all')}
                                        </Button>
                                    </Group>
                                </Group>

                                <Group grow gap="xs">
                                    <Button
                                        color="teal"
                                        onClick={() =>
                                            bulkEnable({ variables: { uuids: selected } })
                                        }
                                        size="xs"
                                        variant="soft"
                                    >
                                        {tr('common.enable', 'Enable')}
                                    </Button>
                                    <Button
                                        color="orange"
                                        onClick={() =>
                                            bulkDisable({ variables: { uuids: selected } })
                                        }
                                        size="xs"
                                        variant="soft"
                                    >
                                        {tr('common.disable', 'Disable')}
                                    </Button>
                                    <Button
                                        color="blue"
                                        onClick={intervalModalHandlers.open}
                                        size="xs"
                                        variant="soft"
                                    >
                                        {tr('srs-lists.feature.set-update-interval', 'Interval')}
                                    </Button>
                                    <Button
                                        color="red"
                                        onClick={handleBulkDelete}
                                        size="xs"
                                        variant="soft"
                                    >
                                        {tr('common.delete', 'Delete')}
                                    </Button>
                                </Group>
                            </Stack>
                        </Paper>
                    )}
                </Transition>
            </Affix>

            <Modal
                centered
                onClose={intervalModalHandlers.close}
                opened={intervalModalOpened}
                size="sm"
                title={tr('srs-lists.feature.set-update-interval', 'Set update interval')}
            >
                <Stack gap="md">
                    <TextInput
                        description={tr(
                            'srs-lists.feature.update-interval-hint',
                            'e.g. 10m, 1h, 12h, 1d, 7d'
                        )}
                        label={tr('srs-lists.feature.update-interval', 'Update interval')}
                        onChange={(event) => setBulkUpdateInterval(event.currentTarget.value)}
                        placeholder="1d"
                        value={bulkUpdateInterval}
                    />
                    <Button fullWidth onClick={applyBulkInterval}>
                        {tr('common.save', 'Save')}
                    </Button>
                </Stack>
            </Modal>
        </>
    )
}
