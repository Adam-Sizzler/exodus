import type { PublicContextDescriptor } from '@dnd-kit/core/dist/store';
import type { ClientRect } from '@dnd-kit/core/dist/types';
interface Arguments {
    measure(element: HTMLElement): ClientRect;
}
export declare function useDragOverlayMeasuring({ measure, }: Arguments): PublicContextDescriptor['dragOverlay'];
export {};
