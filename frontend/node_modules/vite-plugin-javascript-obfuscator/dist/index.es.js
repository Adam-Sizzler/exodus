import { obfuscate } from 'javascript-obfuscator';
import anymatch from 'anymatch';
import { resolve } from 'path';

const defaultIncludeMatcher = [/\.(jsx?|tsx?|cjs|mjs)$/];
const defaultExcludeMatcher = [/node_modules/, /\.nuxt/];
function handleMatcher(matcher) {
    matcher = matcher instanceof Array ? matcher : [matcher];
    return matcher.map((matcher) => {
        if (typeof matcher !== "string") {
            return matcher;
        }
        return resolve(".", matcher).replace(/\\/g, "/");
    });
}
function obfuscatorPlugin(obOptions) {
    let { include, exclude, options } = obOptions || {};
    const consoleLog = (obOptions === null || obOptions === void 0 ? void 0 : obOptions.debugger) ? console.log.bind(console) : () => { };
    options = options || {};
    const includeMatcher = include
        ? handleMatcher(include)
        : defaultIncludeMatcher;
    const excludeMatcher = exclude
        ? handleMatcher(exclude)
        : defaultExcludeMatcher;
    return {
        name: "vite-plugin-javascript-obfuscator",
        enforce: "post",
        apply: (obOptions === null || obOptions === void 0 ? void 0 : obOptions.apply) || (() => true),
        transform(src, id) {
            if (anymatch(excludeMatcher, id, { dot: true })) {
                consoleLog("[::plugin-javascript-obfuscator]::exclude", id);
                return;
            }
            if (anymatch(includeMatcher, id)) {
                consoleLog("[::plugin-javascript-obfuscator]::include matched", id);
                const obfuscationResult = obfuscate(src, options);
                const result = { code: obfuscationResult.getObfuscatedCode() };
                if ((options === null || options === void 0 ? void 0 : options.sourceMap) && (options === null || options === void 0 ? void 0 : options.sourceMapMode) !== "inline") {
                    result.map = obfuscationResult.getSourceMap();
                }
                return result;
            }
            consoleLog(`[::plugin-javascript-obfuscator]::not matched`, id);
        },
    };
}

export { obfuscatorPlugin as default };
