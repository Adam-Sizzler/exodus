import type { Active, Data, DroppableContainer, DraggableNode, Over } from '@dnd-kit/core';
import type { SortableData } from '@dnd-kit/sortable/dist/types/data';
export declare function hasSortableData<T extends Active | Over | DraggableNode | DroppableContainer>(entry: T | null | undefined): entry is T & {
    data: {
        current: Data<SortableData>;
    };
};
