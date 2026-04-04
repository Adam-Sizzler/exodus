import type { ClientRect } from '@dnd-kit/core/dist/types';
export declare function useRect(element: HTMLElement | null, measure?: (element: HTMLElement) => ClientRect, fallbackRect?: ClientRect | null): ClientRect | null;
