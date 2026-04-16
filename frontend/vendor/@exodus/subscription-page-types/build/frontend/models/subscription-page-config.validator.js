"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.cleanLocalizedTexts = exports.validateSvgReferences = exports.validateLocalizedTexts = void 0;
const zod_1 = require("zod");
const constants_1 = require("../constants");
const isLocalizedText = (obj) => {
    if (obj === null || typeof obj !== 'object' || Array.isArray(obj))
        return false;
    const entries = Object.entries(obj);
    if (entries.length === 0)
        return false;
    const hasStringValues = entries.every(([, value]) => typeof value === 'string');
    const hasLanguageCodeKey = entries.some(([key]) => key.length === 2);
    return hasStringValues && hasLanguageCodeKey;
};
const validateLocalizedTexts = (data, requiredLocales, ctx) => {
    const checkLocalizedText = (obj, path) => {
        if (obj === null || typeof obj !== 'object')
            return;
        if (isLocalizedText(obj)) {
            for (const locale of requiredLocales) {
                const value = obj[locale];
                if (!value || value.trim() === '') {
                    ctx.addIssue({
                        code: zod_1.z.ZodIssueCode.custom,
                        message: `Missing required locale '${locale}' at ${path}`,
                        path: [path, locale],
                    });
                }
            }
            return;
        }
        if (Array.isArray(obj)) {
            obj.forEach((item, index) => {
                checkLocalizedText(item, `${path}[${index}]`);
            });
            return;
        }
        for (const [key, value] of Object.entries(obj)) {
            if (typeof value === 'object' && value !== null) {
                checkLocalizedText(value, `${path}.${key}`);
            }
        }
    };
    checkLocalizedText(data.platforms, 'platforms');
    checkLocalizedText(data.baseTranslations, 'baseTranslations');
};
exports.validateLocalizedTexts = validateLocalizedTexts;
const validateSvgReferences = (data, ctx) => {
    const validKeys = new Set(Object.keys(data.svgLibrary));
    const checkSvgRef = (obj, path) => {
        if (obj === null || typeof obj !== 'object')
            return;
        for (const [key, value] of Object.entries(obj)) {
            if (key === 'svgIconKey' && typeof value === 'string') {
                if (!validKeys.has(value)) {
                    ctx.addIssue({
                        code: zod_1.z.ZodIssueCode.custom,
                        message: `Unknown svgIconKey '${value}' at ${path}.${key}. Available: ${[...validKeys].join(', ')}`,
                        path: [path, key],
                    });
                }
            }
            else if (Array.isArray(value)) {
                value.forEach((item, index) => checkSvgRef(item, `${path}.${key}[${index}]`));
            }
            else if (typeof value === 'object' && value !== null) {
                checkSvgRef(value, `${path}.${key}`);
            }
        }
    };
    checkSvgRef(data.platforms, 'platforms');
};
exports.validateSvgReferences = validateSvgReferences;
const cleanLocalizedTexts = (obj, activeLocales) => {
    if (obj === null || typeof obj !== 'object')
        return obj;
    if (isLocalizedText(obj)) {
        const allowedKeys = new Set(activeLocales);
        const cleaned = {};
        for (const [key, value] of Object.entries(obj)) {
            if (typeof value === 'string' && (0, constants_1.isLanguageCode)(key) && allowedKeys.has(key)) {
                const trimmed = value.trim();
                if (trimmed) {
                    cleaned[key] = trimmed;
                }
            }
        }
        return cleaned;
    }
    if (Array.isArray(obj)) {
        return obj.map((item) => (0, exports.cleanLocalizedTexts)(item, activeLocales));
    }
    const result = {};
    for (const [key, value] of Object.entries(obj)) {
        result[key] = (0, exports.cleanLocalizedTexts)(value, activeLocales);
    }
    return result;
};
exports.cleanLocalizedTexts = cleanLocalizedTexts;
