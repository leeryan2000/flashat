export type Comment = {
  id: string;
  postId: string;
  authorUid: string;
  body: string;
  createdAt: number;
};

export type CommentDto = {
  id: string;
  post_id: string;
  author_uid: string;
  body: string;
  created_at: number;
};

export const toComment = (w: CommentDto): Comment => ({
  id: w.id,
  postId: w.post_id,
  authorUid: w.author_uid,
  body: w.body,
  createdAt: w.created_at,
});
