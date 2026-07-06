import {
    closestCenter,
    DndContext,
    DragEndEvent,
    DragOverlay,
    DragStartEvent,
    KeyboardSensor,
    MouseSensor,
    TouchSensor,
    useSensor,
    useSensors
} from '@dnd-kit/core'
import { restrictToVerticalAxis } from '@dnd-kit/modifiers'
import { arrayMove, SortableContext, verticalListSortingStrategy } from '@dnd-kit/sortable'
import { Box, Container, Stack } from '@mantine/core'
import { useListState } from '@mantine/hooks'
import { GetAllNodesCommand } from '@exodus/backend-contract'
import { useWindowVirtualizer } from '@tanstack/react-virtual'
import { memo, useCallback, useEffect, useLayoutEffect, useMemo, useRef, useState } from 'react'

import { queryClient } from '@shared/api'
import { nodesQueryKeys, useGetNodes, useReorderNodes } from '@shared/api/hooks'
import { useIsMobile } from '@shared/hooks'
import { NO_TAG, TagFilterBar } from '@shared/ui'
import { EmptyPageLayout } from '@shared/ui/layouts/empty-page'
import { sToMs } from '@shared/utils/time-utils'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'
import {
    useNodesActiveTag,
    useViewPreferencesStoreActions
} from '@entities/dashboard/view-preferences-store'

import { NodeCardWidget } from '../node-card'
import { NodesSpotlightSearchWidget } from '../nodes-spotlight-search'
import { IProps } from './interfaces'
import styles from './NodesTable.module.css'

export const NodesTableWidget = memo((props: IProps) => {
    const { nodes } = props

    const activeTag = useNodesActiveTag()
    const { setNodesActiveTag } = useViewPreferencesStoreActions()

    const visibleNodes = useMemo(() => {
        if (!nodes) return []
        if (activeTag === null) return nodes
        if (activeTag === NO_TAG) return nodes.filter((node) => (node.tags ?? []).length === 0)
        return nodes.filter((node) => (node.tags ?? []).includes(activeTag))
    }, [nodes, activeTag])

    const [state, handlers] = useListState(visibleNodes)

    const openModalWithData = useModalsStoreOpenWithData()
    const [isPollingEnabled, setIsPollingEnabled] = useState(true)
    const [draggedNode, setDraggedNode] = useState<
        GetAllNodesCommand.Response['response'][number] | null
    >(null)
    const [scrollMargin, setScrollMargin] = useState(0)
    const listRef = useRef<HTMLDivElement | null>(null)
    const prevStateRef = useRef(state)
    const activeTagRef = useRef(activeTag)
    activeTagRef.current = activeTag
    const isMobile = useIsMobile()

    useGetNodes({
        rQueryParams: {
            enabled: isPollingEnabled,
            refetchInterval: isPollingEnabled ? sToMs(5) : false
        }
    })

    const { mutate: reorderNodes } = useReorderNodes({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(nodesQueryKeys.getAllNodes.queryKey, data)
            },
            onError: () => {
                queryClient.invalidateQueries({ queryKey: nodesQueryKeys.getAllNodes.queryKey })
            }
        }
    })

    const virtualizer = useWindowVirtualizer({
        count: state.length,
        estimateSize: () => (isMobile ? 190 : 90),
        overscan: 7,
        scrollMargin,
        getItemKey: (index) => state[index].uuid
    })

    const dataIds = useMemo(() => state.map((node) => node.uuid), [state])

    const sensors = useSensors(
        useSensor(MouseSensor, {
            activationConstraint: {
                distance: 5
            }
        }),
        useSensor(TouchSensor, {
            activationConstraint: {
                delay: 250,
                tolerance: 5
            }
        }),
        useSensor(KeyboardSensor, {})
    )

    useEffect(() => {
        ;(async () => {
            if (!state || state.length === 0) {
                return
            }

            if (activeTagRef.current !== null) {
                prevStateRef.current = state
                return
            }

            const updatedNodes = state.map((node, index) => ({
                uuid: node.uuid,
                viewPosition: index
            }))

            const hasOrderChanged = prevStateRef.current?.some(
                (node, index) => state[index] && node.uuid !== state[index].uuid
            )

            if (hasOrderChanged) {
                reorderNodes({ variables: { nodes: updatedNodes } })
            }

            prevStateRef.current = state
        })()
    }, [state])

    useEffect(() => {
        handlers.setState(visibleNodes)
        prevStateRef.current = visibleNodes
    }, [visibleNodes])

    useLayoutEffect(() => {
        if (listRef.current) {
            setScrollMargin(listRef.current.offsetTop)
        }
    }, [])

    const handleDragStart = useCallback(
        (event: DragStartEvent) => {
            setIsPollingEnabled(false)
            const draggedItem = state.find((item) => item.uuid === event.active.id)
            setDraggedNode(draggedItem || null)
        },
        [state]
    )

    const handleDragEnd = useCallback(
        (event: DragEndEvent) => {
            const { active, over } = event

            if (!over || active.id === over.id) {
                setIsPollingEnabled(true)
                setDraggedNode(null)
                return
            }

            const oldIndex = dataIds.indexOf(String(active.id))
            const newIndex = dataIds.indexOf(String(over.id))

            if (oldIndex !== -1 && newIndex !== -1) {
                const newState = arrayMove(state, oldIndex, newIndex)
                handlers.setState(newState)
            }

            setIsPollingEnabled(true)
            setDraggedNode(null)
        },
        [dataIds, state, handlers]
    )

    const handleDragCancel = useCallback(() => {
        setIsPollingEnabled(true)
        setDraggedNode(null)
    }, [])

    const handleViewNode = (nodeUuid: string) => {
        openModalWithData(MODALS.EDIT_NODE_BY_UUID_MODAL, { nodeUuid })
    }

    if (!nodes) {
        return null
    }

    if (nodes.length === 0) {
        return <EmptyPageLayout />
    }

    return (
        <>
            <TagFilterBar activeTag={activeTag} items={nodes} onChange={setNodesActiveTag} />

            <DndContext
                collisionDetection={closestCenter}
                modifiers={[restrictToVerticalAxis]}
                onDragCancel={handleDragCancel}
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
                        <SortableContext items={dataIds} strategy={verticalListSortingStrategy}>
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
                                                    transform: `translateY(${
                                                        virtualItem.start -
                                                        virtualizer.options.scrollMargin
                                                    }px)`,
                                                    willChange: 'transform'
                                                }}
                                            >
                                                <div className={styles.nodeFadeIn}>
                                                    <NodeCardWidget
                                                        disableReordering={activeTag !== null}
                                                        handleViewNode={handleViewNode}
                                                        isMobile={isMobile}
                                                        node={item}
                                                    />
                                                </div>
                                            </Box>
                                        )
                                    })}
                                </Stack>
                            </Container>
                        </SortableContext>
                    </div>
                </div>
                <DragOverlay>
                    {draggedNode && (
                        <Container fluid pl={0} pr={0}>
                            <NodeCardWidget
                                handleViewNode={handleViewNode}
                                isDragOverlay
                                isMobile={isMobile}
                                node={draggedNode}
                            />
                        </Container>
                    )}
                </DragOverlay>
            </DndContext>
            {nodes && nodes.length > 0 && <NodesSpotlightSearchWidget nodes={nodes} />}
        </>
    )
})
