export declare const SUBSCRIPTION_INFO_BLOCK_VARIANTS: {
    readonly COLLAPSED: "collapsed";
    readonly EXPANDED: "expanded";
    readonly CARDS: "cards";
    readonly HIDDEN: "hidden";
};
export declare const INSTALLATION_GUIDE_BLOCKS_VARIANTS: {
    readonly CARDS: "cards";
    readonly ACCORDION: "accordion";
    readonly MINIMAL: "minimal";
    readonly TIMELINE: "timeline";
};
export type TSubscriptionInfoBlockVariant = (typeof SUBSCRIPTION_INFO_BLOCK_VARIANTS)[keyof typeof SUBSCRIPTION_INFO_BLOCK_VARIANTS];
export declare const SUBSCRIPTION_INFO_BLOCK_VARIANTS_VALUES: ("collapsed" | "expanded" | "cards" | "hidden")[];
export type TInstallationGuideBlockVariant = (typeof INSTALLATION_GUIDE_BLOCKS_VARIANTS)[keyof typeof INSTALLATION_GUIDE_BLOCKS_VARIANTS];
export declare const INSTALLATION_GUIDE_BLOCKS_VARIANTS_VALUES: ("cards" | "accordion" | "minimal" | "timeline")[];
export declare const isInstallationGuideBlockVariant: (value: string) => value is TInstallationGuideBlockVariant;
export declare const isSubscriptionInfoBlockVariant: (value: string) => value is TSubscriptionInfoBlockVariant;
//# sourceMappingURL=ui-blocks-variants.constant.d.ts.map