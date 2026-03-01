import { useEffect, useMemo, useState } from 'react';
import { authApi } from '../api';

function Login({ onAuthenticated }) {
  const [bootstrap, setBootstrap] = useState(null);
  const [username, setUsername] = useState('admin');
  const [password, setPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');

  useEffect(() => {
    let mounted = true;

    const loadBootstrap = async () => {
      try {
        const response = await authApi.bootstrap();
        if (!mounted) {
          return;
        }

        setBootstrap(response);
        if (response?.default_username) {
          setUsername(response.default_username);
        }
      } catch (err) {
        if (mounted) {
          setError(err.message || 'Failed to load login settings');
        }
      }
    };

    loadBootstrap();
    return () => {
      mounted = false;
    };
  }, []);

  const branding = useMemo(() => {
    const fallback = { title: 'V2RS', logoUrl: null };
    if (!bootstrap?.branding_settings) {
      return fallback;
    }
    return {
      title: bootstrap.branding_settings.title || fallback.title,
      logoUrl: bootstrap.branding_settings.logoUrl || null,
    };
  }, [bootstrap]);

  const requiresSetup = bootstrap?.has_admin_configured === false;

  const handleSubmit = async (event) => {
    event.preventDefault();
    setError('');
    setSubmitting(true);

    try {
      if (requiresSetup) {
        if (password !== confirmPassword) {
          setError('Пароли не совпадают');
          setSubmitting(false);
          return;
        }
      }

      const response = requiresSetup
        ? await authApi.setup({ username, password })
        : await authApi.login({ username, password });

      setPassword('');
      setConfirmPassword('');
      onAuthenticated?.(response);
    } catch (err) {
      setError(err.message || (requiresSetup ? 'Не удалось создать учетную запись' : 'Login failed'));
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <div className="login-page">
      <div className="login-card">
        <div className="login-brand">
          {branding.logoUrl ? (
            <img className="login-brand-logo" src={branding.logoUrl} alt="Panel logo" />
          ) : (
            <svg className="login-brand-icon" viewBox="0 0 16 16" fill="none" xmlns="http://www.w3.org/2000/svg" aria-hidden="true">
              <path
                d="M8 1a.75.75 0 0 1 .75.75v12.5a.75.75 0 0 1-1.5 0V1.75A.75.75 0 0 1 8 1Zm6 2a.75.75 0 0 1 .75.75v8.5a.75.75 0 0 1-1.5 0v-8.5A.75.75 0 0 1 14 3ZM5 4a.75.75 0 0 1 .75.75v6.5a.75.75 0 0 1-1.5 0v-6.5A.75.75 0 0 1 5 4Zm6 1a.75.75 0 0 1 .75.75v4.5a.75.75 0 0 1-1.5 0v-4.5A.75.75 0 0 1 11 5ZM2 6a.75.75 0 0 1 .75.75v2.5a.75.75 0 0 1-1.5 0v-2.5A.75.75 0 0 1 2 6Z"
                fill="currentColor"
              />
            </svg>
          )}
          <h1 className="login-brand-title">{branding.title}</h1>
          {requiresSetup && (
            <div className="login-hint">
              Учетная запись администратора не найдена. Создайте первого администратора для входа в панель.
            </div>
          )}
        </div>

        <form className="login-form" onSubmit={handleSubmit}>
          <div className="form-group">
            <label className="form-label" htmlFor="login-username">Имя пользователя</label>
            <input
              id="login-username"
              className="form-input"
              autoComplete="username"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              required
            />
          </div>

          <div className="form-group">
            <label className="form-label" htmlFor="login-password">Пароль</label>
            <input
              id="login-password"
              className="form-input"
              type="password"
              autoComplete={requiresSetup ? 'new-password' : 'current-password'}
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              required
            />
          </div>

          {requiresSetup && (
            <div className="form-group">
              <label className="form-label" htmlFor="login-confirm-password">Подтвердите пароль</label>
              <input
                id="login-confirm-password"
                className="form-input"
                type="password"
                autoComplete="new-password"
                value={confirmPassword}
                onChange={(event) => setConfirmPassword(event.target.value)}
                required
              />
            </div>
          )}

          {error && <div className="login-error">{error}</div>}

          <button className="btn btn-primary btn-block" type="submit" disabled={submitting}>
            {submitting ? (requiresSetup ? 'Создание...' : 'Вход...') : (requiresSetup ? 'Создать и войти' : 'Войти')}
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;
