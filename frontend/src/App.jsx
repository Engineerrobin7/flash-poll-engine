import React, { useState, useEffect } from 'react';
import PollDashboard from './pages/PollDashboard';
import './index.css';

// Simple Toast Provider
export const ToastContext = React.createContext();

function App() {
  const [toasts, setToasts] = useState([]);

  const showToast = (msg) => {
    const id = Math.random();
    setToasts([...toasts, { id, msg }]);
    setTimeout(() => {
      setToasts(prev => prev.filter(t => t.id !== id));
    }, 3000);
  };

  return (
    <ToastContext.Provider value={{ showToast }}>
      <PollDashboard />
      <div className="toast-container">
        {toasts.map(t => (
          <div key={t.id} className="toast">
            {t.msg}
          </div>
        ))}
      </div>
    </ToastContext.Provider>
  );
}

export default App;
