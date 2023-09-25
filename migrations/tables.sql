create table segment_history(
	id serial PRIMARY KEY
	user_id int
	segment_id int
	timestamp timestamp
	operation bool
)

create table users(
	user_id int
)

create table segments(
	id serial
	name varchar(100)
)

create table user_segments(
	user_id int
	segment_id int
	ttl date
)