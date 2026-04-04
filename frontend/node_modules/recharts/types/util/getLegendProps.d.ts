import { ReactNode, ReactElement } from 'react';
import { Props as LegendProps } from 'recharts/types/component/Legend';
import { FormattedGraphicalItem } from 'recharts/types/util/ChartUtils';
export declare const getLegendProps: ({ children, formattedGraphicalItems, legendWidth, legendContent, }: {
    children: ReactNode[];
    formattedGraphicalItems?: Array<FormattedGraphicalItem>;
    legendWidth: number;
    legendContent?: 'children';
}) => null | (LegendProps & {
    item: ReactElement;
});
