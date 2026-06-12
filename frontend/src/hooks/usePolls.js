import { useState, useEffect, useCallback } from 'react';
import { fetchPolls } from '../services/api';

export const usePolls = () => {
  const [polls, setPolls] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);

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

  useEffect(() => {
    refresh();

    // Setup SSE for real-time updates
    const eventSource = new EventSource('http://localhost:8080/api/events');

    eventSource.onmessage = (event) => {
      const updatedPoll = JSON.parse(event.data);
      setPolls(prev => prev.map(p => p.id === updatedPoll.id ? updatedPoll : p));
    };

    eventSource.onerror = (err) => {
      console.error("EventSource failed:", err);
      eventSource.close();
    };

    return () => {
      eventSource.close();
    };
  }, [refresh]);

  const addPoll = (poll) => setPolls(prev => [poll, ...prev]);
  const updatePoll = (poll) => setPolls(prev => prev.map(p => p.id === poll.id ? poll : p));
  const removePoll = (id) => setPolls(prev => prev.filter(p => p.id !== id));

  return { polls, loading, error, refresh, addPoll, updatePoll, removePoll };
};
