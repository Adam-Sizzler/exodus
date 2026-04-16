"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.SubscriptionPageRawConfigSchema = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
const subscription_page_config_validator_1 = require("./subscription-page-config.validator");
const LocalizedTextSchema = zod_1.z
    .record(zod_1.z.string().regex(/^[a-z]{2}$/, 'Language code must be 2 lowercase letters'), zod_1.z.string())
    .refine((obj) => Object.keys(obj).length > 0, {
    message: 'At least one language must be specified',
});
const SvgLibrarySchema = zod_1.z.record(zod_1.z.string().regex(/^[A-Za-z]+$/, { message: 'Only latin characters, no spaces allowed' }), zod_1.z.string());
const ButtonSchema = zod_1.z.object({
    link: zod_1.z.string(),
    type: zod_1.z.nativeEnum(constants_1.BUTTON_TYPES),
    text: LocalizedTextSchema,
    svgIconKey: zod_1.z.string(),
});
const BlockSchema = zod_1.z.object({
    svgIconKey: zod_1.z.string(),
    svgIconColor: zod_1.z
        .string()
        .refine((value) => [
        'blue',
        'cyan',
        'dark',
        'grape',
        'gray',
        'green',
        'indigo',
        'lime',
        'orange',
        'pink',
        'red',
        'teal',
        'violet',
        'yellow',
    ].includes(value) || /^#[0-9a-fA-F]{3,8}$/.test(value), {
        message: 'svgIconColor must be one of the predefined colors or a hex color beginning with #',
    }),
    title: LocalizedTextSchema,
    description: LocalizedTextSchema,
    buttons: zod_1.z.array(ButtonSchema),
});
const PlatformAppSchema = zod_1.z.object({
    name: zod_1.z.string().min(2, 'Name must be at least 2 characters long'),
    svgIconKey: zod_1.z.optional(zod_1.z.string()),
    featured: zod_1.z.boolean(),
    blocks: zod_1.z.array(BlockSchema),
});
const PlatformSchema = zod_1.z.object({
    displayName: LocalizedTextSchema,
    svgIconKey: zod_1.z.string(),
    apps: zod_1.z.array(PlatformAppSchema),
});
const BrandingSettingsSchema = zod_1.z.object({
    title: zod_1.z.string(),
    logoUrl: zod_1.z.string(),
    supportUrl: zod_1.z.string().url(),
});
const UiConfigSchema = zod_1.z.object({
    subscriptionInfoBlockType: zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_INFO_BLOCK_VARIANTS),
    installationGuidesBlockType: zod_1.z.nativeEnum(constants_1.INSTALLATION_GUIDE_BLOCKS_VARIANTS),
});
const SubscriptionPageTranslateKeysSchema = zod_1.z.object({
    installationGuideHeader: LocalizedTextSchema,
    connectionKeysHeader: LocalizedTextSchema,
    linkCopied: LocalizedTextSchema,
    linkCopiedToClipboard: LocalizedTextSchema,
    getLink: LocalizedTextSchema,
    scanQrCode: LocalizedTextSchema,
    scanQrCodeDescription: LocalizedTextSchema,
    copyLink: LocalizedTextSchema,
    name: LocalizedTextSchema,
    status: LocalizedTextSchema,
    active: LocalizedTextSchema,
    inactive: LocalizedTextSchema,
    expires: LocalizedTextSchema,
    bandwidth: LocalizedTextSchema,
    scanToImport: LocalizedTextSchema,
    expiresIn: LocalizedTextSchema,
    expired: LocalizedTextSchema,
    unknown: LocalizedTextSchema,
    indefinitely: LocalizedTextSchema,
});
const BaseSettingsSchema = zod_1.z
    .object({
    metaTitle: zod_1.z.string().default('Subscription'),
    metaDescription: zod_1.z.string().default('Subscription'),
    showConnectionKeys: zod_1.z.boolean().default(false),
    hideGetLinkButton: zod_1.z.boolean().default(false),
})
    .default({
    metaTitle: 'Subscription',
    metaDescription: 'Subscription',
    showConnectionKeys: false,
    hideGetLinkButton: false,
});
exports.SubscriptionPageRawConfigSchema = zod_1.z
    .object({
    version: zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_PAGE_CONFIG_VERSION),
    locales: zod_1.z.array(zod_1.z.enum(constants_1.LANGUAGE_CODES)).min(1, 'At least one locale must be specified'),
    brandingSettings: BrandingSettingsSchema,
    uiConfig: UiConfigSchema,
    baseSettings: BaseSettingsSchema,
    baseTranslations: SubscriptionPageTranslateKeysSchema,
    svgLibrary: SvgLibrarySchema,
    platforms: zod_1.z.record(zod_1.z.nativeEnum(constants_1.SUBSCRIPTION_PAGE_CONFIG_PLATFORM_TYPES), PlatformSchema),
})
    .superRefine((data, ctx) => {
    (0, subscription_page_config_validator_1.validateLocalizedTexts)(data, data.locales, ctx);
    (0, subscription_page_config_validator_1.validateSvgReferences)(data, ctx);
});
