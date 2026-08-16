export type Post = {
  id: string;
  seq: number;
  authorUid: string;
  body: string;
  likesCount: number;
  commentsCount: number;
  likedByMe: boolean;
  createdAt: number;
};

export type PostDto = {
  id: string;
  seq: number;
  author_uid: string;
  body: string;
  likes_count: number;
  comments_count: number;
  liked_by_me: boolean;
  created_at: number;
};

export const toPost = (w: PostDto): Post => ({
  id: w.id,
  seq: w.seq,
  authorUid: w.author_uid,
  body: w.body,
  likesCount: w.likes_count,
  commentsCount: w.comments_count,
  likedByMe: w.liked_by_me,
  createdAt: w.created_at,
});
