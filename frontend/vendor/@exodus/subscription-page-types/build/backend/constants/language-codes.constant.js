"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.getLanguagesArray = exports.isLanguageCode = exports.getLanguageNativeName = exports.getLanguageName = exports.getLanguageInfo = exports.LANGUAGE_CODES = void 0;
const LANGUAGES_LIST = {
    en: { name: 'English', nativeName: 'English', emoji: '🇬🇧' },
    ru: { name: 'Russian', nativeName: 'Русский', emoji: '🇷🇺' },
    zh: { name: 'Chinese', nativeName: '中文', emoji: '🇨🇳' },
    fr: { name: 'French', nativeName: 'Français', emoji: '🇫🇷' },
    fa: { name: 'Persian', nativeName: 'فارسی', emoji: '🇮🇷' },
    uz: { name: 'Uzbek', nativeName: 'Ўзбек', emoji: '🇺🇿' },
    de: { name: 'German', nativeName: 'Deutsch', emoji: '🇩🇪' },
    hi: { name: 'Hindi', nativeName: 'हिन्दी', emoji: '🇮🇳' },
    tr: { name: 'Turkish', nativeName: 'Türkçe', emoji: '🇹🇷' },
    az: { name: 'Azerbaijani', nativeName: 'azərbaycan dili', emoji: '🇦🇿' },
    es: { name: 'Spanish', nativeName: 'Español', emoji: '🇪🇸' },
    vi: { name: 'Vietnamese', nativeName: 'Tiếng Việt', emoji: '🇻🇳' },
    ja: { name: 'Japanese', nativeName: '日本語', emoji: '🇯🇵' },
    be: { name: 'Belarusian', nativeName: 'беларуская мова', emoji: '🇧🇾' },
    uk: { name: 'Ukrainian', nativeName: 'Українська', emoji: '🇺🇦' },
    pt: { name: 'Portuguese', nativeName: 'Português', emoji: '🇵🇹' },
    pl: { name: 'Polish', nativeName: 'Polski', emoji: '🇵🇱' },
    id: { name: 'Indonesian', nativeName: 'Bahasa Indonesia', emoji: '🇮🇩' },
    tk: { name: 'Turkmen', nativeName: 'Türkmençe', emoji: '🇹🇲' },
    th: { name: 'Thai', nativeName: 'ไทย', emoji: '🇹🇭' },
};
exports.LANGUAGE_CODES = Object.keys(LANGUAGES_LIST);
/**
 * Get language info by code
 * @param code - ISO 639-1 language code
 * @returns Language info object with name and nativeName
 */
const getLanguageInfo = (code) => {
    if ((0, exports.isLanguageCode)(code)) {
        return LANGUAGES_LIST[code];
    }
    return null;
};
exports.getLanguageInfo = getLanguageInfo;
/**
 * Get language name (English) by code
 * @param code - ISO 639-1 language code
 * @returns English language name or the code itself if not found
 */
const getLanguageName = (code) => {
    if ((0, exports.isLanguageCode)(code)) {
        return LANGUAGES_LIST[code].name;
    }
    return code;
};
exports.getLanguageName = getLanguageName;
/**
 * Get native language name by code
 * @param code - ISO 639-1 language code
 * @returns Native language name or the code itself if not found
 */
const getLanguageNativeName = (code) => {
    if ((0, exports.isLanguageCode)(code)) {
        return LANGUAGES_LIST[code].nativeName;
    }
    return code;
};
exports.getLanguageNativeName = getLanguageNativeName;
/**
 * Check if a string is a valid language code
 */
const isLanguageCode = (code) => {
    return code in LANGUAGES_LIST;
};
exports.isLanguageCode = isLanguageCode;
/**
 * Get all languages as array for UI selects
 */
const getLanguagesArray = () => {
    return Object.entries(LANGUAGES_LIST).map(([code, info]) => ({
        code,
        ...info,
    }));
};
exports.getLanguagesArray = getLanguagesArray;
