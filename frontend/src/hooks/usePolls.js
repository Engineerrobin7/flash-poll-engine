import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchPolls, API_BASE } from '../services/api';

export const usePolls = () => {
  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const reconnectTimeoutRef = useRef(null);

  const refresh = useCallback(async () => {
    setLoading(true);
    setError(null);
    try {
      const data = await fetchPolls();
      setPolls(data);
    } catch (e) {
      setError(e.message);
    } finally {
      setLoading(false);
    }
  }, []);

  const setupSSE = useCallback(() => {
    // Consistent URL usage
    const eventSource = new EventSource(`${API_BASE}/events`);

    eventSource.onmessage = (event) => {
      try {
        const updatedPoll = JSON.parse(event.data);
        setPolls(prev => {
          const exists = prev.find(p => p.id === updatedPoll.id);
          if (exists) {
            return prev.map(p => p.id === updatedPoll.id ? updatedPoll : p);
          }
          return [updatedPoll, ...prev];
        });
      } catch (err) {
        console.error("SSE parse error", err);
      }
    };

    eventSource.onerror = (err) => {
      console.log("SSE dropped. Reconnecting...");
      eventSource.close();
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = setTimeout(() => setupSSE(), 5000);
    };

    return eventSource;
  }, []);

  useEffect(() => {
    refresh();
    const es = setupSSE();
    return () => {
      if (es) es.close();
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);
    };
  }, [refresh, setupSSE]);

  const addPoll = (poll) => setPolls(prev => [poll, ...prev]);
  const updatePoll = (poll) => setPolls(prev => prev.map(p => p.id === poll.id ? poll : p));
  const removePoll = (id) => setPolls(prev => prev.filter(p => p.id !== id));

  return { polls, loading, error, refresh, addPoll, updatePoll, removePoll };
};
