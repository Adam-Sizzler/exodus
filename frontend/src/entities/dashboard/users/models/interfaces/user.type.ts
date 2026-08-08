import { GetUsersCommand } from '@exodus/backend-contract'

export type User = GetUsersCommand.Response['response']['users'][number]
