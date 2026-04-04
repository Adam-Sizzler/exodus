"use strict";
var __importDefault = (this && this.__importDefault) || function (mod) {
    return (mod && mod.__esModule) ? mod : { "default": mod };
};
Object.defineProperty(exports, "__esModule", { value: true });
exports.extractASTNodes = extractASTNodes;
const traverse_1 = __importDefault(require("@babel/traverse"));
function extractASTNodes(ast) {
    const importDeclarations = [];
    (0, traverse_1.default)(ast, {
        noScope: true, // This is required in order to avoid traverse errors if a variable is redefined (https://github.com/babel/babel/issues/12950#issuecomment-788974837)
        ImportDeclaration(path) {
            const tsModuleParent = path.findParent((p) => p.isTSModuleDeclaration());
            // Do not sort imports inside of typescript module declarations.  See `import-inside-ts-declare.ts` test.
            if (!tsModuleParent) {
                importDeclarations.push(path.node);
            }
        },
    });
    return { importDeclarations };
}
