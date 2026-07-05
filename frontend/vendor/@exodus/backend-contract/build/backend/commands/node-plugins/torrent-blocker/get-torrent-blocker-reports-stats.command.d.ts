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
            }, "strip", z.ZodTypeAny, {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            }, {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            }>;
            topUsers: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                color: z.ZodString;
                username: z.ZodString;
                total: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }, {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }>, "many">;
            topNodes: z.ZodArray<z.ZodObject<{
                uuid: z.ZodString;
                countryCode: z.ZodString;
                color: z.ZodString;
                name: z.ZodString;
                total: z.ZodNumber;
            }, "strip", z.ZodTypeAny, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }, {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }>, "many">;
        }, "strip", z.ZodTypeAny, {
            stats: {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            };
            topUsers: {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
        }, {
            stats: {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            };
            topUsers: {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
        }>;
    }, "strip", z.ZodTypeAny, {
        response: {
            stats: {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            };
            topUsers: {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
        };
    }, {
        response: {
            stats: {
                distinctNodes: number;
                distinctUsers: number;
                totalReports: number;
                reportsLast24Hours: number;
            };
            topUsers: {
                uuid: string;
                username: string;
                total: number;
                color: string;
            }[];
            topNodes: {
                uuid: string;
                total: number;
                countryCode: string;
                color: string;
                name: string;
            }[];
        };
    }>;
    type Response = z.infer<typeof ResponseSchema>;
}
//# sourceMappingURL=get-torrent-blocker-reports-stats.command.d.ts.map