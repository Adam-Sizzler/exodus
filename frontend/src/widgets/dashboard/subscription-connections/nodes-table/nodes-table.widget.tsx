import { RestrictToVerticalAxis } from '@dnd-kit/abstract/modifiers'
import { move } from '@dnd-kit/helpers'
import {
    DragDropProvider,
    DragEndEvent,
    DragOverEvent,
    DragOverlay,
    DragStartEvent
} from '@dnd-kit/react'
import { memo, useCallback, useEffect, useLayoutEffect, useRef, useState } from 'react'
import { useWindowVirtualizer } from '@tanstack/react-virtual'
import { useListState, useMediaQuery } from '@mantine/hooks'
import { Box, Container, em, Stack } from '@mantine/core'

import { MODALS, useModalsStoreOpenWithData } from '@entities/dashboard/modal-store'
import {
    SubscriptionConnectionResponse,
    subscriptionConnectionsQueryKeys,
    useGetSubNodes,
    useReorderSubscriptionConnections
} from '@shared/api/hooks'
import { TbServer } from 'react-icons/tb'
import { EmptyPageLayout } from '@shared/ui/layouts/empty-page'
import { sToMs } from '@shared/utils/time-utils'
import { queryClient } from '@shared/api'

import { NodesSpotlightSearchWidget } from '../nodes-spotlight-search'
import { NodeCardWidget } from '../node-card'
import styles from './NodesTable.module.css'
import { IProps } from './interfaces'

export const NodesTableWidget = memo((props: IProps) => {
    const { nodes } = props
    const [state, handlers] = useListState(nodes || [])

    const openModalWithData = useModalsStoreOpenWithData()
    const [isPollingEnabled, setIsPollingEnabled] = useState(true)
    const [draggedNode, setDraggedNode] = useState<SubscriptionConnectionResponse | null>(null)
    const listRef = useRef<HTMLDivElement | null>(null)
    const parentOffsetRef = useRef(0)
    const prevStateRef = useRef(state)
    const isDraggingRef = useRef(false)
    const dragSnapshotRef = useRef<typeof state | null>(null)
    const isMobile = useMediaQuery(`(max-width: ${em(768)})`)

    useGetSubNodes({
        rQueryParams: {
            enabled: isPollingEnabled,
            refetchInterval: isPollingEnabled ? sToMs(5) : false
        }
    })

    const { mutate: reorderNodes } = useReorderSubscriptionConnections({
        mutationFns: {
            onSuccess: (data) => {
                queryClient.setQueryData(subscriptionConnectionsQueryKeys.getAllNodes.queryKey, data)
            },
            onError: () => {
                queryClient.invalidateQueries({ queryKey: subscriptionConnectionsQueryKeys.getAllNodes.queryKey })
            }
        }
    })

    const virtualizer = useWindowVirtualizer({
        count: state.length,
        estimateSize: () => (isMobile ? 169 : 64),
        overscan: 7,
        scrollMargin: parentOffsetRef.current,
        getItemKey: (index) => state[index]?.uuid ?? index
    })

    useEffect(() => {
        ;(async () => {
            if (!state || state.length === 0) {
                return
            }

            const updatedNodes = state.map((node, index) => ({
                uuid: node.uuid,
                viewPosition: index
            }))

            const hasOrderChanged = prevStateRef.current?.some(
                (node, index) => node.uuid !== state[index]?.uuid
            )

            if (hasOrderChanged) {
                reorderNodes({ variables: { nodes: updatedNodes } })
            }

            prevStateRef.current = state
        })()
    }, [state])

    useEffect(() => {
        handlers.setState(nodes || [])
        prevStateRef.current = nodes || []
    }, [nodes])

    useLayoutEffect(() => {
        parentOffsetRef.current = listRef.current?.offsetTop ?? 0
    }, [])

    const handleDragStart = useCallback(
        (event: DragStartEvent) => {
            setIsPollingEnabled(false)
            isDraggingRef.current = true
            dragSnapshotRef.current = state
            const draggedItem = state.find((item) => item.uuid === event.operation.source?.id)
            setDraggedNode(draggedItem || null)
        },
        [state]
    )

    const handleDragOver = useCallback(
        (event: DragOverEvent) => {
            handlers.setState((prev) => {
                const ids = prev.map((node) => node.uuid)
                const newIds = move(ids, event)
                if (newIds === ids) return prev

                const nodesByUuid = new Map(prev.map((node) => [node.uuid, node]))
                return newIds.map((uuid) => nodesByUuid.get(uuid)!)
            })
        },
        [handlers]
    )

    const handleDragEnd = useCallback(
        (event: DragEndEvent) => {
            isDraggingRef.current = false
            setIsPollingEnabled(true)
            setDraggedNode(null)

            const snapshot = dragSnapshotRef.current
            dragSnapshotRef.current = null

            if (event.canceled) {
                if (snapshot) {
                    prevStateRef.current = snapshot
                    handlers.setState(snapshot)
                }
                return
            }

            handlers.setState((prev) => [...prev])
        },
        [handlers]
    )

    const handleViewNode = (nodeUuid: string) => {
        openModalWithData(MODALS.EDIT_NODE_BY_UUID_MODAL, { nodeUuid })
    }

    if (!nodes) {
        return null
    }

    if (nodes.length === 0) {
        return <EmptyPageLayout icon={<TbServer size={32} />} />
    }

    return (
        <>
            <DragDropProvider
                modifiers={[RestrictToVerticalAxis]}
                onDragEnd={handleDragEnd}
                onDragOver={handleDragOver}
                onDragStart={handleDragStart}
            >
                <div ref={listRef}>
                    <div
                        style={{
                            height: `${virtualizer.getTotalSize()}px`,
                            width: '100%',
                            position: 'relative'
                        }}
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
                                                transform: `translateY(${
                                                    virtualItem.start -
                                                    virtualizer.options.scrollMargin
                                                }px)`
                                            }}
                                        >
                                            <div className={styles.nodeFadeIn}>
                                                <NodeCardWidget
                                                    handleViewNode={handleViewNode}
                                                    index={virtualItem.index}
                                                    isMobile={isMobile}
                                                    node={item}
                                                />
                                            </div>
                                        </Box>
                                    )
                                })}
                            </Stack>
                        </Container>
                    </div>
                </div>
                <DragOverlay>
                    {draggedNode && (
                        <Container fluid pl={0} pr={0}>
                            <NodeCardWidget
                                handleViewNode={handleViewNode}
                                index={0}
                                isDragOverlay
                                isMobile={isMobile}
                                node={draggedNode}
                            />
                        </Container>
                    )}
                </DragOverlay>
            </DragDropProvider>
            {nodes && nodes.length > 0 && <NodesSpotlightSearchWidget nodes={nodes} />}
        </>
    )
})
