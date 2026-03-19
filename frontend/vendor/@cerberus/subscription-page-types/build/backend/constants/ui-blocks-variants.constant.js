"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.isSubscriptionInfoBlockVariant = exports.isInstallationGuideBlockVariant = exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS_VALUES = exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS_VALUES = exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS = exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS = void 0;
exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS = {
    COLLAPSED: 'collapsed',
    EXPANDED: 'expanded',
    CARDS: 'cards',
    HIDDEN: 'hidden',
};
exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS = {
    CARDS: 'cards',
    ACCORDION: 'accordion',
    MINIMAL: 'minimal',
    TIMELINE: 'timeline',
};
exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS_VALUES = Object.values(exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS);
exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS_VALUES = Object.values(exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS);
const isInstallationGuideBlockVariant = (value) => {
    return exports.INSTALLATION_GUIDE_BLOCKS_VARIANTS_VALUES.includes(value);
};
exports.isInstallationGuideBlockVariant = isInstallationGuideBlockVariant;
const isSubscriptionInfoBlockVariant = (value) => {
    return exports.SUBSCRIPTION_INFO_BLOCK_VARIANTS_VALUES.includes(value);
};
exports.isSubscriptionInfoBlockVariant = isSubscriptionInfoBlockVariant;
