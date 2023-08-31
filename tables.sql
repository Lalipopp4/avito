create table segment_history(
	id serial PRIMARY KEY
	user_id int
	segment_name varchar(100)
	timestamp_req 
	operation bool
)

create table users(
	user_id int
)

create tablesegments(
	name varchar(100)
)

create table user_segments(
	user_id int
	segment_name varchar(100)
	ttl date
)