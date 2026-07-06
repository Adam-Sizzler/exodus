import { forwardRef } from 'react'

import classes from './virtualized-dnd-grid.module.css'

export const VirtualizedGridComponents = {
    List: forwardRef<HTMLDivElement, React.HTMLProps<HTMLDivElement>>(
        ({ children, ...props }, ref) => (
            <div {...props} className={classes.list} ref={ref}>
                {children}
            </div>
        )
    ),
    Item: ({ children, ...props }: React.HTMLProps<HTMLDivElement>) => (
        <div {...props} className={classes.item}>
            {children}
        </div>
    )
}
