import { useState } from "react";
import { Heart, MessageCircle } from "lucide-react";
import type { Post } from "../wire/post";
import type { Comment } from "../wire/comment";
import Avatar from "./Avatar";

export type AuthorInfo = { name: string; avatarUrl?: string | null };

const UNKNOWN_AUTHOR: AuthorInfo = { name: "Unknown user", avatarUrl: null };

type PostCardProps = {
  post: Post;
  author: AuthorInfo;
  comments: Comment[];
  commentAuthors: Record<string, AuthorInfo>;
  onLike: (postId: string) => void;
  onLoadComments: (postId: string) => void;
  onAddComment: (postId: string, body: string) => void;
};

export default function PostCard({
  post,
  author,
  comments,
  commentAuthors,
  onLike,
  onLoadComments,
  onAddComment,
}: PostCardProps) {
  const [commentsOpen, setCommentsOpen] = useState(false);
  const [draft, setDraft] = useState("");

  const toggleComments = () => {
    if (!commentsOpen) onLoadComments(post.id);
    setCommentsOpen((o) => !o);
  };

  const submitComment = () => {
    const body = draft.trim();
    if (!body) return;
    onAddComment(post.id, body);
    setDraft("");
  };

  return (
    <div
      className="p-4 rounded-xl border flex flex-col gap-3"
      style={{ background: "var(--sidebar-item)", borderColor: "var(--panel-border)" }}
    >
      <div className="flex items-center gap-3">
        <Avatar name={author.name} avatarUrl={author.avatarUrl} size="sm" />
        <div className="flex-1 min-w-0">
          <p className="font-medium truncate" style={{ color: "var(--text)" }}>
            {author.name}
          </p>
          <p className="text-xs" style={{ color: "var(--text-faint)" }}>
            {new Date(post.createdAt).toLocaleString()}
          </p>
        </div>
      </div>

      <p className="whitespace-pre-wrap" style={{ color: "var(--text)" }}>
        {post.body}
      </p>

      <div className="flex items-center gap-4 text-sm" style={{ color: "var(--text-soft)" }}>
        <button
          onClick={() => onLike(post.id)}
          className="flex items-center gap-1.5 transition"
          style={{ color: post.likedByMe ? "var(--danger-text)" : "var(--text-soft)" }}
        >
          <Heart size={16} fill={post.likedByMe ? "currentColor" : "none"} />
          {post.likesCount}
        </button>
        <button onClick={toggleComments} className="flex items-center gap-1.5 transition">
          <MessageCircle size={16} />
          {post.commentsCount}
        </button>
      </div>

      {commentsOpen && (
        <div className="flex flex-col gap-3 pt-2 border-t" style={{ borderColor: "var(--panel-border)" }}>
          {comments.map((c) => {
            const commentAuthor = commentAuthors[c.authorUid] ?? UNKNOWN_AUTHOR;
            return (
              <div key={c.id} className="flex items-start gap-2">
                <Avatar name={commentAuthor.name} avatarUrl={commentAuthor.avatarUrl} size="xs" />
                <div className="flex-1 min-w-0">
                  <p className="text-xs font-medium" style={{ color: "var(--text)" }}>
                    {commentAuthor.name}
                  </p>
                  <p className="text-sm whitespace-pre-wrap" style={{ color: "var(--text-soft)" }}>
                    {c.body}
                  </p>
                </div>
              </div>
            );
          })}

          <div className="flex items-center gap-2">
            <input
              value={draft}
              onChange={(e) => setDraft(e.target.value)}
              onKeyDown={(e) => e.key === "Enter" && submitComment()}
              placeholder="Write a comment..."
              className="flex-1 rounded-lg px-3 py-1.5 text-sm outline-none focus:ring-2 transition placeholder:text-[color:var(--text-faint)]"
              style={{ background: "var(--surface-muted)", color: "var(--text)", "--tw-ring-color": "var(--primary)" } as React.CSSProperties}
            />
            <button
              onClick={submitComment}
              className="px-3 py-1.5 rounded-lg text-sm font-medium text-white transition"
              style={{ background: "var(--primary)" }}
            >
              Post
            </button>
          </div>
        </div>
      )}
    </div>
  );
}
