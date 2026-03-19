declare const LANGUAGES_LIST: {
    readonly en: {
        readonly name: "English";
        readonly nativeName: "English";
        readonly emoji: "🇬🇧";
    };
    readonly ru: {
        readonly name: "Russian";
        readonly nativeName: "Русский";
        readonly emoji: "🇷🇺";
    };
    readonly zh: {
        readonly name: "Chinese";
        readonly nativeName: "中文";
        readonly emoji: "🇨🇳";
    };
    readonly fr: {
        readonly name: "French";
        readonly nativeName: "Français";
        readonly emoji: "🇫🇷";
    };
    readonly fa: {
        readonly name: "Persian";
        readonly nativeName: "فارسی";
        readonly emoji: "🇮🇷";
    };
    readonly uz: {
        readonly name: "Uzbek";
        readonly nativeName: "Ўзбек";
        readonly emoji: "🇺🇿";
    };
    readonly de: {
        readonly name: "German";
        readonly nativeName: "Deutsch";
        readonly emoji: "🇩🇪";
    };
    readonly hi: {
        readonly name: "Hindi";
        readonly nativeName: "हिन्दी";
        readonly emoji: "🇮🇳";
    };
    readonly tr: {
        readonly name: "Turkish";
        readonly nativeName: "Türkçe";
        readonly emoji: "🇹🇷";
    };
    readonly az: {
        readonly name: "Azerbaijani";
        readonly nativeName: "azərbaycan dili";
        readonly emoji: "🇦🇿";
    };
    readonly es: {
        readonly name: "Spanish";
        readonly nativeName: "Español";
        readonly emoji: "🇪🇸";
    };
    readonly vi: {
        readonly name: "Vietnamese";
        readonly nativeName: "Tiếng Việt";
        readonly emoji: "🇻🇳";
    };
    readonly ja: {
        readonly name: "Japanese";
        readonly nativeName: "日本語";
        readonly emoji: "🇯🇵";
    };
    readonly be: {
        readonly name: "Belarusian";
        readonly nativeName: "беларуская мова";
        readonly emoji: "🇧🇾";
    };
    readonly uk: {
        readonly name: "Ukrainian";
        readonly nativeName: "Українська";
        readonly emoji: "🇺🇦";
    };
    readonly pt: {
        readonly name: "Portuguese";
        readonly nativeName: "Português";
        readonly emoji: "🇵🇹";
    };
    readonly pl: {
        readonly name: "Polish";
        readonly nativeName: "Polski";
        readonly emoji: "🇵🇱";
    };
    readonly id: {
        readonly name: "Indonesian";
        readonly nativeName: "Bahasa Indonesia";
        readonly emoji: "🇮🇩";
    };
    readonly tk: {
        readonly name: "Turkmen";
        readonly nativeName: "Türkmençe";
        readonly emoji: "🇹🇲";
    };
    readonly th: {
        readonly name: "Thai";
        readonly nativeName: "ไทย";
        readonly emoji: "🇹🇭";
    };
};
export type TSubscriptionPageLanguageCode = keyof typeof LANGUAGES_LIST;
export declare const LANGUAGE_CODES: [TSubscriptionPageLanguageCode, ...TSubscriptionPageLanguageCode[]];
export interface ILanguageInfo {
    name: string;
    nativeName: string;
    emoji: string;
}
/**
 * Get language info by code
 * @param code - ISO 639-1 language code
 * @returns Language info object with name and nativeName
 */
export declare const getLanguageInfo: (code: string) => ILanguageInfo | null;
/**
 * Get language name (English) by code
 * @param code - ISO 639-1 language code
 * @returns English language name or the code itself if not found
 */
export declare const getLanguageName: (code: string) => string;
/**
 * Get native language name by code
 * @param code - ISO 639-1 language code
 * @returns Native language name or the code itself if not found
 */
export declare const getLanguageNativeName: (code: string) => string;
/**
 * Check if a string is a valid language code
 */
export declare const isLanguageCode: (code: string) => code is TSubscriptionPageLanguageCode;
/**
 * Get all languages as array for UI selects
 */
export declare const getLanguagesArray: () => Array<{
    code: TSubscriptionPageLanguageCode;
} & ILanguageInfo>;
export {};
//# sourceMappingURL=language-codes.constant.d.ts.map