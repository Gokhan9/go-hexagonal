import { apiClient } from './client'; // Daha önceden oluşturduğumuz temel axios/fetch client'ı

export const authService = {
    login: async (credentials: { email: string; password: string }) => {
        // backend endpointi buraya gelir..
        const response = await apiClient.post('/aıth/login', credentials);
        return response.data; // backend'den dönen veriyi döndür
    }
};