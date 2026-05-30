create extension if not exists "pgcrypto";

create table users (
  id bigserial primary key,
  google_sub varchar(255) not null unique,
  name varchar(255),
  is_actice boolean not null default true,
  created_at timestamptz not null default now(),
  updated_at timestamptz not null default now()
);

create table tasks (
    id uuid primary key default gen_random_uuid(),
    user_id bigint not null references users(id) on delete cascade,
    title varchar(200) not null,
    description text,
    is_completed boolean not null default false,
    created_at timestamptz not null default now(),
    updated_at timestamptz not null default now()
);
