ALTER TABLE citys
DROP INDEX name,
DROP COLUMN name,
ADD city_name VARCHAR(50) NOT NULL UNIQUE;