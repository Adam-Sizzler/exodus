import { UniqueIdentifier, Draggable, Droppable, DragDropManager, DragDropEvents } from '@dnd-kit/abstract';

/**
 * Move an array item to a different position. Returns a new array with the item moved to the new position.
 */
declare function arrayMove<T extends any[]>(array: T, from: number, to: number): T;
/**
 * Move an array item to a different position. Returns a new array with the item moved to the new position.
 */
declare function arraySwap<T extends any[]>(array: T, from: number, to: number): T;
type Items = UniqueIdentifier[] | {
    id: UniqueIdentifier;
}[];
declare function move<T extends Items | Record<UniqueIdentifier, Items>, U extends Draggable, V extends Droppable, W extends DragDropManager<U, V>>(items: T, event: Parameters<DragDropEvents<U, V, W>['dragover'] | DragDropEvents<U, V, W>['dragend']>[0]): T;
declare function swap<T extends Items | Record<UniqueIdentifier, Items>, U extends Draggable, V extends Droppable, W extends DragDropManager<U, V>>(items: T, event: Parameters<DragDropEvents<U, V, W>['dragover'] | DragDropEvents<U, V, W>['dragend']>[0]): T;

export { arrayMove, arraySwap, move, swap };
