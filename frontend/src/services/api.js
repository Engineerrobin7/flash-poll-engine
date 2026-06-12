// Smart URL detector - Handles any Vercel configuration
const getApiBase = () => {
  let url = import.meta.env.VITE_API_URL;
  if (url) {
    if (!url.includes('/api')) {
      url = url.endsWith('/') ? `${url}api` : `${url}/api`;
    }
    return url.endsWith('/') ? url.slice(0, -1) : url;
  }
  return `http://${window.location.hostname}:8080/api`;
};

export const API_BASE = getApiBase();

const handleResponse = async (res) => {
  const contentType = res.headers.get("content-type");
  if (!contentType || !contentType.includes("application/json")) {
    const text = await res.text();
    if (res.status === 404) throw new Error("POLL_EXPIRED_OR_RESET");
    throw new Error(text || `Error ${res.status}`);
  }

  const json = await res.json();
  if (!res.ok) throw new Error(json.error?.message || json.message || 'FAILED');
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
    body: JSON.stringify({ option_id: Number(optionId) }),
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
  const res = await fetch(`${API_BASE}/stats`);
  const json = await handleResponse(res);
  return json.data;
};
