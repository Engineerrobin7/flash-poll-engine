import React, { useState } from 'react';
import PollDashboard from './pages/PollDashboard';
import './index.css';

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
      {/*
        In a full app we'd use react-router, but for this MVP
        we stay on the Dashboard. Even if you open a share link,
        it will load the dashboard where your poll lives.
      */}
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
