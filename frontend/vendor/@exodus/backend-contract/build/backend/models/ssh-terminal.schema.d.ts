import { z } from 'zod';
export declare const SSH_TERMINAL_WS_PATH = "/api/node-ssh/ws";
export declare const SSH_TERMINAL_WS_PROTOCOL = "ex";
export declare const SshOpenMessageSchema: z.ZodObject<{
    t: z.ZodLiteral<"open">;
    host: z.ZodString;
    port: z.ZodNumber;
    username: z.ZodString;
    cols: z.ZodNumber;
    rows: z.ZodNumber;
}, z.core.$strip>;
export declare const SshResizeMessageSchema: z.ZodObject<{
    t: z.ZodLiteral<"resize">;
    cols: z.ZodNumber;
    rows: z.ZodNumber;
}, z.core.$strip>;
export declare const SshIdentitiesReplySchema: z.ZodObject<{
    t: z.ZodLiteral<"identities">;
    id: z.ZodNumber;
    keys: z.ZodArray<z.ZodString>;
}, z.core.$strip>;
export declare const SshSignReplySchema: z.ZodObject<{
    t: z.ZodLiteral<"sign">;
    id: z.ZodNumber;
    signature: z.ZodBase64;
}, z.core.$strip>;
export declare const SshHostKeyReplySchema: z.ZodObject<{
    t: z.ZodLiteral<"hostkey">;
    id: z.ZodNumber;
    accept: z.ZodBoolean;
}, z.core.$strip>;
export declare const SshClientErrorSchema: z.ZodObject<{
    t: z.ZodLiteral<"error">;
    id: z.ZodNumber;
    message: z.ZodString;
}, z.core.$strip>;
/** Browser -> panel. */
export declare const SshClientMessageSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    t: z.ZodLiteral<"open">;
    host: z.ZodString;
    port: z.ZodNumber;
    username: z.ZodString;
    cols: z.ZodNumber;
    rows: z.ZodNumber;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"resize">;
    cols: z.ZodNumber;
    rows: z.ZodNumber;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"identities">;
    id: z.ZodNumber;
    keys: z.ZodArray<z.ZodString>;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"sign">;
    id: z.ZodNumber;
    signature: z.ZodBase64;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"hostkey">;
    id: z.ZodNumber;
    accept: z.ZodBoolean;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"error">;
    id: z.ZodNumber;
    message: z.ZodString;
}, z.core.$strip>], "t">;
/** Panel -> browser. */
export declare const SshServerMessageSchema: z.ZodDiscriminatedUnion<[z.ZodObject<{
    t: z.ZodLiteral<"agent-identities">;
    id: z.ZodNumber;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"agent-sign">;
    id: z.ZodNumber;
    publicKey: z.ZodBase64;
    data: z.ZodBase64;
    hash: z.ZodNullable<z.ZodString>;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"hostkey">;
    id: z.ZodNumber;
    algo: z.ZodString;
    fingerprint: z.ZodString;
    key: z.ZodBase64;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"ready">;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"exit">;
    code: z.ZodNullable<z.ZodNumber>;
    signal: z.ZodNullable<z.ZodString>;
}, z.core.$strip>, z.ZodObject<{
    t: z.ZodLiteral<"error">;
    message: z.ZodString;
}, z.core.$strip>], "t">;
export type TSshOpenMessage = z.infer<typeof SshOpenMessageSchema>;
export type TSshClientMessage = z.infer<typeof SshClientMessageSchema>;
export type TSshServerMessage = z.infer<typeof SshServerMessageSchema>;
//# sourceMappingURL=ssh-terminal.schema.d.ts.map