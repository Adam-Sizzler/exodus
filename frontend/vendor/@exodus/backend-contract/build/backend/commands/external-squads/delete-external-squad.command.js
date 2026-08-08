"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteExternalSquadCommand;
(function (DeleteExternalSquadCommand) {
    DeleteExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.DELETE;
    DeleteExternalSquadCommand.TSQ_url = DeleteExternalSquadCommand.url(':uuid');
    DeleteExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.DELETE(':uuid'), 'delete', 'Delete external squad', { scope: 'delete', kind: 'write' });
    DeleteExternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid().describe('UUID of the external squad'),
    });

})(DeleteExternalSquadCommand || (exports.DeleteExternalSquadCommand = DeleteExternalSquadCommand = {}));
