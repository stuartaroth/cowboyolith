create table users (
   id uuid primary key,
   email text not null,
   is_admin boolean default false,
   inserted_at timestamp default (timezone('utc', now())),
   unique (email)
);

create table pending_user_sessions(
    id uuid primary key,
    user_id uuid not null references users,
    hashed_cookie_token text not null,
    salt text not null,
    ip_address text not null,
    user_agent text not null,
    inserted_at timestamp default (timezone('utc', now()))
);


create table user_sessions (
    id uuid primary key,
    user_id uuid not null references users,
    hashed_cookie_token text not null,
    salt text not null,
    user_agent text not null,
    inserted_at timestamp default (timezone('utc', now()))
);

create table orgs (
    id uuid primary key,
    name text not null,
    inserted_at timestamp default (timezone('utc', now()))
);

create table org_users (
    org_id uuid not null references orgs,
    user_id uuid not null references users,
    unique (org_id, user_id)
);

create table org_groups (
    id uuid primary key,
    org_id uuid not null references orgs,
    name text not null,
    inserted_at timestamp default (timezone('utc', now())),
    unique (org_id, name)
);

create table org_group_users (
    org_group_id uuid not null references org_groups,
    user_id uuid not null references users,
    unique (org_group_id, user_id)
);

create table org_roles_and_permissions (
    id uuid primary key,
    name text not null,
    can_mutate_groups boolean default false,
    can_view_all_users boolean default false,
    can_view_all_groups boolean default false,
    unique (name)
);

create table org_group_roles_and_permissions (
    id uuid primary key,
    name text not null,
    can_update_group boolean default false,
    can_invite_to_group boolean default false,
    unique (name)
);
