import { z } from 'zod';

export const ExternalSquadResponseHeadersAddSchema = z.record(z.string(), z.string());
export const ExternalSquadResponseHeadersRemoveSchema = z.array(z.string());
