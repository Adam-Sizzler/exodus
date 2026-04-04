import { BoxProps, PolymorphicFactory, StylesApiProps } from '../../core';
export type CenterStylesNames = 'root';
export interface CenterProps extends BoxProps, StylesApiProps<CenterFactory> {
    /** Content to center */
    children?: React.ReactNode;
    /** If set, `inline-flex` is used instead of `flex` */
    inline?: boolean;
}
export type CenterFactory = PolymorphicFactory<{
    props: CenterProps;
    defaultRef: HTMLDivElement;
    defaultComponent: 'div';
    stylesNames: CenterStylesNames;
}>;
export declare const Center: (<C = "div">(props: import("../..").PolymorphicComponentProps<C, CenterProps>) => React.ReactElement) & Omit<import("react").FunctionComponent<(CenterProps & {
    component?: any;
} & Omit<Omit<any, "ref">, "component" | keyof CenterProps> & {
    ref?: any;
    renderRoot?: (props: any) => any;
}) | (CenterProps & {
    component: React.ElementType;
    renderRoot?: (props: Record<string, any>) => any;
})>, never> & import("../..").ThemeExtend<{
    props: CenterProps;
    defaultRef: HTMLDivElement;
    defaultComponent: "div";
    stylesNames: CenterStylesNames;
}> & import("../..").ComponentClasses<{
    props: CenterProps;
    defaultRef: HTMLDivElement;
    defaultComponent: "div";
    stylesNames: CenterStylesNames;
}> & import("../..").PolymorphicComponentWithProps<{
    props: CenterProps;
    defaultRef: HTMLDivElement;
    defaultComponent: "div";
    stylesNames: CenterStylesNames;
}> & Record<string, never>;
