"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetExodusHealthCommand = void 0;
const zod_1 = require("zod");
const constants_1 = require("../../constants");
const api_1 = require("../../api");
var GetExodusHealthCommand;
(function (GetExodusHealthCommand) {
    GetExodusHealthCommand.url = api_1.REST_API.SYSTEM.HEALTH;
    GetExodusHealthCommand.TSQ_url = GetExodusHealthCommand.url;
    GetExodusHealthCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SYSTEM_ROUTES.HEALTH, 'get', 'Get Exodus Health');
    GetExodusHealthCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            pm2Stats: zod_1.z.array(zod_1.z.object({
                name: zod_1.z.string(),
                memory: zod_1.z.string(),
                cpu: zod_1.z.string(),
            })),
        }),
    });
})(GetExodusHealthCommand || (exports.GetExodusHealthCommand = GetExodusHealthCommand = {}));
