"use strict";
/**
 * @fileoverview Meta should be followed by `satisfies Meta`
 * @author Tiger Oakes
 */
const utils_1 = require("@typescript-eslint/utils");
const utils_2 = require("../utils");
const create_storybook_rule_1 = require("../utils/create-storybook-rule");
const ast_1 = require("../utils/ast");
module.exports = (0, create_storybook_rule_1.createStorybookRule)({
    name: 'meta-satisfies-type',
    defaultOptions: [],
    meta: {
        type: 'problem',
        fixable: 'code',
        severity: 'error',
        docs: {
            description: 'Meta should use `satisfies Meta`',
            categories: [],
            excludeFromConfig: true,
        },
        messages: {
            metaShouldSatisfyType: 'CSF Meta should use `satisfies` for type safety',
        },
        schema: [],
    },
    create(context) {
        // variables should be defined here
        const sourceCode = context.getSourceCode();
        //----------------------------------------------------------------------
        // Helpers
        //----------------------------------------------------------------------
        const getTextWithParentheses = (node) => {
            // Capture parentheses before and after the node
            let beforeCount = 0;
            let afterCount = 0;
            if (utils_1.ASTUtils.isParenthesized(node, sourceCode)) {
                const bodyOpeningParen = sourceCode.getTokenBefore(node, utils_1.ASTUtils.isOpeningParenToken);
                const bodyClosingParen = sourceCode.getTokenAfter(node, utils_1.ASTUtils.isClosingParenToken);
                if (bodyOpeningParen && bodyClosingParen) {
                    beforeCount = node.range[0] - bodyOpeningParen.range[0];
                    afterCount = bodyClosingParen.range[1] - node.range[1];
                }
            }
            return sourceCode.getText(node, beforeCount, afterCount);
        };
        const getFixer = (meta) => {
            const { parent } = meta;
            if (!parent) {
                return undefined;
            }
            switch (parent.type) {
                // {} as Meta
                case utils_1.AST_NODE_TYPES.TSAsExpression:
                    return (fixer) => [
                        fixer.replaceText(parent, getTextWithParentheses(meta)),
                        fixer.insertTextAfter(parent, ` satisfies ${getTextWithParentheses(parent.typeAnnotation)}`),
                    ];
                // const meta: Meta = {}
                case utils_1.AST_NODE_TYPES.VariableDeclarator: {
                    const { typeAnnotation } = parent.id;
                    if (typeAnnotation) {
                        return (fixer) => [
                            fixer.remove(typeAnnotation),
                            fixer.insertTextAfter(meta, ` satisfies ${getTextWithParentheses(typeAnnotation.typeAnnotation)}`),
                        ];
                    }
                    return undefined;
                }
                default:
                    return undefined;
            }
        };
        // any helper functions should go here or else delete this section
        //----------------------------------------------------------------------
        // Public
        //----------------------------------------------------------------------
        return {
            ExportDefaultDeclaration(node) {
                const meta = (0, utils_2.getMetaObjectExpression)(node, context);
                if (!meta) {
                    return null;
                }
                if (!meta.parent || !(0, ast_1.isTSSatisfiesExpression)(meta.parent)) {
                    context.report({
                        node: meta,
                        messageId: 'metaShouldSatisfyType',
                        fix: getFixer(meta),
                    });
                }
            },
        };
    },
});
