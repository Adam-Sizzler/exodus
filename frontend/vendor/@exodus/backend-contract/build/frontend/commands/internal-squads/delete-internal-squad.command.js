"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteInternalSquadCommand;
(function (DeleteInternalSquadCommand) {
    DeleteInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.DELETE;
    DeleteInternalSquadCommand.TSQ_url = DeleteInternalSquadCommand.url(':uuid');
    DeleteInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.DELETE(':uuid'), 'delete', 'Delete internal squad', { scope: 'delete', kind: 'write' });
    DeleteInternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteInternalSquadCommand || (exports.DeleteInternalSquadCommand = DeleteInternalSquadCommand = {}));
