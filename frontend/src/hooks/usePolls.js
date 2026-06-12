import { useState, useEffect, useCallback, useRef } from 'react';
import { fetchPolls } from '../services/api';

const getEventsUrl = () => {
  if (import.meta.env.VITE_API_URL) return `${import.meta.env.VITE_API_URL}/events`;
  return `http://${window.location.hostname}:8080/api/events`;
};

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
    const eventSource = new EventSource(getEventsUrl());

    eventSource.onmessage = (event) => {
      try {
        const updatedPoll = JSON.parse(event.data);
        setPolls(prev => prev.map(p => p.id === updatedPoll.id ? updatedPoll : p));
      } catch (err) {
        console.error("Failed to parse SSE data:", err);
      }
    };

    eventSource.onerror = (err) => {
      console.error("SSE connection lost. Reconnecting in 5s...", err);
      eventSource.close();

      // Prevent multiple concurrent timeouts
      if (reconnectTimeoutRef.current) clearTimeout(reconnectTimeoutRef.current);

      reconnectTimeoutRef.current = setTimeout(() => {
        setupSSE();
      }, 5000);
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
