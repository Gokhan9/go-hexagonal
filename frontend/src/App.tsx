import { Button } from './components/ui/Button';
import { Card } from './components/ui/Card';
import { Input } from './components/ui/Input';

function App() {
  return (
    <div style={{ padding: '2rem', display: 'flex', justifyContent: 'center', alignItems: 'center', minHeight: '100vh' }}>
      <Card style={{ width: '300px' }}>
        <h2 style={{ marginBottom: '1rem', color: 'var(--color-text-main)' }}>Giriş Yap</h2>
        <div style={{ marginBottom: '1rem' }}>
          <Input label="E-posta" placeholder="ornek@mail.com" />
        </div>
        <Button onClick={() => console.log('Tıklandı!')}>Gönder</Button>
      </Card>
    </div>
  );
  
}

export default App;
