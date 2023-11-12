CREATE SEQUENCE my_serial2 AS integer START 1 OWNED BY sensor_to_sensor_data_connector.id;
ALTER TABLE sensor_to_sensor_data_connector ALTER COLUMN id SET DEFAULT nextval('my_serial2');
