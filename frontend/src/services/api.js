const API_BASE = 'http://localhost:8080/api';

export const fetchPolls = async () => {
  const res = await fetch(`${API_BASE}/polls`);
  if (!res.ok) throw new Error('Failed to fetch polls');
  const json = await res.json();
  return json.data;
};

export const createPoll = async (pollData) => {
  const res = await fetch(`${API_BASE}/polls`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(pollData),
  });
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.error?.message || 'Failed to create poll');
  }
  const json = await res.json();
  return json.data;
};

export const votePoll = async (pollId, optionId) => {
  const res = await fetch(`${API_BASE}/polls/${pollId}/vote`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ option_id: optionId }),
  });
  if (!res.ok) {
    const error = await res.json();
    throw new Error(error.error?.message || 'Failed to vote');
  }
  const json = await res.json();
  return json.data;
};

export const deletePoll = async (pollId) => {
  const res = await fetch(`${API_BASE}/polls/${pollId}`, {
    method: 'DELETE',
  });
  if (!res.ok) throw new Error('Failed to delete poll');
  return pollId;
};
