import type { SensorDescriptor, SensorOptions } from '@dnd-kit/core/dist/sensors/types';
export declare function useSensors(...sensors: (SensorDescriptor<any> | undefined | null)[]): SensorDescriptor<SensorOptions>[];
