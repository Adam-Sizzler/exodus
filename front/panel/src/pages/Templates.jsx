import React, { useEffect, useMemo, useState } from 'react';
import { templatesApi } from '../api';
import { validateWithWasm, formatWithWasm } from '../utils/wasmValidator';
import TilesToolbar from '../components/TilesToolbar';
import useAdaptiveEntityGridColumns from '../components/useAdaptiveEntityGridColumns';
import ConfirmActionModal from '../components/ConfirmActionModal';

const TEMPLATE_TYPES = ['XRAY_JSON', 'MIHOMO', 'STASH', 'CLASH', 'SINGBOX'];
const JSON_TEMPLATE_TYPES = new Set(['XRAY_JSON', 'SINGBOX']);
const TEMPLATE_TYPE_LABELS = {
  XRAY_JSON: 'Xray JSON',
  MIHOMO: 'Mihomo',
  STASH: 'Stash',
  CLASH: 'Clash',
  SINGBOX: 'Sing-box',
};

function Templates({ activeTemplateType }) {
  const { gridRef, columns: gridColumns } = useAdaptiveEntityGridColumns();
  const [templates, setTemplates] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [showModal, setShowModal] = useState(false);
  const [editingTemplate, setEditingTemplate] = useState(null);
  const [templateName, setTemplateName] = useState('');
  const [yamlValue, setYamlValue] = useState('');
  const [jsonValue, setJsonValue] = useState('');
  const [saving, setSaving] = useState(false);
  const [dragIndex, setDragIndex] = useState(null);
  const [openMenuUUID, setOpenMenuUUID] = useState(null);
  const [searchQuery, setSearchQuery] = useState('');
  const [deleteTarget, setDeleteTarget] = useState(null);
  const [deleting, setDeleting] = useState(false);

  const currentTemplateType = TEMPLATE_TYPES.includes(activeTemplateType) ? activeTemplateType : TEMPLATE_TYPES[0];
  const isJSONType = JSON_TEMPLATE_TYPES.has(currentTemplateType);
  const currentTypeLabel = TEMPLATE_TYPE_LABELS[currentTemplateType] || currentTemplateType;
  const currentWasmEngine = currentTemplateType === 'XRAY_JSON'
    ? 'xray'
    : ['MIHOMO', 'CLASH', 'STASH'].includes(currentTemplateType)
      ? 'mihomo'
      : null;

  const loadTemplates = async () => {
    try {
      setLoading(true);
      const data = await templatesApi.getAll(currentTemplateType);
      setTemplates(data.templates || []);
      setError(null);
    } catch (err) {
      setError(err.message || 'Failed to load templates');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadTemplates();
    setShowModal(false);
    setEditingTemplate(null);
  }, [currentTemplateType]);

  useEffect(() => {
    const onDocClick = () => setOpenMenuUUID(null);
    if (openMenuUUID) {
      document.addEventListener('click', onDocClick);
    }
    return () => document.removeEventListener('click', onDocClick);
  }, [openMenuUUID]);

  const templateCount = useMemo(() => templates.length, [templates]);

  const formatJSON = (raw) => {
    if (!String(raw || '').trim()) {
      return '{}';
    }
    return JSON.stringify(JSON.parse(raw), null, 2);
  };

  const formatJson = async () => {
    if (currentWasmEngine) {
      try {
        const source = isJSONType ? (jsonValue || '{}') : (yamlValue || '');
        const formatted = await formatWithWasm(currentWasmEngine, source);
        if (isJSONType) {
          setJsonValue(formatted);
        } else {
          setYamlValue(formatted);
        }
        return;
      } catch (err) {
        console.warn(`[wasm:${currentWasmEngine}] format fallback`, err);
      }
    }

    try {
      setJsonValue(formatJSON(jsonValue));
    } catch (err) {
      alert(`Invalid JSON: ${err.message}`);
    }
  };

  const handleCreate = () => {
    setEditingTemplate(null);
    const typeBase = currentTemplateType.toLowerCase();
    const generatedName = templateCount === 0 ? typeBase : `${typeBase}_${templateCount + 1}`;
    setTemplateName(generatedName);
    if (isJSONType) {
      setJsonValue('{}');
      setYamlValue('');
    } else {
      setYamlValue('');
      setJsonValue('');
    }
    setShowModal(true);
  };

  const handleEdit = (template) => {
    setEditingTemplate(template);
    setTemplateName(template.name || '');
    if (isJSONType) {
      try {
        setJsonValue(formatJSON(template.template_json || '{}'));
      } catch (err) {
        setJsonValue(template.template_json || '{}');
      }
      setYamlValue('');
    } else {
      setYamlValue(template.template_yaml || '');
      setJsonValue('');
    }
    setShowModal(true);
  };

  const handleDelete = (template) => {
    setOpenMenuUUID(null);
    setDeleteTarget({ uuid: template.uuid, name: template.name });
  };

  const handleConfirmDelete = async () => {
    if (!deleteTarget || deleting) {
      return;
    }

    try {
      setDeleting(true);
      await templatesApi.delete(deleteTarget.uuid);
      setDeleteTarget(null);
      await loadTemplates();
    } catch (err) {
      alert(`Failed to delete template: ${err.message}`);
    } finally {
      setDeleting(false);
    }
  };

  const handleCopyUUID = async (template) => {
    try {
      await navigator.clipboard.writeText(template.uuid);
      setOpenMenuUUID(null);
    } catch (err) {
      alert(`Failed to copy UUID: ${err.message}`);
    }
  };

  const handleRename = async (template) => {
    const nextName = prompt('Template name', template.name || '');
    if (nextName === null) {
      return;
    }
    const trimmed = nextName.trim();
    if (!trimmed || trimmed === template.name) {
      setOpenMenuUUID(null);
      return;
    }
    try {
      await templatesApi.update(template.uuid, { name: trimmed });
      setOpenMenuUUID(null);
      await loadTemplates();
    } catch (err) {
      alert(`Failed to rename template: ${err.message}`);
    }
  };

  const handleSave = async (event) => {
    event.preventDefault();

    const payload = {
      name: templateName.trim(),
      template_yaml: isJSONType ? '' : yamlValue,
      template_json: isJSONType ? jsonValue : '',
    };

    try {
      setSaving(true);
      if (editingTemplate) {
        await templatesApi.update(editingTemplate.uuid, payload);
      } else {
        await templatesApi.create({
          ...payload,
          template_type: currentTemplateType,
          view_position: templateCount,
        });
      }
      setShowModal(false);
      await loadTemplates();
    } catch (err) {
      alert(`Failed to save template: ${err.message}`);
    } finally {
      setSaving(false);
    }
  };

  const resetEditorValues = () => {
    if (editingTemplate) {
      setTemplateName(editingTemplate.name || '');
      if (isJSONType) {
        try {
          setJsonValue(formatJSON(editingTemplate.template_json || '{}'));
        } catch (err) {
          setJsonValue(editingTemplate.template_json || '{}');
        }
        setYamlValue('');
      } else {
        setYamlValue(editingTemplate.template_yaml || '');
        setJsonValue('');
      }
      return;
    }

    const typeBase = currentTemplateType.toLowerCase();
    const generatedName = templateCount === 0 ? typeBase : `${typeBase}_${templateCount + 1}`;
    setTemplateName(generatedName);
    if (isJSONType) {
      setJsonValue('{}');
      setYamlValue('');
    } else {
      setYamlValue('');
      setJsonValue('');
    }
  };

  const validateEditorContent = async () => {
    if (currentWasmEngine) {
      try {
        const source = isJSONType ? (jsonValue || '{}') : (yamlValue || '');
        const result = await validateWithWasm(currentWasmEngine, source);
        if (result.ok) {
          if (result.formatted) {
            if (isJSONType) {
              setJsonValue(result.formatted);
            } else {
              setYamlValue(result.formatted);
            }
          }
          alert(result.message || 'Config is valid (WASM).');
        } else {
          alert(result.message || 'Config validation failed.');
        }
        return;
      } catch (err) {
        console.warn(`[wasm:${currentWasmEngine}] validate fallback`, err);
      }
    }

    if (isJSONType) {
      try {
        JSON.parse(jsonValue || '{}');
        alert('JSON is valid.');
      } catch (err) {
        alert(`Invalid JSON: ${err.message}`);
      }
      return;
    }

    alert('YAML validation is not available (WASM module is not connected).');
  };

  const prepareDragPreview = (event) => {
    const clone = event.currentTarget.cloneNode(true);
    const { width, height } = event.currentTarget.getBoundingClientRect();
    const toRem = (value) => `${value / 16}rem`;
    const previewHeight = Math.min(height, 352);
    clone.style.position = 'fixed';
    clone.style.top = '-62.5rem';
    clone.style.left = '-62.5rem';
    clone.style.width = toRem(width);
    clone.style.height = toRem(previewHeight);
    clone.style.maxHeight = toRem(previewHeight);
    clone.style.overflow = 'hidden';
    clone.style.background = 'var(--card-background)';
    clone.style.boxShadow = '0 0.875rem 2.25rem rgba(0, 0, 0, 0.28)';
    clone.style.border = '0.0625rem solid var(--border-color)';
    clone.style.opacity = '1';
    clone.style.pointerEvents = 'none';
    document.body.appendChild(clone);
    event.dataTransfer.setDragImage(clone, 20, 20);
    window.setTimeout(() => {
      if (clone.parentNode) {
        clone.parentNode.removeChild(clone);
      }
    }, 120);
  };

  const handleDragStart = (event, index) => {
    if (!canReorder) return;
    setDragIndex(index);
    event.dataTransfer.effectAllowed = 'move';
    prepareDragPreview(event);
  };

  const handleDrop = async (dropIndex) => {
    if (!canReorder || dragIndex === null || dragIndex === dropIndex) {
      setDragIndex(null);
      return;
    }

    const reordered = [...templates];
    const [moved] = reordered.splice(dragIndex, 1);
    reordered.splice(dropIndex, 0, moved);

    setTemplates(reordered);
    setDragIndex(null);

    try {
      await templatesApi.reorder(reordered.map((template) => template.uuid));
      setTemplates((prev) => prev.map((item, index) => ({ ...item, view_position: index })));
    } catch (err) {
      alert(`Failed to reorder templates: ${err.message}`);
      await loadTemplates();
    }
  };

  const normalizedSearch = searchQuery.trim().toLowerCase();
  const canReorder = normalizedSearch === '';
  const filteredTemplates = templates.filter((template) => {
    if (!normalizedSearch) return true;
    const name = template.name?.toLowerCase() || '';
    const uuid = template.uuid?.toLowerCase() || '';
    return name.includes(normalizedSearch) || uuid.includes(normalizedSearch);
  });

  if (loading) {
    return <div className="loading">Loading...</div>;
  }

  return (
    <div className="page">
      {!showModal ? (
        <>
          <TilesToolbar
            title={`Шаблоны • ${currentTypeLabel}`}
            searchValue={searchQuery}
            onSearchChange={setSearchQuery}
            searchPlaceholder="Поиск шаблонов..."
            onReload={loadTemplates}
            onCreate={handleCreate}
            createTitle="Добавить шаблон"
          />

          {error && (
            <div className="alert alert-error">
              {error}
            </div>
          )}

          {templates.length === 0 ? (
            <div className="empty-state">
              <p>No templates for {currentTemplateType}</p>
              <button className="btn btn-primary" onClick={handleCreate}>Create first template</button>
            </div>
          ) : filteredTemplates.length === 0 ? (
            <div className="empty-state">
              <p>По запросу ничего не найдено</p>
              <button type="button" className="btn btn-secondary" onClick={() => setSearchQuery('')}>
                Сбросить поиск
              </button>
            </div>
          ) : (
            <div
              ref={gridRef}
              className="templates-windows-grid entity-windows-grid"
              style={{
                '--entity-grid-columns': gridColumns,
              }}
            >
              {filteredTemplates.map((template, index) => (
                <div key={template.uuid} className="templates-grid-item">
                  <div
                    className={`templates-item-wrapper card-tile templates-window ${dragIndex === index ? 'dragging' : ''} ${openMenuUUID === template.uuid ? 'menu-open' : ''}`}
                    draggable={canReorder}
                    onDragStart={(e) => handleDragStart(e, index)}
                    onDragOver={(e) => canReorder && e.preventDefault()}
                    onDrop={() => handleDrop(index)}
                    onDragEnd={() => setDragIndex(null)}
                  >
                    <button
                      type="button"
                      className={`templates-drag-handle ${canReorder ? '' : 'disabled'}`}
                      aria-label="Reorder item"
                      title={canReorder ? 'Перетащите для сортировки' : 'Сортировка отключена при поиске'}
                    >
                      <svg viewBox="0 0 24 24" fill="currentColor" aria-hidden="true">
                        <path d="M8.5 7A1.5 1.5 0 1 0 8.5 4a1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zM15.5 7a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3zm0 6.5a1.5 1.5 0 1 0 0-3 1.5 1.5 0 0 0 0 3z" />
                      </svg>
                    </button>
                    <div className="templates-top-accent"></div>
                    <div className="templates-glow-effect"></div>
                    <div className="templates-window-main">
                      <div className="templates-window-header">
                        <div className="templates-icon-wrapper">
                          <button type="button" className="templates-icon-button" tabIndex={-1} aria-hidden="true">
                            <svg fill="currentColor" preserveAspectRatio="xMidYMid meet" viewBox="0 0 35 35" xmlns="http://www.w3.org/2000/svg">
                              <path d="M16.6961 15.2606L16.5825 3.49701C16.5718 2.38439 15.025 2.11843 14.6433 3.16356L11.7279 11.1447C11.6384 11.3898 11.4566 11.5902 11.2213 11.7031L5.66765 14.3687C4.70841 14.8291 5.03635 16.2703 6.10036 16.2703H15.6962C16.2522 16.2703 16.7015 15.8166 16.6961 15.2606Z" />
                              <path d="M18.6471 15.2703V5.88936C18.6471 4.84679 20.0428 4.49998 20.5308 5.4213L23.5833 11.1845C23.7 11.4049 23.8948 11.5737 24.1296 11.6578L31.5829 14.3289C32.6388 14.7073 32.3671 16.2703 31.2455 16.2703H19.6471C19.0948 16.2703 18.6471 15.8226 18.6471 15.2703Z" />
                              <path d="M18.6471 31.4643V19.3784C18.6471 18.8261 19.0948 18.3784 19.6471 18.3784H29.2853C30.3376 18.3784 30.676 19.7947 29.7374 20.2704L24.1129 23.1208C23.889 23.2343 23.716 23.4278 23.6281 23.663L20.5839 31.8141C20.1941 32.8578 18.6471 32.5783 18.6471 31.4643Z" />
                              <path d="M16.7059 28.9873V19.3784C16.7059 18.8261 16.2582 18.3784 15.7059 18.3784H3.83963C2.71522 18.3784 2.44656 19.9473 3.50691 20.3214L11.5457 23.1578C11.7987 23.247 12.0052 23.4342 12.1188 23.6772L14.8 29.4109C15.2531 30.3797 16.7059 30.0568 16.7059 28.9873Z" />
                            </svg>
                          </button>
                        </div>
                        <div className="templates-title-wrap">
                          <h3 className="card-tile-title" title={template.name}>{template.name}</h3>
                          <p className="templates-subtitle">{currentTypeLabel}</p>
                        </div>
                      </div>

                      <div className="templates-window-body">
                        <div className="card-tile-actions templates-actions">
                          <button className="btn btn-primary btn-block templates-edit-btn" onClick={() => handleEdit(template)}>
                            <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                              <path d="M7 7h-1a2 2 0 0 0 -2 2v9a2 2 0 0 0 2 2h9a2 2 0 0 0 2 -2v-1"></path>
                              <path d="M20.385 6.585a2.1 2.1 0 0 0 -2.97 -2.97l-8.415 8.385v3h3l8.385 -8.415z"></path>
                              <path d="M16 5l3 3"></path>
                            </svg>
                            <span>Изменить</span>
                          </button>
                          <div className="templates-menu-wrap" onClick={(e) => e.stopPropagation()}>
                            <button
                              className={`templates-menu-control ${openMenuUUID === template.uuid ? 'open' : ''}`}
                              onClick={(e) => {
                                e.stopPropagation();
                                setOpenMenuUUID(openMenuUUID === template.uuid ? null : template.uuid);
                              }}
                              title="More actions"
                              aria-label="More actions"
                            >
                              <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" aria-hidden="true">
                                <path d="M6 9l6 6l6 -6"></path>
                              </svg>
                            </button>
                            {openMenuUUID === template.uuid && (
                              <div className="templates-dropdown-menu">
                                <button className="templates-dropdown-item" onClick={() => handleCopyUUID(template)}>
                                  Copy UUID
                                </button>
                                <button className="templates-dropdown-item" onClick={() => handleRename(template)}>
                                  Rename
                                </button>
                                <button className="templates-dropdown-item danger" onClick={() => handleDelete(template)}>
                                  Delete
                                </button>
                              </div>
                            )}
                          </div>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              ))}
            </div>
          )}
        </>
      ) : (
        <div className="right-editor-inline" aria-label="Template editor panel">
          <div className="right-editor-modal right-editor-inline-window">
            <div className="modal-body-fullscreen right-editor-layout">
              <div className="right-editor-card right-editor-head">
                <div className="right-editor-head-left">
                  <button type="button" className="editor-round-icon editor-round-icon-cyan" title="Template">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M14 2H6a2 2 0 0 0 -2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2 -2V8z" />
                      <polyline points="14 2 14 8 20 8" />
                      <line x1="16" y1="13" x2="8" y2="13" />
                      <line x1="16" y1="17" x2="8" y2="17" />
                    </svg>
                  </button>
                  <div className="right-editor-title-stack">
                    <h4 className="right-editor-title">{editingTemplate?.name || templateName || 'Default'}</h4>
                    <span className="right-editor-subtitle">{currentTypeLabel}</span>
                  </div>
                </div>
                <div className="right-editor-head-actions">
                  <button type="button" className="editor-round-icon editor-round-icon-lime" title="Validate" onClick={validateEditorContent}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M12 17h.01" />
                      <path d="M12 3a9 9 0 0 1 9 9a9 9 0 1 1 -9 -9z" />
                      <path d="M9.09 9a3 3 0 1 1 5.82 1c0 2-3 3-3 3" />
                    </svg>
                  </button>
                  <button type="button" className="editor-round-icon editor-round-icon-gray" title="Reset" onClick={resetEditorValues}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <path d="M3 10h11a4 4 0 1 1 0 8h-1" />
                      <path d="M6 14l-4-4l4-4" />
                    </svg>
                  </button>
                  <button type="button" className="editor-round-icon editor-round-icon-default" title="Свернуть" onClick={() => setShowModal(false)}>
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <line x1="5" y1="12" x2="19" y2="12" />
                    </svg>
                  </button>
                </div>
              </div>

              <div className="right-editor-card right-editor-main">
                <form id="template-editor-form" className="config-editor-form right-editor-form" onSubmit={handleSave}>
                  {isJSONType ? (
                    <div className="form-group config-editor-group">
                      <label className="form-label">Template JSON *</label>
                      <textarea
                        className="form-textarea templates-json-editor-full"
                        value={jsonValue}
                        onChange={(e) => setJsonValue(e.target.value)}
                      />
                    </div>
                  ) : (
                    <div className="form-group config-editor-group">
                      <label className="form-label">Template YAML *</label>
                      <textarea
                        className="form-textarea templates-yaml-editor-full"
                        value={yamlValue}
                        onChange={(e) => setYamlValue(e.target.value)}
                      />
                    </div>
                  )}
                </form>
              </div>

              <div className="right-editor-card right-editor-footer">
                <button type="submit" className="btn btn-primary" form="template-editor-form" disabled={saving}>
                  {saving ? 'Saving...' : 'Сохранить'}
                </button>
                <div className="right-editor-footer-actions">
                  <button type="button" className="editor-round-icon editor-round-icon-default" title="Actions">
                    <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                      <line x1="4" y1="6" x2="20" y2="6" />
                      <line x1="7" y1="12" x2="20" y2="12" />
                      <line x1="10" y1="18" x2="20" y2="18" />
                    </svg>
                  </button>
                  <button
                    type="button"
                    className="btn btn-secondary"
                    onClick={formatJson}
                    disabled={!isJSONType && !currentWasmEngine}
                  >
                    Форматировать
                  </button>
                  <button type="button" className="btn btn-secondary" onClick={() => setShowModal(false)}>
                    Закрыть
                  </button>
                </div>
              </div>
            </div>
          </div>
        </div>
      )}

      <ConfirmActionModal
        open={Boolean(deleteTarget)}
        title="Удалить"
        message={`Вы уверены, что хотите удалить шаблон "${deleteTarget?.name || ''}"? Это действие нельзя отменить.`}
        confirmLabel="Удалить"
        cancelLabel="Закрыть"
        loading={deleting}
        onClose={() => {
          if (!deleting) {
            setDeleteTarget(null);
          }
        }}
        onConfirm={handleConfirmDelete}
      />
    </div>
  );
}

export default Templates;
