import type { Coordinates, ClientRect } from '@dnd-kit/core/dist/types';
export declare function createRectAdjustmentFn(modifier: number): (rect: ClientRect, ...adjustments: Coordinates[]) => ClientRect;
export declare const getAdjustedRect: (rect: ClientRect, ...adjustments: Coordinates[]) => ClientRect;
