import React, { useState, useContext, useEffect } from 'react';
import { votePoll, deletePoll } from '../services/api';
import { ToastContext } from '../App';

const PollCard = ({ poll, onUpdate, onDelete, isHighlighted }) => {
  const [loading, setLoading] = useState(false);
  const [votedId, setVotedId] = useState(null);
  const [showStats, setShowStats] = useState(false);
  const [pulse, setPulse] = useState(false);
  const { showToast } = useContext(ToastContext);

  // Trigger pulse when poll data changes from external source (SSE)
  useEffect(() => {
    if (poll.total_votes > 0) {
      setPulse(true);
      const timer = setTimeout(() => setPulse(false), 1000);
      return () => clearTimeout(timer);
    }
  }, [poll.total_votes]);

  const castVote = async (oid) => {
    if (votedId || loading) return;
    setLoading(true);
    try {
      const updated = await votePoll(poll.id, oid);
      onUpdate(updated);
      setVotedId(oid);
      setShowStats(true);
      showToast("VOTE REGISTERED // ANALYTICS UPDATED");
    } catch (e) {
      if (e.message === "POLL_EXPIRED_OR_RESET") {
        showToast("SIGNAL EXPIRED: BACKEND RESET. REFRESHING...");
        setTimeout(() => window.location.reload(), 2000);
      } else {
        showToast(`VOTE FAILED: ${e.message.toUpperCase()}`);
      }
    } finally {
      setLoading(false);
    }
  };

  const remove = async () => {
    if (!confirm("PURGE THIS DATA?")) return;
    try {
      await deletePoll(poll.id);
      onDelete(poll.id);
      showToast("POLL PURGED FROM DATABASE");
    } catch (e) {
      showToast(`ERROR: ${e.message.toUpperCase()}`);
    }
  };

  const copyLink = () => {
    const url = `${window.location.origin}/?id=${poll.id}`;
    navigator.clipboard.writeText(url);
    showToast("SIGNAL LINK COPIED TO CLIPBOARD");
  };

  const isRevealed = votedId || showStats;

  return (
    <div
      className="poll-card"
      style={isHighlighted ? { border: '6px solid var(--tech)', transform: 'scale(1.02)', boxShadow: '12px 12px 0px 0px #000' } : {}}
    >
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
        <div style={{display: 'flex', alignItems: 'center', gap: '8px'}}>
          <span className={`tag tag-${poll.category}`}>{poll.category} {isHighlighted && "// TARGET"}</span>
          {pulse && <span className="tag" style={{background: '#76ff03', animation: 'blink 0.5s infinite'}}>LIVE UPDATE</span>}
        </div>
        <div style={{display: 'flex', gap: '8px'}}>
          <button className="tag" onClick={copyLink} style={{background: '#fff', cursor: 'pointer'}}>SHARE</button>
          <button className="tag" onClick={remove} style={{background: '#000', color: '#fff', cursor: 'pointer'}}>PURGE</button>
        </div>
      </div>

      <h3 className="poll-question">{poll.question}</h3>

      <div className="options-container">
        {poll.options.map((opt) => (
          <button
            key={opt.id}
            className="option-btn"
            onClick={() => castVote(opt.id)}
            disabled={loading || votedId}
          >
            {isRevealed && (
              <div className="progress-fill" style={{ width: `${opt.percentage}%` }} />
            )}

            <div style={{ position: 'relative', zIndex: 1, display: 'flex', justifyContent: 'space-between', width: '100%' }}>
              <div style={{display: 'flex', alignItems: 'center', gap: '10px'}}>
                {votedId === opt.id && <span style={{fontSize: '18px'}}>🎯</span>}
                <span>{opt.option_text}</span>
              </div>

              {isRevealed && (
                <span style={{ fontWeight: '900' }}>{Math.round(opt.percentage)}%</span>
              )}
            </div>
          </button>
        ))}
      </div>

      <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', marginTop: '20px', borderTop: '3px solid #000', paddingTop: '15px' }}>
        <div style={{ fontWeight: '800', fontSize: '12px', textTransform: 'uppercase' }}>
          {poll.total_votes} SIGNALS // {new Date(poll.created_at).toLocaleDateString()}
        </div>

        {!votedId && (
          <button
            className="tag"
            style={{ cursor: 'pointer', background: isRevealed ? '#eee' : '#fff' }}
            onClick={() => setShowStats(!showStats)}
          >
            {showStats ? 'HIDE ANALYTICS' : 'VIEW ANALYTICS'}
          </button>
        )}
      </div>
    </div>
  );
};

export default PollCard;
