import type { PointerActivationConstraint } from '@dnd-kit/core/dist/sensors';
import type { Active, Over } from '@dnd-kit/core/dist/store';
import type { Collision } from '@dnd-kit/core/dist/utilities/algorithms';
import type { Coordinates, Translate } from '@dnd-kit/core/dist/types/coordinates';
import type { UniqueIdentifier } from '@dnd-kit/core/dist/types';
interface DragEvent {
    activatorEvent: Event;
    active: Active;
    collisions: Collision[] | null;
    delta: Translate;
    over: Over | null;
}
/**
 * Fired if a pending drag was aborted before it started.
 * Only meaningful in the context of activation constraints.
 **/
export interface DragAbortEvent {
    id: UniqueIdentifier;
}
/**
 * Fired when a drag is about to start pending activation constraints.
 * @note For pointer events, it will be fired repeatedly with updated
 * coordinates when pointer is moved until the drag starts.
 */
export interface DragPendingEvent {
    id: UniqueIdentifier;
    constraint: PointerActivationConstraint;
    initialCoordinates: Coordinates;
    offset?: Coordinates | undefined;
}
export interface DragStartEvent extends Pick<DragEvent, 'active' | 'activatorEvent'> {
}
export interface DragMoveEvent extends DragEvent {
}
export interface DragOverEvent extends DragMoveEvent {
}
export interface DragEndEvent extends DragEvent {
}
export interface DragCancelEvent extends DragEndEvent {
}
export {};
