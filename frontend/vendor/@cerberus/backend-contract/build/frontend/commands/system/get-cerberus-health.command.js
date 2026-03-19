"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetCerberusHealthCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetCerberusHealthCommand;
(function (GetCerberusHealthCommand) {
    GetCerberusHealthCommand.url = api_1.REST_API.SYSTEM.HEALTH;
    GetCerberusHealthCommand.TSQ_url = GetCerberusHealthCommand.url;
    GetCerberusHealthCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.HEALTH, 'get', 'Get Cerberus Health');
    GetCerberusHealthCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            pm2Stats: zod_1.z.array(zod_1.z.object({
                name: zod_1.z.string(),
                memory: zod_1.z.string(),
                cpu: zod_1.z.string(),
            })),
        }),
    });
})(GetCerberusHealthCommand || (exports.GetCerberusHealthCommand = GetCerberusHealthCommand = {}));
