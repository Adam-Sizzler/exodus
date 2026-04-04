import type { DraggableNode, DraggableNodes } from '@dnd-kit/core/dist/store';
import type { UniqueIdentifier } from '@dnd-kit/core/dist/types';
export declare function useCachedNode(draggableNodes: DraggableNodes, id: UniqueIdentifier | null): DraggableNode['node']['current'];
