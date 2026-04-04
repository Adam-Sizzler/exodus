import type { UniqueIdentifier } from '@dnd-kit/core/dist/types';
import type { DroppableContainer } from '@dnd-kit/core/dist/store/types';
declare type Identifier = UniqueIdentifier | null | undefined;
export declare class DroppableContainersMap extends Map<UniqueIdentifier, DroppableContainer> {
    get(id: Identifier): DroppableContainer | undefined;
    toArray(): DroppableContainer[];
    getEnabled(): DroppableContainer[];
    getNodeFor(id: Identifier): HTMLElement | undefined;
}
export {};
