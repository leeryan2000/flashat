import { useMemo, useState } from "react";
import { Rss } from "lucide-react";
import { usePosts } from "../context/PostContext";
import { useAuth } from "../context/AuthContext";
import { useChat } from "../context/ChatContext";
import PostCard, { type AuthorInfo } from "../components/PostCard";

export default function Posts() {
  const { posts, comments, isLoading, hasMore, loadMore, createPost, toggleLike, loadComments, addComment } =
    usePosts();
  const { user } = useAuth();
  const { friendships } = useChat();
  const [draft, setDraft] = useState("");

  // The posts service only knows author_uid (separate database, no
  // access to the users table) — since the feed is friends-only, every
  // author is either the viewer or an already-loaded friend, so this is
  // resolved client-side instead of a third gRPC call.
  const authorsByUid = useMemo(() => {
    const map: Record<string, AuthorInfo> = {};
    if (user) map[user.uid] = { name: user.name, avatarUrl: user.user_avatar_url };
    for (const f of friendships) {
      if (f.status === "accepted") map[f.uid] = { name: f.name, avatarUrl: f.avatarUrl };
    }
    return map;
  }, [user, friendships]);

  const submitPost = () => {
    const body = draft.trim();
    if (!body) return;
    createPost(body);
    setDraft("");
  };

  return (
    <div className="flex h-screen bg-black">
      <main className="flex-1 h-full overflow-y-auto" style={{ background: "var(--sidebar-bg)", color: "var(--text)" }}>
        <div className="max-w-2xl mx-auto p-6 flex flex-col gap-6">
          <h1 className="text-2xl font-bold flex items-center gap-3">
            <Rss className="w-6 h-6" style={{ color: "var(--primary)" }} />
            Posts
          </h1>

          {/* Composer */}
          <div
            className="p-4 rounded-xl border flex flex-col gap-3"
            style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
          >
            <textarea
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              placeholder="Share your daily stuff..."
              rows={3}
              className="w-full rounded-lg px-3 py-2 text-sm outline-none focus:ring-2 transition resize-none placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--surface-muted)", color: "var(--text)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
            <button
              onClick={submitPost}
              disabled={!draft.trim()}
              className="self-end px-4 py-2 rounded-xl font-medium text-sm text-white transition disabled:opacity-50"
              style={{ background: "var(--primary)" }}
            >
              Post
            </button>
          </div>

          {/* Feed */}
          <div className="flex flex-col gap-4">
            {posts.order.length === 0 && !isLoading ? (
              <div className="text-center py-10 italic" style={{ color: "var(--text-faint)" }}>
                No posts yet. Be the first to share something!
              </div>
            ) : (
              posts.order.map((id) => {
                const post = posts.entities[id];
                const author = authorsByUid[post.authorUid] ?? { name: "Unknown user", avatarUrl: null };
                const slice = comments[id];
                const postComments = slice ? slice.order.map((cid) => slice.entities[cid]) : [];
                return (
                  <PostCard
                    key={id}
                    post={post}
                    author={author}
                    comments={postComments}
                    commentAuthors={authorsByUid}
                    onLike={toggleLike}
                    onLoadComments={loadComments}
                    onAddComment={addComment}
                  />
                );
              })
            )}

            {hasMore && posts.order.length > 0 && (
              <button
                onClick={loadMore}
                disabled={isLoading}
                className="self-center text-sm font-medium transition disabled:opacity-50"
                style={{ color: "var(--primary)" }}
              >
                {isLoading ? "Loading..." : "Load more"}
              </button>
            )}
          </div>
        </div>
      </main>
    </div>
  );
}
