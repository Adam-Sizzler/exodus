"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetOttCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetOttCommand;
(function (GetOttCommand) {
    GetOttCommand.url = api_1.REST_API.API_TOKENS.OTT;
    GetOttCommand.TSQ_url = GetOttCommand.url;
    GetOttCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.API_TOKENS_ROUTES.OTT, 'post', 'Get short-lived token for accessing backend tools (Swagger, Scalar, Bull Board)', { scope: 'ott', kind: 'write' });
    GetOttCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            ott: zod_1.z.string(),
        }),
    });
})(GetOttCommand || (exports.GetOttCommand = GetOttCommand = {}));
