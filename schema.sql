CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date TEXT(8),
	title TEXT(100),
	comment TEXT(500),
	repeat TEXT(128)
);
CREATE INDEX scheduler_date_IDX ON scheduler (date);