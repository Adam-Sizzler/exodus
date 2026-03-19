import { z } from 'zod';
declare const LocalizedTextSchema: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
declare const SvgLibrarySchema: z.ZodRecord<z.ZodString, z.ZodString>;
declare const ButtonSchema: z.ZodObject<{
    link: z.ZodString;
    type: z.ZodNativeEnum<{
        readonly EXTERNAL: "external";
        readonly SUBSCRIPTION_LINK: "subscriptionLink";
        readonly COPY_BUTTON: "copyButton";
    }>;
    text: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    svgIconKey: z.ZodString;
}, "strip", z.ZodTypeAny, {
    type: "external" | "subscriptionLink" | "copyButton";
    svgIconKey: string;
    link: string;
    text: Record<string, string>;
}, {
    type: "external" | "subscriptionLink" | "copyButton";
    svgIconKey: string;
    link: string;
    text: Record<string, string>;
}>;
declare const BlockSchema: z.ZodObject<{
    svgIconKey: z.ZodString;
    svgIconColor: z.ZodEffects<z.ZodString, string, string>;
    title: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    description: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    buttons: z.ZodArray<z.ZodObject<{
        link: z.ZodString;
        type: z.ZodNativeEnum<{
            readonly EXTERNAL: "external";
            readonly SUBSCRIPTION_LINK: "subscriptionLink";
            readonly COPY_BUTTON: "copyButton";
        }>;
        text: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        svgIconKey: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        type: "external" | "subscriptionLink" | "copyButton";
        svgIconKey: string;
        link: string;
        text: Record<string, string>;
    }, {
        type: "external" | "subscriptionLink" | "copyButton";
        svgIconKey: string;
        link: string;
        text: Record<string, string>;
    }>, "many">;
}, "strip", z.ZodTypeAny, {
    svgIconKey: string;
    svgIconColor: string;
    title: Record<string, string>;
    description: Record<string, string>;
    buttons: {
        type: "external" | "subscriptionLink" | "copyButton";
        svgIconKey: string;
        link: string;
        text: Record<string, string>;
    }[];
}, {
    svgIconKey: string;
    svgIconColor: string;
    title: Record<string, string>;
    description: Record<string, string>;
    buttons: {
        type: "external" | "subscriptionLink" | "copyButton";
        svgIconKey: string;
        link: string;
        text: Record<string, string>;
    }[];
}>;
declare const PlatformAppSchema: z.ZodObject<{
    name: z.ZodString;
    svgIconKey: z.ZodOptional<z.ZodString>;
    featured: z.ZodBoolean;
    blocks: z.ZodArray<z.ZodObject<{
        svgIconKey: z.ZodString;
        svgIconColor: z.ZodEffects<z.ZodString, string, string>;
        title: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        description: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        buttons: z.ZodArray<z.ZodObject<{
            link: z.ZodString;
            type: z.ZodNativeEnum<{
                readonly EXTERNAL: "external";
                readonly SUBSCRIPTION_LINK: "subscriptionLink";
                readonly COPY_BUTTON: "copyButton";
            }>;
            text: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
            svgIconKey: z.ZodString;
        }, "strip", z.ZodTypeAny, {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }, {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        svgIconKey: string;
        svgIconColor: string;
        title: Record<string, string>;
        description: Record<string, string>;
        buttons: {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }[];
    }, {
        svgIconKey: string;
        svgIconColor: string;
        title: Record<string, string>;
        description: Record<string, string>;
        buttons: {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }[];
    }>, "many">;
}, "strip", z.ZodTypeAny, {
    name: string;
    featured: boolean;
    blocks: {
        svgIconKey: string;
        svgIconColor: string;
        title: Record<string, string>;
        description: Record<string, string>;
        buttons: {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }[];
    }[];
    svgIconKey?: string | undefined;
}, {
    name: string;
    featured: boolean;
    blocks: {
        svgIconKey: string;
        svgIconColor: string;
        title: Record<string, string>;
        description: Record<string, string>;
        buttons: {
            type: "external" | "subscriptionLink" | "copyButton";
            svgIconKey: string;
            link: string;
            text: Record<string, string>;
        }[];
    }[];
    svgIconKey?: string | undefined;
}>;
declare const PlatformSchema: z.ZodObject<{
    displayName: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    svgIconKey: z.ZodString;
    apps: z.ZodArray<z.ZodObject<{
        name: z.ZodString;
        svgIconKey: z.ZodOptional<z.ZodString>;
        featured: z.ZodBoolean;
        blocks: z.ZodArray<z.ZodObject<{
            svgIconKey: z.ZodString;
            svgIconColor: z.ZodEffects<z.ZodString, string, string>;
            title: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
            description: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
            buttons: z.ZodArray<z.ZodObject<{
                link: z.ZodString;
                type: z.ZodNativeEnum<{
                    readonly EXTERNAL: "external";
                    readonly SUBSCRIPTION_LINK: "subscriptionLink";
                    readonly COPY_BUTTON: "copyButton";
                }>;
                text: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
                svgIconKey: z.ZodString;
            }, "strip", z.ZodTypeAny, {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }, {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }, {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        name: string;
        featured: boolean;
        blocks: {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }[];
        svgIconKey?: string | undefined;
    }, {
        name: string;
        featured: boolean;
        blocks: {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }[];
        svgIconKey?: string | undefined;
    }>, "many">;
}, "strip", z.ZodTypeAny, {
    svgIconKey: string;
    displayName: Record<string, string>;
    apps: {
        name: string;
        featured: boolean;
        blocks: {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }[];
        svgIconKey?: string | undefined;
    }[];
}, {
    svgIconKey: string;
    displayName: Record<string, string>;
    apps: {
        name: string;
        featured: boolean;
        blocks: {
            svgIconKey: string;
            svgIconColor: string;
            title: Record<string, string>;
            description: Record<string, string>;
            buttons: {
                type: "external" | "subscriptionLink" | "copyButton";
                svgIconKey: string;
                link: string;
                text: Record<string, string>;
            }[];
        }[];
        svgIconKey?: string | undefined;
    }[];
}>;
declare const BrandingSettingsSchema: z.ZodObject<{
    title: z.ZodString;
    logoUrl: z.ZodString;
    supportUrl: z.ZodString;
}, "strip", z.ZodTypeAny, {
    title: string;
    logoUrl: string;
    supportUrl: string;
}, {
    title: string;
    logoUrl: string;
    supportUrl: string;
}>;
declare const UiConfigSchema: z.ZodObject<{
    subscriptionInfoBlockType: z.ZodNativeEnum<{
        readonly COLLAPSED: "collapsed";
        readonly EXPANDED: "expanded";
        readonly CARDS: "cards";
        readonly HIDDEN: "hidden";
    }>;
    installationGuidesBlockType: z.ZodNativeEnum<{
        readonly CARDS: "cards";
        readonly ACCORDION: "accordion";
        readonly MINIMAL: "minimal";
        readonly TIMELINE: "timeline";
    }>;
}, "strip", z.ZodTypeAny, {
    subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
    installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
}, {
    subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
    installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
}>;
declare const SubscriptionPageTranslateKeysSchema: z.ZodObject<{
    installationGuideHeader: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    connectionKeysHeader: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    linkCopied: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    linkCopiedToClipboard: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    getLink: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    scanQrCode: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    scanQrCodeDescription: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    copyLink: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    name: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    status: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    active: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    inactive: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    expires: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    bandwidth: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    scanToImport: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    expiresIn: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    expired: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    unknown: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    indefinitely: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
}, "strip", z.ZodTypeAny, {
    unknown: Record<string, string>;
    status: Record<string, string>;
    name: Record<string, string>;
    installationGuideHeader: Record<string, string>;
    connectionKeysHeader: Record<string, string>;
    linkCopied: Record<string, string>;
    linkCopiedToClipboard: Record<string, string>;
    getLink: Record<string, string>;
    scanQrCode: Record<string, string>;
    scanQrCodeDescription: Record<string, string>;
    copyLink: Record<string, string>;
    active: Record<string, string>;
    inactive: Record<string, string>;
    expires: Record<string, string>;
    bandwidth: Record<string, string>;
    scanToImport: Record<string, string>;
    expiresIn: Record<string, string>;
    expired: Record<string, string>;
    indefinitely: Record<string, string>;
}, {
    unknown: Record<string, string>;
    status: Record<string, string>;
    name: Record<string, string>;
    installationGuideHeader: Record<string, string>;
    connectionKeysHeader: Record<string, string>;
    linkCopied: Record<string, string>;
    linkCopiedToClipboard: Record<string, string>;
    getLink: Record<string, string>;
    scanQrCode: Record<string, string>;
    scanQrCodeDescription: Record<string, string>;
    copyLink: Record<string, string>;
    active: Record<string, string>;
    inactive: Record<string, string>;
    expires: Record<string, string>;
    bandwidth: Record<string, string>;
    scanToImport: Record<string, string>;
    expiresIn: Record<string, string>;
    expired: Record<string, string>;
    indefinitely: Record<string, string>;
}>;
export declare const SubscriptionPageRawConfigSchema: z.ZodEffects<z.ZodObject<{
    version: z.ZodNativeEnum<{
        readonly 1: "1";
    }>;
    locales: z.ZodArray<z.ZodEnum<["en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th", ...("en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th")[]]>, "many">;
    brandingSettings: z.ZodObject<{
        title: z.ZodString;
        logoUrl: z.ZodString;
        supportUrl: z.ZodString;
    }, "strip", z.ZodTypeAny, {
        title: string;
        logoUrl: string;
        supportUrl: string;
    }, {
        title: string;
        logoUrl: string;
        supportUrl: string;
    }>;
    uiConfig: z.ZodObject<{
        subscriptionInfoBlockType: z.ZodNativeEnum<{
            readonly COLLAPSED: "collapsed";
            readonly EXPANDED: "expanded";
            readonly CARDS: "cards";
            readonly HIDDEN: "hidden";
        }>;
        installationGuidesBlockType: z.ZodNativeEnum<{
            readonly CARDS: "cards";
            readonly ACCORDION: "accordion";
            readonly MINIMAL: "minimal";
            readonly TIMELINE: "timeline";
        }>;
    }, "strip", z.ZodTypeAny, {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    }, {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    }>;
    baseSettings: z.ZodDefault<z.ZodObject<{
        metaTitle: z.ZodDefault<z.ZodString>;
        metaDescription: z.ZodDefault<z.ZodString>;
        showConnectionKeys: z.ZodDefault<z.ZodBoolean>;
        hideGetLinkButton: z.ZodDefault<z.ZodBoolean>;
    }, "strip", z.ZodTypeAny, {
        metaTitle: string;
        metaDescription: string;
        showConnectionKeys: boolean;
        hideGetLinkButton: boolean;
    }, {
        metaTitle?: string | undefined;
        metaDescription?: string | undefined;
        showConnectionKeys?: boolean | undefined;
        hideGetLinkButton?: boolean | undefined;
    }>>;
    baseTranslations: z.ZodObject<{
        installationGuideHeader: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        connectionKeysHeader: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        linkCopied: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        linkCopiedToClipboard: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        getLink: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        scanQrCode: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        scanQrCodeDescription: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        copyLink: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        name: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        status: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        active: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        inactive: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        expires: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        bandwidth: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        scanToImport: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        expiresIn: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        expired: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        unknown: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        indefinitely: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
    }, "strip", z.ZodTypeAny, {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    }, {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    }>;
    svgLibrary: z.ZodRecord<z.ZodString, z.ZodString>;
    platforms: z.ZodRecord<z.ZodNativeEnum<{
        readonly IOS: "ios";
        readonly ANDROID: "android";
        readonly LINUX: "linux";
        readonly MACOS: "macos";
        readonly WINDOWS: "windows";
        readonly ANDROID_TV: "androidTV";
        readonly APPLE_TV: "appleTV";
    }>, z.ZodObject<{
        displayName: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
        svgIconKey: z.ZodString;
        apps: z.ZodArray<z.ZodObject<{
            name: z.ZodString;
            svgIconKey: z.ZodOptional<z.ZodString>;
            featured: z.ZodBoolean;
            blocks: z.ZodArray<z.ZodObject<{
                svgIconKey: z.ZodString;
                svgIconColor: z.ZodEffects<z.ZodString, string, string>;
                title: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
                description: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
                buttons: z.ZodArray<z.ZodObject<{
                    link: z.ZodString;
                    type: z.ZodNativeEnum<{
                        readonly EXTERNAL: "external";
                        readonly SUBSCRIPTION_LINK: "subscriptionLink";
                        readonly COPY_BUTTON: "copyButton";
                    }>;
                    text: z.ZodEffects<z.ZodRecord<z.ZodString, z.ZodString>, Record<string, string>, Record<string, string>>;
                    svgIconKey: z.ZodString;
                }, "strip", z.ZodTypeAny, {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }, {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }>, "many">;
            }, "strip", z.ZodTypeAny, {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }, {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }, {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }>, "many">;
    }, "strip", z.ZodTypeAny, {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }, {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }>>;
}, "strip", z.ZodTypeAny, {
    platforms: Partial<Record<"ios" | "android" | "linux" | "macos" | "windows" | "androidTV" | "appleTV", {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }>>;
    baseTranslations: {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    };
    version: "1";
    locales: ("en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th")[];
    brandingSettings: {
        title: string;
        logoUrl: string;
        supportUrl: string;
    };
    uiConfig: {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    };
    baseSettings: {
        metaTitle: string;
        metaDescription: string;
        showConnectionKeys: boolean;
        hideGetLinkButton: boolean;
    };
    svgLibrary: Record<string, string>;
}, {
    platforms: Partial<Record<"ios" | "android" | "linux" | "macos" | "windows" | "androidTV" | "appleTV", {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }>>;
    baseTranslations: {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    };
    version: "1";
    locales: ("en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th")[];
    brandingSettings: {
        title: string;
        logoUrl: string;
        supportUrl: string;
    };
    uiConfig: {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    };
    svgLibrary: Record<string, string>;
    baseSettings?: {
        metaTitle?: string | undefined;
        metaDescription?: string | undefined;
        showConnectionKeys?: boolean | undefined;
        hideGetLinkButton?: boolean | undefined;
    } | undefined;
}>, {
    platforms: Partial<Record<"ios" | "android" | "linux" | "macos" | "windows" | "androidTV" | "appleTV", {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }>>;
    baseTranslations: {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    };
    version: "1";
    locales: ("en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th")[];
    brandingSettings: {
        title: string;
        logoUrl: string;
        supportUrl: string;
    };
    uiConfig: {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    };
    baseSettings: {
        metaTitle: string;
        metaDescription: string;
        showConnectionKeys: boolean;
        hideGetLinkButton: boolean;
    };
    svgLibrary: Record<string, string>;
}, {
    platforms: Partial<Record<"ios" | "android" | "linux" | "macos" | "windows" | "androidTV" | "appleTV", {
        svgIconKey: string;
        displayName: Record<string, string>;
        apps: {
            name: string;
            featured: boolean;
            blocks: {
                svgIconKey: string;
                svgIconColor: string;
                title: Record<string, string>;
                description: Record<string, string>;
                buttons: {
                    type: "external" | "subscriptionLink" | "copyButton";
                    svgIconKey: string;
                    link: string;
                    text: Record<string, string>;
                }[];
            }[];
            svgIconKey?: string | undefined;
        }[];
    }>>;
    baseTranslations: {
        unknown: Record<string, string>;
        status: Record<string, string>;
        name: Record<string, string>;
        installationGuideHeader: Record<string, string>;
        connectionKeysHeader: Record<string, string>;
        linkCopied: Record<string, string>;
        linkCopiedToClipboard: Record<string, string>;
        getLink: Record<string, string>;
        scanQrCode: Record<string, string>;
        scanQrCodeDescription: Record<string, string>;
        copyLink: Record<string, string>;
        active: Record<string, string>;
        inactive: Record<string, string>;
        expires: Record<string, string>;
        bandwidth: Record<string, string>;
        scanToImport: Record<string, string>;
        expiresIn: Record<string, string>;
        expired: Record<string, string>;
        indefinitely: Record<string, string>;
    };
    version: "1";
    locales: ("en" | "ru" | "zh" | "fr" | "fa" | "uz" | "de" | "hi" | "tr" | "az" | "es" | "vi" | "ja" | "be" | "uk" | "pt" | "pl" | "id" | "tk" | "th")[];
    brandingSettings: {
        title: string;
        logoUrl: string;
        supportUrl: string;
    };
    uiConfig: {
        subscriptionInfoBlockType: "collapsed" | "expanded" | "cards" | "hidden";
        installationGuidesBlockType: "cards" | "accordion" | "minimal" | "timeline";
    };
    svgLibrary: Record<string, string>;
    baseSettings?: {
        metaTitle?: string | undefined;
        metaDescription?: string | undefined;
        showConnectionKeys?: boolean | undefined;
        hideGetLinkButton?: boolean | undefined;
    } | undefined;
}>;
export type TSubscriptionPageSvgLibrary = z.infer<typeof SvgLibrarySchema>;
export type TSubscriptionPageRawConfig = z.infer<typeof SubscriptionPageRawConfigSchema>;
export type TSubscriptionPageBrandingSettings = z.infer<typeof BrandingSettingsSchema>;
export type TSubscriptionPagePlatformSchema = z.infer<typeof PlatformSchema>;
export type TSubscriptionPagePlatformKey = keyof TSubscriptionPageRawConfig['platforms'];
export type TSubscriptionPageAppConfig = z.infer<typeof PlatformAppSchema>;
export type TSubscriptionPageBlockConfig = z.infer<typeof BlockSchema>;
export type TSubscriptionPageButtonConfig = z.infer<typeof ButtonSchema>;
export type TSubscriptionPageLocalizedText = z.infer<typeof LocalizedTextSchema>;
export type TSubscriptionPageUiConfig = z.infer<typeof UiConfigSchema>;
export type TSubscriptionPageTranslateKeys = z.infer<typeof SubscriptionPageTranslateKeysSchema>;
export type TSubscriptionPageBaseTranslationKeys = keyof TSubscriptionPageTranslateKeys;
export {};
//# sourceMappingURL=subscription-page-config.schema.d.ts.map