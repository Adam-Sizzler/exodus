import type { PointerEvent } from 'react';
import type { SensorProps } from '@dnd-kit/core/dist/sensors/types';
import { AbstractPointerSensor, AbstractPointerSensorOptions } from '@dnd-kit/core/dist/sensors/pointer/AbstractPointerSensor';
export interface PointerSensorOptions extends AbstractPointerSensorOptions {
}
export declare type PointerSensorProps = SensorProps<PointerSensorOptions>;
export declare class PointerSensor extends AbstractPointerSensor {
    constructor(props: PointerSensorProps);
    static activators: {
        eventName: "onPointerDown";
        handler: ({ nativeEvent: event }: PointerEvent, { onActivation }: PointerSensorOptions) => boolean;
    }[];
}
