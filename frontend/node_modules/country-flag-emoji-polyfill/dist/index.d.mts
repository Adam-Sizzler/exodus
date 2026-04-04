/**
 * Injects a style element into the HEAD with a web font with country flags,
 * iff the browser does support emojis but not country flags.
 *
 * @param fontName - Override the default font name ("Twemoji Country Flags")
 * @param fontUrl - Override the font URL (defaults to a jsdeliver-hosted)
 *
 * @returns true if the web font was loaded (ie the browser does not support country flags)
 */
declare function polyfillCountryFlagEmojis(fontName?: string, fontUrl?: string): boolean;

export { polyfillCountryFlagEmojis };
