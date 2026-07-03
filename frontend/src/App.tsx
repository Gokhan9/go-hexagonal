import { getUserWallets } from './api/wallet-service'

function App() {

    // Veriyi Çek
    const { data: wallets, isLoading, isError, error } = getUserWallets();

    if (isLoading) {
        return (
            <div className="App">
              <h1>HEXAGONAL GO & REACT</h1>
              <p>Cüzdanlar yükleniyor.......</p>
            </div>
        );
    }
    
    // Hata varsa, hatayı göster
    if (isError) {
        return (
            <div className="App">
                <h1>HEXAGONAL GO & REACT</h1>
                <p>Error :{(error as Error).message}</p>
            </div>
        );
    }

    // Cüzdanlar yüklendiğinde, cüzdanları listele
    return (
        <div className="App">
            <h1>Hexagonal Go & React</h1>
            <h2>Cüzdanlarım</h2>
            {wallets && wallets.length > 0 ? (
                <ul>
                    {wallets.map((wallet) => (
                        <li key={wallet.id}>
                            <strong>{wallet.currency}</strong>: {wallet.balance}
                        </li>
                    ))}

                </ul>
            ) : (
                <p>Henüz cüzdanınız yok.</p>
            )}
        </div>
    )
}

export default App;