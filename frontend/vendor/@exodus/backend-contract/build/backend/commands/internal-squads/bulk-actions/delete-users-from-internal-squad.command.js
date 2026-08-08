"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteUsersFromInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var DeleteUsersFromInternalSquadCommand;
(function (DeleteUsersFromInternalSquadCommand) {
    DeleteUsersFromInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.BULK_ACTIONS.REMOVE_USERS;
    DeleteUsersFromInternalSquadCommand.TSQ_url = DeleteUsersFromInternalSquadCommand.url(':uuid');
    DeleteUsersFromInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_USERS(':uuid'), 'delete', 'Delete users from internal squad', { scope: 'remove-users', kind: 'write' });
    DeleteUsersFromInternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });

})(DeleteUsersFromInternalSquadCommand || (exports.DeleteUsersFromInternalSquadCommand = DeleteUsersFromInternalSquadCommand = {}));
