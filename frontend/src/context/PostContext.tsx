import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { useAuth } from "./AuthContext";
import { api } from "../api/api";
import { toPost, type Post, type PostDto } from "../wire/post";
import { toComment, type Comment, type CommentDto } from "../wire/comment";
import { generateUUID } from "../utils/uuid";

// Mirrors ChatContext's normalized {order, entities} shape — order here
// is newest-first (seq desc), feed semantics, unlike message threads
// which read oldest-first.
export type PostState = {
  order: string[];
  entities: Record<string, Post>;
};

export type CommentSlice = {
  order: string[];
  entities: Record<string, Comment>;
};

interface PostContextType {
  posts: PostState;
  comments: Record<string, CommentSlice>; // postId -> CommentSlice
  isLoading: boolean;
  hasMore: boolean;
  loadMore: () => void;
  createPost: (body: string) => Promise<void>;
  toggleLike: (postId: string) => void;
  loadComments: (postId: string) => void;
  addComment: (postId: string, body: string) => void;
}

const PostContext = createContext<PostContextType | null>(null);

const PAGE_LIMIT = 20;

export function PostProvider({ children }: { children: React.ReactNode }) {
  const { user } = useAuth();
  const [posts, setPosts] = useState<PostState>({ order: [], entities: {} });
  const [comments, setComments] = useState<Record<string, CommentSlice>>({});
  const [isLoading, setIsLoading] = useState(false);
  const [hasMore, setHasMore] = useState(true);

  const mergePosts = useCallback((list: Post[]) => {
    setPosts((prev) => {
      const entities = { ...prev.entities };
      const newIds: string[] = [];
      for (const p of list) {
        if (!entities[p.id]) newIds.push(p.id);
        entities[p.id] = p;
      }
      const order = [...new Set([...prev.order, ...newIds])].sort(
        (a, b) => entities[b].seq - entities[a].seq
      );
      return { entities, order };
    });
  }, []);

  const loadMore = useCallback(() => {
    setIsLoading((currentlyLoading) => {
      if (currentlyLoading) return currentlyLoading;

      (async () => {
        try {
          const last = posts.order[posts.order.length - 1];
          const beforeSeq = last ? posts.entities[last]?.seq : undefined;
          const query =
            beforeSeq !== undefined
              ? `?limit=${PAGE_LIMIT}&before_seq=${beforeSeq}`
              : `?limit=${PAGE_LIMIT}`;
          const dtos = await api<PostDto[]>(`/posts/feed${query}`);
          const list = (dtos || []).map(toPost);
          if (list.length < PAGE_LIMIT) setHasMore(false);
          mergePosts(list);
        } catch (err) {
          console.error("Error loading feed:", err);
        } finally {
          setIsLoading(false);
        }
      })();

      return true;
    });
  }, [posts, mergePosts]);

  useEffect(() => {
    loadMore();
    // Only ever want the initial page on mount — loadMore captures
    // pagination cursor state itself for subsequent calls.
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, []);

  const createPost = useCallback(
    async (body: string) => {
      if (!user) return;
      const tempId = `pending-${generateUUID()}`;
      const optimisticPost: Post = {
        id: tempId,
        seq: Number.MAX_SAFE_INTEGER, // sorts to the top until replaced by the real seq
        authorUid: user.uid,
        body,
        likesCount: 0,
        commentsCount: 0,
        likedByMe: false,
        createdAt: Date.now(),
      };

      setPosts((prev) => ({
        entities: { ...prev.entities, [tempId]: optimisticPost },
        order: [tempId, ...prev.order],
      }));

      try {
        const dto = await api<PostDto>("/posts", {
          method: "POST",
          body: JSON.stringify({ body }),
        });
        const post = toPost(dto);
        setPosts((prev) => {
          const { [tempId]: _removed, ...restEntities } = prev.entities;
          return {
            entities: { ...restEntities, [post.id]: post },
            order: prev.order.map((id) => (id === tempId ? post.id : id)),
          };
        });
      } catch (err) {
        console.error("Error creating post:", err);
        setPosts((prev) => {
          const { [tempId]: _removed, ...restEntities } = prev.entities;
          return { entities: restEntities, order: prev.order.filter((id) => id !== tempId) };
        });
      }
    },
    [user]
  );

  const toggleLike = useCallback(
    (postId: string) => {
      const before = posts.entities[postId];
      if (!before) return;

      setPosts((prev) => {
        const post = prev.entities[postId];
        if (!post) return prev;
        const likedByMe = !post.likedByMe;
        const likesCount = post.likesCount + (likedByMe ? 1 : -1);
        return { ...prev, entities: { ...prev.entities, [postId]: { ...post, likedByMe, likesCount } } };
      });

      (async () => {
        try {
          const resp = await api<{ liked: boolean; likes_count: number }>(`/posts/${postId}/like`, {
            method: "POST",
          });
          setPosts((prev) => {
            const post = prev.entities[postId];
            if (!post) return prev;
            return {
              ...prev,
              entities: { ...prev.entities, [postId]: { ...post, likedByMe: resp.liked, likesCount: resp.likes_count } },
            };
          });
        } catch (err) {
          console.error("Error toggling like:", err);
          setPosts((prev) => ({ ...prev, entities: { ...prev.entities, [postId]: before } }));
        }
      })();
    },
    [posts]
  );

  const loadComments = useCallback((postId: string) => {
    (async () => {
      try {
        const dtos = await api<CommentDto[]>(`/posts/${postId}/comments?limit=100`);
        const list = (dtos || []).map(toComment);
        setComments((prev) => ({
          ...prev,
          [postId]: {
            order: list.map((c) => c.id),
            entities: Object.fromEntries(list.map((c) => [c.id, c])),
          },
        }));
      } catch (err) {
        console.error("Error loading comments:", err);
      }
    })();
  }, []);

  const addComment = useCallback(
    (postId: string, body: string) => {
      if (!user) return;
      const tempId = `pending-${generateUUID()}`;
      const optimisticComment: Comment = {
        id: tempId,
        postId,
        authorUid: user.uid,
        body,
        createdAt: Date.now(),
      };

      setComments((prev) => {
        const slice = prev[postId] ?? { order: [], entities: {} };
        return {
          ...prev,
          [postId]: { order: [...slice.order, tempId], entities: { ...slice.entities, [tempId]: optimisticComment } },
        };
      });
      setPosts((prev) => {
        const post = prev.entities[postId];
        if (!post) return prev;
        return { ...prev, entities: { ...prev.entities, [postId]: { ...post, commentsCount: post.commentsCount + 1 } } };
      });

      (async () => {
        try {
          const dto = await api<CommentDto>(`/posts/${postId}/comments`, {
            method: "POST",
            body: JSON.stringify({ body }),
          });
          const comment = toComment(dto);
          setComments((prev) => {
            const slice = prev[postId];
            if (!slice) return prev;
            const { [tempId]: _removed, ...restEntities } = slice.entities;
            return {
              ...prev,
              [postId]: {
                order: slice.order.map((id) => (id === tempId ? comment.id : id)),
                entities: { ...restEntities, [comment.id]: comment },
              },
            };
          });
        } catch (err) {
          console.error("Error adding comment:", err);
          setComments((prev) => {
            const slice = prev[postId];
            if (!slice) return prev;
            const { [tempId]: _removed, ...restEntities } = slice.entities;
            return { ...prev, [postId]: { order: slice.order.filter((id) => id !== tempId), entities: restEntities } };
          });
          setPosts((prev) => {
            const post = prev.entities[postId];
            if (!post) return prev;
            return {
              ...prev,
              entities: { ...prev.entities, [postId]: { ...post, commentsCount: Math.max(0, post.commentsCount - 1) } },
            };
          });
        }
      })();
    },
    [user]
  );

  const value = useMemo(
    () => ({ posts, comments, isLoading, hasMore, loadMore, createPost, toggleLike, loadComments, addComment }),
    [posts, comments, isLoading, hasMore, loadMore, createPost, toggleLike, loadComments, addComment]
  );

  return <PostContext.Provider value={value}>{children}</PostContext.Provider>;
}

export function usePosts() {
  const context = useContext(PostContext);
  if (!context) {
    throw new Error("usePosts must be used within a PostProvider");
  }
  return context;
}
