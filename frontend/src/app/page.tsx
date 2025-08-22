import { revalidatePath } from "next/cache";
import { headers } from "next/headers";

type Note = {
  id: number;
  content: string;
  createdAt: string;
};

async function getBaseUrl(): Promise<string | undefined> {
  const envBase = process.env.NEXT_PUBLIC_APP_URL;
  if (envBase) return envBase;
  const h = await headers();
  const host = h.get("x-forwarded-host") ?? h.get("host");
  const proto = h.get("x-forwarded-proto") ?? (process.env.NODE_ENV === "production" ? "https" : "http");
  return host ? `${proto}://${host}` : undefined;
}

async function getNotes(): Promise<Note[]> {
  const base = await getBaseUrl();
  if (!base) return [];

  try {
    const resp = await fetch(`${base}/api/notes`, { cache: "no-store" });
    if (!resp.ok) return [];
    return await resp.json();
  } catch {
    return [];
  }
}

async function addNote(formData: FormData) {
  "use server";
  const content = (formData.get("content") || "").toString().trim();
  if (!content) return;

  // Prefer routing through our API route to keep a single integration point
  const base = await getBaseUrl();
  if (!base) return;

  await fetch(`${base}/api/notes`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ content }),
  });

  revalidatePath("/");
}

export default async function Home() {
  const notes = await getNotes();

  return (
    <div className="min-h-screen bg-background text-foreground font-sans">
      <main className="mx-auto max-w-2xl p-6 sm:p-10 space-y-8">
        <header className="space-y-2">
          <h1 className="text-2xl font-semibold tracking-tight">Notes</h1>
          <p className="text-sm text-black/60 dark:text-white/60">
            バックエンドのノート一覧を表示し、フォームから追加します。
          </p>
        </header>

        <form action={addNote} className="space-y-3 rounded-lg border border-black/10 dark:border-white/10 p-4 bg-white/60 dark:bg-black/20 backdrop-blur">
          <label htmlFor="content" className="block text-sm font-medium">
            新しいノート
          </label>
          <textarea
            id="content"
            name="content"
            placeholder="メモを書いてください..."
            className="w-full rounded-md border border-black/10 dark:border-white/15 bg-white dark:bg-black/30 px-3 py-2 text-sm outline-none focus:ring-2 focus:ring-foreground/30 min-h-[84px]"
          />
          <div className="flex items-center justify-end">
            <button
              type="submit"
              className="inline-flex items-center rounded-md bg-foreground text-background px-4 py-2 text-sm font-medium shadow hover:opacity-90 active:translate-y-[1px]"
            >
              追加する
            </button>
          </div>
        </form>

        <section className="space-y-3">
          <h2 className="text-lg font-medium">最新のノート</h2>
          {notes.length === 0 ? (
            <p className="text-sm text-black/60 dark:text-white/60">ノートがありません。</p>
          ) : (
            <ul className="space-y-3">
              {notes.map((n) => (
                <li key={n.id} className="rounded-lg border border-black/10 dark:border-white/10 p-4 bg-white/60 dark:bg-black/20">
                  <p className="whitespace-pre-wrap text-sm">{n.content}</p>
                  <p className="mt-2 text-[12px] text-black/50 dark:text-white/50">
                    {new Date(n.createdAt).toLocaleString()}
                  </p>
                </li>
              ))}
            </ul>
          )}
        </section>
      </main>
    </div>
  );
}
