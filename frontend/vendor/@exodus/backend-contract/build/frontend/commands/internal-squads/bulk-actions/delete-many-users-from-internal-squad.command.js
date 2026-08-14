"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteManyUsersFromInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var DeleteManyUsersFromInternalSquadCommand;
(function (DeleteManyUsersFromInternalSquadCommand) {
    DeleteManyUsersFromInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.BULK_ACTIONS.REMOVE_MANY_USERS;
    DeleteManyUsersFromInternalSquadCommand.TSQ_url = DeleteManyUsersFromInternalSquadCommand.url(':uuid');
    DeleteManyUsersFromInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_MANY_USERS(':uuid'), 'delete', 'Delete many users from internal squad', { scope: 'remove-many-users', kind: 'write' });
    DeleteManyUsersFromInternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    DeleteManyUsersFromInternalSquadCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(1000),
    });
})(DeleteManyUsersFromInternalSquadCommand || (exports.DeleteManyUsersFromInternalSquadCommand = DeleteManyUsersFromInternalSquadCommand = {}));
