import type { Sensor, SensorDescriptor, SensorOptions } from '@dnd-kit/core/dist/sensors/types';
export declare function useSensor<T extends SensorOptions>(sensor: Sensor<T>, options?: T): SensorDescriptor<T>;
