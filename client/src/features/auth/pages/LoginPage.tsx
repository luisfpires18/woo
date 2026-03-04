import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { Button } from '../../../components/Button/Button';
import { Input } from '../../../components/Input/Input';
import { Card } from '../../../components/Card/Card';
import { ApiRequestError } from '../../../services/api';
import styles from './LoginPage.module.css';

export function LoginPage() {
  const login = useAuthStore((s) => s.login);
  const navigate = useNavigate();

  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      await login({ email, password });
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
          label="Email"
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          placeholder="your@email.com"
          required
          autoComplete="email"
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
