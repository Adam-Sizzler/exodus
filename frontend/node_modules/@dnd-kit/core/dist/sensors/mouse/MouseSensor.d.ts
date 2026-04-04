import type { MouseEvent } from 'react';
import type { SensorProps } from '@dnd-kit/core/dist/sensors/types';
import { AbstractPointerSensor, AbstractPointerSensorOptions } from '@dnd-kit/core/dist/sensors/pointer';
export interface MouseSensorOptions extends AbstractPointerSensorOptions {
}
export declare type MouseSensorProps = SensorProps<MouseSensorOptions>;
export declare class MouseSensor extends AbstractPointerSensor {
    constructor(props: MouseSensorProps);
    static activators: {
        eventName: "onMouseDown";
        handler: ({ nativeEvent: event }: MouseEvent, { onActivation }: MouseSensorOptions) => boolean;
    }[];
}
