ALTER TABLE movies ADD CONSTRAINT movies_runtime_check check (runtime >=0);

ALTER TABLE movies ADD CONSTRAINT movies_year_check check (year between 1888 and date_part('year', now()));

ALTER TABLE movies ADD CONSTRAINT genres_length_check check (array_length(genres,1) between 1 and 5 );