import type { ClientRect } from '@dnd-kit/core/dist/types';
import type { DroppableMeasuring } from '@dnd-kit/core/dist/hooks/utilities';
export declare type MeasuringFunction = (node: HTMLElement) => ClientRect;
interface Measuring {
    measure: MeasuringFunction;
}
export interface DraggableMeasuring extends Measuring {
}
export interface DragOverlayMeasuring extends Measuring {
}
export interface MeasuringConfiguration {
    draggable?: Partial<DraggableMeasuring>;
    droppable?: Partial<DroppableMeasuring>;
    dragOverlay?: Partial<DragOverlayMeasuring>;
}
export {};
