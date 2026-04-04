"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
const unitsConverters = require('./converters');
const rem = unitsConverters.remStrict;
module.exports = () => {
    return {
        postcssPlugin: 'postcss-auto-rem',
        Declaration: (decl) => {
            if (!decl.value.includes('px')) {
                return;
            }
            if (decl.prop === 'content') {
                return;
            }
            decl.value = rem(decl.value);
        },
    };
};
module.exports.postcss = true;
