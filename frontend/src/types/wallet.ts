export interface CreateWalletRequest {
    owner: string;
    currency: string;
}

export interface WalletResponse {
    id: string;
    ownerid: string;
    balance: number; //api'den float gelir.
    currency: string;
    created_at: string; // ISO string olarak gelir.
}

export interface TransactionRequest {
    amount: number;
    transaction_id: string;
}