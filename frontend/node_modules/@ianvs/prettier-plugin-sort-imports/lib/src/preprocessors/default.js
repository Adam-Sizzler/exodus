"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.defaultPreprocessor = defaultPreprocessor;
const preprocessor_1 = require("./preprocessor");
function defaultPreprocessor(code, options) {
    if (options.filepath?.endsWith('.vue'))
        return code;
    return (0, preprocessor_1.preprocessor)(code, { options });
}
