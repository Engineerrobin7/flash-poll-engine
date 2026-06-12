import React from 'react';
import PollCard from '../components/PollCard';
import CreatePollForm from '../components/CreatePollForm';
import { usePolls } from '../hooks/usePolls';

const PollDashboard = () => {
  const { polls, loading, error, refresh, addPoll, updatePoll, removePoll } = usePolls();

  return (
    <div className="container">
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
        <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-end', marginBottom: '32px' }}>
          <div>
            <h2 style={{ margin: 0, textTransform: 'uppercase', fontSize: '2.5rem', fontWeight: '900' }}>Active</h2>
            <div style={{height: '6px', background: '#000', width: '60px'}}></div>
          </div>
          <button onClick={refresh} className="tag" style={{cursor: 'pointer'}}>RE-SYNC</button>
        </div>

        {loading && polls.length === 0 ? (
            <div style={{ textAlign: 'center', padding: '100px', fontWeight: '900', fontSize: '24px' }}>INITIALIZING...</div>
        ) : polls.length === 0 ? (
          <div style={{ textAlign: 'center', padding: '80px', border: '4px dashed #000', background: '#fff' }}>
            <h3 style={{ textTransform: 'uppercase', fontSize: '2rem' }}>Station Empty</h3>
            <p style={{fontWeight: '700'}}>Waiting for the first poll signal.</p>
          </div>
        ) : (
          <div style={{display: 'flex', flexDirection: 'column', gap: '20px'}}>
            {polls.map(p => (
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
        <div style={{fontWeight: '900', fontSize: '18px'}}>FLASH-POLL ENGINE v1.2</div>
        <div style={{fontSize: '12px', fontWeight: '700'}}>LE MICI ENGINEERING TECHNICAL ASSESSMENT // 2026</div>
      </footer>
    </div>
  );
};

export default PollDashboard;
