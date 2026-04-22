create table users (
  id bigserial primary key,
  google_sub varchar(255) not null unique,
  name varchar(255),
  is_actice boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table tasks (
    id bigserial primary key,
    user_id bigdecimal references users(id),
    content varchar(256) not null,
    is_completed boolean not null default false
)