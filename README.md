# Prerequisite
Prepare a `postgresql` table for testing
```
create table benchmark (
	id bigserial primary key,
	name varchar,
	address varchar,
	status varchar,
	created_on timestamp default current_timestamp,
	updated_on timestamp default current_timestamp
);
```
# Build program
### Get dependencies
```
dep init
dep ensure
```
### Build
```
go build
```

# Run Program
#### Benchmarking insert command with prepared statements
```
./benchmark -mode=insert -stmt=true
```
#### Benchmarking insert command without prepared statements
```
./benchmark -mode=insert -stmt=false
```
#### Benchmarking update command with prepared statements
```
./benchmark -mode=update -stmt=true
```
#### Benchmarking update command without prepared statements
```
./benchmark -mode=update -stmt=false
```

# Result
n-th attempt | #1 | #2 | #3 | #4 | #5 | #6 | #7 | #8 | #9 | #10 |
--- | --- | --- | --- |--- |--- |--- |--- |--- |--- |---
Insert with prepared statements (s) | 3.84 | 3.58 | 3.49 | 3.56 | 3.47 | 3.54 | 3.60 | 3.56 | 3.21 | 3.59 |
Insert withou prepared statements (s) | 4.01 | 3.94 | 3.97 | 3.82 | 4.10 | 4.00 | 4.06 | 3.95 | 4.09 | 3.24 |
Update with prepared statements (ms) | 325.82 | 330.90 | 324.23 | 313.86 | 334.92 | 328.92 | 328.90 | 316.56 | 329.58 | 369.39 |
Update without prepared statements (ms) | 496.42 | 443.11 | 438.18 | 454.77 | 447.60 | 484.58 | 441.03 | 444.77 | 433.84 | 441.09 |
