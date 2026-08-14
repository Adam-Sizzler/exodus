"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.DeleteSnippetCommand = void 0;
const zod_1 = require("zod");
const api_1 = require("../../api");
const constants_1 = require("../../constants");
var DeleteSnippetCommand;
(function (DeleteSnippetCommand) {
    DeleteSnippetCommand.url = api_1.REST_API.SNIPPETS.DELETE;
    DeleteSnippetCommand.TSQ_url = DeleteSnippetCommand.url;
    DeleteSnippetCommand.endpointDetails = (0, constants_1.getEndpointDetails)(api_1.SNIPPETS_ROUTES.DELETE, 'delete', 'Delete snippet', { scope: 'delete', kind: 'write' });
    DeleteSnippetCommand.RequestBodySchema = zod_1.z.object({
        name: zod_1.z
            .string()
            .min(2, 'Name must be at least 2 characters')
            .max(255, 'Name must be less than 255 characters')
            .regex(/^[A-Za-z0-9_\s-]+$/, 'Name can only contain letters, numbers, underscores, dashes and spaces'),
    });
})(DeleteSnippetCommand || (exports.DeleteSnippetCommand = DeleteSnippetCommand = {}));
