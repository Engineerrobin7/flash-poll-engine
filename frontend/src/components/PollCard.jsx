import React, { useState, useContext } from 'react';
import { votePoll, deletePoll } from '../services/api';
import { ToastContext } from '../App';

const PollCard = ({ poll, onUpdate, onDelete }) => {
  const [loading, setLoading] = useState(false);
  const [votedId, setVotedId] = useState(null);
  const [showStats, setShowStats] = useState(false);
  const { showToast } = useContext(ToastContext);

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
      showToast("SIGNAL LOST: VOTE FAILED");
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
      showToast("ERROR: PURGE FAILED");
    }
  };

  const copyLink = () => {
    const url = `${window.location.origin}/?id=${poll.id}`;
    navigator.clipboard.writeText(url);
    showToast("SIGNAL LINK COPIED TO CLIPBOARD");
  };

  const isRevealed = votedId || showStats;

  return (
    <div className="poll-card" role="article" aria-label={`Poll: ${poll.question}`}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '15px' }}>
        <span className={`tag tag-${poll.category}`}>{poll.category}</span>
        <div style={{display: 'flex', gap: '8px'}}>
          <button className="tag" onClick={copyLink} style={{background: '#fff', cursor: 'pointer'}}>SHARE</button>
          <button className="tag" onClick={remove} style={{background: '#000', color: '#fff', cursor: 'pointer'}} aria-label="Delete Poll">PURGE</button>
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
                {votedId === opt.id && <span style={{fontSize: '20px'}}>🎯</span>}
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
          {poll.total_votes} SIGNALS RECORDED // {new Date(poll.created_at).toLocaleDateString()}
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
