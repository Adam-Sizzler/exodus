import type { DndMonitorEvent } from '@dnd-kit/core/dist/components/DndMonitor/types';
export declare function useDndMonitorProvider(): readonly [({ type, event }: DndMonitorEvent) => void, (listener: any) => () => boolean];
