import { GetSnippetsCommand } from '@exodus/backend-contract'

export type TSnippet = GetSnippetsCommand.Response['response']['snippets'][number]
