export declare const BUTTON_TYPES: {
    readonly EXTERNAL: "external";
    readonly SUBSCRIPTION_LINK: "subscriptionLink";
    readonly COPY_BUTTON: "copyButton";
};
export type TButtonType = (typeof BUTTON_TYPES)[keyof typeof BUTTON_TYPES];
export declare const BUTTON_TYPES_VALUES: ("external" | "subscriptionLink" | "copyButton")[];
export declare const isButtonType: (value: string) => value is TButtonType;
//# sourceMappingURL=button-types.constant.d.ts.map