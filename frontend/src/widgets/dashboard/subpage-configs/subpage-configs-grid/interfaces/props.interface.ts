import { GetSubpageConfigsCommand } from '@exodus/backend-contract'

export interface IProps {
    configs: GetSubpageConfigsCommand.Response['response']['configs']
}
