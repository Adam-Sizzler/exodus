import React, { useEffect } from 'react';

function ConfirmActionModal({
  open,
  title = 'Удалить',
  message = 'Вы уверены, что хотите выполнить это действие? Это действие нельзя отменить.',
  confirmLabel = 'Удалить',
  cancelLabel = 'Закрыть',
  onConfirm,
  onClose,
  loading = false,
}) {
  useEffect(() => {
    if (!open) {
      return undefined;
    }

    const onKeyDown = (event) => {
      if (event.key === 'Escape') {
        onClose?.();
      }
    };

    document.addEventListener('keydown', onKeyDown);
    return () => document.removeEventListener('keydown', onKeyDown);
  }, [open, onClose]);

  if (!open) {
    return null;
  }

  return (
    <div className="modal-overlay confirm-action-overlay" onClick={onClose}>
      <section
        className="modal confirm-action-modal"
        onClick={(event) => event.stopPropagation()}
        role="dialog"
        tabIndex={-1}
        aria-modal="true"
        aria-labelledby="confirm-action-modal-title"
        aria-describedby="confirm-action-modal-body"
      >
        <header className="confirm-action-header">
          <h2 id="confirm-action-modal-title">{title}</h2>
          <button
            type="button"
            className="confirm-action-close"
            onClick={onClose}
            aria-label="Закрыть"
          >
            <svg viewBox="0 0 15 15" fill="none" aria-hidden="true">
              <path d="M11.782 4.032a.813.813 0 0 0 0-1.149a.813.813 0 0 0-1.149 0L7.5 6.015L4.367 2.883a.813.813 0 0 0-1.149 0a.813.813 0 0 0 0 1.149L6.351 7.165l-3.133 3.133a.813.813 0 0 0 0 1.149a.813.813 0 0 0 1.149 0L7.5 8.314l3.133 3.133a.813.813 0 0 0 1.149 0a.813.813 0 0 0 0-1.149L8.649 7.165z" fill="currentColor" />
            </svg>
          </button>
        </header>
        <div className="confirm-action-body" id="confirm-action-modal-body">
          <p className="confirm-action-message">{message}</p>
          <div className="confirm-action-footer">
            <button
              type="button"
              className="confirm-action-btn confirm-action-btn-cancel"
              onClick={onClose}
              disabled={loading}
            >
              {cancelLabel}
            </button>
            <button
              type="button"
              className="confirm-action-btn confirm-action-btn-danger"
              onClick={onConfirm}
              disabled={loading}
            >
              {loading ? 'Удаление...' : confirmLabel}
            </button>
          </div>
        </div>
      </section>
    </div>
  );
}

export default ConfirmActionModal;
