"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.EvaluateVaultCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var EvaluateVaultCommand;
(function (EvaluateVaultCommand) {
    EvaluateVaultCommand.url = api_1.REST_API.NODE_SSH.EVALUATE_VAULT;
    EvaluateVaultCommand.TSQ_url = EvaluateVaultCommand.url;
    EvaluateVaultCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.NODE_SSH_ROUTES.EVALUATE_VAULT, 'post', 'Oblivious evaluation step for unlocking the SSH key vault', { scope: 'node-ssh', kind: 'write' });
    EvaluateVaultCommand.RequestBodySchema = zod_1.z.object({
        blinded: zod_1.z.base64().max(128),
    });
    EvaluateVaultCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            evaluated: zod_1.z.base64(),
        }),
    });
})(EvaluateVaultCommand || (exports.EvaluateVaultCommand = EvaluateVaultCommand = {}));
