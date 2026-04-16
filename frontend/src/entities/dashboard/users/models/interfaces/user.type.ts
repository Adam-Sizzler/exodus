import { GetAllUsersCommand } from '@exodus/backend-contract'

export type User = GetAllUsersCommand.Response['response']['users'][number]
