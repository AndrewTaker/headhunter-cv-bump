package database

var tables = `
	create table if not exists users (
		id text primary key,
		first_name text,
		last_name text,
		middle_name text
	);

	create table if not exists tokens (
		access_token text,
		refresh_token text,
		expires_in integer,
		code text unique,
		user_id text unique,
		
		foreign key (user_id) references users(id) on delete cascade
	);

	create table if not exists resumes (
		id text primary key unique,
		alternate_url text,
		title text,
		created_at text,
		updated_at text,
		user_id text,
		is_scheduled integer not null default 0,

		foreign key (user_id) references users(id) on delete cascade
	);

	create table if not exists scheduler (
		user_id text,
		resume_id text,
		resume_title text,
		timestamp text,
		error text
	);
`
