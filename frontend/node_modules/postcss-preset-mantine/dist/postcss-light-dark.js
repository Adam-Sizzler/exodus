"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const postcss_1 = require("postcss");
const FUNCTION = 'light-dark(';
function splitStringAtCharacter(character, search) {
    let characterIndex = 0;
    let openedParentheses = 0;
    while (characterIndex < search.length &&
        (search[characterIndex] !== character || openedParentheses)) {
        if (search[characterIndex] === '(') {
            openedParentheses += 1;
        }
        if (search[characterIndex] === ')') {
            openedParentheses -= 1;
        }
        characterIndex += 1;
    }
    return [search.slice(0, characterIndex), search.slice(characterIndex + 1)];
}
function getLightDarkValue(value) {
    const [prefix, ...search] = value.split(FUNCTION);
    if (!search.length) {
        return { light: value, dark: value };
    }
    const [macro, suffix] = splitStringAtCharacter(')', search.join(FUNCTION));
    const [light, dark] = splitStringAtCharacter(',', macro);
    const parsedSuffix = getLightDarkValue(suffix);
    return {
        light: prefix + getLightDarkValue(light.trim()).light + parsedSuffix.light,
        dark: prefix + getLightDarkValue(dark.trim()).dark + parsedSuffix.dark,
    };
}
module.exports = () => {
    return {
        postcssPlugin: 'postcss-light-dark',
        Once(root) {
            root.walkDecls((decl) => {
                const { value } = decl;
                const regex = /\blight-dark\b/;
                if (regex.test(value)) {
                    const { light: lightVal, dark: darkVal } = getLightDarkValue(value);
                    const darkMixin = (0, postcss_1.atRule)({ name: 'mixin', params: 'dark' });
                    darkMixin.append(decl.clone({ value: darkVal }));
                    decl.parent?.insertAfter(decl, darkMixin);
                    decl.parent?.insertAfter(decl, decl.clone({ value: lightVal }));
                    decl.remove();
                }
            });
        },
    };
};
module.exports.postcss = true;
