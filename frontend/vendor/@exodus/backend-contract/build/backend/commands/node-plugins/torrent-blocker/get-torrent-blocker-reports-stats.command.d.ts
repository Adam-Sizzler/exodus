import { z } from 'zod';
export declare namespace GetTorrentBlockerReportsStatsCommand {
    const url: "/api/node-plugins/torrent-blocker/stats";
    const TSQ_url: "/api/node-plugins/torrent-blocker/stats";
    const endpointDetails: import("../../../constants").EndpointDetails;
    const ResponseSchema: z.ZodObject<{
        response: z.ZodObject<{
            stats: z.ZodObject<{
                distinctNodes: z.ZodNumber;
                distinctUsers: z.ZodNumber;
                totalReports: z.ZodNumber;
                reportsLast24Hours: z.ZodNumber;
            }, z.core.$strip>;
            topUsers: z.ZodArray<z.ZodObject<{
                userId: z.ZodNumber;
                color: z.ZodString;
                username: z.ZodString;
                total: z.ZodNumber;
            }, z.core.$strip>>;
            topNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodUUID;
                countryCode: z.ZodString;
                color: z.ZodString;
                name: z.ZodString;
                total: z.ZodNumber;
            }, z.core.$strip>>;
        }, z.core.$strip>;
    }, z.core.$strip>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-torrent-blocker-reports-stats.command.d.ts.map