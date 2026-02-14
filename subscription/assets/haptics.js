(function () {
    const TARGET_SELECTORS = ['button', 'a', '.option', '.select-trigger', '.tab-btn', '.icon-btn', '.btn'];

    function triggerHaptic(type = 'light') {
        if (typeof navigator.vibrate === 'function') {
            const patterns = {
                'light': [3, 1, 3],
                'medium': [8],
                'heavy': [15]
            };
            
            navigator.vibrate(patterns[type] || [10]);
            return;
        }
    }

    function handleEvent(e) {
        const target = e.target.closest(TARGET_SELECTORS.join(','));
        if (target) {
            const hapticType = target.dataset.haptic || 'light';
            triggerHaptic(hapticType);
        }
    }

    // document.addEventListener('touchstart', handleEvent, { passive: true });
    document.addEventListener('click', handleEvent, { passive: true });

    window.triggerHaptic = triggerHaptic;
})();