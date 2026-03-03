import { BrowserRouter, Routes, Route } from 'react-router-dom';

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </BrowserRouter>
  );
}

function Home() {
  return (
    <div style={{ textAlign: 'center', marginTop: '4rem' }}>
      <h1 style={{ fontFamily: 'Cinzel, serif' }}>Weapons of Order</h1>
      <p style={{ fontFamily: 'EB Garamond, serif', fontSize: '1.2rem' }}>
        The forge awaits.
      </p>
    </div>
  );
}

export default App;
