import { Tabs } from '@mantine/core'

import classes from './tabs.module.css'

export default {
    Tabs: Tabs.extend({
        classNames: {
            list: classes.list,
            tab: classes.tab,
            tabLabel: classes.tabLabel
        },
        defaultProps: {
            variant: 'unstyled'
        }
    })
}
