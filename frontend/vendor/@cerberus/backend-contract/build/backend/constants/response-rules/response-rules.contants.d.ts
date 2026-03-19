export declare const RESPONSE_RULES_CONFIG_VERSION: {
    readonly 1: "1";
};
export declare const RESPONSE_RULES_OPERATORS: {
    readonly AND: "AND";
    readonly OR: "OR";
};
export declare const RESPONSE_RULE_CONDITION_TYPES: {
    readonly HEADER: "HEADER";
};
export declare const RESPONSE_RULES_CONDITION_OPERATORS: {
    readonly EQUALS: "EQUALS";
    readonly NOT_EQUALS: "NOT_EQUALS";
    readonly CONTAINS: "CONTAINS";
    readonly NOT_CONTAINS: "NOT_CONTAINS";
    readonly STARTS_WITH: "STARTS_WITH";
    readonly NOT_STARTS_WITH: "NOT_STARTS_WITH";
    readonly ENDS_WITH: "ENDS_WITH";
    readonly NOT_ENDS_WITH: "NOT_ENDS_WITH";
    readonly REGEX: "REGEX";
    readonly NOT_REGEX: "NOT_REGEX";
};
export declare const RESPONSE_RULES_RESPONSE_TYPES: {
    readonly BROWSER: "BROWSER";
    readonly BLOCK: "BLOCK";
    readonly STATUS_CODE_404: "STATUS_CODE_404";
    readonly STATUS_CODE_451: "STATUS_CODE_451";
    readonly SOCKET_DROP: "SOCKET_DROP";
    readonly XRAY_JSON: "XRAY_JSON";
    readonly XRAY_BASE64: "XRAY_BASE64";
    readonly MIHOMO: "MIHOMO";
    readonly STASH: "STASH";
    readonly CLASH: "CLASH";
    readonly SINGBOX: "SINGBOX";
};
export type TResponseRulesResponseType = [keyof typeof RESPONSE_RULES_RESPONSE_TYPES][number];
export declare const RESPONSE_RULES_RESPONSE_TYPES_VALUES: ("STASH" | "SINGBOX" | "MIHOMO" | "XRAY_JSON" | "CLASH" | "XRAY_BASE64" | "BROWSER" | "BLOCK" | "STATUS_CODE_404" | "STATUS_CODE_451" | "SOCKET_DROP")[];
export type TResponseRulesResponseTypeKeys = (typeof RESPONSE_RULES_RESPONSE_TYPES)[keyof typeof RESPONSE_RULES_RESPONSE_TYPES];
export type TResponseRulesConditionOperator = [
    keyof typeof RESPONSE_RULES_CONDITION_OPERATORS
][number];
export declare const RESPONSE_RULES_CONDITION_OPERATORS_VALUES: ("EQUALS" | "NOT_EQUALS" | "CONTAINS" | "NOT_CONTAINS" | "STARTS_WITH" | "NOT_STARTS_WITH" | "ENDS_WITH" | "NOT_ENDS_WITH" | "REGEX" | "NOT_REGEX")[];
export type TResponseRulesConditionOperatorKeys = (typeof RESPONSE_RULES_CONDITION_OPERATORS)[keyof typeof RESPONSE_RULES_CONDITION_OPERATORS];
export declare const RESPONSE_RULES_RESPONSE_TYPES_DESCRIPTION: {
    readonly XRAY_JSON: "Return **subscription** in XRAY-JSON format. (Using `Xray Json` template)";
    readonly XRAY_BASE64: "Return **subscription** in BASE64 encoded string. Compatible with most client application with Xray core.";
    readonly MIHOMO: "Return **subscription** in Mihomo format. (Using `Mihomo` template)";
    readonly STASH: "Return **subscription** in Stash format. (Using `Stash` template)";
    readonly CLASH: "Return **subscription** in Clash format. (Using `Clash` template) Useful for client application that use Legacy Clash core.";
    readonly SINGBOX: "Return **subscription** in Singbox format. (Using `Singbox` template) Format which is used by Singbox client application.";
    readonly BROWSER: "Return **subscription** as browser format. The same as on `/info` route.";
    readonly BLOCK: "**Drop** request and return `403` status code.";
    readonly STATUS_CODE_404: "**Drop** request and return `404` status code.";
    readonly STATUS_CODE_451: "**Drop** request and return `451` status code.";
    readonly SOCKET_DROP: "**Drop** the socket connection.";
};
export declare const RESPONSE_RULES_CONDITION_OPERATORS_DESCRIPTION: {
    readonly EQUALS: "Performs an exact, comparison between the header value and specified string. `string === value`";
    readonly NOT_EQUALS: "Ensures the header value does not exactly match the specified string. `string !== value`";
    readonly CONTAINS: "Checks if the header value contains the specified string as a substring. `string.includes()`";
    readonly NOT_CONTAINS: "Verifies the header value does not contain the specified string as a substring. `!string.includes()`";
    readonly STARTS_WITH: "Validates that the header value begins with the specified string. `string.startsWith()`";
    readonly NOT_STARTS_WITH: "Validates that the header value does not begin with the specified string. `!string.startsWith()`";
    readonly ENDS_WITH: "Confirms the header value ends with the specified string. `string.endsWith()`";
    readonly NOT_ENDS_WITH: "Confirms the header value does not end with the specified string. `!string.endsWith()`";
    readonly REGEX: "Evaluates if the header value matches the specified regular expression pattern. `regex.test()`";
    readonly NOT_REGEX: "Evaluates if the header value does not match the specified regular expression pattern. `!regex.test()`";
};
//# sourceMappingURL=response-rules.contants.d.ts.map