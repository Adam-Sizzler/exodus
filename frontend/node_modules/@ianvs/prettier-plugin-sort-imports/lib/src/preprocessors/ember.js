"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.emberPreprocessor = emberPreprocessor;
const glimmer_content_tag_1 = require("../utils/glimmer-content-tag");
const preprocessor_1 = require("./preprocessor");
function emberPreprocessor(code, options) {
    let parseableCode = code;
    const templates = (0, glimmer_content_tag_1.extractTemplates)(code);
    for (const template of templates) {
        parseableCode = (0, glimmer_content_tag_1.preprocessTemplateRange)(template, parseableCode);
    }
    const sorted = (0, preprocessor_1.preprocessor)(code, { parseableCode, options });
    return sorted;
}
