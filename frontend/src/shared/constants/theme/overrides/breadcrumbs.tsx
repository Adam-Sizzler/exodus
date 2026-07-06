// oxlint-disable
import { Breadcrumbs, px } from '@mantine/core'
import { GoDotFill as BreadcrumbsSeparator } from 'react-icons/go'

export default {
    Breadcrumbs: Breadcrumbs.extend({
        defaultProps: {
            separator: (
                <BreadcrumbsSeparator
                    color="var(--mantine-color-dimmed)"
                    opacity={0.4}
                    size={px('0.5rem')}
                />
            )
        }
    })
}
