import React, { useState } from 'react';
import { createPoll } from '../services/api';

const CreatePollForm = ({ onCreated }) => {
  const [q, setQ] = useState('');
  const [cat, setCat] = useState('Tech');
  const [opts, setOpts] = useState(['', '']);
  const [busy, setBusy] = useState(false);

  const save = async (e) => {
    e.preventDefault();
    const cleanOpts = opts.filter(o => o.trim());
    if (cleanOpts.length < 2) return alert("MINIMUM 2 CHOICES REQUIRED.");

    setBusy(true);
    try {
      const data = await createPoll({
        question: q.trim(),
        category: cat,
        options: cleanOpts
      });
      onCreated(data);
      setQ('');
      setOpts(['', '']);
    } catch (err) {
      alert("BROADCAST FAILED: " + err.message);
    } finally {
      setBusy(false);
    }
  };

  return (
    <div>
      <h2 style={{marginTop: 0, textTransform: 'uppercase', fontSize: '2rem', fontWeight: '900'}}>New Signal</h2>
      <form onSubmit={save}>
        <div style={{marginBottom: '20px'}}>
          <label style={{fontWeight: '900', fontSize: '12px'}}>QUESTION (5-280 CHARS)</label>
          <input
            type="text"
            className="form-input"
            value={q}
            onChange={e => setQ(e.target.value)}
            required
            minLength={5}
            maxLength={280}
            placeholder="What needs a decision?"
          />
        </div>

        <div style={{marginBottom: '20px'}}>
          <label style={{fontWeight: '900', fontSize: '12px'}}>CATEGORY</label>
          <select className="form-input" value={cat} onChange={e => setCat(e.target.value)}>
            <option>Tech</option>
            <option>Business</option>
            <option>Design</option>
          </select>
        </div>

        <div style={{marginBottom: '30px'}}>
          <label style={{fontWeight: '900', fontSize: '12px'}}>OPTIONS</label>
          {opts.map((o, i) => (
            <div key={i} style={{display: 'flex', gap: '8px', marginBottom: '10px'}}>
              <input
                type="text"
                className="form-input"
                style={{marginBottom: 0}}
                value={o}
                onChange={e => {
                  const n = [...opts];
                  n[i] = e.target.value;
                  setOpts(n);
                }}
                required
                maxLength={80}
                placeholder={`Option ${i+1}`}
              />
              {opts.length > 2 && (
                <button type="button" className="tag" onClick={() => setOpts(opts.filter((_, idx) => idx !== i))} style={{background: '#ff4d4d', cursor: 'pointer'}}>X</button>
              )}
            </div>
          ))}

          {opts.length < 10 && (
            <button type="button" className="tag" onClick={() => setOpts([...opts, ''])} style={{marginTop: '10px', cursor: 'pointer', background: '#fff'}}>+ ADD OPTION</button>
          )}
        </div>

        <button type="submit" className="btn-launch" disabled={busy}>
          {busy ? 'BROADCASTING...' : 'LAUNCH POLL'}
        </button>
      </form>
    </div>
  );
};

export default CreatePollForm;
