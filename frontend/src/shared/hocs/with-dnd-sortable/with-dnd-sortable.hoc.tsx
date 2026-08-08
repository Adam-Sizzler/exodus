import { OptimisticSortingPlugin } from '@dnd-kit/dom/sortable'
import { useSortable } from '@dnd-kit/react/sortable'
import { ActionIcon } from '@mantine/core'
import { createContext, CSSProperties, forwardRef, useContext } from 'react'
import { RiDraggable } from 'react-icons/ri'

import classes from './with-dnd-sortable.module.css'

export const DndSortableIndexContext = createContext(0)

interface WithDndSortableProps {
    children: React.ReactNode
    dragHandlePosition?: 'bottom-left' | 'bottom-right' | 'top-left' | 'top-right'
    id: string
    isDragOverlay?: boolean
    showDragHandle?: boolean
}

export const WithDndSortable = forwardRef<HTMLDivElement, WithDndSortableProps>(
    (props, externalRef) => {
        const {
            id,
            isDragOverlay = false,
            children,
            showDragHandle = true,
            dragHandlePosition = 'top-right'
        } = props

        const index = useContext(DndSortableIndexContext)

        const sortable = useSortable({
            id,
            index,
            disabled: isDragOverlay,
            plugins: (defaults) => defaults.filter((plugin) => plugin !== OptimisticSortingPlugin)
        })

        const isDragging = !isDragOverlay && sortable.isDragging
        const { ref, handleRef } = sortable

        const style: CSSProperties = {
            opacity: isDragging ? 0 : 1,
            position: 'relative'
        }

        const dragHandleClasses = {
            'top-left': classes.dragHandleTopLeft,
            'top-right': classes.dragHandleTopRight,
            'bottom-left': classes.dragHandleBottomLeft,
            'bottom-right': classes.dragHandleBottomRight
        }

        return (
            <div ref={isDragOverlay ? externalRef : ref} style={style}>
                {showDragHandle && (
                    <ActionIcon
                        className={`${classes.dragHandle} ${dragHandleClasses[dragHandlePosition]}`}
                        ref={isDragOverlay ? undefined : handleRef}
                        size="lg"
                        variant="transparent"
                    >
                        <RiDraggable size={24} />
                    </ActionIcon>
                )}
                {children}
            </div>
        )
    }
)

WithDndSortable.displayName = 'WithDndSortable'
