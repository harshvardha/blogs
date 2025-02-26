-- +goose Up
alter table likes add constraint userID_FK foreign key(user_id) references users(id) on delete cascade,
add constraint blogID_FK foreign key(blog_id) references blogs(id) on delete cascade;

-- +goose Down
alter table likes drop constraint userID_FK, drop constraint blogID_FK;