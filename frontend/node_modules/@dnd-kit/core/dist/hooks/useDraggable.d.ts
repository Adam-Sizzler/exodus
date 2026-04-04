/// <reference types="react" />
import { Transform } from '@dnd-kit/utilities';
import { Data } from '@dnd-kit/core/dist/store';
import type { UniqueIdentifier } from '@dnd-kit/core/dist/types';
import { SyntheticListenerMap } from '@dnd-kit/core/dist/hooks/utilities';
export interface UseDraggableArguments {
    id: UniqueIdentifier;
    data?: Data;
    disabled?: boolean;
    attributes?: {
        role?: string;
        roleDescription?: string;
        tabIndex?: number;
    };
}
export interface DraggableAttributes {
    role: string;
    tabIndex: number;
    'aria-disabled': boolean;
    'aria-pressed': boolean | undefined;
    'aria-roledescription': string;
    'aria-describedby': string;
}
export declare type DraggableSyntheticListeners = SyntheticListenerMap | undefined;
export declare function useDraggable({ id, data, disabled, attributes, }: UseDraggableArguments): {
    active: import("@dnd-kit/core/dist/store").Active | null;
    activatorEvent: Event | null;
    activeNodeRect: import("@dnd-kit/core/dist/types").ClientRect | null;
    attributes: DraggableAttributes;
    isDragging: boolean;
    listeners: SyntheticListenerMap | undefined;
    node: import("react").MutableRefObject<HTMLElement | null>;
    over: import("@dnd-kit/core/dist/store").Over | null;
    setNodeRef: (element: HTMLElement | null) => void;
    setActivatorNodeRef: (element: HTMLElement | null) => void;
    transform: Transform | null;
};
