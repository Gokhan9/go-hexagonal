import { useState } from 'react';
import { authService } from '../api/auth-service';

export const useAuth = () => {
    const [isLoading, setIsLoading] = useState(false);
    const [error, setError] = useState<string | null>(null);

    const login = async (email: string, password: string) => {
        setIsLoading(true);
        setError(null);

        try {
            const data = await authService.login({ email, password });
            // token kaydetme ve kullanıcı yönlendirme mantığı
            return data;
        } catch (err) {
            setError('Giriş başarısız oldu. Lütfen tekrar deneyin.');
            throw err;
        } finally {
            setIsLoading(false);
        }
    };

    return { login, isLoading, error };
}