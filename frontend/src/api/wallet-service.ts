import {useQuery, useMutation, useQueryClient } from '@tanstack/react-query'
import { apiClient} from './client'
import { CreateWalletRequest, WalletResponse } from '../types/wallet'

// Cüzdanları çektik.
export const getUserWallets = () => {

    return useQuery({
        queryKey: ['wallets'],
        queryFn: async () => {
            const { data } = await apiClient.get<WalletResponse[]>('/wallets');
            return data;
        },
    });
};

// cüzdan oluşturma işlemi için mutation hook'u.
export const createWallet = () => {
    
    const queryClient = useQueryClient();

    return useMutation({
        mutationFn: async (newWallet: CreateWalletRequest) => {
            const { data} = await apiClient.post<WalletResponse>('/wallets', newWallet);
            return data;
        },
        onSuccess: () => {
            // başarılı olursa, cüzdan listesini güncelle.
            queryClient.invalidateQueries({ queryKey: ['wallets']});
        },
    });
};