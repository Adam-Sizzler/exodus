"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.GetHostsTagsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var GetHostsTagsCommand;
(function (GetHostsTagsCommand) {
    GetHostsTagsCommand.url = api_1.REST_API.HOSTS.TAGS.GET;
    GetHostsTagsCommand.TSQ_url = GetHostsTagsCommand.url;
    GetHostsTagsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.TAGS.GET, 'get', 'Get tags of hosts', { scope: 'list-tags', kind: 'read' });
    GetHostsTagsCommand.ResponseSchema = zod_1.z.object({
        response: zod_1.z.object({
            tags: zod_1.z.array(zod_1.z.string()),
        }),
    });
})(GetHostsTagsCommand || (exports.GetHostsTagsCommand = GetHostsTagsCommand = {}));
