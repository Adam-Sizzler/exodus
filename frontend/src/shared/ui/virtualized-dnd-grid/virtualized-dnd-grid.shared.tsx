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
import { restrictToWindowEdges } from '@dnd-kit/modifiers'
import { arrayMove, rectSortingStrategy, SortableContext } from '@dnd-kit/sortable'
import { Box } from '@mantine/core'
import { ReactNode, useEffect, useState } from 'react'
import { VirtuosoGrid } from 'react-virtuoso'

import classes from './virtualized-dnd-grid.module.css'
import { VirtualizedGridComponents } from './virtualized-grid-components'

interface VirtualizedDndGridProps<T extends { uuid: string }> {
    enableDnd?: boolean
    items: T[]
    onReorder?: (items: T[]) => void
    renderDragOverlay?: (item: T) => ReactNode
    renderItem: (item: T, index: number) => ReactNode
    style?: React.CSSProperties
    useWindowScroll?: boolean
}

export function VirtualizedDndGrid<T extends { uuid: string }>(props: VirtualizedDndGridProps<T>) {
    const {
        items: initialItems,
        renderItem,
        renderDragOverlay,
        onReorder,
        useWindowScroll = true,
        style,
        enableDnd = true
    } = props

    const [items, setItems] = useState(initialItems)
    const [activeId, setActiveId] = useState<null | UniqueIdentifier>(null)

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
        setItems(initialItems)
    }, [initialItems])

    const handleDragStart = (event: DragStartEvent) => {
        setActiveId(event.active.id)
    }

    const handleDragEnd = (event: DragEndEvent) => {
        const { active, over } = event

        if (over && active.id !== over.id) {
            setItems((items) => {
                const oldIndex = items.findIndex((item) => item.uuid === active.id)
                const newIndex = items.findIndex((item) => item.uuid === over.id)

                const newItems = arrayMove(items, oldIndex, newIndex)
                onReorder?.(newItems)
                return newItems
            })
        }

        setActiveId(null)
    }

    const handleDragCancel = () => {
        setActiveId(null)
    }

    const draggedItem = activeId ? items.find((item) => item.uuid === activeId) : null

    const itemContent = (index: number) => {
        const item = items[index]
        if (!item) return null

        return (
            <div className={classes.itemWrapper} key={item.uuid}>
                {renderItem(item, index)}
            </div>
        )
    }

    if (!enableDnd) {
        return (
            <Box style={style}>
                <VirtuosoGrid
                    components={VirtualizedGridComponents}
                    itemContent={itemContent}
                    totalCount={items.length}
                    useWindowScroll={useWindowScroll}
                />
            </Box>
        )
    }

    return (
        <DndContext
            collisionDetection={closestCenter}
            modifiers={[restrictToWindowEdges]}
            onDragCancel={handleDragCancel}
            onDragEnd={handleDragEnd}
            onDragStart={handleDragStart}
            sensors={sensors}
        >
            <SortableContext items={items.map((item) => item.uuid)} strategy={rectSortingStrategy}>
                <Box style={style}>
                    <VirtuosoGrid
                        components={VirtualizedGridComponents}
                        itemContent={itemContent}
                        overscan={{
                            main: 10,
                            reverse: 10
                        }}
                        totalCount={items.length}
                        useWindowScroll={useWindowScroll}
                    />
                </Box>
            </SortableContext>

            <DragOverlay>
                {draggedItem && renderDragOverlay ? (
                    <Box style={{ width: '100%' }}>{renderDragOverlay(draggedItem)}</Box>
                ) : null}
            </DragOverlay>
        </DndContext>
    )
}
