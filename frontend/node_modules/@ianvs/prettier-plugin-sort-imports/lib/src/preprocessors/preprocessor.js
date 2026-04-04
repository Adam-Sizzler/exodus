"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.preprocessor = preprocessor;
const parser_1 = require("@babel/parser");
const extract_ast_nodes_1 = require("../utils/extract-ast-nodes");
const get_code_from_ast_1 = require("../utils/get-code-from-ast");
const get_sorted_nodes_1 = require("../utils/get-sorted-nodes");
const normalize_plugin_options_1 = require("../utils/normalize-plugin-options");
/**
 *
 * @param originalCode The raw source code from the file being processed
 * @param parseableCode This is code that has been transformed (if necessary) so that babel can parse it, but is the same size as the original code
 * @param options PrettierOptions
 * @returns
 */
function preprocessor(originalCode, { options, parseableCode, }) {
    const { plugins, ...remainingOptions } = (0, normalize_plugin_options_1.examineAndNormalizePluginOptions)(options);
    const parserOptions = {
        sourceType: 'module',
        attachComment: true,
        errorRecovery: true,
        allowReturnOutsideFunction: true,
        allowNewTargetOutsideFunction: true,
        allowSuperOutsideMethod: true,
        allowUndeclaredExports: true,
        plugins,
    };
    // short-circuit if importOrder is an empty array (can be used to disable plugin)
    if (!remainingOptions.importOrder.length) {
        return originalCode;
    }
    let ast;
    try {
        ast = (0, parser_1.parse)(parseableCode ?? originalCode, parserOptions);
    }
    catch (err) {
        // Don't throw warning messages if this is a codeblock in markdown.  (Prettier doesn't either)
        if (options.parentParser !== 'markdown' &&
            options.parentParser !== 'mdx') {
            console.error(' [error] [prettier-plugin-sort-imports]: import sorting aborted due to parsing error:\n%s', err);
        }
        return originalCode;
    }
    const directives = ast.program.directives;
    const interpreter = ast.program.interpreter;
    const { importDeclarations: allOriginalImportNodes } = (0, extract_ast_nodes_1.extractASTNodes)(ast);
    // short-circuit if there are no import declarations
    if (allOriginalImportNodes.length === 0) {
        return originalCode;
    }
    const nodesToOutput = (0, get_sorted_nodes_1.getSortedNodes)(allOriginalImportNodes, remainingOptions);
    return (0, get_code_from_ast_1.getCodeFromAst)({
        nodesToOutput,
        allOriginalImportNodes,
        originalCode,
        directives,
        interpreter,
    });
}
