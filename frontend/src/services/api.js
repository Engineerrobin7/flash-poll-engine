// Dynamic API configuration
const getApiBase = () => {
  let url = import.meta.env.VITE_API_URL;
  if (url) {
    return url.endsWith('/') ? url.slice(0, -1) : url;
  }
  return `http://${window.location.hostname}:8080/api`;
};

const API_BASE = getApiBase();

// Helper to handle responses safely
const handleResponse = async (res) => {
  const contentType = res.headers.get("content-type");

  // If not JSON, it's likely a server error or rate limit text
  if (!contentType || !contentType.includes("application/json")) {
    const text = await res.text();
    throw new Error(text || `Server Error: ${res.status}`);
  }

  const json = await res.json();

  if (!res.ok) {
    throw new Error(json.error?.message || json.message || 'API Request Failed');
  }

  return json;
};

export const fetchPolls = async () => {
  const res = await fetch(`${API_BASE}/polls`);
  const json = await handleResponse(res);
  return json.data;
};

export const createPoll = async (pollData) => {
  const res = await fetch(`${API_BASE}/polls`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(pollData),
  });
  const json = await handleResponse(res);
  return json.data;
};

export const votePoll = async (pollId, optionId) => {
  const res = await fetch(`${API_BASE}/polls/${pollId}/vote`, {
    method: 'PATCH',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ option_id: optionId }),
  });
  const json = await handleResponse(res);
  return json.data;
};

export const deletePoll = async (pollId) => {
  const res = await fetch(`${API_BASE}/polls/${pollId}`, {
    method: 'DELETE',
  });
  await handleResponse(res);
  return pollId;
};

export const fetchStats = async () => {
  const res = await fetch(`${API_BASE.replace('/api', '')}/api/stats`);
  const json = await handleResponse(res);
  return json.data;
};
