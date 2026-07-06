import { useState } from 'react';
import { Card } from '../components/ui/Card';
import { Input } from '../components/ui/Input';
import { Button } from '../components/ui/Button';
import { useAuth } from '../hooks/useAuth'
import './LoginPage.css';

export const LoginPage = () => {
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    // Hook'u çağırıyoruz
    const { login, isLoading, error } = useAuth();

    const handleSubmit = async (e: React.FormEvent) => {
        e.preventDefault();
        try {
            await login(email, password);
            console.log('Login successful');
            // TODO: Başarılı giriş sonrası yönlendirme yap (örneğin: navigate('/dashboard'))
        } catch (err) {
            console.error('Login failed', err);
        }
    };

    return (
        <div className="login-page">
            <Card className="login-card">
                
                <h2>Giriş Yap</h2>
                {error && <p style={{ color: 'red', marginTop: '10px' }}>{error}</p>}

                <form onSubmit={handleSubmit} className="login-form">
                    <Input 
                        label="E-posta" 
                        type="email" 
                        value={email} 
                        onChange={(e) => setEmail(e.target.value)} 
                        required
                        disabled={isLoading}
                    />
                    <Input
                    label="Şifre"
                    type="password"
                    name="password"
                    value={password}
                    onChange={(e) => setPassword(e.target.value)}
                    required
                    disabled={isLoading}
                    />
                    <Button type="submit" disabled={isLoading}>
                        {isLoading ? 'Giriş Yapılıyor...' : 'Giriş Yap'}
                    </Button>
                </form>
            </Card>
        </div>
    );
}
