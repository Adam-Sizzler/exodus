"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.TestSrrMatcherCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../../constants");
const models_1 = require("../../../models");
const api_1 = require("../../../api");
var TestSrrMatcherCommand;
(function (TestSrrMatcherCommand) {
    TestSrrMatcherCommand.url = api_1.REST_API.SYSTEM.TESTERS.SRR_MATCHER;
    TestSrrMatcherCommand.TSQ_url = TestSrrMatcherCommand.url;
    TestSrrMatcherCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.TESTERS.SRR_MATCHER, 'post', 'Test SRR Matcher');
    TestSrrMatcherCommand.RequestSchema = zod_1.z.object({
        responseRules: models_1.ResponseRulesConfigSchema,
    });
    TestSrrMatcherCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            matched: zod_1.z.boolean(),
            responseType: zod_1.z.nativeEnum(constants_1.RESPONSE_RULES_RESPONSE_TYPES),
            matchedRule: zod_1.z.nullable(models_1.ResponseRuleSchema),
            inputHeaders: zod_1.z.record(zod_1.z.string(), zod_1.z.string()),
            outputHeaders: zod_1.z.record(zod_1.z.string(), zod_1.z.string()),
        }),
    });
})(TestSrrMatcherCommand || (exports.TestSrrMatcherCommand = TestSrrMatcherCommand = {}));
