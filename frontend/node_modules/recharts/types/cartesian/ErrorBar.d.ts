/**
 * @fileOverview Render a group of error bar
 */
import React, { SVGProps } from 'react';
import { Props as XAxisProps } from 'recharts/types/cartesian/XAxis';
import { Props as YAxisProps } from 'recharts/types/cartesian/YAxis';
import { D3Scale, DataKey } from 'recharts/types/util/types';
import { BarRectangleItem } from 'recharts/types/cartesian/Bar';
import { LinePointItem } from 'recharts/types/cartesian/Line';
import { ScatterPointItem } from 'recharts/types/cartesian/Scatter';
export interface ErrorBarDataItem {
    x: number;
    y: number;
    value: number;
    errorVal?: number[] | number;
}
export type ErrorBarDataPointFormatter = (entry: BarRectangleItem | LinePointItem | ScatterPointItem, dataKey: DataKey<any>) => ErrorBarDataItem;
interface InternalErrorBarProps {
    xAxis?: Omit<XAxisProps, 'scale'> & {
        scale: D3Scale<string | number>;
    };
    yAxis?: Omit<YAxisProps, 'scale'> & {
        scale: D3Scale<string | number>;
    };
    data?: any[];
    layout?: 'horizontal' | 'vertical';
    dataPointFormatter?: ErrorBarDataPointFormatter;
    /** The offset between central and the given coordinate, often set by <Bar/> */
    offset?: number;
}
interface ErrorBarProps extends InternalErrorBarProps {
    dataKey: DataKey<any>;
    /** the width of the error bar ends */
    width?: number;
    /**
     * Only used for ScatterChart with error bars in two directions.
     * Only accepts a value of "x" or "y" and makes the error bars lie in that direction.
     */
    direction?: 'x' | 'y';
}
export type Props = SVGProps<SVGLineElement> & ErrorBarProps;
export declare class ErrorBar extends React.Component<Props> {
    static defaultProps: {
        stroke: string;
        strokeWidth: number;
        width: number;
        offset: number;
        layout: string;
    };
    static displayName: string;
    render(): React.JSX.Element;
}
export {};
