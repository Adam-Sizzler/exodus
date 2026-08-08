"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteConfigProfileCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteConfigProfileCommand;
(function (DeleteConfigProfileCommand) {
    DeleteConfigProfileCommand.url = api_1.REST_API.CONFIG_PROFILES.DELETE;
    DeleteConfigProfileCommand.TSQ_url = DeleteConfigProfileCommand.url(':uuid');
    DeleteConfigProfileCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.CONFIG_PROFILES_ROUTES.DELETE(':uuid'), 'delete', 'Delete config profile', { scope: 'delete', kind: 'write' });
    DeleteConfigProfileCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid().describe('UUID of the config profile'),
    });

})(DeleteConfigProfileCommand || (exports.DeleteConfigProfileCommand = DeleteConfigProfileCommand = {}));
