import z from 'zod';

export const ResponseRuleSettingsSchema = z
    .object({
        disableSubscriptionAccessByPath: z.boolean().optional().meta({
            title: 'Disable Subscription Access by Path',
            markdownDescription:
                "Usually, a user's subscription may also be available via additional paths such as **/json**, **/stash**, or **/mihomo**. If this flag is set to **true**, access via these additional paths will be disabled.",
        }),
    })
    .optional()
    .meta({
        title: 'Response Rule Settings',
        markdownDescription: 'Settings for the **response rules** config. Optional.',
    });
