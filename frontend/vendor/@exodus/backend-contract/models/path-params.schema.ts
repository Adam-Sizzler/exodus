import { z } from 'zod';

export const numberParamSchema = z.coerce.number().positive();
