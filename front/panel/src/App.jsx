import React, { useState, useEffect } from 'react';
import Dashboard from './pages/Dashboard';
import Nodes from './pages/Nodes';
import Users from './pages/Users';
import Settings from './pages/Settings';
import ConfigProfiles from './pages/ConfigProfiles';
import InternalSquads from './pages/InternalSquads';
import Hosts from './pages/Hosts';
import Templates from './pages/Templates';
import SubscriptionSettings from './pages/SubscriptionSettings';
import Login from './pages/Login';
import { authApi, healthCheck } from './api';
import './App.css';

const TEMPLATE_TYPES = ['XRAY_JSON', 'MIHOMO', 'STASH', 'CLASH', 'SINGBOX'];
const BACKEND_RECHECK_INTERVAL_MS = 5000;

function BackendUnavailableLoader() {
  return (
    <div className="backend-unavailable-loader">
      <div className="backend-unavailable-loader-track">
        <div className="backend-unavailable-loader-bar" />
      </div>
    </div>
  );
}

function App() {
  const [isAuthenticated, setIsAuthenticated] = useState(null);
  const [authChecking, setAuthChecking] = useState(true);
  const [adminProfile, setAdminProfile] = useState(null);
  const [branding, setBranding] = useState({
    title: 'V2RS',
    logoUrl: null,
  });
  const [activePage, setActivePage] = useState(() => {
    const saved = localStorage.getItem('v2ray-active-page') || 'dashboard';
    if (saved.startsWith('settings-subscription-')) {
      return 'settings-subscription';
    }
    return saved;
  });
  const [sidebarOpen, setSidebarOpen] = useState(true);
  const [backendStatus, setBackendStatus] = useState('checking');
  const [templatesExpanded, setTemplatesExpanded] = useState(() => {
    const saved = localStorage.getItem('v2ray-active-page') || 'dashboard';
    return saved.startsWith('templates-');
  });

  useEffect(() => {
    localStorage.setItem('v2ray-active-page', activePage);
  }, [activePage]);

  useEffect(() => {
    if (activePage.startsWith('templates-')) {
      setTemplatesExpanded(true);
    }
  }, [activePage]);

  useEffect(() => {
    let mounted = true;

    const restoreSession = async () => {
      try {
        const session = await authApi.me();
        if (!mounted) {
          return;
        }
        setIsAuthenticated(true);
        setAdminProfile(session?.admin || null);
        if (session?.branding_settings) {
          setBranding({
            title: session.branding_settings.title || 'V2RS',
            logoUrl: session.branding_settings.logoUrl || null,
          });
        }
      } catch {
        if (mounted) {
          setIsAuthenticated(false);
          setAdminProfile(null);
        }
      } finally {
        if (mounted) {
          setAuthChecking(false);
        }
      }
    };

    restoreSession();
    return () => {
      mounted = false;
    };
  }, []);

  useEffect(() => {
    const handleAuthExpired = () => {
      setIsAuthenticated(false);
      setAdminProfile(null);
      setBackendStatus('online');
    };

    const handleBrandingUpdated = (event) => {
      const nextBranding = event?.detail || {};
      setBranding((prev) => ({
        title: nextBranding.title || prev.title || 'V2RS',
        logoUrl: nextBranding.logoUrl === undefined ? prev.logoUrl : nextBranding.logoUrl,
      }));
    };

    window.addEventListener('auth-expired', handleAuthExpired);
    window.addEventListener('panel-branding-updated', handleBrandingUpdated);
    return () => {
      window.removeEventListener('auth-expired', handleAuthExpired);
      window.removeEventListener('panel-branding-updated', handleBrandingUpdated);
    };
  }, []);

  useEffect(() => {
    let isMounted = true;

    const updateBackendStatus = (nextStatus) => {
      if (!isMounted) {
        return;
      }
      setBackendStatus((currentStatus) =>
        currentStatus === nextStatus ? currentStatus : nextStatus
      );
    };

    const checkBackend = async () => {
      try {
        const health = await healthCheck();
        updateBackendStatus(health?.status === 'ok' ? 'online' : 'offline');
      } catch {
        updateBackendStatus('offline');
      }
    };

    checkBackend();
    const intervalId = window.setInterval(checkBackend, BACKEND_RECHECK_INTERVAL_MS);

    return () => {
      isMounted = false;
      window.clearInterval(intervalId);
    };
  }, []);

  const isTemplatesPage = activePage.startsWith('templates-');
  const activeTemplateType = isTemplatesPage ? activePage.replace('templates-', '') : '';
  const isBackendAvailable = backendStatus === 'online';

  const renderPage = () => {
    switch (activePage) {
      case 'dashboard':
        return <Dashboard />;
      case 'nodes':
        return <Nodes />;
      case 'users':
        return <Users />;
      case 'config-profiles':
        return <ConfigProfiles />;
      case 'internal-squads':
        return <InternalSquads />;
      case 'hosts':
        return <Hosts />;
      case 'settings-subscription':
        return <SubscriptionSettings />;
      case 'settings':
        return <Settings />;
      default:
        if (activePage.startsWith('templates-')) {
          return <Templates activeTemplateType={activeTemplateType} />;
        }
        return <Dashboard />;
    }
  };

  const handleNavClick = (page, closeSidebarOnMobile = true) => {
    setActivePage(page);
    if (closeSidebarOnMobile && window.innerWidth <= 1024) {
      setSidebarOpen(false);
    }
  };

  const handleTemplatesClick = (event) => {
    event.preventDefault();
    const nextExpanded = !templatesExpanded;
    setTemplatesExpanded(nextExpanded);
    if (!isTemplatesPage) {
      handleNavClick(`templates-${TEMPLATE_TYPES[0]}`, false);
    }
  };

  const handleTemplateTypeClick = (type) => {
    setTemplatesExpanded(true);
    handleNavClick(`templates-${type}`);
  };

  const toggleSidebar = () => {
    setSidebarOpen(!sidebarOpen);
  };

  const handleAuthenticated = (payload) => {
    setIsAuthenticated(true);
    setAuthChecking(false);
    setAdminProfile(payload?.admin || null);
    if (payload?.branding_settings) {
      setBranding({
        title: payload.branding_settings.title || 'V2RS',
        logoUrl: payload.branding_settings.logoUrl || null,
      });
    }
  };

  const handleLogout = async () => {
    try {
      await authApi.logout();
    } catch (error) {
      console.error('Logout failed:', error);
    } finally {
      setIsAuthenticated(false);
      setAdminProfile(null);
      setActivePage('dashboard');
    }
  };

  if (authChecking) {
    return (
      <div className="app app-auth-loading">
        <BackendUnavailableLoader />
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Login onAuthenticated={handleAuthenticated} />;
  }

  return (
    <div className={`app ${sidebarOpen ? 'sidebar-open' : 'sidebar-closed'}`}>
      <div
        className={`sidebar-overlay ${sidebarOpen ? 'open' : ''}`}
        onClick={() => setSidebarOpen(false)}
      />

      <aside className={`sidebar ${sidebarOpen ? 'open' : ''}`}>
        <div className="sidebar-brand">
          {branding.logoUrl ? (
            <img className="logo-image" src={branding.logoUrl} alt="Panel logo" />
          ) : (
            <svg className="logo-icon" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
              <path
                d="M8 1a.75.75 0 0 1 .75.75v12.5a.75.75 0 0 1-1.5 0V1.75A.75.75 0 0 1 8 1Zm6 2a.75.75 0 0 1 .75.75v8.5a.75.75 0 0 1-1.5 0v-8.5A.75.75 0 0 1 14 3ZM5 4a.75.75 0 0 1 .75.75v6.5a.75.75 0 0 1-1.5 0v-6.5A.75.75 0 0 1 5 4Zm6 1a.75.75 0 0 1 .75.75v4.5a.75.75 0 0 1-1.5 0v-4.5A.75.75 0 0 1 11 5ZM2 6a.75.75 0 0 1 .75.75v2.5a.75.75 0 0 1-1.5 0v-2.5A.75.75 0 0 1 2 6Z"
                fill="currentColor"
              />
            </svg>
          )}
          <div className="sidebar-brand-text">
            <span className="sidebar-brand-title">{branding.title || 'V2RS'}</span>
          </div>
        </div>
        <nav className="nav">
          <div className="nav-section">
            <h6 className="nav-section-title">Обзор</h6>
            <ul className="nav-list">
              <li className={`nav-item ${activePage === 'dashboard' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('dashboard'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <rect x="3" y="3" width="7" height="7"/>
                    <rect x="14" y="3" width="7" height="7"/>
                    <rect x="14" y="14" width="7" height="7"/>
                    <rect x="3" y="14" width="7" height="7"/>
                  </svg>
                  <span>Главная</span>
                </a>
              </li>
            </ul>
          </div>

          <div className="nav-section">
            <div className="nav-divider" />
            <h6 className="nav-section-title">Управление</h6>
            <ul className="nav-list">
              <li className={`nav-item ${activePage === 'users' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('users'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <path d="M23 21v-2a4 4 0 00-3-3.87"/>
                    <path d="M16 3.13a4 4 0 010 7.75"/>
                  </svg>
                  <span>Пользователи</span>
                </a>
              </li>
              <li className={`nav-item ${activePage === 'internal-squads' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('internal-squads'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M17 21v-2a4 4 0 00-4-4H5a4 4 0 00-4 4v2"/>
                    <circle cx="9" cy="7" r="4"/>
                    <path d="M23 21v-2a4 4 0 00-3-3.87"/>
                    <path d="M16 3.13a4 4 0 010 7.75"/>
                  </svg>
                  <span>Внутренние сквады</span>
                </a>
              </li>
              <li className={`nav-item ${activePage === 'config-profiles' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('config-profiles'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                    <line x1="16" y1="13" x2="8" y2="13"/>
                    <line x1="16" y1="17" x2="8" y2="17"/>
                    <polyline points="10 9 9 9 8 9"/>
                  </svg>
                  <span>Профили</span>
                </a>
              </li>
              <li className={`nav-item ${activePage === 'hosts' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('hosts'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
                  </svg>
                  <span>Хосты</span>
                </a>
              </li>
              <li className={`nav-item ${activePage === 'nodes' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('nodes'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="10"/>
                    <line x1="2" y1="12" x2="22" y2="12"/>
                    <path d="M12 2a15.3 15.3 0 014 10 15.3 15.3 0 01-4 10 15.3 15.3 0 01-4-10 15.3 15.3 0 014-10z"/>
                  </svg>
                  <span>Ноды</span>
                </a>
              </li>
              <li className={`nav-item ${activePage === 'settings' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('settings'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
                  </svg>
                  <span>Настройки панели</span>
                </a>
              </li>
            </ul>
          </div>

          <div className="nav-section">
            <div className="nav-divider" />
            <h6 className="nav-section-title">Подписка</h6>
            <ul className="nav-list">
              <li className={`nav-item ${activePage === 'settings-subscription' ? 'active' : ''}`}>
                <a href="#" className="nav-link" onClick={(e) => { e.preventDefault(); handleNavClick('settings-subscription'); }}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <circle cx="12" cy="12" r="3"/>
                    <path d="M19.4 15a1.65 1.65 0 00.33 1.82l.06.06a2 2 0 010 2.83 2 2 0 01-2.83 0l-.06-.06a1.65 1.65 0 00-1.82-.33 1.65 1.65 0 00-1 1.51V21a2 2 0 01-2 2 2 2 0 01-2-2v-.09A1.65 1.65 0 009 19.4a1.65 1.65 0 00-1.82.33l-.06.06a2 2 0 01-2.83 0 2 2 0 010-2.83l.06-.06a1.65 1.65 0 00.33-1.82 1.65 1.65 0 00-1.51-1H3a2 2 0 01-2-2 2 2 0 012-2h.09A1.65 1.65 0 004.6 9a1.65 1.65 0 00-.33-1.82l-.06-.06a2 2 0 010-2.83 2 2 0 012.83 0l.06.06a1.65 1.65 0 001.82.33H9a1.65 1.65 0 001-1.51V3a2 2 0 012-2 2 2 0 012 2v.09a1.65 1.65 0 001 1.51 1.65 1.65 0 001.82-.33l.06-.06a2 2 0 012.83 0 2 2 0 010 2.83l-.06.06a1.65 1.65 0 00-.33 1.82V9a1.65 1.65 0 001.51 1H21a2 2 0 012 2 2 2 0 01-2 2h-.09a1.65 1.65 0 00-1.51 1z"/>
                  </svg>
                  <span>Настройки подписки</span>
                </a>
              </li>
              <li className="nav-item">
                <a href="#" className="nav-link nav-link-with-arrow" onClick={handleTemplatesClick}>
                  <svg className="nav-icon" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <path d="M14 2H6a2 2 0 00-2 2v16a2 2 0 002 2h12a2 2 0 002-2V8z"/>
                    <polyline points="14 2 14 8 20 8"/>
                    <line x1="16" y1="13" x2="8" y2="13"/>
                    <line x1="16" y1="17" x2="8" y2="17"/>
                  </svg>
                  <span>Шаблоны</span>
                  <svg className={`nav-arrow ${templatesExpanded ? 'open' : ''}`} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <polyline points="6 9 12 15 18 9"/>
                  </svg>
                </a>
                {templatesExpanded && (
                  <ul className="nav-sublist">
                    {TEMPLATE_TYPES.map((type) => (
                      <li key={type} className={`nav-subitem ${activeTemplateType === type ? 'active' : ''}`}>
                        <a href="#" className="nav-sublink" onClick={(e) => { e.preventDefault(); handleTemplateTypeClick(type); }}>
                          {type}
                        </a>
                      </li>
                    ))}
                  </ul>
                )}
              </li>
            </ul>
          </div>
        </nav>
      </aside>

      <div className="app-main">
        <header className="header">
          <div className="header-content">
            <div className="header-left">
              <button
                className={`shell-burger ${sidebarOpen ? 'open' : ''}`}
                onClick={toggleSidebar}
                title={sidebarOpen ? 'Close sidebar' : 'Open sidebar'}
                aria-label={sidebarOpen ? 'Close sidebar' : 'Open sidebar'}
              >
                {sidebarOpen ? (
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <line x1="18" y1="6" x2="6" y2="18"/>
                    <line x1="6" y1="6" x2="18" y2="18"/>
                  </svg>
                ) : (
                  <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                    <line x1="3" y1="12" x2="21" y2="12"/>
                    <line x1="3" y1="6" x2="21" y2="6"/>
                    <line x1="3" y1="18" x2="21" y2="18"/>
                  </svg>
                )}
              </button>
            </div>
            <div className="header-actions">
              <span className="header-user-chip header-hide-mobile">{adminProfile?.username || 'admin'}</span>
              <a
                className="header-icon-btn header-social-btn"
                href="https://t.me/v2rs"
                target="_blank"
                rel="noopener noreferrer"
                title="Telegram"
                aria-label="Telegram"
              >
                <svg viewBox="0 0 256 256" fill="currentColor">
                  <path d="M231.49,23.16a13,13,0,0,0-13.23-2.26L15.6,100.21a18.22,18.22,0,0,0,3.12,34.86L68,144.74V200a20,20,0,0,0,34.4,13.88l22.67-23.51L162.35,223a20,20,0,0,0,32.7-10.54L235.67,35.91A13,13,0,0,0,231.49,23.16ZM139.41,77.52,77.22,122.09l-34.43-6.75ZM92,190.06V161.35l15,13.15Zm81.16,10.52L99.28,135.81,205.59,59.63Z"></path>
                </svg>
              </a>
              <a
                className="header-icon-btn header-social-btn"
                href="https://github.com/v2rs/panel"
                target="_blank"
                rel="noopener noreferrer"
                title="GitHub"
                aria-label="GitHub"
              >
                <svg viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M9 19c-4.3 1.4 -4.3 -2.5 -6 -3m12 5v-3.5c0 -1 .1 -1.4 -.5 -2c2.8 -.3 5.5 -1.4 5.5 -6a4.6 4.6 0 0 0 -1.3 -3.2a4.2 4.2 0 0 0 -.1 -3.2s-1.1 -.3 -3.5 1.3a12.3 12.3 0 0 0 -6.2 0c-2.4 -1.6 -3.5 -1.3 -3.5 -1.3a4.2 4.2 0 0 0 -.1 3.2a4.6 4.6 0 0 0 -1.3 3.2c0 4.6 2.7 5.7 5.5 6c-.6 .6 -.6 1.2 -.5 2v3.5"></path>
                </svg>
              </a>
              <span className="header-version header-hide-mobile">
                <svg viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
                  <path d="M8 1a.75.75 0 0 1 .75.75v12.5a.75.75 0 0 1-1.5 0V1.75A.75.75 0 0 1 8 1Zm6 2a.75.75 0 0 1 .75.75v8.5a.75.75 0 0 1-1.5 0v-8.5A.75.75 0 0 1 14 3ZM5 4a.75.75 0 0 1 .75.75v6.5a.75.75 0 0 1-1.5 0v-6.5A.75.75 0 0 1 5 4ZM11 5a.75.75 0 0 1 .75.75v4.5a.75.75 0 0 1-1.5 0v-4.5A.75.75 0 0 1 11 5ZM2 6a.75.75 0 0 1 .75.75v2.5a.75.75 0 0 1-1.5 0v-2.5A.75.75 0 0 1 2 6Z" fill="currentColor" />
                </svg>
                <span>v2.6.1</span>
              </span>
              <button className="header-icon-btn header-logout-btn" type="button" title="Выйти" aria-label="Выйти" onClick={handleLogout}>
                <svg viewBox="0 0 256 256" fill="currentColor">
                  <path d="M216,120a8,8,0,0,1-8,8H85.94l37.65,37.66a8,8,0,0,1-11.32,11.31l-51.31-51.31a8,8,0,0,1,0-11.32l51.31-51.31a8,8,0,0,1,11.32,11.31L85.94,112H208A8,8,0,0,1,216,120ZM144,40a8,8,0,0,0,0,16h40V184H144a8,8,0,0,0,0,16h48a8,8,0,0,0,8-8V48a8,8,0,0,0-8-8Z"></path>
                </svg>
              </button>
            </div>
          </div>
        </header>

        <main className="main-content">
          <div className="content-wrapper">
            {isBackendAvailable ? renderPage() : <BackendUnavailableLoader />}
          </div>
        </main>
      </div>
    </div>
  );
}

export default App;
