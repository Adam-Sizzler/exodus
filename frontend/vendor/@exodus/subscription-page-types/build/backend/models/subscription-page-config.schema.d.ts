import { z } from 'zod';
declare const LocalizedTextSchema: z.ZodRecord<z.ZodString, z.ZodString>;
declare const SvgLibrarySchema: z.ZodRecord<z.ZodString, z.ZodString>;
declare const ButtonSchema: z.ZodObject<{
    link: z.ZodString;
    type: z.ZodEnum<{
        readonly EXTERNAL: "external";
        readonly SUBSCRIPTION_LINK: "subscriptionLink";
        readonly COPY_BUTTON: "copyButton";
    }>;
    text: z.ZodRecord<z.ZodString, z.ZodString>;
    svgIconKey: z.ZodString;
}, z.core.$strip>;
declare const BlockSchema: z.ZodObject<{
    svgIconKey: z.ZodString;
    svgIconColor: z.ZodString;
    title: z.ZodRecord<z.ZodString, z.ZodString>;
    description: z.ZodRecord<z.ZodString, z.ZodString>;
    buttons: z.ZodArray<z.ZodObject<{
        link: z.ZodString;
        type: z.ZodEnum<{
            readonly EXTERNAL: "external";
            readonly SUBSCRIPTION_LINK: "subscriptionLink";
            readonly COPY_BUTTON: "copyButton";
        }>;
        text: z.ZodRecord<z.ZodString, z.ZodString>;
        svgIconKey: z.ZodString;
    }, z.core.$strip>>;
}, z.core.$strip>;
declare const PlatformAppSchema: z.ZodObject<{
    name: z.ZodString;
    svgIconKey: z.ZodOptional<z.ZodString>;
    featured: z.ZodBoolean;
    blocks: z.ZodArray<z.ZodObject<{
        svgIconKey: z.ZodString;
        svgIconColor: z.ZodString;
        title: z.ZodRecord<z.ZodString, z.ZodString>;
        description: z.ZodRecord<z.ZodString, z.ZodString>;
        buttons: z.ZodArray<z.ZodObject<{
            link: z.ZodString;
            type: z.ZodEnum<{
                readonly EXTERNAL: "external";
                readonly SUBSCRIPTION_LINK: "subscriptionLink";
                readonly COPY_BUTTON: "copyButton";
            }>;
            text: z.ZodRecord<z.ZodString, z.ZodString>;
            svgIconKey: z.ZodString;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
declare const PlatformSchema: z.ZodObject<{
    displayName: z.ZodRecord<z.ZodString, z.ZodString>;
    svgIconKey: z.ZodString;
    apps: z.ZodArray<z.ZodObject<{
        name: z.ZodString;
        svgIconKey: z.ZodOptional<z.ZodString>;
        featured: z.ZodBoolean;
        blocks: z.ZodArray<z.ZodObject<{
            svgIconKey: z.ZodString;
            svgIconColor: z.ZodString;
            title: z.ZodRecord<z.ZodString, z.ZodString>;
            description: z.ZodRecord<z.ZodString, z.ZodString>;
            buttons: z.ZodArray<z.ZodObject<{
                link: z.ZodString;
                type: z.ZodEnum<{
                    readonly EXTERNAL: "external";
                    readonly SUBSCRIPTION_LINK: "subscriptionLink";
                    readonly COPY_BUTTON: "copyButton";
                }>;
                text: z.ZodRecord<z.ZodString, z.ZodString>;
                svgIconKey: z.ZodString;
            }, z.core.$strip>>;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
declare const BrandingSettingsSchema: z.ZodObject<{
    title: z.ZodString;
    logoUrl: z.ZodString;
    supportUrl: z.ZodURL;
}, z.core.$strip>;
declare const UiConfigSchema: z.ZodObject<{
    subscriptionInfoBlockType: z.ZodEnum<{
        readonly COLLAPSED: "collapsed";
        readonly EXPANDED: "expanded";
        readonly CARDS: "cards";
        readonly HIDDEN: "hidden";
    }>;
    installationGuidesBlockType: z.ZodEnum<{
        readonly CARDS: "cards";
        readonly ACCORDION: "accordion";
        readonly MINIMAL: "minimal";
        readonly TIMELINE: "timeline";
    }>;
}, z.core.$strip>;
declare const SubscriptionPageTranslateKeysSchema: z.ZodObject<{
    installationGuideHeader: z.ZodRecord<z.ZodString, z.ZodString>;
    connectionKeysHeader: z.ZodRecord<z.ZodString, z.ZodString>;
    linkCopied: z.ZodRecord<z.ZodString, z.ZodString>;
    linkCopiedToClipboard: z.ZodRecord<z.ZodString, z.ZodString>;
    getLink: z.ZodRecord<z.ZodString, z.ZodString>;
    scanQrCode: z.ZodRecord<z.ZodString, z.ZodString>;
    scanQrCodeDescription: z.ZodRecord<z.ZodString, z.ZodString>;
    copyLink: z.ZodRecord<z.ZodString, z.ZodString>;
    name: z.ZodRecord<z.ZodString, z.ZodString>;
    status: z.ZodRecord<z.ZodString, z.ZodString>;
    active: z.ZodRecord<z.ZodString, z.ZodString>;
    inactive: z.ZodRecord<z.ZodString, z.ZodString>;
    expires: z.ZodRecord<z.ZodString, z.ZodString>;
    bandwidth: z.ZodRecord<z.ZodString, z.ZodString>;
    scanToImport: z.ZodRecord<z.ZodString, z.ZodString>;
    expiresIn: z.ZodRecord<z.ZodString, z.ZodString>;
    expired: z.ZodRecord<z.ZodString, z.ZodString>;
    unknown: z.ZodRecord<z.ZodString, z.ZodString>;
    indefinitely: z.ZodRecord<z.ZodString, z.ZodString>;
}, z.core.$strip>;
export declare const SubscriptionPageRawConfigSchema: z.ZodObject<{
    version: z.ZodEnum<{
        readonly 1: "1";
    }>;
    locales: z.ZodArray<z.ZodEnum<{
        en: "en";
        ru: "ru";
        zh: "zh";
        fr: "fr";
        fa: "fa";
        uz: "uz";
        de: "de";
        hi: "hi";
        tr: "tr";
        az: "az";
        es: "es";
        vi: "vi";
        ja: "ja";
        be: "be";
        uk: "uk";
        pt: "pt";
        pl: "pl";
        id: "id";
        tk: "tk";
        th: "th";
    }>>;
    brandingSettings: z.ZodObject<{
        title: z.ZodString;
        logoUrl: z.ZodString;
        supportUrl: z.ZodURL;
    }, z.core.$strip>;
    uiConfig: z.ZodObject<{
        subscriptionInfoBlockType: z.ZodEnum<{
            readonly COLLAPSED: "collapsed";
            readonly EXPANDED: "expanded";
            readonly CARDS: "cards";
            readonly HIDDEN: "hidden";
        }>;
        installationGuidesBlockType: z.ZodEnum<{
            readonly CARDS: "cards";
            readonly ACCORDION: "accordion";
            readonly MINIMAL: "minimal";
            readonly TIMELINE: "timeline";
        }>;
    }, z.core.$strip>;
    baseSettings: z.ZodDefault<z.ZodObject<{
        metaTitle: z.ZodDefault<z.ZodString>;
        metaDescription: z.ZodDefault<z.ZodString>;
        showConnectionKeys: z.ZodDefault<z.ZodBoolean>;
        hideGetLinkButton: z.ZodDefault<z.ZodBoolean>;
    }, z.core.$strip>>;
    baseTranslations: z.ZodObject<{
        installationGuideHeader: z.ZodRecord<z.ZodString, z.ZodString>;
        connectionKeysHeader: z.ZodRecord<z.ZodString, z.ZodString>;
        linkCopied: z.ZodRecord<z.ZodString, z.ZodString>;
        linkCopiedToClipboard: z.ZodRecord<z.ZodString, z.ZodString>;
        getLink: z.ZodRecord<z.ZodString, z.ZodString>;
        scanQrCode: z.ZodRecord<z.ZodString, z.ZodString>;
        scanQrCodeDescription: z.ZodRecord<z.ZodString, z.ZodString>;
        copyLink: z.ZodRecord<z.ZodString, z.ZodString>;
        name: z.ZodRecord<z.ZodString, z.ZodString>;
        status: z.ZodRecord<z.ZodString, z.ZodString>;
        active: z.ZodRecord<z.ZodString, z.ZodString>;
        inactive: z.ZodRecord<z.ZodString, z.ZodString>;
        expires: z.ZodRecord<z.ZodString, z.ZodString>;
        bandwidth: z.ZodRecord<z.ZodString, z.ZodString>;
        scanToImport: z.ZodRecord<z.ZodString, z.ZodString>;
        expiresIn: z.ZodRecord<z.ZodString, z.ZodString>;
        expired: z.ZodRecord<z.ZodString, z.ZodString>;
        unknown: z.ZodRecord<z.ZodString, z.ZodString>;
        indefinitely: z.ZodRecord<z.ZodString, z.ZodString>;
    }, z.core.$strip>;
    svgLibrary: z.ZodRecord<z.ZodString, z.ZodString>;
    platforms: z.ZodRecord<z.ZodEnum<{
        readonly IOS: "ios";
        readonly ANDROID: "android";
        readonly LINUX: "linux";
        readonly MACOS: "macos";
        readonly WINDOWS: "windows";
        readonly ANDROID_TV: "androidTV";
        readonly APPLE_TV: "appleTV";
    }> & z.core.$partial, z.ZodObject<{
        displayName: z.ZodRecord<z.ZodString, z.ZodString>;
        svgIconKey: z.ZodString;
        apps: z.ZodArray<z.ZodObject<{
            name: z.ZodString;
            svgIconKey: z.ZodOptional<z.ZodString>;
            featured: z.ZodBoolean;
            blocks: z.ZodArray<z.ZodObject<{
                svgIconKey: z.ZodString;
                svgIconColor: z.ZodString;
                title: z.ZodRecord<z.ZodString, z.ZodString>;
                description: z.ZodRecord<z.ZodString, z.ZodString>;
                buttons: z.ZodArray<z.ZodObject<{
                    link: z.ZodString;
                    type: z.ZodEnum<{
                        readonly EXTERNAL: "external";
                        readonly SUBSCRIPTION_LINK: "subscriptionLink";
                        readonly COPY_BUTTON: "copyButton";
                    }>;
                    text: z.ZodRecord<z.ZodString, z.ZodString>;
                    svgIconKey: z.ZodString;
                }, z.core.$strip>>;
            }, z.core.$strip>>;
        }, z.core.$strip>>;
    }, z.core.$strip>>;
}, z.core.$strip>;
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