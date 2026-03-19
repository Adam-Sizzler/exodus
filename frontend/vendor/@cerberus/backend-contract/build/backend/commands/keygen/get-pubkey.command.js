"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetPubKeyCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetPubKeyCommand;
(function (GetPubKeyCommand) {
    GetPubKeyCommand.url = api_1.REST_API.KEYGEN.GET;
    GetPubKeyCommand.TSQ_url = GetPubKeyCommand.url;
    GetPubKeyCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.KEYGEN_ROUTES.GET, 'get', 'Get SSL_CERT for Cerberus Node');
    GetPubKeyCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            pubKey: zod_1.z.string(),
        }),
    });
})(GetPubKeyCommand || (exports.GetPubKeyCommand = GetPubKeyCommand = {}));
