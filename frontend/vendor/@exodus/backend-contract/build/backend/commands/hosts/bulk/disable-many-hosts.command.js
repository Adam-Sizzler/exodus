"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BulkDisableHostsCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../../api");
const constants_1 = require("../../../constants");
var BulkDisableHostsCommand;
(function (BulkDisableHostsCommand) {
    BulkDisableHostsCommand.url = api_1.REST_API.HOSTS.BULK.DISABLE_HOSTS;
    BulkDisableHostsCommand.TSQ_url = BulkDisableHostsCommand.url;
    BulkDisableHostsCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.HOSTS_ROUTES.BULK.DISABLE_HOSTS, 'post', 'Disable hosts by UUIDs', { scope: 'bulk-disable', kind: 'write' });
    BulkDisableHostsCommand.RequestBodySchema = zod_1.z.object({
        uuids: zod_1.z.array(zod_1.z.uuid()),
    });
})(BulkDisableHostsCommand || (exports.BulkDisableHostsCommand = BulkDisableHostsCommand = {}));
