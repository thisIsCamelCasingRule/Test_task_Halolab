CREATE SEQUENCE my_serial AS integer START 1 OWNED BY sensor_data.id;
ALTER TABLE sensor_data ALTER COLUMN id SET DEFAULT nextval('my_serial');
