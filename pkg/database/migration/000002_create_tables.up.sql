create table if not exists sensor (
    id integer,
    group_name text,
    x float,
    y float,
    z float,
    primary key (id, group_name)
);

create table if not exists sensor_to_sensor_data_connector (
  id integer primary key ,
  sensor_id integer,
  group_name text,
  foreign key (sensor_id, group_name) references sensor(id, group_name)
);

create type fish_specie  as enum ('Atlantic Blue Tuna', 'Atlantic Cod', 'Blue Marlin', 'Coelacanth', 'Yellow Tuna');

create table if not exists sensor_data (
    id integer,
    sensor_id integer,
    temperature float,
    transparency integer,
    fish_specie fish_specie,
    fish_amount integer,
    timestamp timestamp,
    primary key (id),
    foreign key (sensor_id) references sensor_to_sensor_data_connector(id)
);