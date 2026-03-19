import { z } from 'zod';
import { TSubscriptionPageLanguageCode } from '../constants';
export declare const validateLocalizedTexts: (data: {
    locales: TSubscriptionPageLanguageCode[];
    platforms: Record<string, unknown>;
    uiConfig: unknown;
    baseTranslations: unknown;
}, requiredLocales: TSubscriptionPageLanguageCode[], ctx: z.RefinementCtx) => void;
export declare const validateSvgReferences: (data: {
    svgLibrary: Record<string, string>;
    platforms: Record<string, unknown>;
}, ctx: z.RefinementCtx) => void;
export declare const cleanLocalizedTexts: <T>(obj: T, activeLocales: TSubscriptionPageLanguageCode[]) => T;
//# sourceMappingURL=subscription-page-config.validator.d.ts.map