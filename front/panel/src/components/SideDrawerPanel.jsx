import { useEffect } from 'react';

function SideDrawerPanel({
  open,
  onClose,
  title,
  subtitle,
  icon,
  children,
  footer,
  headerActions,
  width = '50rem',
  className = '',
}) {
  useEffect(() => {
    if (!open) return undefined;

    const previousOverflow = document.body.style.overflow;
    document.body.style.overflow = 'hidden';

    const onKeyDown = (event) => {
      if (event.key === 'Escape') {
        onClose?.();
      }
    };

    document.addEventListener('keydown', onKeyDown);

    return () => {
      document.body.style.overflow = previousOverflow;
      document.removeEventListener('keydown', onKeyDown);
    };
  }, [open, onClose]);

  return (
    <div
      className={`side-drawer-root ${open ? 'open' : ''}`}
      aria-hidden={!open}
      style={{ '--side-drawer-width': width }}
    >
      <button
        type="button"
        className="side-drawer-overlay"
        aria-label="Close panel"
        onClick={onClose}
      />

      <aside
        className={`side-drawer-panel ${className}`.trim()}
        role="dialog"
        aria-modal="true"
        aria-label={title || 'Side panel'}
      >
        <header className="side-drawer-header">
          <div className="side-drawer-title-wrap">
            {icon ? <div className="side-drawer-title-icon">{icon}</div> : null}
            <div className="side-drawer-title-stack">
              {title ? <h2 className="side-drawer-title">{title}</h2> : null}
              {subtitle ? <p className="side-drawer-subtitle">{subtitle}</p> : null}
            </div>
          </div>
          <div className="side-drawer-header-actions">
            {headerActions}
            <button
              type="button"
              className="side-drawer-close"
              onClick={onClose}
              aria-label="Close panel"
            >
              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                <line x1="18" y1="6" x2="6" y2="18" />
                <line x1="6" y1="6" x2="18" y2="18" />
              </svg>
            </button>
          </div>
        </header>

        <div className="side-drawer-body">{children}</div>

        {footer ? <footer className="side-drawer-footer">{footer}</footer> : null}
      </aside>
    </div>
  );
}

export default SideDrawerPanel;
