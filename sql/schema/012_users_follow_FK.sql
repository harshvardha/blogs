-- +goose Up
alter table users_follow add constraint follower_FK foreign key(follower_id) references users(id) on delete cascade,
add constraint following_FK foreign key(following_id) references users(id) on delete cascade;

-- +goose Down
alter table users_follow drop constraint follower_FK, drop constraint following_FK;