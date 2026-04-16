"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION = exports.RESPONSE_RULES_RESPONSE_TYPES_DESCRIPTION = exports.RESPONSE_RULES_CONDITION_OPERATORS_VALUES = exports.RESPONSE_RULES_RESPONSE_TYPES_VALUES = exports.RESPONSE_RULES_RESPONSE_TYPES = exports.RESPONSE_RULES_CONDITION_OPERATORS = exports.RESPONSE_RULE_CONDITION_TYPES = exports.RESPONSE_RULES_OPERATORS = exports.RESPONSE_RULES_CONFIG_VERSION = void 0;
const subscription_template_1 = require("../subscription-template");
exports.RESPONSE_RULES_CONFIG_VERSION = {
    1: '1',
};
exports.RESPONSE_RULES_OPERATORS = {
    AND: 'AND',
    OR: 'OR',
};
exports.RESPONSE_RULE_CONDITION_TYPES = {
    HEADER: 'HEADER',
};
exports.RESPONSE_RULES_CONDITION_OPERATORS = {
    EQUALS: 'EQUALS',
    NOT_EQUALS: 'NOT_EQUALS',
    CONTAINS: 'CONTAINS',
    NOT_CONTAINS: 'NOT_CONTAINS',
    STARTS_WITH: 'STARTS_WITH',
    NOT_STARTS_WITH: 'NOT_STARTS_WITH',
    ENDS_WITH: 'ENDS_WITH',
    NOT_ENDS_WITH: 'NOT_ENDS_WITH',
    REGEX: 'REGEX',
    NOT_REGEX: 'NOT_REGEX',
};
exports.RESPONSE_RULES_RESPONSE_TYPES = Object.assign(Object.assign({}, subscription_template_1.SUBSCRIPTION_TEMPLATE_TYPE), { BROWSER: 'BROWSER', BLOCK: 'BLOCK', STATUS_CODE_404: 'STATUS_CODE_404', STATUS_CODE_451: 'STATUS_CODE_451', SOCKET_DROP: 'SOCKET_DROP' });
exports.RESPONSE_RULES_RESPONSE_TYPES_VALUES = Object.values(exports.RESPONSE_RULES_RESPONSE_TYPES);
exports.RESPONSE_RULES_CONDITION_OPERATORS_VALUES = Object.values(exports.RESPONSE_RULES_CONDITION_OPERATORS);
exports.RESPONSE_RULES_RESPONSE_TYPES_DESCRIPTION = {
    [exports.RESPONSE_RULES_RESPONSE_TYPES.XRAY_JSON]: 'Return **subscription** in XRAY-JSON format. (Using `Xray Json` template)',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.XRAY_BASE64]: 'Return **subscription** in BASE64 encoded string. Compatible with most client application with Xray core.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.MIHOMO]: 'Return **subscription** in Mihomo format. (Using `Mihomo` template)',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.STASH]: 'Return **subscription** in Stash format. (Using `Stash` template)',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.CLASH]: 'Return **subscription** in Clash format. (Using `Clash` template) Useful for client application that use Legacy Clash core.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.SINGBOX]: 'Return **subscription** in Singbox format. (Using `Singbox` template) Format which is used by Singbox client application.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.BROWSER]: 'Return **subscription** as browser format. The same as on `/info` route.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.BLOCK]: '**Drop** request and return `403` status code.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.STATUS_CODE_404]: '**Drop** request and return `404` status code.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.STATUS_CODE_451]: '**Drop** request and return `451` status code.',
    [exports.RESPONSE_RULES_RESPONSE_TYPES.SOCKET_DROP]: '**Drop** the socket connection.',
};
exports.RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION = {
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.EQUALS]: 'Performs an exact, comparison between the header value and specified string. `string === value`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.NOT_EQUALS]: 'Ensures the header value does not exactly match the specified string. `string !== value`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.CONTAINS]: 'Checks if the header value contains the specified string as a substring. `string.includes()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.NOT_CONTAINS]: 'Verifies the header value does not contain the specified string as a substring. `!string.includes()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.STARTS_WITH]: 'Validates that the header value begins with the specified string. `string.startsWith()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.NOT_STARTS_WITH]: 'Validates that the header value does not begin with the specified string. `!string.startsWith()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.ENDS_WITH]: 'Confirms the header value ends with the specified string. `string.endsWith()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.NOT_ENDS_WITH]: 'Confirms the header value does not end with the specified string. `!string.endsWith()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.REGEX]: 'Evaluates if the header value matches the specified regular expression pattern. `regex.test()`',
    [exports.RESPONSE_RULES_CONDITION_OPERATORS.NOT_REGEX]: 'Evaluates if the header value does not match the specified regular expression pattern. `!regex.test()`',
};
