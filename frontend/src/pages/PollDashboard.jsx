import React, { useState, useEffect } from 'react';
import PollCard from '../components/PollCard';
import CreatePollForm from '../components/CreatePollForm';
import { usePolls } from '../hooks/usePolls';
import { fetchStats } from '../services/api';

const PollDashboard = () => {
  const { polls, loading, error, refresh, addPoll, updatePoll, removePoll } = usePolls();
  const [filter, setFilter] = useState('ALL');
  const [stats, setStats] = useState({ total_polls: 0, total_votes: 0 });

  useEffect(() => {
    const loadStats = async () => {
      try {
        const data = await fetchStats();
        setStats(data);
      } catch (e) {
        console.error("Stats fetch failed:", e);
      }
    };
    loadStats();
    // Refresh stats every 10 seconds or when polls change
    const interval = setInterval(loadStats, 10000);
    return () => clearInterval(interval);
  }, [polls]);

  const filteredPolls = filter === 'ALL'
    ? polls
    : polls.filter(p => p.category.toUpperCase() === filter);

  const categories = ['ALL', 'TECH', 'BUSINESS', 'DESIGN'];

  return (
    <div className="container" style={{paddingTop: '100px'}}>
      <div className="ticker-wrap">
        <div className="ticker">
          <span style={{marginRight: '100px'}}>NETWORK STATUS: ONLINE // TOTAL POLLS: {stats.total_polls} // TOTAL SIGNALS RECORDED: {stats.total_votes} // LIVE FEED ACTIVE //</span>
          <span style={{marginRight: '100px'}}>NETWORK STATUS: ONLINE // TOTAL POLLS: {stats.total_polls} // TOTAL SIGNALS RECORDED: {stats.total_votes} // LIVE FEED ACTIVE //</span>
          <span>NETWORK STATUS: ONLINE // TOTAL POLLS: {stats.total_polls} // TOTAL SIGNALS RECORDED: {stats.total_votes} // LIVE FEED ACTIVE //</span>
        </div>
      </div>

      <header style={{textAlign: 'center', marginBottom: '60px'}}>
        <div className="live-indicator">
          <div className="pulse-dot"></div>
          ENGINE CONNECTED // LIVE
        </div>
        <h1>FlashPoll</h1>
        <p style={{fontWeight: '800', textTransform: 'uppercase', fontSize: '14px'}}>High-Velocity Decision Engine</p>
      </header>

      <section style={{background: '#fff', border: '4px solid #000', padding: '30px', boxShadow: '8px 8px 0px 0px #000', marginBottom: '80px'}}>
        <CreatePollForm onCreated={addPoll} />
      </section>

      <div style={{ marginTop: '40px' }}>
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', flexWrap: 'wrap', gap: '20px', marginBottom: '32px' }}>
          <div>
            <h2 style={{ margin: 0, textTransform: 'uppercase', fontSize: '2.5rem', fontWeight: '900' }}>Active</h2>
            <div style={{height: '6px', background: '#000', width: '60px'}}></div>
          </div>

          <div style={{display: 'flex', gap: '8px'}}>
            {categories.map(cat => (
              <button
                key={cat}
                onClick={() => setFilter(cat)}
                className="tag"
                style={{
                  cursor: 'pointer',
                  background: filter === cat ? '#000' : '#fff',
                  color: filter === cat ? '#fff' : '#000',
                  transition: '0.2s'
                }}
              >
                {cat}
              </button>
            ))}
          </div>
        </div>

        {loading && polls.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '100px', fontWeight: '900', fontSize: '24px' }}>INITIALIZING...</div>
        ) : filteredPolls.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '80px', border: '4px dashed #000', background: '#fff' }}>
            <h3 style={{ textTransform: 'uppercase', fontSize: '1.5rem' }}>No {filter !== 'ALL' ? filter : ''} Signals Found</h3>
            <p style={{fontWeight: '700'}}>Adjust your filters or launch a new poll.</p>
          </div>
        ) : (
          <div style={{display: 'flex', flexDirection: 'column', gap: '20px'}}>
            {filteredPolls.map(p => (
              <PollCard
                key={p.id}
                poll={p}
                onUpdate={updatePoll}
                onDelete={removePoll}
              />
            ))}
          </div>
        )}
      </div>

      <footer style={{ marginTop: '120px', textAlign: 'center', padding: '40px 0', borderTop: '4px solid #000' }}>
        <div style={{fontWeight: '900', fontSize: '18px'}}>FLASH-POLL ENGINE v1.5</div>
        <div style={{fontSize: '12px', fontWeight: '700'}}>LE MICI ENGINEERING TECHNICAL ASSESSMENT // 2026</div>
      </footer>
    </div>
  );
};

export default PollDashboard;
