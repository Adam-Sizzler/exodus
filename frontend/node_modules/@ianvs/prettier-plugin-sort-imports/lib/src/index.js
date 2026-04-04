"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.printers = exports.parsers = exports.options = void 0;
const parser_babel_1 = require("prettier/parser-babel");
const parser_flow_1 = require("prettier/parser-flow");
const parser_html_1 = require("prettier/parser-html");
const parser_typescript_1 = require("prettier/parser-typescript");
const constants_1 = require("./constants");
const default_1 = require("./preprocessors/default");
const ember_1 = require("./preprocessors/ember");
const vue_1 = require("./preprocessors/vue");
exports.options = {
    importOrder: {
        type: 'string',
        category: 'Global',
        array: true,
        default: [{ value: constants_1.DEFAULT_IMPORT_ORDER }],
        description: 'An array of regex strings for the import sort order.',
    },
    importOrderParserPlugins: {
        type: 'string',
        category: 'Global',
        array: true,
        // By default, we add ts and jsx as parsers but if users define something
        // we take that option
        default: [{ value: ['typescript', 'jsx'] }],
        description: 'A list of babel parser plugins for special syntax.',
    },
    importOrderTypeScriptVersion: {
        type: 'string',
        category: 'Global',
        default: '1.0.0',
        description: 'Version of TypeScript in use in the project.  Determines some output syntax when using TypeScript.',
    },
    importOrderCaseSensitive: {
        type: 'boolean',
        category: 'Global',
        default: false,
        description: 'Should capitalization be considered when sorting imports?',
    },
    importOrderSafeSideEffects: {
        type: 'string',
        category: 'Global',
        array: true,
        default: [{ value: [] }],
        description: 'Array of globs for side-effect-only imports that are considered safe to sort.',
    },
};
const getEmberPlugin = () => {
    try {
        const emberPlugin = require('prettier-plugin-ember-template-tag');
        return emberPlugin;
    }
    catch {
        throw new Error('prettier-plugin-ember-template-tag could not be loaded.  Is it installed?');
    }
};
const getOxcPlugin = () => {
    try {
        const oxcPlugin = require('@prettier/plugin-oxc');
        return oxcPlugin;
    }
    catch {
        throw new Error('@prettier/plugin-oxc could not be loaded.  Is it installed?');
    }
};
exports.parsers = {
    babel: {
        ...parser_babel_1.parsers.babel,
        preprocess: default_1.defaultPreprocessor,
    },
    'babel-ts': {
        ...parser_babel_1.parsers['babel-ts'],
        preprocess: default_1.defaultPreprocessor,
    },
    get 'ember-template-tag'() {
        const emberPlugin = getEmberPlugin();
        return {
            ...emberPlugin.parsers['ember-template-tag'],
            preprocess: ember_1.emberPreprocessor,
        };
    },
    flow: {
        ...parser_flow_1.parsers.flow,
        preprocess: default_1.defaultPreprocessor,
    },
    get oxc() {
        const oxcPlugin = getOxcPlugin();
        return {
            ...oxcPlugin.parsers.oxc,
            preprocess: default_1.defaultPreprocessor,
        };
    },
    get 'oxc-ts'() {
        const oxcPlugin = getOxcPlugin();
        return {
            ...oxcPlugin.parsers['oxc-ts'],
            preprocess: default_1.defaultPreprocessor,
        };
    },
    typescript: {
        ...parser_typescript_1.parsers.typescript,
        preprocess: default_1.defaultPreprocessor,
    },
    vue: {
        ...parser_html_1.parsers.vue,
        preprocess: vue_1.vuePreprocessor,
    },
};
exports.printers = {
    get 'estree-oxc'() {
        const oxcPlugin = getOxcPlugin();
        return oxcPlugin.printers['estree-oxc'];
    },
};
