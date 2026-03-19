import {
    closestCenter,
    DndContext,
    DragEndEvent,
    DragOverlay,
    DragStartEvent,
    KeyboardSensor,
    MouseSensor,
    TouchSensor,
    UniqueIdentifier,
    useSensor,
    useSensors
} from '@dnd-kit/core'
import { arrayMove, SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import {
    PiEmpty,
    PiInfo,
    PiListChecks,
    PiMagnifyingGlass,
    PiProhibit,
    PiPulse,
    PiTag,
    PiWarningCircle
} from 'react-icons/pi'
import { CSS } from '@dnd-kit/utilities'
import {
    Accordion,
    ActionIcon,
    ActionIconGroup,
    Affix,
    Badge,
    Box,
    Button,
    Card,
    Checkbox,
    Container,
    Divider,
    Group,
    Modal,
    Paper,
    Select,
    Stack,
    Text,
    TextInput,
    Textarea,
    ThemeIcon,
    Title,
    Tooltip,
    Transition
} from '@mantine/core'
import { useDisclosure, useMediaQuery } from '@mantine/hooks'
import { useWindowVirtualizer } from '@tanstack/react-virtual'
import { modals } from '@mantine/modals'
import { useTranslation } from 'node_modules/react-i18next'
import { useCallback, useEffect, useMemo, useRef, useState } from 'react'
import { TbPlus, TbRefresh, TbTagStarred } from 'react-icons/tb'
import { useSortable } from '@dnd-kit/sortable'
import { restrictToVerticalAxis } from '@dnd-kit/modifiers'
import { HiFilter } from 'react-icons/hi'
import { RiDraggable } from 'react-icons/ri'
import { motion } from 'framer-motion'
import cx from 'clsx'

import {
    QueryKeys,
    useBulkDeleteSRSLists,
    useBulkDisableSRSLists,
    useBulkEnableSRSLists,
    useBulkSetIntervalSRSLists,
    useCheckSRSLists,
    useCreateSRSLists,
    useReorderSRSLists,
    useUpdateSRSList
} from '@shared/api/hooks'
import { UniversalSpotlightActionIconShared } from '@shared/ui/universal-spotlight'
import { UniversalSpotlightContentShared } from '@shared/ui/universal-spotlight/universal-spotlight-content.shared'
import { queryClient } from '@shared/api/query-client'
import { DataTableShared } from '@shared/ui/table'
import { Page } from '@shared/ui/page'

import classes from './SRSCard.module.css'

type SRSList = {
    uuid: string
    tag: string
    url: string
    updateInterval: string
    fileName: string
    shortName: string
    viewPosition: number
    isEnabled: boolean
    isAvailable: boolean
    lastError?: null | string
}

interface Props {
    srsLists: SRSList[]
}

function SRSListsIcon(props: { size?: number }) {
    const { size = 20 } = props

    return (
        <svg fill="currentColor" height={size} viewBox="0 0 256 256" width={size}>
            <path d="M224,128a8,8,0,0,1-8,8H128a8,8,0,0,1,0-16h88A8,8,0,0,1,224,128ZM128,72h88a8,8,0,0,0,0-16H128a8,8,0,0,0,0,16Zm88,112H128a8,8,0,0,0,0,16h88a8,8,0,0,0,0-16ZM82.34,42.34,56,68.69,45.66,58.34A8,8,0,0,0,34.34,69.66l16,16a8,8,0,0,0,11.32,0l32-32A8,8,0,0,0,82.34,42.34Zm0,64L56,132.69,45.66,122.34a8,8,0,0,0-11.32,11.32l16,16a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Zm0,64L56,196.69,45.66,186.34a8,8,0,0,0-11.32,11.32l16,16a8,8,0,0,0,11.32,0l32-32a8,8,0,0,0-11.32-11.32Z" />
        </svg>
    )
}

function SRSListCard(props: {
    highlightedUuid: null | string
    isDragOverlay?: boolean
    isSelected: boolean
    item: SRSList
    onEdit: (item: SRSList) => void
    onSelect: () => void
    selectedNameUuid: null | string
    selectedTag: null | string
    selectedURLUuid: null | string
}) {
    const {
        item,
        isSelected,
        onSelect,
        onEdit,
        isDragOverlay = false,
        selectedTag,
        selectedNameUuid,
        selectedURLUuid,
        highlightedUuid
    } = props

    const { t } = useTranslation()
    const [isHovered, setIsHovered] = useState(false)

    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const isFiltered =
        (!!selectedTag && item.tag !== selectedTag) ||
        (!!selectedNameUuid && item.uuid !== selectedNameUuid) ||
        (!!selectedURLUuid && item.uuid !== selectedURLUuid)

    const { attributes, listeners, setNodeRef, transform, transition, isDragging } = useSortable({
        id: item.uuid,
        disabled: isDragOverlay || isFiltered
    })

    const stateColor = !item.isEnabled ? 'gray' : item.isAvailable ? 'teal' : 'red'
    const stateIcon = !item.isEnabled ? (
        <PiProhibit size={16} style={{ color: 'var(--mantine-color-gray-6)' }} />
    ) : item.isAvailable ? (
        <PiPulse size={16} style={{ color: 'var(--mantine-color-teal-6)' }} />
    ) : (
        <PiWarningCircle size={18} style={{ color: 'var(--mantine-color-red-3)' }} />
    )

    return (
        <Box
            className={cx(classes.item, {
                [classes.itemDragging]: isDragging || isHovered,
                [classes.selectedItem]: isSelected,
                [classes.filteredItem]: isFiltered,
                [classes.highlightedItem]: highlightedUuid === item.uuid,
                [classes.itemAvailable]: item.isEnabled && item.isAvailable,
                [classes.itemDisabled]: !item.isEnabled,
                [classes.itemUnavailable]: item.isEnabled && !item.isAvailable
            })}
            data-dnd-overlay={isDragOverlay}
            ref={isDragOverlay ? undefined : setNodeRef}
            style={{
                transform: CSS.Transform.toString(transform),
                transition,
                opacity: isDragging ? 0 : 1,
                zIndex: isDragging ? 1000 : 'auto',
                position: 'relative'
            }}
        >
            <Group gap="md" w="100%" wrap="nowrap">
                <Group gap="xs" wrap="nowrap">
                    <Checkbox checked={isSelected} onChange={onSelect} size="md" />
                    <Box
                        {...(isDragOverlay ? {} : attributes)}
                        {...(isDragOverlay ? {} : listeners)}
                        className={classes.dragHandle}
                    >
                        <RiDraggable color="white" size="24px" />
                    </Box>
                </Group>

                <Box
                    className={classes.contentArea}
                    flex={1}
                    onClick={() => onEdit(item)}
                    onMouseEnter={() => setIsHovered(true)}
                    onMouseLeave={() => setIsHovered(false)}
                    onTouchEnd={() => setIsHovered(false)}
                    onTouchStart={() => setIsHovered(true)}
                >
                    <Group gap="md" justify="space-between" wrap="nowrap">
                        <Group
                            flex={1}
                            gap="sm"
                            miw={0}
                            style={{ overflow: 'hidden' }}
                            wrap="nowrap"
                        >
                            <Tooltip
                                label={
                                    !item.isEnabled
                                        ? tr('srs-lists.feature.disabled', 'Disabled')
                                        : item.isAvailable
                                          ? tr('srs-lists.feature.available', 'Available')
                                          : item.lastError
                                            ? `${tr('srs-lists.feature.unavailable', 'Unavailable')}. ${tr('srs-lists.feature.last-error', 'Error')}: ${item.lastError}`
                                            : tr('srs-lists.feature.unavailable', 'Unavailable')
                                }
                                multiline
                                maw={460}
                                withArrow
                            >
                                <ActionIcon
                                    color={stateColor}
                                    size="lg"
                                    style={{ cursor: 'default', flexShrink: 0 }}
                                    variant="light"
                                >
                                    {stateIcon}
                                </ActionIcon>
                            </Tooltip>

                            <Group gap="md" style={{ flexShrink: 1, minWidth: 0 }} wrap="nowrap">
                                <Text fw={600} style={{ flexShrink: 0 }} truncate>
                                    {item.tag || item.shortName}
                                </Text>
                                <Text
                                    c="dimmed"
                                    className={classes.hostAddress}
                                    style={{ flexShrink: 1, minWidth: 0 }}
                                    truncate
                                >
                                    {item.url}
                                </Text>
                            </Group>
                        </Group>

                        <Group className={classes.rightMeta} gap="md" wrap="nowrap">
                            <Badge className={classes.intervalBadge} color="gray" leftSection={<PiTag size={12} />} size="md" variant="outline">
                                {item.updateInterval}
                            </Badge>
                        </Group>
                    </Group>
                </Box>
            </Group>
        </Box>
    )
}

export function SRSListsPageComponent(props: Props) {
    const { srsLists } = props
    const { t } = useTranslation()

    const tr = (key: string, defaultValue: string) => t(key, { defaultValue })

    const [opened, handlers] = useDisclosure(false)
    const [intervalModalOpened, intervalModalHandlers] = useDisclosure(false)
    const [editModalOpened, editModalHandlers] = useDisclosure(false)

    const [state, setState] = useState<SRSList[]>(srsLists || [])
    const [draggedItem, setDraggedItem] = useState<null | SRSList>(null)
    const [selected, setSelected] = useState<string[]>([])
    const [editingItem, setEditingItem] = useState<null | SRSList>(null)

    const [selectedTag, setSelectedTag] = useState<null | string>(null)
    const [selectedNameUuid, setSelectedNameUuid] = useState<null | string>(null)
    const [selectedURLUuid, setSelectedURLUuid] = useState<null | string>(null)
    const [highlightedUuid, setHighlightedUuid] = useState<null | string>(null)

    const [urlsText, setUrlsText] = useState('')
    const [updateInterval, setUpdateInterval] = useState('1d')
    const [createEnabled, setCreateEnabled] = useState(true)
    const [bulkUpdateInterval, setBulkUpdateInterval] = useState('1d')
    const [editURL, setEditURL] = useState('')
    const [editUpdateInterval, setEditUpdateInterval] = useState('1d')
    const [editEnabled, setEditEnabled] = useState(true)

    const listRef = useRef<HTMLDivElement | null>(null)
    const isMobile = useMediaQuery('(max-width: 48em)')

    const virtualizer = useWindowVirtualizer({
        count: state.length,
        estimateSize: () => 68,
        overscan: 5,
        scrollMargin: listRef.current?.offsetTop ?? 0,
        getItemKey: (index) => state[index].uuid
    })

    const dataIds = useRef<UniqueIdentifier[]>([])
    dataIds.current = state.map((item) => item.uuid)

    const { mutate: createSRSLists, isPending: isCreating } = useCreateSRSLists({
        mutationFns: {
            onSuccess: async () => {
                handlers.close()
                setUrlsText('')
                setUpdateInterval('1d')
                setCreateEnabled(true)
                await queryClient.refetchQueries({ queryKey: QueryKeys.srsLists.getSRSLists.queryKey })
            }
        }
    })

    const { mutate: reorderSRSLists } = useReorderSRSLists({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(QueryKeys.srsLists.getSRSLists.queryKey, data)
            }
        }
    })

    const mutationRefetch = async () => {
        await queryClient.refetchQueries({ queryKey: QueryKeys.srsLists.getSRSLists.queryKey })
        setSelected([])
    }

    const { mutate: bulkDelete } = useBulkDeleteSRSLists({ mutationFns: { onSuccess: mutationRefetch } })
    const { mutate: bulkEnable } = useBulkEnableSRSLists({ mutationFns: { onSuccess: mutationRefetch } })
    const { mutate: bulkDisable } = useBulkDisableSRSLists({ mutationFns: { onSuccess: mutationRefetch } })
    const { mutate: bulkSetInterval } = useBulkSetIntervalSRSLists({
        mutationFns: {
            onSuccess: async () => {
                intervalModalHandlers.close()
                await mutationRefetch()
            }
        }
    })

    const { mutate: checkLists, isPending: isChecking } = useCheckSRSLists({
        mutationFns: {
            onSuccess: async () => {
                await queryClient.refetchQueries({ queryKey: QueryKeys.srsLists.getSRSLists.queryKey })
            }
        }
    })

    const { mutate: updateSRSList, isPending: isUpdating } = useUpdateSRSList({
        mutationFns: {
            onSuccess: async () => {
                editModalHandlers.close()
                setEditingItem(null)
                await mutationRefetch()
            }
        }
    })

    useEffect(() => {
        setState(srsLists || [])
    }, [srsLists])

    useEffect(() => {
        setSelected((prev) => prev.filter((id) => srsLists.some((item) => item.uuid === id)))
    }, [srsLists])

    useEffect(() => {
        if (highlightedUuid) {
            const timeout = setTimeout(() => setHighlightedUuid(null), 2000)
            return () => clearTimeout(timeout)
        }

        return undefined
    }, [highlightedUuid])

    const sensors = useSensors(
        useSensor(MouseSensor, { activationConstraint: { distance: 5 } }),
        useSensor(TouchSensor, { activationConstraint: { delay: 250, tolerance: 5 } }),
        useSensor(KeyboardSensor, {})
    )

    const selectedCount = selected.length

    const selectedSet = useMemo(() => new Set(selected), [selected])

    const tagOptions = useMemo(() => {
        const uniq = Array.from(new Set(state.map((item) => item.tag).filter(Boolean)))
        return [{ value: '', label: tr('srs-lists.feature.all-tags', 'All tags') }, ...uniq.map((value) => ({ value, label: value }))]
    }, [state])

    const searchNameOptions = useMemo(
        () =>
            state.map((item) => ({
                value: item.uuid,
                label: `${item.shortName} (${item.tag})`
            })),
        [state]
    )

    const searchURLOptions = useMemo(
        () =>
            state.map((item) => ({
                value: item.uuid,
                label: item.url
            })),
        [state]
    )

    const toggleSelect = useCallback((uuid: string) => {
        setSelected((prev) =>
            prev.includes(uuid) ? prev.filter((value) => value !== uuid) : [...prev, uuid]
        )
    }, [])

    const selectAll = () => setSelected(state.map((item) => item.uuid))
    const clearSelection = () => setSelected([])

    const handleSearchSelect = (value: null | string) => {
        setSelectedNameUuid(value)
        if (!value) return
        const index = state.findIndex((item) => item.uuid === value)
        if (index >= 0) {
            virtualizer.scrollToIndex(index, { align: 'center', behavior: 'smooth' })
            setHighlightedUuid(value)
        }
    }

    const handleURLSearchSelect = (value: null | string) => {
        setSelectedURLUuid(value)
        if (!value) return
        const index = state.findIndex((item) => item.uuid === value)
        if (index >= 0) {
            virtualizer.scrollToIndex(index, { align: 'center', behavior: 'smooth' })
            setHighlightedUuid(value)
        }
    }

    const handleTagFilter = (value: null | string) => setSelectedTag(value || null)

    const onSubmitCreate = () => {
        const urls = urlsText
            .split(/\n|,|\s+/g)
            .map((value) => value.trim())
            .filter(Boolean)

        if (urls.length === 0) {
            return
        }

        createSRSLists({
            variables: {
                urls,
                updateInterval: updateInterval.trim() || '1d',
                isEnabled: createEnabled
            }
        })
    }

    const handleDragStart = useCallback(
        (event: DragStartEvent) => {
            const found = state.find((item) => item.uuid === event.active.id)
            setDraggedItem(found || null)
        },
        [state]
    )

    const handleDragEnd = useCallback(
        (event: DragEndEvent) => {
            const { active, over } = event

            if (!over || active.id === over.id) {
                setDraggedItem(null)
                return
            }

            const oldIndex = dataIds.current.indexOf(active.id)
            const newIndex = dataIds.current.indexOf(over.id)

            if (oldIndex !== -1 && newIndex !== -1) {
                const next = arrayMove(state, oldIndex, newIndex)
                setState(next)
                reorderSRSLists({
                    variables: {
                        items: next.map((item, index) => ({ uuid: item.uuid, viewPosition: index }))
                    }
                })
            }

            setDraggedItem(null)
        },
        [state, reorderSRSLists]
    )

    const deleteSelected = () => {
        if (selectedCount === 0) return

        modals.openConfirmModal({
            centered: true,
            title: tr('srs-lists.feature.delete-selected-title', 'Delete selected SRS lists'),
            children: `${tr('srs-lists.feature.selected-count', 'Selected')}: ${selectedCount}`,
            labels: {
                confirm: tr('common.delete', 'Delete'),
                cancel: tr('common.cancel', 'Cancel')
            },
            confirmProps: { color: 'red' },
            onConfirm: () => bulkDelete({ variables: { uuids: selected } })
        })
    }

    const applyBulkInterval = () => {
        const value = bulkUpdateInterval.trim()
        if (!value || selectedCount === 0) return

        bulkSetInterval({ variables: { uuids: selected, updateInterval: value } })
    }

    const openEditModal = (item: SRSList) => {
        setEditingItem(item)
        setEditURL(item.url)
        setEditUpdateInterval(item.updateInterval || '1d')
        setEditEnabled(item.isEnabled)
        editModalHandlers.open()
    }

    const onSubmitEdit = () => {
        if (!editingItem) return

        const url = editURL.trim()
        if (!url) return

        updateSRSList({
            variables: {
                uuid: editingItem.uuid,
                url,
                updateInterval: editUpdateInterval.trim() || '1d',
                isEnabled: editEnabled
            }
        })
    }

    const spotlightActions = state.map((item) => ({
        id: item.uuid,
        label: item.shortName,
        keywords: [item.url, item.tag, item.fileName, item.updateInterval],
        onClick: () => {
            const index = state.findIndex((row) => row.uuid === item.uuid)
            if (index >= 0) {
                virtualizer.scrollToIndex(index, { align: 'center', behavior: 'smooth' })
                setHighlightedUuid(item.uuid)
            }
        },
        description: item.url
    }))

    return (
        <Page title={tr('constants.srs-lists', 'SRS Lists')}>
            <Stack gap="md">
                <DataTableShared.Container>
                    <DataTableShared.Title
                        actions={
                            <Group grow preventGrowOverflow={false} wrap="wrap">
                                <UniversalSpotlightActionIconShared />

                                <ActionIconGroup>
                                    <Tooltip label={tr('common.update', 'Update')} withArrow>
                                        <ActionIcon
                                            loading={isChecking}
                                            onClick={() => checkLists({ variables: {} })}
                                            size="input-md"
                                            variant="light"
                                        >
                                            <TbRefresh size="24px" />
                                        </ActionIcon>
                                    </Tooltip>
                                </ActionIconGroup>

                                <ActionIconGroup>
                                    <Tooltip
                                        label={tr('srs-lists.feature.add-links', 'Add links')}
                                        withArrow
                                    >
                                        <ActionIcon
                                            color="teal"
                                            onClick={handlers.open}
                                            size="input-md"
                                            variant="light"
                                        >
                                            <TbPlus size="24px" />
                                        </ActionIcon>
                                    </Tooltip>
                                </ActionIconGroup>
                            </Group>
                        }
                        icon={<PiListChecks size={24} />}
                        title={tr('constants.srs-lists', 'SRS Lists')}
                    />

                    <DataTableShared.Content>
                        <Accordion
                            bg="linear-gradient(135deg, var(--mantine-color-dark-6) 0%, var(--mantine-color-dark-7) 100%)"
                            variant="filled"
                        >
                            <Accordion.Item value="filters">
                                <Accordion.Control component="a">
                                    <Group align="center" gap="md" wrap="nowrap">
                                        <ActionIcon color="gray" size="input-sm" variant="light">
                                            <HiFilter size={20} />
                                        </ActionIcon>
                                        <Title fw={500} fz="md" order={4}>
                                            {tr('srs-lists.feature.filters', 'Filters')}
                                        </Title>
                                    </Group>
                                </Accordion.Control>
                                <Accordion.Panel>
                                    <Stack gap="md">
                                        <Group grow preventGrowOverflow={false} wrap="wrap">
                                            <Select
                                                clearable
                                                data={tagOptions}
                                                leftSection={<TbTagStarred size="16px" />}
                                                leftSectionPointerEvents="none"
                                                onChange={handleTagFilter}
                                                placeholder={tr('srs-lists.feature.filter-by-tags', 'Filter by tags')}
                                                size="sm"
                                                value={selectedTag || ''}
                                            />

                                            <Select
                                                clearable
                                                data={searchNameOptions}
                                                leftSection={<PiMagnifyingGlass size={16} />}
                                                onChange={handleSearchSelect}
                                                placeholder={tr('srs-lists.feature.search-by-name', 'Search by name')}
                                                searchable
                                                value={selectedNameUuid}
                                            />

                                            <Select
                                                clearable
                                                data={searchURLOptions}
                                                leftSection={<PiMagnifyingGlass size={16} />}
                                                onChange={handleURLSearchSelect}
                                                placeholder={tr('srs-lists.feature.search-by-url', 'Search by URL')}
                                                searchable
                                                value={selectedURLUuid}
                                            />
                                        </Group>
                                    </Stack>
                                </Accordion.Panel>
                            </Accordion.Item>
                        </Accordion>
                    </DataTableShared.Content>
                </DataTableShared.Container>

                {state.length === 0 && (
                    <Card p="xl" withBorder>
                        <Stack align="center" gap="md">
                            <PiEmpty size={48} style={{ opacity: 0.5 }} />
                            <div>
                                <Title c="dimmed" order={4} ta="center">
                                    {tr('srs-lists.feature.no-srs-lists-found', 'SRS списки не найдены')}
                                </Title>
                                <Text c="dimmed" mt="xs" size="sm" ta="center">
                                    {tr(
                                        'srs-lists.feature.create-your-first-srs-list-to-get-started',
                                        'Создайте свой первый SRS список, чтобы начать работу',
                                    )}
                                </Text>
                            </div>
                        </Stack>
                    </Card>
                )}

                {state.length > 0 && (
                    <DndContext
                        collisionDetection={closestCenter}
                        modifiers={[restrictToVerticalAxis]}
                        onDragEnd={handleDragEnd}
                        onDragStart={handleDragStart}
                        sensors={sensors}
                    >
                        <div ref={listRef}>
                            <div
                                style={{
                                    height: `${virtualizer.getTotalSize()}px`,
                                    width: '100%',
                                    position: 'relative'
                                }}
                            >
                                <SortableContext
                                    items={dataIds.current}
                                    strategy={verticalListSortingStrategy}
                                >
                                    <Container fluid>
                                        <Stack gap={0}>
                                            {virtualizer.getVirtualItems().map((virtualItem) => {
                                                const item = state[virtualItem.index]
                                                if (!item) return null

                                                return (
                                                    <Box
                                                        data-index={virtualItem.index}
                                                        key={item.uuid}
                                                        style={{
                                                            position: 'absolute',
                                                            marginLeft: isMobile ? '0px' : '16px',
                                                            marginRight: isMobile ? '0px' : '16px',
                                                            top: 0,
                                                            left: 0,
                                                            right: 0,
                                                            transform: `translateY(${virtualItem.start - virtualizer.options.scrollMargin}px)`
                                                        }}
                                                    >
                                                        <motion.div
                                                            animate={{ opacity: 1 }}
                                                            exit={{ opacity: 0 }}
                                                            initial={{ opacity: 0 }}
                                                            transition={{ duration: 0.1 }}
                                                        >
                                                            <SRSListCard
                                                                highlightedUuid={highlightedUuid}
                                                                isSelected={selectedSet.has(item.uuid)}
                                                                item={item}
                                                                onEdit={openEditModal}
                                                                onSelect={() => toggleSelect(item.uuid)}
                                                                selectedNameUuid={selectedNameUuid}
                                                                selectedTag={selectedTag}
                                                                selectedURLUuid={selectedURLUuid}
                                                            />
                                                        </motion.div>
                                                    </Box>
                                                )
                                            })}
                                        </Stack>
                                    </Container>
                                </SortableContext>
                            </div>
                        </div>

                        <DragOverlay>
                            {draggedItem && (
                                <Container fluid pl={0} pr={0}>
                                    <SRSListCard
                                        highlightedUuid={highlightedUuid}
                                        isDragOverlay
                                        isSelected={selectedSet.has(draggedItem.uuid)}
                                        item={draggedItem}
                                        onEdit={openEditModal}
                                        onSelect={() => toggleSelect(draggedItem.uuid)}
                                        selectedNameUuid={selectedNameUuid}
                                        selectedTag={selectedTag}
                                        selectedURLUuid={selectedURLUuid}
                                    />
                                </Container>
                            )}
                        </DragOverlay>
                    </DndContext>
                )}
            </Stack>

            <UniversalSpotlightContentShared actions={spotlightActions} />

            <Affix position={{ bottom: 20, right: 20 }} zIndex={100}>
                <Transition mounted={selectedCount > 0} transition="slide-up">
                    {(styles) => (
                        <Paper
                            p="md"
                            shadow="md"
                            style={{ ...styles, width: '300px', maxWidth: '1200px', margin: '0 auto' }}
                            withBorder
                        >
                            <Stack>
                                <Group justify="flex-start">
                                    <Group justify="center" w="100%">
                                        <Badge color="blue" size="lg">
                                            {tr('srs-lists.feature.selected-count', 'Selected')}: {selectedCount}
                                        </Badge>
                                    </Group>
                                    <Group grow justify="apart" preventGrowOverflow={false} wrap="wrap">
                                        <Button onClick={clearSelection} variant="subtle">
                                            {tr('multi-select-hosts.feature.clear-selection', 'Clear selection')}
                                        </Button>
                                        <Button onClick={selectAll} variant="subtle">
                                            {tr('multi-select-hosts.feature.select-all', 'Select all')}
                                        </Button>
                                    </Group>
                                </Group>

                                <Group grow justify="apart" preventGrowOverflow={false} wrap="wrap">
                                    <Button
                                        color="green"
                                        onClick={() => bulkEnable({ variables: { uuids: selected } })}
                                    >
                                        {tr('common.enable', 'Enable')}
                                    </Button>
                                    <Button
                                        color="gray"
                                        onClick={() => bulkDisable({ variables: { uuids: selected } })}
                                    >
                                        {tr('common.disable', 'Disable')}
                                    </Button>
                                    <Button color="cyan" onClick={intervalModalHandlers.open}>
                                        {tr('srs-lists.feature.set-update-interval', 'Set update interval')}
                                    </Button>
                                    <Button color="red" onClick={deleteSelected}>
                                        {tr('common.delete', 'Delete')}
                                    </Button>
                                </Group>
                            </Stack>
                        </Paper>
                    )}
                </Transition>
            </Affix>

            <Modal
                onClose={handlers.close}
                opened={opened}
                size={792}
                styles={{ body: { minHeight: 648 } }}
                title={
                    <Group gap="sm" wrap="nowrap">
                        <ThemeIcon size="lg" variant="gradient-teal">
                            <SRSListsIcon size={20} />
                        </ThemeIcon>
                        <Title c="white" order={4}>
                            {tr('srs-lists.feature.add-links', 'Add links')}
                        </Title>
                    </Group>
                }
            >
                <Stack>
                    <Card
                        style={{
                            background: 'rgba(255, 255, 255, 0.02)',
                            border: '1px solid rgba(255, 255, 255, 0.08)',
                            padding: 'var(--mantine-spacing-md)'
                        }}
                        withBorder
                    >
                        <Stack>
                            <Group gap="sm" wrap="nowrap">
                                <ThemeIcon size="lg" variant="gradient-violet">
                                    <PiInfo size={18} />
                                </ThemeIcon>
                                <Title c="white" order={5}>
                                    {tr('srs-lists.feature.direct-link-required', 'Direct link required')}
                                </Title>
                            </Group>
                            <Divider opacity={0.3} />
                            <Text c="dimmed" size="sm">
                                {tr(
                                    'srs-lists.feature.direct-url-hint',
                                    'Use a direct .srs file URL (example: https://raw.githubusercontent.com/.../ruleset.srs). Links with /blob/ are pages and will not work as direct file sources.'
                                )}
                            </Text>
                            <Textarea
                                autosize
                                description={tr(
                                    'srs-lists.feature.multiple-urls-hint',
                                    'One or multiple URLs separated by spaces, commas or new lines'
                                )}
                                label={tr('srs-lists.feature.urls', 'URLs')}
                                minRows={10}
                                onChange={(event) => setUrlsText(event.currentTarget.value)}
                                placeholder="https://.../ruleset.srs"
                                value={urlsText}
                            />
                        </Stack>
                    </Card>
                    <Card
                        style={{
                            background: 'rgba(255, 255, 255, 0.02)',
                            border: '1px solid rgba(255, 255, 255, 0.08)',
                            padding: 'var(--mantine-spacing-md)'
                        }}
                        withBorder
                    >
                        <Stack gap="sm">
                            <Group gap="sm" wrap="nowrap">
                                <ThemeIcon size="lg" variant="gradient-orange">
                                    <svg
                                        fill="none"
                                        height="20"
                                        stroke="currentColor"
                                        strokeLinecap="round"
                                        strokeLinejoin="round"
                                        strokeWidth="2"
                                        viewBox="0 0 24 24"
                                        width="20"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path d="M10.325 4.317c.426 -1.756 2.924 -1.756 3.35 0a1.724 1.724 0 0 0 2.573 1.066c1.543 -.94 3.31 .826 2.37 2.37a1.724 1.724 0 0 0 1.065 2.572c1.756 .426 1.756 2.924 0 3.35a1.724 1.724 0 0 0 -1.066 2.573c.94 1.543 -.826 3.31 -2.37 2.37a1.724 1.724 0 0 0 -2.572 1.065c-.426 1.756 -2.924 1.756 -3.35 0a1.724 1.724 0 0 0 -2.573 -1.066c-1.543 .94 -3.31 -.826 -2.37 -2.37a1.724 1.724 0 0 0 -1.065 -2.572c-1.756 -.426 -1.756 -2.924 0 -3.35a1.724 1.724 0 0 0 1.066 -2.573c-.94 -1.543 .826 -3.31 2.37 -2.37c1 .608 2.296 .07 2.572 -1.065z" />
                                        <path d="M9 12a3 3 0 1 0 6 0a3 3 0 0 0 -6 0" />
                                    </svg>
                                </ThemeIcon>
                                <Title c="white" order={5}>
                                    {tr('srs-lists.feature.settings', 'Settings')}
                                </Title>
                            </Group>
                            <Divider opacity={0.3} />
                            <TextInput
                                description={tr(
                                    'srs-lists.feature.update-interval-hint',
                                    'Supported by sing-box duration parser, for example: 10m, 1h, 12h, 1d, 7d.'
                                )}
                                label={tr('srs-lists.feature.update-interval', 'Update interval')}
                                onChange={(event) => setUpdateInterval(event.currentTarget.value)}
                                placeholder="1d"
                                value={updateInterval}
                            />
                            <Select
                                data={[
                                    { value: 'enabled', label: tr('srs-lists.feature.enabled', 'Enabled') },
                                    { value: 'disabled', label: tr('srs-lists.feature.disabled', 'Disabled') }
                                ]}
                                label={tr('srs-lists.feature.status', 'Status')}
                                onChange={(value) => setCreateEnabled(value !== 'disabled')}
                                value={createEnabled ? 'enabled' : 'disabled'}
                            />
                        </Stack>
                    </Card>

                    <Button
                        fullWidth
                        leftSection={<SRSListsIcon size={16} />}
                        loading={isCreating}
                        onClick={onSubmitCreate}
                    >
                        {tr('common.save', 'Save')}
                    </Button>
                </Stack>
            </Modal>

            <Modal
                onClose={() => {
                    editModalHandlers.close()
                    setEditingItem(null)
                }}
                opened={editModalOpened}
                size={660}
                title={
                    <Group gap="sm" wrap="nowrap">
                        <ThemeIcon size="lg" variant="gradient-blue">
                            <SRSListsIcon size={20} />
                        </ThemeIcon>
                        <Title c="white" order={4}>
                            {tr('srs-lists.feature.edit-list', 'Edit SRS list')}
                        </Title>
                    </Group>
                }
            >
                <Stack>
                    <TextInput
                        label={tr('srs-lists.feature.url', 'URL')}
                        onChange={(event) => setEditURL(event.currentTarget.value)}
                        placeholder="https://.../ruleset.srs"
                        value={editURL}
                    />
                    <TextInput
                        label={tr('srs-lists.feature.update-interval', 'Update interval')}
                        onChange={(event) => setEditUpdateInterval(event.currentTarget.value)}
                        placeholder="1d"
                        value={editUpdateInterval}
                    />
                    <Select
                        data={[
                            { value: 'enabled', label: tr('srs-lists.feature.enabled', 'Enabled') },
                            { value: 'disabled', label: tr('srs-lists.feature.disabled', 'Disabled') }
                        ]}
                        label={tr('srs-lists.feature.status', 'Status')}
                        onChange={(value) => setEditEnabled(value !== 'disabled')}
                        value={editEnabled ? 'enabled' : 'disabled'}
                    />

                    <Button
                        fullWidth
                        leftSection={<SRSListsIcon size={16} />}
                        loading={isUpdating}
                        onClick={onSubmitEdit}
                    >
                        {tr('common.save', 'Save')}
                    </Button>
                </Stack>
            </Modal>

            <Modal
                onClose={intervalModalHandlers.close}
                opened={intervalModalOpened}
                title={tr('srs-lists.feature.set-update-interval', 'Set update interval')}
            >
                <Stack>
                    <TextInput
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
        </Page>
    )
}
