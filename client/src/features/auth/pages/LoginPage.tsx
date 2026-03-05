import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { Button } from '../../../components/Button/Button';
import { Input } from '../../../components/Input/Input';
import { Card } from '../../../components/Card/Card';
import { ApiRequestError } from '../../../services/api';
import styles from './LoginPage.module.css';

export function LoginPage() {
  const doLogin = useAuthStore((s) => s.login);
  const navigate = useNavigate();

  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await doLogin({ login, password });
      navigate('/');
    } catch (err) {
      if (err instanceof ApiRequestError) {
        setError(err.message);
      } else {
        setError('An unexpected error occurred');
      }
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card>
      <form onSubmit={handleSubmit} className={styles.form}>
        <h2 className={styles.heading}>Sign In</h2>

        {error && <p className={styles.error}>{error}</p>}

        <Input
          label="Username or Email"
          type="text"
          value={login}
          onChange={(e) => setLogin(e.target.value)}
          placeholder="username or email"
          required
          autoComplete="username"
        />

        <Input
          label="Password"
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          placeholder="Enter your password"
          required
          autoComplete="current-password"
        />

        <Button type="submit" loading={loading} size="lg">
          Enter the Realm
        </Button>

        <p className={styles.link}>
          No account yet?{' '}
          <Link to="/register">Join the battle</Link>
        </p>
      </form>
    </Card>
  );
}
