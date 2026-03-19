"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.AddUsersToExternalSquadCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var AddUsersToExternalSquadCommand;
(function (AddUsersToExternalSquadCommand) {
    AddUsersToExternalSquadCommand.url = api_1.REST_API.EXTERNAL_SQUADS.BULK_ACTIONS.ADD_USERS;
    AddUsersToExternalSquadCommand.TSQ_url = AddUsersToExternalSquadCommand.url(':uuid');
    AddUsersToExternalSquadCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.EXTERNAL_SQUADS_ROUTES.BULK_ACTIONS.ADD_USERS(':uuid'), 'post', 'Add all users to external squad');
    AddUsersToExternalSquadCommand.RequestSchema = zod_1.z.object({
        uuid: zod_1.z.string().uuid(),
    });
    AddUsersToExternalSquadCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            eventSent: zod_1.z.boolean(),
        }),
    });
})(AddUsersToExternalSquadCommand || (exports.AddUsersToExternalSquadCommand = AddUsersToExternalSquadCommand = {}));
