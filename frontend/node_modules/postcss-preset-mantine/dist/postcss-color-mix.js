"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const getValueRegexp = (name) => new RegExp('\\b' + name + '\\(([^()]+)\\)', 'g');
const getVarRegexp = (name) => new RegExp('\\b' + name + '\\(([^()]*\\([^()]*\\)[^()]*)+\\)', 'g');
function replaceValues(root, fn, replace) {
    root.replaceValues(getValueRegexp(fn), { fast: `${fn}(` }, (_, values) => replace(values));
    root.replaceValues(getVarRegexp(fn), { fast: `${fn}(` }, (_, values) => replace(values));
}
function getParsedColor(input) {
    if (typeof input !== 'string') {
        return null;
    }
    const color = input.trim();
    const lastCommaIndex = color.lastIndexOf(',');
    if (lastCommaIndex === -1) {
        return null;
    }
    const rawPayload = color.slice(lastCommaIndex + 1).trim();
    const payload = rawPayload.endsWith('%')
        ? Number(rawPayload.slice(0, -1)) / 100
        : Number(color.slice(lastCommaIndex + 1));
    if (Number.isNaN(payload)) {
        return null;
    }
    return {
        color: color.slice(0, lastCommaIndex).trim(),
        payload: Math.max(0, Math.min(1, payload)),
    };
}
function alpha(input) {
    const parsed = getParsedColor(input);
    if (!parsed) {
        return input;
    }
    if (parsed.payload === 1) {
        return parsed.color;
    }
    if (parsed.payload === 0) {
        return 'transparent';
    }
    const mixPercentage = (1 - parsed.payload) * 100;
    return `color-mix(in srgb, ${parsed.color}, transparent ${mixPercentage}%)`;
}
function lighten(input) {
    const parsed = getParsedColor(input);
    if (!parsed) {
        return input;
    }
    return `color-mix(in srgb, ${parsed.color}, white ${parsed.payload * 100}%)`;
}
function darken(input) {
    const parsed = getParsedColor(input);
    if (!parsed) {
        return input;
    }
    return `color-mix(in srgb, ${parsed.color}, black ${parsed.payload * 100}%)`;
}
module.exports = () => {
    return {
        postcssPlugin: 'postcss-mantine-color-mix',
        Once(root) {
            replaceValues(root, 'alpha', alpha);
            replaceValues(root, 'lighten', lighten);
            replaceValues(root, 'darken', darken);
        },
    };
};
module.exports.postcss = true;
