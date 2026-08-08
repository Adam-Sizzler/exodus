"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AddManyUsersToInternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var AddManyUsersToInternalSquadCommand;
(function (AddManyUsersToInternalSquadCommand) {
    AddManyUsersToInternalSquadCommand.url = api_1.REST_API.INTERNAL_SQUADS.BULK_ACTIONS.ADD_MANY_USERS;
    AddManyUsersToInternalSquadCommand.TSQ_url = AddManyUsersToInternalSquadCommand.url(':uuid');
    AddManyUsersToInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.ADD_MANY_USERS(':uuid'), 'post', 'Add many users to internal squad', { scope: 'add-many-users', kind: 'write' });
    AddManyUsersToInternalSquadCommand.RequestParamSchema = zod_1.z.object({
        uuid: zod_1.z.uuid(),
    });
    AddManyUsersToInternalSquadCommand.RequestBodySchema = zod_1.z.object({
        userIds: zod_1.z.array(zod_1.z.number()).min(1).max(1000),
    });

})(AddManyUsersToInternalSquadCommand || (exports.AddManyUsersToInternalSquadCommand = AddManyUsersToInternalSquadCommand = {}));
