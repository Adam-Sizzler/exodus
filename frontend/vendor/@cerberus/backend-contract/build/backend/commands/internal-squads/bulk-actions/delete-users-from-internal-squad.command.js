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
    DeleteUsersFromInternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.INTERNAL_SQUADS_ROUTES.BULK_ACTIONS.REMOVE_USERS(':uuid'), 'delete', 'Delete users from internal squad');
    DeleteUsersFromInternalSquadCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    DeleteUsersFromInternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(DeleteUsersFromInternalSquadCommand || (exports.DeleteUsersFromInternalSquadCommand = DeleteUsersFromInternalSquadCommand = {}));
