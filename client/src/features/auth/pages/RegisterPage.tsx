import { useState } from 'react';
import { Link, useNavigate } from 'react-router-dom';
import { useAuthStore } from '../../../stores/authStore';
import { Button } from '../../../components/Button/Button';
import { Input } from '../../../components/Input/Input';
import { Card } from '../../../components/Card/Card';
import { ApiRequestError } from '../../../services/api';
import styles from './RegisterPage.module.css';

export function RegisterPage() {
  const register = useAuthStore((s) => s.register);
  const navigate = useNavigate();

  const [username, setUsername] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    setLoading(true);
    try {
      await register({ username, email, password });
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
        <h2 className={styles.heading}>Create Account</h2>

        {error && <p className={styles.error}>{error}</p>}

        <Input
          label="Username"
          type="text"
          value={username}
          onChange={(e) => setUsername(e.target.value)}
          placeholder="3-20 characters"
          required
          autoComplete="username"
          minLength={3}
          maxLength={20}
        />

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
          placeholder="At least 8 characters"
          required
          autoComplete="new-password"
          minLength={8}
        />

        <Button type="submit" loading={loading} size="lg">
          Create Account
        </Button>

        <p className={styles.link}>
          Already have an account?{' '}
          <Link to="/login">Sign in</Link>
        </p>
      </form>
    </Card>
  );
}
