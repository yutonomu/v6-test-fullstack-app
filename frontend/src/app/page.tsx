'use client';
import { useEffect, useState } from 'react';

type Note = {
  id: number;
  content: string;
  createdAt: string;
};

export default function Home() {
  const [content, setContent] = useState('');
  const [status, setStatus] = useState<string | null>(null);
  const [notes, setNotes] = useState<Note[]>([]);

  async function fetchNotes() {
    try {
      const res = await fetch('/api/notes');
      const data = await res.json();
      if (!res.ok) throw new Error(data?.error || '取得に失敗しました');
      setNotes(data as Note[]);
    } catch (e: any) {
      setStatus(e.message || '取得に失敗しました');
    }
  }

  useEffect(() => {
    fetchNotes();
  }, []);

  async function onSubmit(e: React.FormEvent) {
    e.preventDefault();
    setStatus('送信中...');
    try {
      const res = await fetch('/api/notes', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ content }),
      });
      const data = await res.json();
      if (!res.ok) throw new Error(data?.error || '送信失敗');
      setStatus(`保存しました ID: ${data.id}`);
      setContent('');
      // 最新の一覧を取得して表示
      fetchNotes();
    } catch (err: any) {
      setStatus(err.message);
    }
  }

  return (
    <main style={{ padding: 24 }}>
      <h1>超シンプル入力 → Go → Postgres</h1>
      <form onSubmit={onSubmit} style={{ marginTop: 12 }}>
        <input
          value={content}
          onChange={(e) => setContent(e.target.value)}
          placeholder="メモ内容（content）"
          style={{ padding: 8, width: 320 }}
        />
        <button type="submit" style={{ marginLeft: 8 }}>保存</button>
      </form>
      {status && <p style={{ marginTop: 12 }}>{status}</p>}

      <section style={{ marginTop: 24 }}>
        <h2>保存済みノート（最新順）</h2>
        {notes.length === 0 ? (
          <p style={{ color: '#666' }}>まだデータがありません</p>
        ) : (
          <ul style={{ marginTop: 8, paddingLeft: 16 }}>
            {notes.map((n) => (
              <li key={n.id} style={{ marginBottom: 6 }}>
                <span>{n.content}</span>
                <small style={{ marginLeft: 8, color: '#666' }}>
                  {new Date(n.createdAt).toLocaleString()}
                </small>
              </li>
            ))}
          </ul>
        )}
      </section>
      <p style={{ marginTop: 24, color: '#666' }}>
        フロントは <code>/api/notes</code> にPOSTします（NextのrewriteでGoサーバへプロキシ）。
      </p>
    </main>
  );
}
