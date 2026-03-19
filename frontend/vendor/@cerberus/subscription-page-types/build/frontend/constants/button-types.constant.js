"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isButtonType = exports.BUTTON_TYPES_VALUES = exports.BUTTON_TYPES = void 0;
exports.BUTTON_TYPES = {
    EXTERNAL: 'external',
    SUBSCRIPTION_LINK: 'subscriptionLink',
    COPY_BUTTON: 'copyButton',
};
exports.BUTTON_TYPES_VALUES = Object.values(exports.BUTTON_TYPES);
const isButtonType = (value) => {
    return exports.BUTTON_TYPES_VALUES.includes(value);
};
exports.isButtonType = isButtonType;
