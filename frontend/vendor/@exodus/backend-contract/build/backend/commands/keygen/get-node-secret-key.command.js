"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetNodeSecretKeyCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var GetNodeSecretKeyCommand;
(function (GetNodeSecretKeyCommand) {
    GetNodeSecretKeyCommand.url = api_1.REST_API.KEYGEN.GET;
    GetNodeSecretKeyCommand.TSQ_url = GetNodeSecretKeyCommand.url;
    GetNodeSecretKeyCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.KEYGEN_ROUTES.GET, 'get', 'Get SECRET_KEY for Exodus Node', { scope: 'get', kind: 'read' });
    GetNodeSecretKeyCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            secretKey: zod_1.z.string(),
            grpcToken: zod_1.z.string(),
        }),
    });
})(GetNodeSecretKeyCommand || (exports.GetNodeSecretKeyCommand = GetNodeSecretKeyCommand = {}));
