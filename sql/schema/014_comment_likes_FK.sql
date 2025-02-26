-- +goose Up
alter table comment_likes add constraint userID_FK foreign key(user_id) references users(id) on delete cascade,
add constraint commentID_FK foreign key(comment_id) references comments(id) on delete cascade;

-- +goose Down
alter table comment_likes drop constraint userID_FK, drop constraint commentID_FK;